package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/bearer/bearer/cmd/bearer/build"
	"github.com/bearer/bearer/pkg/commands"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

func main() {
	dir := "./docs/_data"
	if _, err := os.Stat(dir); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	err := writeDocs(commands.NewApp(build.Version, build.CommitSHA), dir)

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

func writeDocs(cmd *cobra.Command, dir string) error {
	// Exit if there's an error
	for _, c := range cmd.Commands() {
		if !c.IsAvailableCommand() || c.IsAdditionalHelpTopicCommand() {
			continue
		}
		if c.HasSubCommands() {
			for _, subCmd := range c.Commands() {
				if err := writeDocs(subCmd, dir); err != nil {
					return err
				}
			}
			continue
		}
		if err := writeDocs(c, dir); err != nil {
			return err
		}
	}

	// create a file to use
	basename := "bearer.yaml"

	if cmd.CommandPath() != "" {
		basename = fmt.Sprintf("bearer%s.yaml", strings.Replace(cmd.CommandPath(), " ", "_", -1))
	}

	filename := filepath.Join(dir, basename)
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	// given the file
	err = doc.GenYaml(cmd, f)
	if err != nil {
		return err
	}

	// add aliases
	aliases := fmt.Sprintf("aliases: %s\n", strings.Join(cmd.Aliases, ", "))

	if _, err := f.WriteString(aliases); err != nil {
		return err
	}
	return nil
}
