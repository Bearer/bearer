package commands

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/bearer/bearer/pkg/flag"
	"github.com/bearer/bearer/pkg/util/ignore"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var bold = color.New(color.Bold).SprintFunc()
var morePrefix = color.HiBlackString("├─ ")
var lastPrefix = color.HiBlackString("└─ ")

func NewIgnoreCommand() *cobra.Command {
	usageTemplate := `
Usage: bearer ignore <command> [flags]

Available Commands:
    add              Add an ignored fingerprint
    show             Show an ignored fingerprint
    migrate          Migrate ignored fingerprints

Examples:
    # Add an ignored fingerprint to your bearer.ignore file
    $ bearer ignore add <fingerprint> --author Mish --comment "investigate this"

    # Show the details of an ignored fingerprint from your bearer.ignore file
    $ bearer ignore show <fingerprint>

    # Migrate existing ignored (excluded) fingerprints from bearer.yml file
    # to bearer.ignore
    $ bearer ignore migrate

`

	cmd := &cobra.Command{
		Use:           "ignore [subcommand] <fingerprint>",
		Short:         "Manage ignored fingerprints",
		Args:          cobra.NoArgs,
		SilenceErrors: false,
		SilenceUsage:  false,
	}

	cmd.AddCommand(
		newIgnoreShowCommand(),
		newIgnoreAddCommand(),
		newIgnoreMigrateCommand(),
	)

	cmd.SetUsageTemplate(usageTemplate)

	return cmd
}

func newIgnoreShowCommand() *cobra.Command {
	var IgnoreShowFlags = &flag.Flags{
		IgnoreFlagGroup: flag.NewIgnoreFlagGroup(),
	}
	cmd := &cobra.Command{
		Use:   "show <fingerprint>",
		Short: "Show an ignored fingerprint",
		Example: `# Show the details of an ignored fingerprint from your bearer.ignore file
$ bearer ignore show <fingerprint>`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := IgnoreShowFlags.Bind(cmd); err != nil {
				return fmt.Errorf("flag bind error: %w", err)
			}

			if len(args) == 0 {
				return cmd.Help()
			}

			options, err := IgnoreShowFlags.ToOptions(args)
			if err != nil {
				return fmt.Errorf("flag error: %s", err)
			}

			ignoredFingerprints, fileExists, err := ignore.GetIgnoredFingerprints(options.IgnoreOptions.BearerIgnoreFile)
			if err != nil {
				cmd.Printf("Issue loading ignored fingerprints from bearer.ignore file: %s", err)
				return nil
			}
			if !fileExists {
				cmd.Printf("bearer.ignore file not found. Perhaps you need to use --bearer-ignore-file to specify the path to bearer.ignore?\n")
				return nil
			}
			fingerprintId := args[0]
			selectedIgnoredFingerprint, ok := ignoredFingerprints[fingerprintId]
			if !ok {
				cmd.Printf("Ignored fingerprint '%s' was not found in bearer.ignore file\n", fingerprintId)
				return nil
			}
			cmd.Print("\n")
			cmd.Print(displayIgnoredEntryTextString(fingerprintId, selectedIgnoredFingerprint))
			cmd.Print("\n")
			return nil
		},
		SilenceErrors: false,
		SilenceUsage:  false,
	}
	IgnoreShowFlags.AddFlags(cmd)
	cmd.SetUsageTemplate(fmt.Sprintf(scanTemplate, IgnoreShowFlags.Usages(cmd)))

	return cmd
}

func newIgnoreAddCommand() *cobra.Command {
	var IgnoreAddFlags = &flag.Flags{
		IgnoreFlagGroup:    flag.NewIgnoreFlagGroup(),
		IgnoreAddFlagGroup: flag.NewIgnoreAddFlagGroup(),
	}
	cmd := &cobra.Command{
		Use:   "add <fingerprint>",
		Short: "Add an ignored fingerprint",
		Example: `# Add an ignored fingerprint to your bearer.ignore file
$ bearer ignore add <fingerprint> --author Mish --comment "Possible false positive"`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := IgnoreAddFlags.Bind(cmd); err != nil {
				return fmt.Errorf("flag bind error: %w", err)
			}

			if len(args) == 0 {
				return cmd.Help()
			}

			options, err := IgnoreAddFlags.ToOptions(args)
			if err != nil {
				return fmt.Errorf("flag error: %s", err)
			}
			fingerprintId := args[0]
			fingerprintEntry := ignore.IgnoredFingerprint{
				Author:  options.IgnoreAddOptions.Author,
				Comment: options.IgnoreAddOptions.Comment,
			}
			fingerprintsToIgnore := map[string]ignore.IgnoredFingerprint{
				fingerprintId: fingerprintEntry,
			}

			ignoredFingerprints, fileExists, err := ignore.GetIgnoredFingerprints(options.IgnoreOptions.BearerIgnoreFile)
			if err != nil {
				return fmt.Errorf("error retrieving existing ignores: %s", err)
			}

			if !fileExists {
				cmd.Printf("\nCreating bearer.ignore file...\n")
			}

			if mergeErr := ignore.MergeIgnoredFingerprints(fingerprintsToIgnore, ignoredFingerprints, options.IgnoreAddOptions.Force); mergeErr != nil {
				// handle expected error (duplicate entry in bearer.ignore)
				cmd.Printf("Error: %s\n", mergeErr.Error())
				return nil
			}

			if err := writeIgnoreFile(ignoredFingerprints, options.IgnoreOptions.BearerIgnoreFile); err != nil {
				return err
			}

			cmd.Print("Fingerprint added to bearer.ignore:\n\n")
			cmd.Print(displayIgnoredEntryTextString(fingerprintId, ignoredFingerprints[fingerprintId]))
			cmd.Print("\n")
			return nil
		},
		SilenceErrors: false,
		SilenceUsage:  false,
	}
	IgnoreAddFlags.AddFlags(cmd)
	cmd.SetUsageTemplate(fmt.Sprintf(scanTemplate, IgnoreAddFlags.Usages(cmd)))

	return cmd
}

