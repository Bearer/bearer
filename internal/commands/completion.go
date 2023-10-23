package commands

import (
	"os"

	"github.com/spf13/cobra"
)

func NewCompletionCommand() *cobra.Command {
	usageTemplate := `
Usage: bearer completion [command]

Available Commands:
    bash        Generate the autocompletion script for bash
    fish        Generate the autocompletion script for fish
    powershell  Generate the autocompletion script for powershell
    zsh         Generate the autocompletion script for zsh
`

	cmd := &cobra.Command{
		Use:                   "completion [command]",
		Short:                 "Generate the autocompletion script for the your shell.",
		SilenceErrors:         false,
		SilenceUsage:          false,
		DisableFlagsInUseLine: true,
		ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
		Args:                  cobra.MatchAll(cobra.ExactArgs(1)),
		Run: func(cmd *cobra.Command, args []string) {
			switch args[0] {
			case "bash":
				cmd.Root().GenBashCompletion(os.Stdout) //nolint:errcheck
			case "zsh":
				cmd.Root().GenZshCompletion(os.Stdout) //nolint:errcheck
			case "fish":
				cmd.Root().GenFishCompletion(os.Stdout, true) //nolint:errcheck
			case "powershell":
				cmd.Root().GenPowerShellCompletionWithDesc(os.Stdout) //nolint:errcheck
			}
		},
	}

	cmd.SetUsageTemplate(usageTemplate)

	return cmd
}
