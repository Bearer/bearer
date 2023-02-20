package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/bearer/curio/cmd/curio/build"
	"github.com/bearer/curio/pkg/commands"
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
		if err := writeDocs(c, dir); err != nil {
			return err
		}
	}

	// create a file to use
	basename := "bearer.yaml"

	if cmd.CommandPath() != "" {
		basename = strings.Replace(cmd.CommandPath(), " ", "bearer_", -1) + ".yaml"
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
