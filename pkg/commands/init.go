package commands

import (
	"fmt"

	"github.com/bearer/curio/pkg/commands/process/settings"
	"github.com/bearer/curio/pkg/util/output"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/xerrors"
)

func NewInitCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "writes default config in curio.yml",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := scanFlags.BindForConfigInit(cmd); err != nil {
				return fmt.Errorf("flag bind error: %w", err)
			}

			options, err := scanFlags.ToOptions(args, outputWriter)
			if err != nil {
				return xerrors.Errorf("flag error: %w", err)
			}
			globalSettings, err := settings.FromOptions(options)
			if err != nil {
				return err
			}
			viper.Set(settings.CustomDetectorKey, globalSettings.CustomDetector)

			viper.SetConfigFile("./curio.yml")
			err = viper.WriteConfig()
			if err != nil {
				return err
			}

			output.StdErrLogger().Msgf("created: curio.yml (default configuration file)")
			return nil
		},
	}

	scanFlags.AddFlags(cmd)

	return cmd
}
