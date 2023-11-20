package commands

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/bearer/bearer/internal/commands/debugprofile"
	"github.com/bearer/bearer/internal/commands/process/orchestrator/worker"
	"github.com/bearer/bearer/internal/flag"
	"github.com/bearer/bearer/internal/util/output"
)

func NewProcessingWorkerCommand() *cobra.Command {
	flags := flag.Flags{flag.WorkerFlagGroup}

	cmd := &cobra.Command{
		Use:   "processing-worker [flags] PATH",
		Short: "start scan processing server",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := flags.Bind(cmd); err != nil {
				return fmt.Errorf("flag bind error: %w", err)
			}

			output.Setup(cmd, output.SetupRequest{
				LogLevel:  viper.GetString(flag.LogLevelFlag.ConfigName),
				Quiet:     viper.GetBool(flag.QuietFlag.ConfigName),
				ProcessID: viper.GetString(flag.WorkerIDFlag.ConfigName),
			})

			if viper.GetBool(flag.DebugProfileFlag.ConfigName) {
				debugprofile.Start()
			}

			options, err := flags.ToOptions(args)
			if err != nil {
				return fmt.Errorf("flag error: %s", err)
			}

			log.Debug().Msgf("running scan worker on port `%s`", options.WorkerOptions.Port)
			return worker.Start(options.WorkerOptions.ParentProcessID, options.WorkerOptions.Port)
		},
		Hidden:        true,
		SilenceErrors: true,
		SilenceUsage:  true,
	}
	cmd.SetFlagErrorFunc(func(command *cobra.Command, err error) error {
		if err := command.Help(); err != nil {
			return err
		}
		command.Println() // add empty line after list of flags
		return err
	})
	flags.AddFlags(cmd)

	return cmd
}
