package commands

import (
	"fmt"

	"github.com/bearer/bearer/pkg/commands/artifact"
	"github.com/bearer/bearer/pkg/commands/debugprofile"
	"github.com/bearer/bearer/pkg/flag"
	"github.com/bearer/bearer/pkg/util/file"
	"github.com/bearer/bearer/pkg/util/output"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
	IgnoreFlagGroup:  flag.NewIgnoreFlagGroup(),
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
				return fmt.Errorf("flag bind error: %w", err)
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := ScanFlags.Bind(cmd); err != nil {
				return fmt.Errorf("flag bind error: %w", err)
			}

			logLevel := viper.GetString(flag.LogLevelFlag.ConfigName)
			if viper.GetBool(flag.DebugFlag.ConfigName) {
				logLevel = flag.DebugLogLevel
			}

			output.Setup(cmd, output.SetupRequest{
				LogLevel:  logLevel,
				Quiet:     viper.GetBool(flag.QuietFlag.ConfigName),
				ProcessID: "main",
			})

			if viper.GetBool(flag.DebugProfileFlag.ConfigName) {
				debugprofile.Start()
			}

			configPath := viper.GetString(flag.ConfigFileFlag.ConfigName)
			var defaultConfigPath = ""
			if len(args) > 0 {
				defaultConfigPath = file.GetFullFilename(args[0], configPath)
			}

			var loadFileMessage string
			if err := readConfig(configPath); err != nil {
				if err := readConfig(defaultConfigPath); err != nil {
					loadFileMessage = "Couldn't find any config file"
				} else {
					loadFileMessage = fmt.Sprintf("Loading default config file %s", defaultConfigPath)
				}
			} else {
				loadFileMessage = fmt.Sprintf("Loading config file %s", configPath)
			}

			options, err := ScanFlags.ToOptions(args)
			if err != nil {
				return fmt.Errorf("flag error: %s", err)
			}

			log.Debug().Msgf(loadFileMessage)

			if len(args) == 0 {
				return cmd.Help()
			} else {
				options.Target = args[0]
			}

			cmd.SilenceUsage = true

			err = artifact.Run(cmd.Context(), options)
			debugprofile.Stop()
			return err
		},
		SilenceErrors: false,
		SilenceUsage:  false,
	}

	ScanFlags.AddFlags(cmd)
	cmd.SetUsageTemplate(fmt.Sprintf(scanTemplate, ScanFlags.Usages(cmd)))

	return cmd
}

func readConfig(configFile string) error {
	viper.SetConfigType("yaml")
	viper.SetConfigFile(configFile)
	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	return nil
}
