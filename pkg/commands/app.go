package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/bearer/curio/pkg/commands/artifact"
	"github.com/bearer/curio/pkg/commands/process/worker"
	"github.com/bearer/curio/pkg/flag"
	"github.com/bearer/curio/pkg/types"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/xerrors"
)

// VersionInfo holds the curio version
type VersionInfo struct {
	Version string `json:",omitempty"`
}

const (
	usageTemplate = `Usage:{{if .Runnable}}
  {{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}
  {{.CommandPath}} [command]{{end}}{{if gt (len .Aliases) 0}}
Aliases:
  {{.NameAndAliases}}{{end}}{{if .HasExample}}
Examples:
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}
Available Commands:{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}
%s
Global Flags:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}
Additional help topics:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}
Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}
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
func NewApp(version string) *cobra.Command {
	globalFlags := flag.NewGlobalFlagGroup()
	rootCmd := NewRootCommand(version, globalFlags)
	rootCmd.AddCommand(
		NewProcessingServerCommand(),
		NewScanCommand(globalFlags),
		NewConfigCommand(globalFlags),
		NewVersionCommand(globalFlags),
	)

	return rootCmd
}

func initConfig(configFile string) error {
	// Read from config
	viper.SetConfigFile(configFile)
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			// log.Logger.Debugf("config file %q not found", configFile)
			return nil
		}
		return xerrors.Errorf("config file %q loading error: %s", configFile, err)
	}
	// log.Logger.Infof("Loaded %s", configFile)
	return nil
}

func NewProcessingServerCommand() *cobra.Command {
	flags := &flag.Flags{
		ProcessFlagGroup: flag.NewProcessGroup(),
	}

	cmd := &cobra.Command{
		Use:   "processing-worker",
		Short: "start scan processing server",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Debug().Msgf("started scan processing")
			if err := flags.Bind(cmd); err != nil {
				return fmt.Errorf("flag bind error: %w", err)
			}

			options, err := flags.ProcessFlagGroup.ToOptions()
			if err != nil {
				return fmt.Errorf("options binding error: %w", err)
			}

			log.Debug().Msgf("running scan worker on port `%s`", options.Port)

			return worker.Start(options.Port)
		},
		Hidden:        true,
		SilenceErrors: true,
		SilenceUsage:  true,
	}
	cmd.SetFlagErrorFunc(flagErrorFunc)
	flags.AddFlags(cmd)

	return cmd
}

func NewRootCommand(version string, globalFlags *flag.GlobalFlagGroup) *cobra.Command {
	var versionFormat string
	cmd := &cobra.Command{
		Use:   "curio [global flags] command [flags] target",
		Short: "Unified data security scanner",
		Long:  "Scanner for Git repositories",
		Example: `  # Scan local repository
  $ curio scan <repository>`,
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
		},
		Args: cobra.NoArgs,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			cmd.SetOut(outputWriter)

			// Set the Curio version here so that we can override version printer.
			cmd.Version = version

			// viper.BindPFlag cannot be called in init().
			// cf. https://github.com/spf13/cobra/issues/875
			//     https://github.com/spf13/viper/issues/233
			if err := globalFlags.Bind(cmd.Root()); err != nil {
				return xerrors.Errorf("flag bind error: %w", err)
			}

			// The config path is needed for config initialization.
			// It needs to be obtained before ToOptions().
			configPath := viper.GetString(flag.ConfigFileFlag.ConfigName)

			// Configure environment variables and config file
			// It cannot be called in init() because it must be called after viper.BindPFlags.
			if err := initConfig(configPath); err != nil {
				return err
			}

			// globalOptions := globalFlags.ToOptions()

			// Initialize logger
			// if err := log.InitLogger(globalOptions.Debug, globalOptions.Quiet); err != nil {
			// 	return err
			// }

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			globalOptions := globalFlags.ToOptions()
			if globalOptions.ShowVersion {
				// Customize version output
				showVersion(versionFormat, version, outputWriter)
			} else {
				return cmd.Help()
			}
			return nil
		},
	}

	// Add version format flag, only json is supported
	cmd.Flags().StringVarP(&versionFormat, flag.FormatFlag.Name, flag.FormatFlag.Shorthand, "", "version format (json)")

	globalFlags.AddFlags(cmd)

	return cmd
}

func NewScanCommand(globalFlags *flag.GlobalFlagGroup) *cobra.Command {
	// reportFlagGroup := flag.NewReportFlagGroup()
	// reportFlagGroup.ReportFormat = nil // TODO: support --report summary

	flags := &flag.Flags{
		// CacheFlagGroup:  flag.NewCacheFlagGroup(),
		// ReportFlagGroup: reportFlagGroup,
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
			options, err := flags.ToOptions(cmd.Version, args, globalFlags, outputWriter)
			if err != nil {
				return xerrors.Errorf("flag error: %w", err)
			}

			if options.Target == "" {
				return fmt.Errorf("PATH is required")
			}

			return artifact.Run(cmd.Context(), options, artifact.TargetFilesystem)
		},
		SilenceErrors: false,
		SilenceUsage:  false,
	}

	cmd.SetFlagErrorFunc(flagErrorFunc)
	flags.AddFlags(cmd)
	cmd.SetUsageTemplate(fmt.Sprintf(usageTemplate, flags.Usages(cmd)))

	return cmd
}

func NewRepositoryCommand(globalFlags *flag.GlobalFlagGroup) *cobra.Command {
	// reportFlagGroup := flag.NewReportFlagGroup()
	// reportFlagGroup.ReportFormat = nil // TODO: support --report summary

	repoFlags := &flag.Flags{
		ScanFlagGroup: flag.NewScanFlagGroup(),
		RepoFlagGroup: flag.NewRepoFlagGroup(),
	}

	cmd := &cobra.Command{
		Use:     "repository [flags] REPO_URL",
		Aliases: []string{"repo"},
		Short:   "Scan a remote repository",
		Example: `  # Scan your remote git repository
  $ curio repo https://github.com/curio/curio-ci-test`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := repoFlags.Bind(cmd); err != nil {
				return xerrors.Errorf("flag bind error: %w", err)
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := repoFlags.Bind(cmd); err != nil {
				return xerrors.Errorf("flag bind error: %w", err)
			}
			options, err := repoFlags.ToOptions(cmd.Version, args, globalFlags, outputWriter)
			if err != nil {
				return xerrors.Errorf("flag error: %w", err)
			}
			return artifact.Run(cmd.Context(), options, artifact.TargetRepository)
		},
		SilenceErrors: true,
		SilenceUsage:  true,
	}
	cmd.SetFlagErrorFunc(flagErrorFunc)
	repoFlags.AddFlags(cmd)
	cmd.SetUsageTemplate(fmt.Sprintf(usageTemplate, repoFlags.Usages(cmd)))

	return cmd
}

