package commands

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/bearer/bearer/internal/flag"
	"github.com/bearer/bearer/internal/util/output"
	"github.com/bearer/bearer/internal/version_check"
)

func NewVersionCommand(version string, commitSHA string) *cobra.Command {
	var flags = flag.Flags{flag.GeneralFlagGroup}
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := flags.Bind(cmd); err != nil {
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

			meta, err := version_check.GetVersionMeta(cmd.Context(), make([]string, 0))
			if err != nil {
				log.Debug().Msgf("failed: %s", err)
			} else {
				version_check.DisplayBinaryVersionWarning(meta, false)
			}
			cmd.Printf("bearer version %s, build %s\n", version, commitSHA)
			return nil
		},
	}
	cmd.SetUsageFunc(func(cmd *cobra.Command) error {
		return nil
	})
	flags.AddFlags(cmd)
	return cmd
}
