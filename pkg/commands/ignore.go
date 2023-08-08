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

var IgnoreFlags = &flag.Flags{
	IgnoreFlagGroup:  flag.NewIgnoreFlagGroup(),
	GeneralFlagGroup: flag.NewGeneralFlagGroup(),
}

func NewIgnoreCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "ignore [subcommand] <fingerprint>",
		Short:         "Manage ignored fingerprints",
		Example:       `  # Add a fingerprint to, or show a fingerprint from, your bearer.ignore file`,
		Args:          cobra.NoArgs,
		SilenceErrors: false,
		SilenceUsage:  false,
	}

	cmd.AddCommand(
		newIgnoreShowCommand(),
		newIgnoreAddCommand(),
		newIgnoreMigrateCommand(),
	)

	return cmd
}

func newIgnoreShowCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "show <fingerprint>",
		Short: "Show an ignored fingerprint",
		Example: `  # Show the details of an ignored fingerprint from your bearer.ignore file
$ bearer ignore show <fingerprint>`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				cmd.Printf("No fingerprint given. Please provide a fingerprint with the command: $ bearer ignore show <fingerprint>\n")
				return nil
			}

			ignoredFingerprints, err := ignore.GetIgnoredFingerprints(nil)
			if err != nil {
				cmd.Printf("Issue loading ignored fingerprints from bearer.ignore file")
				return nil
			}

			selectedIgnoredFingerprint, ok := ignoredFingerprints[args[0]]
			if !ok {
				cmd.Printf("Ignored fingerprint '%s' was not found in bearer.ignore file", args[0])
				return nil
			}

			cmd.Printf("\nIgnored At: %s\nAuthor: %s\n", selectedIgnoredFingerprint.IgnoredAt, selectedIgnoredFingerprint.Author)
			cmd.Printf("Comment: %s\n\n", selectedIgnoredFingerprint.Comment)

			return nil
		},
		SilenceErrors: false,
		SilenceUsage:  false,
	}
}

func newIgnoreAddCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add <fingerprint>",
		Short: "Add an ignored fingerprint",
		Example: `  # Add an ignored fingerprint to your bearer.ignore file
$ bearer ignore add <fingerprint> --author Mish --comment "Possible false positive"`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := IgnoreFlags.Bind(cmd); err != nil {
				return fmt.Errorf("flag bind error: %w", err)
			}

			if len(args) == 0 {
				cmd.Printf("No fingerprint given. Please provide a fingerprint with the command: $ bearer ignore add <fingerprint>\n")
				return nil
			}

			options, err := IgnoreFlags.ToOptions(args)
			if err != nil {
				return fmt.Errorf("flag error: %s", err)
			}

			fingerprintsToIgnore := map[string]ignore.IgnoredFingerprint{
				args[0]: {
					Author:  options.IgnoreOptions.Author,
					Comment: options.IgnoreOptions.Comment,
				},
			}

			if err = ignore.AddToIgnoreFile(fingerprintsToIgnore, options.IgnoreOptions.Force); err != nil {
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
	IgnoreFlags.AddFlags(cmd)
	cmd.SetUsageTemplate(fmt.Sprintf(scanTemplate, IgnoreFlags.Usages(cmd)))

	return cmd
}

func newIgnoreMigrateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "Migrate ignored fingerprints from bearer.yml to bearer.ignore",
		Example: `  # Migrate existing ignored (excluded) fingerprints from bearer.yml file to bearer.ignore
$ bearer ignore migrate`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := IgnoreFlags.Bind(cmd); err != nil {
				return fmt.Errorf("flag bind error: %w", err)
			}

			options, err := IgnoreFlags.ToOptions(args)
			if err != nil {
				return fmt.Errorf("flag error: %s", err)
			}

			ignoredFingerprintsFromConfig, err := getIgnoredFingerprintsFromConfig(args)
			if err != nil {
				return err
			}

			if err = ignore.AddToIgnoreFile(ignoredFingerprintsFromConfig, options.IgnoreOptions.Force); err != nil {
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
	IgnoreFlags.AddFlags(cmd)
	cmd.SetUsageTemplate(fmt.Sprintf(scanTemplate, IgnoreFlags.Usages(cmd)))

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