func NewConfigCommand(globalFlags *flag.GlobalFlagGroup) *cobra.Command {

	scanFlags := &flag.ScanFlagGroup{
		SkipPathFlag: &flag.SkipPathFlag,
		FilePatterns: &flag.FilePatternsFlag,
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
			options, err := configFlags.ToOptions(cmd.Version, args, globalFlags, outputWriter)
			if err != nil {
				return xerrors.Errorf("flag error: %w", err)
			}

			// Disable OS and language analyzers

			// Scan only for misconfigurations
			options.SecurityChecks = []string{types.SecurityCheckConfig}

			return artifact.Run(cmd.Context(), options, artifact.TargetFilesystem)
		},
		SilenceErrors: true,
		SilenceUsage:  true,
	}
	cmd.SetFlagErrorFunc(flagErrorFunc)
	configFlags.AddFlags(cmd)
	cmd.SetUsageTemplate(fmt.Sprintf(usageTemplate, configFlags.Usages(cmd)))

	return cmd
}

func NewVersionCommand(globalFlags *flag.GlobalFlagGroup) *cobra.Command {
	var versionFormat string
	cmd := &cobra.Command{
		Use:   "version [flags]",
		Short: "Print the version",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			showVersion(versionFormat, cmd.Version, outputWriter)

			return nil
		},
		SilenceErrors: true,
		SilenceUsage:  true,
	}
	cmd.SetFlagErrorFunc(flagErrorFunc)

	// Add version format flag, only json is supported
	cmd.Flags().StringVarP(&versionFormat, flag.FormatFlag.Name, flag.FormatFlag.Shorthand, "", "version format (json)")

	return cmd
}

func showVersion(outputFormat, version string, outputWriter io.Writer) {
	switch outputFormat {
	case "json":
		b, _ := json.Marshal(VersionInfo{
			Version: version,
		})
		fmt.Fprint(outputWriter, string(b))
	default:
		output := fmt.Sprintf("Version: %s\n", version)
		fmt.Fprint(outputWriter, output)
	}
}

// show help on using the command when an invalid flag is encountered
func flagErrorFunc(command *cobra.Command, err error) error {
	if err := command.Help(); err != nil {
		return err
	}
	command.Println() // add empty line after list of flags
	return err
}
