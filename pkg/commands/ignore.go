package commands

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/bearer/bearer/pkg/flag"
	"github.com/bearer/bearer/pkg/util/ignore"
	ignoretypes "github.com/bearer/bearer/pkg/util/ignore/types"
	"github.com/bearer/bearer/pkg/util/output"
)

var migratedIgnoreComment = "migrated from bearer.yml"

func NewIgnoreCommand() *cobra.Command {
	usageTemplate := `
Usage: bearer ignore <command> [flags]

Available Commands:
    add              Add an ignored fingerprint
    show             Show an ignored fingerprint
    remove           Remove an ignored fingerprint
    migrate          Migrate ignored fingerprints

Examples:
    # Add an ignored fingerprint to your ignore file
    $ bearer ignore add <fingerprint> --author Mish --comment "investigate this"

    # Show the details of an ignored fingerprint from your ignore file
    $ bearer ignore show <fingerprint>

    # Remove an ignored fingerprint from your ignore file
    $ bearer ignore remove <fingerprint>

    # Migrate existing ignored (excluded) fingerprints from bearer.yml file
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
		newIgnoreRemoveCommand(),
		newIgnoreMigrateCommand(),
	)

	cmd.SetUsageTemplate(usageTemplate)

	return cmd
}

func newIgnoreShowCommand() *cobra.Command {
	var IgnoreShowFlags = flag.Flags{
		flag.GeneralFlagGroup,
		flag.IgnoreShowFlagGroup,
	}
	cmd := &cobra.Command{
		Use:   "show <fingerprint>",
		Short: "Show an ignored fingerprint",
		Example: `# Show the details of an ignored fingerprint from your ignore file
$ bearer ignore show <fingerprint>`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := IgnoreShowFlags.Bind(cmd); err != nil {
				return fmt.Errorf("flag bind error: %w", err)
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			setLogLevel(cmd)

			options, err := IgnoreShowFlags.ToOptions(args)
			if err != nil {
				return fmt.Errorf("flag error: %s", err)
			}

			if len(args) == 0 && !options.All {
				return cmd.Help()
			}

			ignoredFingerprints, ignoreFilepath, fileExists, err := ignore.GetIgnoredFingerprints(options.IgnoreFile, nil)
			if err != nil {
				cmd.Printf("Issue loading ignored fingerprints from %s: %s", err, ignoreFilepath)
				return nil
			}
			if !fileExists {
				cmd.Printf("Ignore file not found. Perhaps you need to use --ignore-file to specify the path to ignore?\n")
				return nil
			}

			cmd.Print("\n")
			if options.All {
				// show all fingerprints sorted by date
				keys := make([]string, 0, len(ignoredFingerprints))
				for key := range ignoredFingerprints {
					keys = append(keys, key)
				}

				sort.SliceStable(keys, func(i, j int) bool {
					return ignoredFingerprints[keys[i]].IgnoredAt < ignoredFingerprints[keys[j]].IgnoredAt
				})

				for _, k := range keys {
					cmd.Print(ignore.DisplayIgnoredEntryTextString(k, ignoredFingerprints[k], options.NoColor))
					cmd.Print("\n\n")
				}
			} else {
				// show a specific fingerprint
				fingerprintId := args[0]
				selectedIgnoredFingerprint, ok := ignoredFingerprints[fingerprintId]
				if !ok {
					cmd.Printf("Ignored fingerprint '%s' was not found in ignore file\n", fingerprintId)
					return nil
				}
				cmd.Print(ignore.DisplayIgnoredEntryTextString(fingerprintId, selectedIgnoredFingerprint, options.NoColor))
			}
			cmd.Print("\n\n")
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
	IgnoreShowFlags := flag.Flags{
		flag.IgnoreAddFlagGroup,
		flag.GeneralFlagGroup,
	}

	cmd := &cobra.Command{
		Use:   "add <fingerprint>",
		Short: "Add an ignored fingerprint",
		Example: `# Add an ignored fingerprint to your ignore file
$ bearer ignore add <fingerprint> --author Mish --comment "Possible false positive"`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := IgnoreShowFlags.Bind(cmd); err != nil {
				return fmt.Errorf("flag bind error: %w", err)
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return cmd.Help()
			}

			setLogLevel(cmd)

			options, err := IgnoreShowFlags.ToOptions(args)
			if err != nil {
				return fmt.Errorf("flag error: %s", err)
			}

			// create initial entry
			fingerprintId := args[0]
			var fingerprintEntry ignoretypes.IgnoredFingerprint
			fingerprintsToIgnore := map[string]ignoretypes.IgnoredFingerprint{
				fingerprintId: fingerprintEntry,
			}

			ignoredFingerprints, ignoreFilepath, fileExists, err := ignore.GetIgnoredFingerprints(options.IgnoreFile, nil)
			if err != nil {
				return fmt.Errorf("error retrieving existing ignores: %s", err)
			}

			// check for merge conflicts
			if mergeErr := ignore.MergeIgnoredFingerprints(fingerprintsToIgnore, ignoredFingerprints, options.IgnoreAddOptions.Force); mergeErr != nil {
				// handle expected error (duplicate entry in ignore)
				cmd.Printf("Error: %s\n", mergeErr.Error())
				return nil
			}

			// ensure ignored at is set
			fingerprintEntry = ignoredFingerprints[fingerprintId]

			// add additional information to entry
			if options.Author != "" {
				fingerprintEntry.Author = &options.Author
			} else {
				if author, err := ignore.GetAuthor(); err == nil {
					fingerprintEntry.Author = author
				}
			}
			if options.FalsePositive {
				fingerprintEntry.FalsePositive = options.FalsePositive
			} else {
				fingerprintEntry.FalsePositive = requestConfirmation("Is this finding a false positive?")
				cmd.Printf("\n")
			}
			if options.Comment != "" {
				fingerprintEntry.Comment = &options.Comment
			} else {
				reader := bufio.NewReader(os.Stdin)
				fmt.Print("Add a comment or press enter to continue: ")
				input, _ := reader.ReadString('\n')
				comment := strings.TrimSuffix(input, "\n")
				if comment != "" {
					fingerprintEntry.Comment = &comment
				}
				cmd.Printf("\n")
			}

			// update entry to include additional information
			ignoredFingerprints[fingerprintId] = fingerprintEntry

			if !fileExists {
				cmd.Printf("\nCreating ignore file...\n")
			}

			if err := writeIgnoreFile(ignoredFingerprints, ignoreFilepath); err != nil {
				return err
			}

			cmd.Print("Fingerprint added to ignore file:\n\n")
			cmd.Print(ignore.DisplayIgnoredEntryTextString(fingerprintId, ignoredFingerprints[fingerprintId], options.NoColor))
			cmd.Print("\n\n")
			return nil
		},
		SilenceErrors: false,
		SilenceUsage:  false,
	}
	IgnoreShowFlags.AddFlags(cmd)
	cmd.SetUsageTemplate(fmt.Sprintf(scanTemplate, IgnoreShowFlags.Usages(cmd)))

	return cmd
}

func newIgnoreRemoveCommand() *cobra.Command {
	flags := flag.Flags{flag.GeneralFlagGroup}

	cmd := &cobra.Command{
		Use:   "remove <fingerprint>",
		Short: "Remove an ignored fingerprint",
		Example: `# Remove an ignored fingerprint from your ignore file
$ bearer ignore remove <fingerprint>`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := flags.Bind(cmd); err != nil {
				return fmt.Errorf("flag bind error: %w", err)
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return cmd.Help()
			}

			setLogLevel(cmd)

			options, err := flags.ToOptions(args)
			if err != nil {
				return fmt.Errorf("flag error: %s", err)
			}

			ignoredFingerprints, ignoreFilepath, fileExists, err := ignore.GetIgnoredFingerprints(options.IgnoreFile, nil)
			if err != nil {
				return fmt.Errorf("error retrieving existing ignores: %s", err)
			}
			if !fileExists {
				cmd.Printf("Ignore file not found. Perhaps you need to use --ignore-file to specify the path?\n")
				return nil
			}

			fingerprintId := args[0]
			removedFingerprint, ok := ignoredFingerprints[fingerprintId]
			if !ok {
				cmd.Printf("Ignored fingerprint '%s' was not found in ignore file\n", fingerprintId)
				return nil
			}

			delete(ignoredFingerprints, fingerprintId)
			if err := writeIgnoreFile(ignoredFingerprints, ignoreFilepath); err != nil {
				return err
			}

			cmd.Print("Fingerprint successfully removed from ignore file:\n\n")
			cmd.Print(ignore.DisplayIgnoredEntryTextString(fingerprintId, removedFingerprint, options.NoColor))
			cmd.Print("\n\n")
			return nil
		},
		SilenceErrors: false,
		SilenceUsage:  false,
	}

	flags.AddFlags(cmd)
	cmd.SetUsageTemplate(fmt.Sprintf(scanTemplate, flags.Usages(cmd)))

	return cmd
}

func newIgnoreMigrateCommand() *cobra.Command {
	flags := flag.Flags{
		flag.GeneralFlagGroup,
		flag.IgnoreMigrateFlagGroup,
	}
	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "Migrate ignored fingerprints from bearer.yml to ignore file",
		Example: `# Migrate existing ignored (excluded) fingerprints from bearer.yml file to ignore file
$ bearer ignore migrate`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := flags.Bind(cmd); err != nil {
				return fmt.Errorf("flag bind error: %w", err)
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			setLogLevel(cmd)

			options, err := flags.ToOptions(args)
			if err != nil {
				return fmt.Errorf("flag error: %s", err)
			}

			configFilePath, _, err := readConfig(args)
			if err != nil {
				return fmt.Errorf("error reading config: %s\nPerhaps you need to use --config-file to specify the config path?", err.Error())
			}
			fingerprintsToMigrate := getIgnoredFingerprintsFromConfig(configFilePath)

			ignoredFingerprints, ignoreFilepath, fileExists, err := ignore.GetIgnoredFingerprints(options.IgnoreFile, nil)
			if err != nil {
				return fmt.Errorf("error retrieving existing ignores: %s", err)
			}

			migratedIgnoredCount := len(fingerprintsToMigrate)
			skippedIgnoresToMigrate := ""
			cmd.Printf("Found %d ignores in:\n\t%s\n", migratedIgnoredCount, configFilePath)

			if !fileExists {
				cmd.Printf("\nCreating ignore file...\n")
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

			cmd.Printf("Added %d ignores to:\n\t%s\n", migratedIgnoredCount, ignoreFilepath)

			if skippedIgnoresToMigrate != "" {
				cmd.Printf("\nThe following ignores already exist in the ignore file:\n")
				cmd.Printf("%s", skippedIgnoresToMigrate)
				cmd.Printf("\nTo overwrite these entries, use --force\n")
			}

			// either no duplicate entries at this point or --force is true so we can ignore merge error
			_ = ignore.MergeIgnoredFingerprints(fingerprintsToMigrate, ignoredFingerprints, options.IgnoreMigrateOptions.Force)
			return writeIgnoreFile(ignoredFingerprints, ignoreFilepath)
		},
		SilenceErrors: false,
		SilenceUsage:  false,
	}

	flags.AddFlags(cmd)
	cmd.SetUsageTemplate(fmt.Sprintf(scanTemplate, flags.Usages(cmd)))

	return cmd
}

func setLogLevel(cmd *cobra.Command) {
	logLevel := viper.GetString(flag.LogLevelFlag.ConfigName)
	if viper.GetBool(flag.DebugFlag.ConfigName) {
		logLevel = flag.DebugLogLevel
	}
	output.Setup(cmd, output.SetupRequest{
		LogLevel: logLevel,
	})
}

func writeIgnoreFile(ignoredFingerprints map[string]ignoretypes.IgnoredFingerprint, ignoreFilePath string) error {
	data, err := json.MarshalIndent(ignoredFingerprints, "", "  ")
	if err != nil {
		// failed to marshall data
		return err
	}

	return os.WriteFile(ignoreFilePath, data, 0644)
}

func getIgnoredFingerprintsFromConfig(configPath string) (ignoredFingerprintsFromConfig map[string]ignoretypes.IgnoredFingerprint) {
	ignoredFingerprintsFromConfig = make(map[string]ignoretypes.IgnoredFingerprint)

	for _, fingerprint := range viper.GetStringSlice("report.exclude-fingerprint") {
		ignoredFingerprintsFromConfig[fingerprint] = ignoretypes.IgnoredFingerprint{
			Comment: &migratedIgnoreComment,
		}
	}

	return ignoredFingerprintsFromConfig
}

func requestConfirmation(s string) bool {
	r := bufio.NewReader(os.Stdin)

	for i := 0; true; i++ {
		if i > 0 {
			fmt.Printf("Please enter y or n\n")
		}
		fmt.Printf("%s [Y/n]: ", s)

		response, _ := r.ReadString('\n')
		response = strings.ToLower(strings.TrimSpace(response))

		if response == "y" || response == "yes" {
			return true
		}

		if response == "" {
			// Enter key defaults to Y
			return true
		}

		if response == "n" || response == "no" {
			return false
		}

		continue
	}

	return false
}
