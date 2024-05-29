package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/bearer/bearer/cmd/bearer/build"
	"github.com/bearer/bearer/pkg/engine"
)

// NewApp is the factory method to return CLI
func NewApp(version string, commitSHA string, engine engine.Engine) *cobra.Command {
	rootCmd := NewRootCommand()
	rootCmd.AddCommand(
		NewCompletionCommand(),
		NewProcessingWorkerCommand(),
		NewInitCommand(),
		NewScanCommand(engine),
		NewIgnoreCommand(),
		NewVersionCommand(version, commitSHA),
	)

	return rootCmd
}

func NewRootCommand() *cobra.Command {
	usageTemplate := `
[0;1;34;94m▄▄▄▄▄[0m
[0;34m█[0m    [0;34m█[0m  [0;34m▄▄▄[0m    [0;37m▄▄▄[0m    [0;37m▄[0m [0;37m▄▄[0m   [0;37m▄▄[0;1;30;90m▄[0m    [0;1;30;90m▄[0m [0;1;30;90m▄▄[0m
[0;34m█▄▄▄▄▀[0m [0;37m█▀[0m  [0;37m█[0m  [0;37m▀[0m   [0;37m█[0m   [0;37m█[0;1;30;90m▀[0m  [0;1;30;90m▀[0m [0;1;30;90m█▀[0m  [0;1;30;90m█[0m   [0;1;30;90m█▀[0m  [0;1;34;94m▀[0m
[0;37m█[0m    [0;37m█[0m [0;37m█▀▀▀▀[0m  [0;37m▄[0;1;30;90m▀▀▀█[0m   [0;1;30;90m█[0m     [0;1;30;90m█▀▀[0;1;34;94m▀▀[0m   [0;1;34;94m█[0m
[0;37m█▄▄▄▄▀[0m [0;1;30;90m▀█▄▄▀[0m  [0;1;30;90m▀▄▄▀█[0m   [0;1;30;90m█[0m     [0;1;34;94m▀█▄▄▀[0m   [0;1;34;94m█[0m

Scan your source code to discover, filter and prioritize security and privacy risks.

Usage: bearer <command> [flags]

Available Commands:
	completion        Generate the autocompletion script for your shell
	scan              Scan a directory or file
	init              Write the default config to bearer.yml
	ignore            Manage ignored fingerprints
	version           Print the version

Examples:
	# Scan local directory or file and output security risks
	$ bearer scan <path> --scanner=sast,secrets

	# Scan current directory and output the privacy report to a file
	$ bearer scan --report privacy --output <output-path> .

	# Scan local directory and output details about the underlying
	# detections and classifications
	$ bearer scan . --report dataflow

Learn More:
	Bearer is a code security tool that scans your source code to
	identify OWASP/CWE top security risks, privacy impact, and sensitive dataflows

	For more examples, tutorials, and to learn more about the project
	visit https://docs.bearer.com
`
	version := fmt.Sprintf("%s, build %s", build.Version, build.CommitSHA)
	cmd := &cobra.Command{
		Use:     "bearer",
		Args:    cobra.NoArgs,
		Version: version,
	}
	cmd.SetUsageTemplate(usageTemplate)

	return cmd
}
