package commands

import (
	"errors"
	"fmt"

	"github.com/bearer/bearer/pkg/flag"
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
	var IgnoreShowFlags = &flag.Flags{
		IgnoreFlagGroup: flag.NewIgnoreFlagGroup(),
	}
	cmd := &cobra.Command{
		Use:   "show <fingerprint>",
		Short: "Show an ignored fingerprint",
		Example: `# Show the details of an ignored fingerprint from your bearer.ignore file
$ bearer ignore show <fingerprint>`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return cmd.Help()
			}

			options, err := IgnoreShowFlags.ToOptions(args)
			if err != nil {
				return fmt.Errorf("flag error: %s", err)
			}

			ignoredFingerprints, err := ignore.GetIgnoredFingerprints(options.IgnoreOptions.BearerIgnoreFile)
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

			existingIgnoredFingerprints, fileExists, err := ignore.GetExistingIgnoredFingerprints(options.IgnoreOptions.BearerIgnoreFile)
			if err != nil {
				return fmt.Errorf("error retrieving existing ignores: %s", err)
			}

			if !fileExists {
				cmd.Printf("\nCreating bearer.ignore file...")
			}

			if err := ignore.AddToIgnoreFile(
				fingerprintsToIgnore,
				existingIgnoredFingerprints,
				options.IgnoreOptions.BearerIgnoreFile,
				options.IgnoreAddOptions.Force,
			); err != nil {
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
			ignoredFingerprintsFromConfig, err := getIgnoredFingerprintsFromConfig(configFilePath)
			if err != nil {
				cmd.Printf("Error: %s, perhaps you need to use --config-file to specify the config path?\n", err.Error())
				return nil
			}

			existingIgnoredFingerprints, fileExists, err := ignore.GetExistingIgnoredFingerprints(options.IgnoreOptions.BearerIgnoreFile)
			if err != nil {
				return fmt.Errorf("error retrieving existing ignores: %s", err)
			}

			entrysInConfig := len(ignoredFingerprintsFromConfig)
			entrysInConfigSkipped := ""
			cmd.Printf("Found %d ignores in:\n\t%s\n", entrysInConfig, configFilePath)

			if !fileExists {
				cmd.Printf("Creating bearer.ignore file\n")
			}

			if !options.IgnoreMigrateOptions.Force {
				for key := range existingIgnoredFingerprints {
					if _, ok := ignoredFingerprintsFromConfig[key]; ok {
						entrysInConfig--
						entrysInConfigSkipped += fmt.Sprintf("- %s\n", key)
						delete(ignoredFingerprintsFromConfig, key)
					}
				}
			}

			cmd.Printf("Added %d ignores to:\n\t%s\n", entrysInConfig, options.IgnoreOptions.BearerIgnoreFile)

			if entrysInConfigSkipped != "" {
				cmd.Printf("\nThe following items where already ignored:\n")
				cmd.Printf(entrysInConfigSkipped)
				cmd.Printf("\nTo overwrite these entries in the ignore file, use --force\n")
			}

			if err = ignore.AddToIgnoreFile(ignoredFingerprintsFromConfig, existingIgnoredFingerprints, options.IgnoreOptions.BearerIgnoreFile, options.IgnoreMigrateOptions.Force); err != nil {
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
