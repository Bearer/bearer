package commands

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/bearer/curio/pkg/commands/artifact"
	config "github.com/bearer/curio/pkg/commands/process/settings"
	"github.com/bearer/curio/pkg/commands/process/worker"
	"github.com/bearer/curio/pkg/flag"
	"github.com/bearer/curio/pkg/util/output"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"golang.org/x/xerrors"
)

// VersionInfo holds the curio version
type VersionInfo struct {
	Version string `json:",omitempty"`
}

const (
	scanTemplate = `Usage:{{if .Runnable}}
  {{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}
  {{.CommandPath}} [command]{{end}}{{if gt (len .Aliases) 0}}
Aliases:
  {{.NameAndAliases}}{{end}}{{if .HasExample}}
Examples:
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}
Available Commands:{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}

{{if .HasAvailableLocalFlags}}
%s
{{end}}
`
)

var (
	outputWriter io.Writer = os.Stdout
)

// SetOut overrides the destination for messages
func SetOut(out io.Writer) {
	outputWriter = out
}

// NewApp is the factory method to return Curio CLI
func NewApp(version string, commitSHA string) *cobra.Command {
	rootCmd := NewRootCommand()
	rootCmd.AddCommand(
		NewProcessingServerCommand(),
		NewScanCommand(),
		NewConfigCommand(),
		NewVersionCommand(version, commitSHA),
	)

	return rootCmd
}

func NewProcessingServerCommand() *cobra.Command {
	flags := &flag.Flags{
		ProcessFlagGroup: flag.NewProcessGroup(),
		ScanFlagGroup:    flag.NewScanFlagGroup(),
	}

	cmd := &cobra.Command{
		Use:   "processing-worker [flags] PATH",
		Short: "start scan processing server",
		RunE: func(cmd *cobra.Command, args []string) error {
			var config config.Config
			if err := json.Unmarshal([]byte(args[0]), &config); err != nil {
				return err
			}

			if err := flags.Bind(cmd); err != nil {
				return fmt.Errorf("flag bind error: %w", err)
			}

			generalOptions, err := flags.ToOptions(cmd.Version, args, outputWriter)
			if err != nil {
				return fmt.Errorf("options binding error: %w", err)
			}

			output.Setup(generalOptions)

			processOptions, err := flags.ProcessFlagGroup.ToOptions()
			if err != nil {
				return fmt.Errorf("options binding error: %w", err)
			}

			log.Debug().Msgf("started scan processing")
			log.Debug().Msgf("running scan worker on port `%s`", processOptions.Port)

			return worker.Start(processOptions.Port, config)
		},
		Hidden:        true,
		SilenceErrors: true,
		SilenceUsage:  true,
	}
	cmd.SetFlagErrorFunc(flagErrorFunc)
	flags.AddFlags(cmd)

	return cmd
}

func NewRootCommand() *cobra.Command {
	usageTemplate := `Curio is a tool for scanning policy breaches

Scan Example:
	# Scan local repository
	$ curio scan <repository>

Available Commands:
	scan              Scan git repository
	init              Inits default configuration file (curio.yml)
	config            Scan config files for misconfigurations
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

func NewScanCommand() *cobra.Command {

	flags := &flag.Flags{
		ScanFlagGroup:   flag.NewScanFlagGroup(),
		WorkerFlagGroup: flag.NewWorkerFlagGroup(),
		ReportFlagGroup: flag.NewReportFlagGroup(),
	}

	cmd := &cobra.Command{
		Use:     "scan [flags] PATH",
		Aliases: []string{"s"},
		Short:   "Scan git repository",
		Example: `  # Scan a local project including language-specific files
  $ curio s /path/to/your_project
  # Scan a single file
  $ curio s ./curio-ci-test/Pipfile.lock`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := flags.Bind(cmd); err != nil {
				return xerrors.Errorf("flag bind error: %w", err)
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := flags.Bind(cmd); err != nil {
				return xerrors.Errorf("flag bind error: %w", err)
			}
			options, err := flags.ToOptions(cmd.Version, args, outputWriter)
			if err != nil {
				return xerrors.Errorf("flag error: %w", err)
			}

			output.Setup(options)

			if options.Target == "" {
				return cmd.Help()
			}

			return artifact.Run(cmd.Context(), options, artifact.TargetFilesystem)
		},
		SilenceErrors: false,
		SilenceUsage:  false,
	}

	flags.AddFlags(cmd)
	cmd.SetUsageTemplate(fmt.Sprintf(scanTemplate, flags.Usages(cmd)))

	return cmd
}

func NewConfigCommand() *cobra.Command {
	scanFlags := &flag.ScanFlagGroup{
		SkipPathFlag: &flag.SkipPathFlag,
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
			options, err := configFlags.ToOptions(cmd.Version, args, outputWriter)
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

func NewVersionCommand(version string, commitSHA string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			output.DefaultLogger().Msgf("curio version: %s\nsha: %s", version, commitSHA)
			return nil
		},
	}
	cmd.SetFlagErrorFunc(flagErrorFunc)

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
