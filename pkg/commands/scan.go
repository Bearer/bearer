package commands

import (
	"fmt"

	"github.com/bearer/bearer/pkg/commands/artifact"
	"github.com/bearer/bearer/pkg/flag"
	"github.com/bearer/bearer/pkg/util/file"
	"github.com/bearer/bearer/pkg/util/output"
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

var ScanFlags = &flag.Flags{
	ScanFlagGroup:    flag.NewScanFlagGroup(),
	RuleFlagGroup:    flag.NewRuleFlagGroup(),
	ReportFlagGroup:  flag.NewReportFlagGroup(),
	GeneralFlagGroup: flag.NewGeneralFlagGroup(),
}

func NewScanCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "scan [flags] <path>",
		Aliases: []string{"s"},
		Short:   "Scan a directory or file",
		Example: `  # Scan a local project, including language-specific files
  $ bearer scan /path/to/your_project`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := ScanFlags.Bind(cmd); err != nil {
				return xerrors.Errorf("flag bind error: %w", err)
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := ScanFlags.Bind(cmd); err != nil {
				return xerrors.Errorf("flag bind error: %w", err)
			}

			configPath := viper.GetString(flag.ConfigFileFlag.ConfigName)
			var defaultConfigPath = ""
			if len(args) == 1 {
				defaultConfigPath = file.GetFullFilename(args[0], configPath)
			}

			var loadFileMessage string
			if err := readConfig(configPath); err != nil {
				if err := readConfig(defaultConfigPath); err != nil {
					loadFileMessage = fmt.Sprintf("Couldn't find config file %s or %s", configPath, defaultConfigPath)
				} else {
					loadFileMessage = fmt.Sprintf("Loading default config file %s", defaultConfigPath)
				}
			} else {
				loadFileMessage = fmt.Sprintf("Loading config file %s", configPath)
			}

			options, err := ScanFlags.ToOptions(args)
			if err != nil {
				return xerrors.Errorf("flag error: %w", err)
			}

			if !options.Quiet {
				output.StdErrLogger().Msgf(loadFileMessage)
			}

			output.Setup(cmd, options)

			if options.Target == "" {
				return cmd.Help()
			}

			cmd.SilenceUsage = true

			return artifact.Run(cmd.Context(), options, artifact.TargetFilesystem)
		},
		SilenceErrors: false,
		SilenceUsage:  false,
	}

	ScanFlags.AddFlags(cmd)
	cmd.SetUsageTemplate(fmt.Sprintf(scanTemplate, ScanFlags.Usages(cmd)))

	return cmd
}

func readConfig(configFile string) error {
	if configFile == "" {
		return nil
	}

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
