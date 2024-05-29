package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewInitCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Generates a default config to `bearer.yml`",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := ScanFlags.BindForConfigInit(NewScanCommand(nil)); err != nil {
				return fmt.Errorf("flag bind error: %w", err)
			}

			viper.SetConfigFile("./bearer.yml")
			err := viper.WriteConfig()
			if err != nil {
				return err
			}

			cmd.PrintErrln("Created: bearer.yml (default configuration file)")
			return nil
		},
	}

	cmd.SetUsageFunc(func(cmd *cobra.Command) error {
		return nil
	})

	return cmd
}
