package commands

import (
	"fmt"

	"github.com/bearer/curio/pkg/commands/process/settings"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/xerrors"
)

func NewInitCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Generates a default config to `curio.yml`",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := scanFlags.BindForConfigInit(NewScanCommand()); err != nil {
				return fmt.Errorf("flag bind error: %w", err)
			}

			options, err := scanFlags.ToOptions(args)
			if err != nil {
				return xerrors.Errorf("flag error: %w", err)
			}
			globalSettings, err := settings.FromOptions(options)
			if err != nil {
				return err
			}
			viper.Set(settings.CustomDetectorKey, globalSettings.CustomDetector)
			viper.Set(settings.PoliciesKey, globalSettings.Policies)

			viper.SetConfigFile("./curio.yml")
			err = viper.WriteConfig()
			if err != nil {
				return err
			}

			cmd.PrintErrln("Created: curio.yml (default configuration file)")
			return nil
		},
	}

	cmd.SetUsageFunc(func(cmd *cobra.Command) error {
		return nil
	})

	return cmd
}
