package commands

import (
	"github.com/spf13/cobra"
)

func NewVersionCommand(version string, commitSHA string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.Printf("bearer version: %s\nsha: %s\n", version, commitSHA)
			return nil
		},
	}
	cmd.SetUsageFunc(func(cmd *cobra.Command) error {
		return nil
	})
	return cmd
}
