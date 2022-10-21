package commands

import (
	"github.com/bearer/curio/pkg/util/output"
	"github.com/spf13/cobra"
)

func NewVersionCommand(version string, commitSHA string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			output.StdOutLogger().Msgf("curio version: %s\nsha: %s", version, commitSHA)
			return nil
		},
	}
	cmd.SetFlagErrorFunc(flagErrorFunc)

	return cmd
}