func newIgnoreMigrateCommand() *cobra.Command {
	IgnoreMigrateFlags := &flag.Flags{
		IgnoreFlagGroup:        flag.NewIgnoreFlagGroup(),
		IgnoreMigrateFlagGroup: flag.NewIgnoreMigrateFlagGroup(),
	}
	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "Migrate ignored fingerprints from bearer.yml to bearer.ignore",
		Example: `# Migrate existing ignored (excluded) fingerprints from bearer.yml file to bearer.ignore
$ bearer ignore migrate`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := IgnoreMigrateFlags.Bind(cmd); err != nil {
				return fmt.Errorf("flag bind error: %w", err)
			}

			options, err := IgnoreMigrateFlags.ToOptions(args)
			if err != nil {
				return fmt.Errorf("flag error: %s", err)
			}
			configFilePath := viper.GetString(flag.ConfigFileFlag.ConfigName)
			fingerprintsToMigrate, err := getIgnoredFingerprintsFromConfig(configFilePath)
			if err != nil {
				return fmt.Errorf("error reading config: %s\nPerhaps you need to use --config-file to specify the config path?", err.Error())
			}

			ignoredFingerprints, fileExists, err := ignore.GetIgnoredFingerprints(options.IgnoreOptions.BearerIgnoreFile)
			if err != nil {
				return fmt.Errorf("error retrieving existing ignores: %s", err)
			}

			migratedIgnoredCount := len(fingerprintsToMigrate)
			skippedIgnoresToMigrate := ""
			cmd.Printf("Found %d ignores in:\n\t%s\n", migratedIgnoredCount, configFilePath)

			if !fileExists {
				cmd.Printf("\nCreating bearer.ignore file...\n")
			}

			if !options.IgnoreMigrateOptions.Force {
				for key := range ignoredFingerprints {
					if _, ok := fingerprintsToMigrate[key]; ok {
						migratedIgnoredCount--
						skippedIgnoresToMigrate += fmt.Sprintf("- %s\n", key)
						delete(fingerprintsToMigrate, key)
					}
				}
			}

			cmd.Printf("Added %d ignores to:\n\t%s\n", migratedIgnoredCount, options.IgnoreOptions.BearerIgnoreFile)

			if skippedIgnoresToMigrate != "" {
				cmd.Printf("\nThe following ignores already exist in the bearer.ignore file:\n")
				cmd.Printf(skippedIgnoresToMigrate)
				cmd.Printf("\nTo overwrite these entries, use --force\n")
			}

			// either no duplicate entries at this point or --force is true so we can ignore merge error
			_ = ignore.MergeIgnoredFingerprints(fingerprintsToMigrate, ignoredFingerprints, options.IgnoreMigrateOptions.Force)
			return writeIgnoreFile(ignoredFingerprints, options.IgnoreOptions.BearerIgnoreFile)
		},
		SilenceErrors: false,
		SilenceUsage:  false,
	}
	IgnoreMigrateFlags.AddFlags(cmd)
	cmd.SetUsageTemplate(fmt.Sprintf(scanTemplate, IgnoreMigrateFlags.Usages(cmd)))

	return cmd
}

func writeIgnoreFile(ignoredFingerprints map[string]ignore.IgnoredFingerprint, bearerIgnoreFilePath string) error {
	data, err := json.MarshalIndent(ignoredFingerprints, "", "  ")
	if err != nil {
		// failed to marshall data
		return err
	}

	return os.WriteFile(bearerIgnoreFilePath, data, 0644)
}

func getIgnoredFingerprintsFromConfig(configPath string) (ignoredFingerprintsFromConfig map[string]ignore.IgnoredFingerprint, err error) {
	ignoredFingerprintsFromConfig = make(map[string]ignore.IgnoredFingerprint)

	if err := readConfig(configPath); err != nil {
		return ignoredFingerprintsFromConfig, err
	}

	for _, fingerprint := range viper.GetStringSlice("report.exclude-fingerprint") {
		ignoredFingerprintsFromConfig[fingerprint] = ignore.IgnoredFingerprint{
			Comment: "migrated from bearer.yml",
		}
	}

	return ignoredFingerprintsFromConfig, nil
}

func displayIgnoredEntryTextString(fingerprintId string, entry ignore.IgnoredFingerprint) string {
	prefix := morePrefix
	result := fmt.Sprintf(bold(color.HiBlueString("%s \n")), fingerprintId)

	if entry.Author == "" && entry.Comment == "" {
		prefix = lastPrefix
	}
	result += fmt.Sprintf("%sIgnored At: %s\n", prefix, bold(entry.IgnoredAt))

	if entry.Author != "" {
		if entry.Comment == "" {
			prefix = lastPrefix
		}

		result += fmt.Sprintf("%sAuthor: %s\n", prefix, bold(entry.Author))
	}

	if entry.Comment != "" {
		result += fmt.Sprintf("%sComment: %s\n", lastPrefix, bold(entry.Comment))
	}

	return result
}
