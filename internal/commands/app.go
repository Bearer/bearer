package commands

import (
	"fmt"

	"github.com/bearer/bearer/cmd/bearer/build"
	"github.com/bearer/bearer/internal/commands/artifact"
	"github.com/bearer/bearer/internal/flag"
	"github.com/spf13/cobra"
)

// VersionInfo holds the bearer version
type VersionInfo struct {
	Version string `json:",omitempty" yaml:",omitempty"`
}

// NewApp is the factory method to return CLI
func NewApp(version string, commitSHA string) *cobra.Command {
	rootCmd := NewRootCommand()
	rootCmd.AddCommand(
		NewCompletionCommand(),
		NewProcessingWorkerCommand(),
		NewInitCommand(),
		NewScanCommand(),
		NewIgnoreCommand(),
		NewVersionCommand(version, commitSHA),
	)

	return rootCmd
}

func NewRootCommand() *cobra.Command {
	usageTemplate := `
[0;1;34;94m▄▄▄▄▄[0m
[0;34m█[0m    [0;34m█[0m  [0;34m▄▄▄[0m    [0;37m▄▄▄[0m    [0;37m▄[0m [0;37m▄▄[0m   [0;37m▄▄[0;1;30;90m▄[0m    [0;1;30;90m▄[0m [0;1;30;90m▄▄[0m
[0;34m█▄▄▄▄▀[0m [0;37m█▀[0m  [0;37m█[0m  [0;37m▀[0m   [0;37m█[0m   [0;37m█[0;1;30;90m▀[0m  [0;1;30;90m▀[0m [0;1;30;90m█▀[0m  [0;1;30;90m█[0m   [0;1;30;90m█▀[0m  [0;1;34;94m▀[0m
[0;37m█[0m    [0;37m█[0m [0;37m█▀▀▀▀[0m  [0;37m▄[0;1;30;90m▀▀▀█[0m   [0;1;30;90m█[0m     [0;1;30;90m█▀▀[0;1;34;94m▀▀[0m   [0;1;34;94m█[0m
[0;37m█▄▄▄▄▀[0m [0;1;30;90m▀█▄▄▀[0m  [0;1;30;90m▀▄▄▀█[0m   [0;1;30;90m█[0m     [0;1;34;94m▀█▄▄▀[0m   [0;1;34;94m█[0m

Scan your source code to discover, filter and prioritize security and privacy risks.

Usage: bearer <command> [flags]

Available Commands:
	completion        Generate the autocompletion script for your shell
	scan              Scan a directory or file
	init              Write the default config to bearer.yml
	ignore            Manage ignored fingerprints
	version           Print the version

Examples:
	# Scan local directory or file and output security risks
	$ bearer scan <path> --scanner=sast,secrets

	# Scan current directory and output the privacy report to a file
	$ bearer scan --report privacy --output <output-path> .

	# Scan local directory and output details about the underlying
	# detections and classifications
	$ bearer scan . --report dataflow

Learn More:
	Bearer is a code security tool that scans your source code to
	identify OWASP/CWE top security risks, privacy impact, and sensitive dataflows

	For more examples, tutorials, and to learn more about the project
	visit https://docs.bearer.com
`
	version := fmt.Sprintf("%s, build %s", build.Version, build.CommitSHA)
	cmd := &cobra.Command{
		Use:     "bearer",
		Args:    cobra.NoArgs,
		Version: version,
	}
	cmd.SetUsageTemplate(usageTemplate)

	return cmd
}

func NewConfigCommand() *cobra.Command {
	scanFlags := &flag.ScanFlagGroup{
		SkipPathFlag:                &flag.SkipPathFlag,
		DisableDomainResolutionFlag: &flag.DisableDomainResolutionFlag,
		DomainResolutionTimeoutFlag: &flag.DomainResolutionTimeoutFlag,
		InternalDomainsFlag:         &flag.InternalDomainsFlag,
		ContextFlag:                 &flag.ContextFlag,
	}

	configFlags := &flag.Flags{
		ScanFlagGroup: scanFlags,
	}

	cmd := &cobra.Command{
		Use:     "config [flags] DIR",
		Aliases: []string{"conf"},
		Short:   "Scan config files for misconfigurations",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := configFlags.Bind(cmd); err != nil {
				return fmt.Errorf("flag bind error: %w", err)
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := configFlags.Bind(cmd); err != nil {
				return fmt.Errorf("flag bind error: %w", err)
			}
			options, err := configFlags.ToOptions(args)
			if err != nil {
				return fmt.Errorf("flag error: %w", err)
			}

			return artifact.Run(cmd.Context(), options)
		},
		SilenceErrors: true,
		SilenceUsage:  true,
	}
	cmd.SetFlagErrorFunc(flagErrorFunc)
	configFlags.AddFlags(cmd)
	cmd.SetUsageTemplate(fmt.Sprintf(scanTemplate, configFlags.Usages(cmd)))

	return cmd
}

// show help on using the command when an invalid flag is encountered
func flagErrorFunc(command *cobra.Command, err error) error {
	if err := command.Help(); err != nil {
		return err
	}
	command.Println() // add empty line after list of flags
	return err
}
