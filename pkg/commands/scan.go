package commands

import (
	"fmt"

	"github.com/bearer/curio/pkg/commands/artifact"
	"github.com/bearer/curio/pkg/flag"
	"github.com/bearer/curio/pkg/util/output"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/xerrors"
)

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

var scanFlags = &flag.Flags{
	ScanFlagGroup:    flag.NewScanFlagGroup(),
	PolicyFlagGroup:  flag.NewPolicyFlagGroup(),
	WorkerFlagGroup:  flag.NewWorkerFlagGroup(),
	ReportFlagGroup:  flag.NewReportFlagGroup(),
	GeneralFlagGroup: flag.NewGeneralFlagGroup(),
}

func NewScanCommand() *cobra.Command {

	cmd := &cobra.Command{
		Use:     "scan [flags] <path>",
		Aliases: []string{"s"},
		Short:   "Scan a directory or file",
		Example: `  # Scan a local project, including language-specific files
  $ curio scan /path/to/your_project`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := scanFlags.Bind(cmd); err != nil {
				return xerrors.Errorf("flag bind error: %w", err)
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := scanFlags.Bind(cmd); err != nil {
				return xerrors.Errorf("flag bind error: %w", err)
			}

			configPath := viper.GetString(flag.ConfigFileFlag.ConfigName)

			if err := readConfig(configPath); err != nil {
				return err
			}

			options, err := scanFlags.ToOptions(args)
			if err != nil {
				return xerrors.Errorf("flag error: %w", err)
			}

			if !options.Quiet && configPath != "" {
				output.StdErrLogger().Msgf("Loaded %s configuration file", configPath)
			}

			output.Setup(cmd, options)

			if options.Target == "" {
				return cmd.Help()
			}

			return artifact.Run(cmd.Context(), options, artifact.TargetFilesystem)
		},
		SilenceErrors: false,
		SilenceUsage:  false,
	}

	scanFlags.AddFlags(cmd)
	cmd.SetUsageTemplate(fmt.Sprintf(scanTemplate, scanFlags.Usages(cmd)))

	return cmd
}

func readConfig(configFile string) error {
	viper.SetConfigType("yaml")
	viper.SetConfigFile(configFile)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil
		}

		return fmt.Errorf("config file %q loading error: %s", configFile, err)
	}

	return nil
}
