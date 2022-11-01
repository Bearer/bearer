package commands

import (
	"fmt"

	"github.com/bearer/curio/pkg/commands/artifact"
	"github.com/bearer/curio/pkg/flag"
	"github.com/spf13/cobra"
	"golang.org/x/xerrors"
)

// VersionInfo holds the curio version
type VersionInfo struct {
	Version string `json:",omitempty"`
}

// NewApp is the factory method to return Curio CLI
func NewApp(version string, commitSHA string) *cobra.Command {
	rootCmd := NewRootCommand()
	rootCmd.AddCommand(
		NewProcessingWorkerCommand(),
		NewInitCommand(),
		NewScanCommand(),
		NewConfigCommand(),
		NewVersionCommand(version, commitSHA),
	)

	return rootCmd
}

func NewRootCommand() *cobra.Command {
	usageTemplate := `Curio is a tool for scanning policy breaches

Scan Example:
	# Scan local repository
	$ curio scan <repository>

Available Commands:
	scan              Scan git repository
	init              Writes default config to curio.yml
	version           Print the version
`

	cmd := &cobra.Command{
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
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
				return xerrors.Errorf("flag bind error: %w", err)
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := configFlags.Bind(cmd); err != nil {
				return xerrors.Errorf("flag bind error: %w", err)
			}
			options, err := configFlags.ToOptions(args, cmd.OutOrStdout())
			if err != nil {
				return xerrors.Errorf("flag error: %w", err)
			}

			return artifact.Run(cmd.Context(), options, artifact.TargetFilesystem)
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
