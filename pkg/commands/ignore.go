package commands

import (
	"errors"
	"fmt"

	"github.com/bearer/bearer/pkg/flag"
	"github.com/bearer/bearer/pkg/util/file"
	"github.com/bearer/bearer/pkg/util/ignore"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

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
	return &cobra.Command{
		Use:   "show <fingerprint>",
		Short: "Show an ignored fingerprint",
		Example: `# Show the details of an ignored fingerprint from your bearer.ignore file
$ bearer ignore show <fingerprint>`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return cmd.Help()
			}

			ignoredFingerprints, err := ignore.GetIgnoredFingerprints(nil)
			if err != nil {
				cmd.Printf("Issue loading ignored fingerprints from bearer.ignore file: %s", err)
				return nil
			}
			fingerprintId := args[0]
			selectedIgnoredFingerprint, ok := ignoredFingerprints[fingerprintId]
			if !ok {
				cmd.Printf("Ignored fingerprint '%s' was not found in bearer.ignore file\n", fingerprintId)
				return nil
			}
			cmd.Print("\n")
			cmd.Print(ignore.DisplayIgnoredEntryTextString(fingerprintId, selectedIgnoredFingerprint))
			cmd.Print("\n")
			return nil
		},
		SilenceErrors: false,
		SilenceUsage:  false,
	}
}

func newIgnoreAddCommand() *cobra.Command {
	var IgnoreAddFlags = &flag.Flags{
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

			if err = ignore.AddToIgnoreFile(fingerprintsToIgnore, options.IgnoreAddOptions.Force); err != nil {
				target := &ignore.DuplicateIgnoredFingerprintError{}
				if errors.As(err, &target) {
					// handle expected error (duplicate entry in bearer.ignore)
					cmd.Printf("Error: %s\n", err.Error())
					return nil
				}
				return err
			}

			cmd.Print("fingerprint added to bearer.ignore:\n\n")
			cmd.Print(ignore.DisplayIgnoredEntryTextString(fingerprintId, fingerprintEntry))
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

			ignoredFingerprintsFromConfig, err := getIgnoredFingerprintsFromConfig(args)
			if err != nil {
				// handle expected error (duplicate entry in bearer.ignore)
				cmd.Printf("%s\n", err.Error())
				return nil
			}

			if err = ignore.AddToIgnoreFile(ignoredFingerprintsFromConfig, options.IgnoreMigrateOptions.Force); err != nil {
				target := &ignore.DuplicateIgnoredFingerprintError{}
				if errors.As(err, &target) {
					// handle expected error (duplicate entry in bearer.ignore)
					cmd.Printf("%s\n", err.Error())
					return nil
				}
				return err
			}

			return nil
		},
		SilenceErrors: false,
		SilenceUsage:  false,
	}
	IgnoreMigrateFlags.AddFlags(cmd)
	cmd.SetUsageTemplate(fmt.Sprintf(scanTemplate, IgnoreMigrateFlags.Usages(cmd)))

	return cmd
}

func getIgnoredFingerprintsFromConfig(args []string) (ignoredFingerprintsFromConfig map[string]ignore.IgnoredFingerprint, err error) {
	ignoredFingerprintsFromConfig = make(map[string]ignore.IgnoredFingerprint)

	configPath := viper.GetString(flag.ConfigFileFlag.ConfigName)
	var defaultConfigPath = ""
	if len(args) > 0 {
		defaultConfigPath = file.GetFullFilename(args[0], configPath)
	}

	if err := readConfig(configPath); err != nil {
		if err := readConfig(defaultConfigPath); err != nil {
			return ignoredFingerprintsFromConfig, err
		}
		return ignoredFingerprintsFromConfig, err
	}

	for _, fingerprint := range viper.GetStringSlice("report.exclude-fingerprint") {
		ignoredFingerprintsFromConfig[fingerprint] = ignore.IgnoredFingerprint{}
	}

	return ignoredFingerprintsFromConfig, nil
}
