package commands

import (
	"fmt"

	"github.com/bearer/bearer/internal/commands/artifact"
	"github.com/bearer/bearer/internal/commands/debugprofile"
	"github.com/bearer/bearer/internal/flag"
	"github.com/bearer/bearer/internal/util/file"
	"github.com/bearer/bearer/internal/util/output"
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

var ScanFlags = flag.Flags{
	flag.ReportFlagGroup,
	flag.RuleFlagGroup,
	flag.ScanFlagGroup,
	flag.RepositoryFlagGroup,
	flag.GeneralFlagGroup,
}

func NewScanCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "scan [flags] <path>",
		Aliases: []string{"s"},
		Short:   "Scan a directory or file",
		Example: `  # Scan a local project, including language-specific files
  $ bearer scan /path/to/your_project`,
		RunE: func(cmd *cobra.Command, args []string) error {
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

			_, loadFileMessage, _ := readConfig(args)
			log.Debug().Msgf(loadFileMessage)

			options, err := ScanFlags.ToOptions(args)
			if err != nil {
				return fmt.Errorf("flag error: %s", err)
			}

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
	ScanFlags.Bind(cmd) // nolint:errcheck
	cmd.SetUsageTemplate(fmt.Sprintf(scanTemplate, ScanFlags.Usages(cmd)))

	return cmd
}

func readConfig(args []string) (string, string, error) {
	configPath := viper.GetString(flag.ConfigFileFlag.ConfigName)
	var loadFileMessage string
	if err := readConfigFromPath(configPath); err != nil {
		// load from default
		var configPath = ""
		if len(args) > 0 {
			configPath = file.GetFullFilename(args[0], configPath)
		}
		if err := readConfigFromPath(configPath); err != nil {
			return configPath, "Couldn't find any config file", err
		} else {
			loadFileMessage = fmt.Sprintf("Loading default config file %s", configPath)
		}
	} else {
		loadFileMessage = fmt.Sprintf("Loading config file %s", configPath)
	}

	return configPath, loadFileMessage, nil
}

func readConfigFromPath(configFile string) error {
	viper.SetConfigType("yaml")
	viper.SetConfigFile(configFile)
	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	return nil
}
