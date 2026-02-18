package commands

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/bearer/bearer/pkg/commands/debugprofile"
	"github.com/bearer/bearer/pkg/commands/process/orchestrator/worker"
	"github.com/bearer/bearer/pkg/engine"
	"github.com/bearer/bearer/pkg/flag"
	"github.com/bearer/bearer/pkg/util/output"
)

func NewProcessingWorkerCommand(engine engine.Engine) *cobra.Command {
	flags := flag.Flags{flag.WorkerFlagGroup}

	cmd := &cobra.Command{
		Use:   "processing-worker [flags] PATH",
		Short: "start scan processing server",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := flags.Bind(cmd); err != nil {
				return fmt.Errorf("flag bind error: %w", err)
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
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

			log.Debug().Msgf("running scan worker on port `%s`", options.Port)
			return worker.Start(options.ParentProcessID, options.Port, engine)
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
