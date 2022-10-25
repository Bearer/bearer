package commands

import (
	"encoding/json"
	"fmt"

	config "github.com/bearer/curio/pkg/commands/process/settings"
	"github.com/bearer/curio/pkg/commands/process/worker"
	"github.com/bearer/curio/pkg/flag"
	"github.com/bearer/curio/pkg/util/output"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func NewProcessingWorkerCommand() *cobra.Command {
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

			generalOptions, err := flags.ToOptions(args, outputWriter)
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
