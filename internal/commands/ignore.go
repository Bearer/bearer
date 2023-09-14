package commands

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/bearer/bearer/internal/flag"
	"github.com/bearer/bearer/internal/report/output/saas"
	"github.com/bearer/bearer/internal/util/ignore"
	ignoretypes "github.com/bearer/bearer/internal/util/ignore/types"
	"github.com/bearer/bearer/internal/util/output"
)

var migratedIgnoreComment = "migrated from bearer.yml"

func NewIgnoreCommand() *cobra.Command {
	usageTemplate := `
Usage: bearer ignore <command> [flags]

Available Commands:
    add              Add an ignored fingerprint
    show             Show an ignored fingerprint
    remove           Remove an ignored fingerprint
    pull             Pull ignored fingerprints from Cloud
    migrate          Migrate ignored fingerprints

Examples:
    # Add an ignored fingerprint to your bearer.ignore file
    $ bearer ignore add <fingerprint> --author Mish --comment "investigate this"

    # Show the details of an ignored fingerprint from your bearer.ignore file
    $ bearer ignore show <fingerprint>

    # Remove an ignored fingerprint from your bearer.ignore file
    $ bearer ignore remove <fingerprint>

    # Pull ignored fingerprints from the Cloud (requires API key)
    $ bearer ignore pull /path/to/your_project --api-key=XXXXX

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
		newIgnoreRemoveCommand(),
		newIgnorePullCommand(),
		newIgnoreMigrateCommand(),
	)

	cmd.SetUsageTemplate(usageTemplate)

	return cmd
}

func newIgnoreShowCommand() *cobra.Command {
	var IgnoreShowFlags = &flag.Flags{
		GeneralFlagGroup:    flag.NewGeneralFlagGroup(),
		IgnoreShowFlagGroup: flag.NewIgnoreShowFlagGroup(),
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

			setLogLevel(cmd)

			options, err := IgnoreShowFlags.ToOptions(args)
			if err != nil {
				return fmt.Errorf("flag error: %s", err)
			}

			if len(args) == 0 && !options.IgnoreShowOptions.All {
				return cmd.Help()
			}

			ignoredFingerprints, fileExists, err := ignore.GetIgnoredFingerprints(options.GeneralOptions.IgnoreFile, nil)
			if err != nil {
				cmd.Printf("Issue loading ignored fingerprints from bearer.ignore file: %s", err)
				return nil
			}
			if !fileExists {
				cmd.Printf("bearer.ignore file not found. Perhaps you need to use --bearer-ignore-file to specify the path to bearer.ignore?\n")
				return nil
			}

			cmd.Print("\n")
			if options.IgnoreShowOptions.All {
				// show all fingerprints
				for fingerprintId, fingerprint := range ignoredFingerprints {
					cmd.Print(ignore.DisplayIgnoredEntryTextString(fingerprintId, fingerprint, options.GeneralOptions.NoColor))
					cmd.Print("\n\n")
				}
			} else {
				// show a specific fingerprint
				fingerprintId := args[0]
				selectedIgnoredFingerprint, ok := ignoredFingerprints[fingerprintId]
				if !ok {
					cmd.Printf("Ignored fingerprint '%s' was not found in bearer.ignore file\n", fingerprintId)
					return nil
				}
				cmd.Print(ignore.DisplayIgnoredEntryTextString(fingerprintId, selectedIgnoredFingerprint, options.GeneralOptions.NoColor))
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
	var IgnoreAddFlags = &flag.Flags{
		GeneralFlagGroup:   flag.NewGeneralFlagGroup(),
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

			setLogLevel(cmd)

			options, err := IgnoreAddFlags.ToOptions(args)
			if err != nil {
				return fmt.Errorf("flag error: %s", err)
			}

			// create initial entry
			fingerprintId := args[0]
			var fingerprintEntry ignoretypes.IgnoredFingerprint
			fingerprintsToIgnore := map[string]ignoretypes.IgnoredFingerprint{
				fingerprintId: fingerprintEntry,
			}

			ignoredFingerprints, fileExists, err := ignore.GetIgnoredFingerprints(options.GeneralOptions.IgnoreFile, nil)
			if err != nil {
				return fmt.Errorf("error retrieving existing ignores: %s", err)
			}

			// check for merge conflicts
			if mergeErr := ignore.MergeIgnoredFingerprints(fingerprintsToIgnore, ignoredFingerprints, options.IgnoreAddOptions.Force); mergeErr != nil {
				// handle expected error (duplicate entry in bearer.ignore)
				cmd.Printf("Error: %s\n", mergeErr.Error())
				return nil
			}

			// ensure ignored at is set
			fingerprintEntry = ignoredFingerprints[fingerprintId]

			// add additional information to entry
			if options.IgnoreAddOptions.Author != "" {
				fingerprintEntry.Author = &options.IgnoreAddOptions.Author
			} else {
				if author, err := ignore.GetAuthor(); err == nil {
					fingerprintEntry.Author = author
				}
			}
			if options.IgnoreAddOptions.FalsePositive {
				fingerprintEntry.FalsePositive = options.IgnoreAddOptions.FalsePositive
			} else {
				fingerprintEntry.FalsePositive = requestConfirmation("Is this finding a false positive?")
				cmd.Printf("\n")
			}
			if options.IgnoreAddOptions.Comment != "" {
				fingerprintEntry.Comment = &options.IgnoreAddOptions.Comment
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
				cmd.Printf("\nCreating bearer.ignore file...\n")
			}

			if err := writeIgnoreFile(ignoredFingerprints, options.GeneralOptions.IgnoreFile); err != nil {
				return err
			}

			cmd.Print("Fingerprint added to bearer.ignore:\n\n")
			cmd.Print(ignore.DisplayIgnoredEntryTextString(fingerprintId, ignoredFingerprints[fingerprintId], options.GeneralOptions.NoColor))
			cmd.Print("\n\n")
			return nil
		},
		SilenceErrors: false,
		SilenceUsage:  false,
	}
	IgnoreAddFlags.AddFlags(cmd)
	cmd.SetUsageTemplate(fmt.Sprintf(scanTemplate, IgnoreAddFlags.Usages(cmd)))

	return cmd
}

func newIgnoreRemoveCommand() *cobra.Command {
	var IgnoreRemoveFlags = &flag.Flags{
		GeneralFlagGroup: flag.NewGeneralFlagGroup(),
	}
	cmd := &cobra.Command{
		Use:   "remove <fingerprint>",
		Short: "Remove an ignored fingerprint",
		Example: `# Remove an ignored fingerprint from your bearer.ignore file
$ bearer ignore remove <fingerprint>`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := IgnoreRemoveFlags.Bind(cmd); err != nil {
				return fmt.Errorf("flag bind error: %w", err)
			}

			if len(args) == 0 {
				return cmd.Help()
			}

			setLogLevel(cmd)

			options, err := IgnoreRemoveFlags.ToOptions(args)
			if err != nil {
				return fmt.Errorf("flag error: %s", err)
			}

			ignoredFingerprints, fileExists, err := ignore.GetIgnoredFingerprints(options.GeneralOptions.IgnoreFile, nil)
			if err != nil {
				return fmt.Errorf("error retrieving existing ignores: %s", err)
			}
			if !fileExists {
				cmd.Printf("bearer.ignore file not found. Perhaps you need to use --bearer-ignore-file to specify the path to bearer.ignore?\n")
				return nil
			}

			fingerprintId := args[0]
			removedFingerprint, ok := ignoredFingerprints[fingerprintId]
			if !ok {
				cmd.Printf("Ignored fingerprint '%s' was not found in bearer.ignore file\n", fingerprintId)
				return nil
			}

			delete(ignoredFingerprints, fingerprintId)
			if err := writeIgnoreFile(ignoredFingerprints, options.GeneralOptions.IgnoreFile); err != nil {
				return err
			}

			cmd.Print("Fingerprint successfully removed from bearer.ignore:\n\n")
			cmd.Print(ignore.DisplayIgnoredEntryTextString(fingerprintId, removedFingerprint, options.GeneralOptions.NoColor))
			cmd.Print("\n\n")
			return nil
		},
		SilenceErrors: false,
		SilenceUsage:  false,
	}
	IgnoreRemoveFlags.AddFlags(cmd)
	cmd.SetUsageTemplate(fmt.Sprintf(scanTemplate, IgnoreRemoveFlags.Usages(cmd)))

	return cmd
}

func newIgnorePullCommand() *cobra.Command {
	var IgnorePullFlags = &flag.Flags{
		GeneralFlagGroup: flag.NewGeneralFlagGroup(),
	}
	cmd := &cobra.Command{
		Use:   "pull",
		Short: "Pull ignored fingerprints from Cloud",
		Example: `# Pull ignored fingerprints from the Cloud (requires API key)
$ bearer ignore pull /path/to/your_project --api-key=XXXXX`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := IgnorePullFlags.Bind(cmd); err != nil {
				return fmt.Errorf("flag bind error: %w", err)
			}

			setLogLevel(cmd)

			options, err := IgnorePullFlags.ToOptions(args)
			if err != nil {
				return fmt.Errorf("flag error: %s", err)
			}

			if len(args) == 0 {
				return cmd.Help()
			} else {
				options.Target = args[0]
			}

			// confirm overwrite if bearer.ignore file exists
			bearerIgnoreFilePath := options.GeneralOptions.IgnoreFile
			fileExists := true
			info, err := os.Stat(bearerIgnoreFilePath)
			if os.IsNotExist(err) {
				fileExists = false
			} else {
				if info.IsDir() {
					return fmt.Errorf("bearer-ignore-file path %s is a dir not a file", bearerIgnoreFilePath)
				}
			}
			if err != nil && fileExists {
				return fmt.Errorf("file error: %s", err)
			}

			if fileExists {
				overwriteApproved := requestConfirmation("Warning: this action will overwrite your current bearer.ignore file. Continue?")
				cmd.Printf("\n")
				if !overwriteApproved {
					cmd.Printf("Okay, pull cancelled!\n")
					return nil
				}
			}

			// get project full name
			vcsInfo, err := saas.GetVCSInfo(options.Target)
			if err != nil {
				return fmt.Errorf("error fetching project info: %s", err)
			}

			data, err := options.GeneralOptions.Client.FetchIgnores(vcsInfo.FullName, []string{})
			if err != nil {
				return fmt.Errorf("cloud error: %s", err)
			}

			if !data.ProjectFound {
				// no project
				cmd.Printf("Project %s not found in Cloud. Pull cancelled.", vcsInfo.FullName)
				return nil
			}

			cloudIgnoresCount := len(data.CloudIgnoredFingerprints)
			if cloudIgnoresCount == 0 {
				// project found but no ignores
				cmd.Printf("No ignores for project %s found in the Cloud. Pull cancelled", vcsInfo.FullName)
				return nil
			}

			// project found and we have ignores - write to bearer.ignore
			cmd.Printf("Pulling %d ignores from the Cloud:\n", cloudIgnoresCount)
			for fingerprintId, fingerprint := range data.CloudIgnoredFingerprints {
				if fingerprint.Comment == nil {
					cmd.Printf("\t- %s\n", fingerprintId)
				} else {
					cmd.Printf("\t- %s (%s)\n", fingerprintId, *fingerprint.Comment)
				}
			}
			cmd.Printf("\n")

			if err = writeIgnoreFile(data.CloudIgnoredFingerprints, bearerIgnoreFilePath); err != nil {
				return fmt.Errorf("error writing to file: %s", err)
			}

			cmd.Printf("Pull successful! To view updated ignore file, run: bearer ignore show --all\n")
			return nil
		},
		SilenceErrors: false,
		SilenceUsage:  false,
	}
	IgnorePullFlags.AddFlags(cmd)
	cmd.SetUsageTemplate(fmt.Sprintf(scanTemplate, IgnorePullFlags.Usages(cmd)))

	return cmd
}

func newIgnoreMigrateCommand() *cobra.Command {
	IgnoreMigrateFlags := &flag.Flags{
		GeneralFlagGroup:       flag.NewGeneralFlagGroup(),
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

			setLogLevel(cmd)

			options, err := IgnoreMigrateFlags.ToOptions(args)
			if err != nil {
				return fmt.Errorf("flag error: %s", err)
			}

			configFilePath, _, err := readConfig(args)
			if err != nil {
				return fmt.Errorf("error reading config: %s\nPerhaps you need to use --config-file to specify the config path?", err.Error())
			}
			fingerprintsToMigrate := getIgnoredFingerprintsFromConfig(configFilePath)

			ignoredFingerprints, fileExists, err := ignore.GetIgnoredFingerprints(options.GeneralOptions.IgnoreFile, nil)
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

			cmd.Printf("Added %d ignores to:\n\t%s\n", migratedIgnoredCount, options.GeneralOptions.IgnoreFile)

			if skippedIgnoresToMigrate != "" {
				cmd.Printf("\nThe following ignores already exist in the bearer.ignore file:\n")
				cmd.Printf(skippedIgnoresToMigrate)
				cmd.Printf("\nTo overwrite these entries, use --force\n")
			}

			// either no duplicate entries at this point or --force is true so we can ignore merge error
			_ = ignore.MergeIgnoredFingerprints(fingerprintsToMigrate, ignoredFingerprints, options.IgnoreMigrateOptions.Force)
			return writeIgnoreFile(ignoredFingerprints, options.GeneralOptions.IgnoreFile)
		},
		SilenceErrors: false,
		SilenceUsage:  false,
	}
	IgnoreMigrateFlags.AddFlags(cmd)
	cmd.SetUsageTemplate(fmt.Sprintf(scanTemplate, IgnoreMigrateFlags.Usages(cmd)))

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

func writeIgnoreFile(ignoredFingerprints map[string]ignoretypes.IgnoredFingerprint, bearerIgnoreFilePath string) error {
	if bearerIgnoreFilePath == "" {
		bearerIgnoreFilePath = ignore.DefaultIgnoreFilepath
	}

	data, err := json.MarshalIndent(ignoredFingerprints, "", "  ")
	if err != nil {
		// failed to marshall data
		return err
	}

	return os.WriteFile(bearerIgnoreFilePath, data, 0644)
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
