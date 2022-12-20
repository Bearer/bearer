package main

import (
	"os"

	"github.com/bearer/curio/cmd/curio/build"

	"github.com/bearer/curio/pkg/commands"
	"github.com/bearer/curio/pkg/util/output"
)

func main() {
	if err := run(); err != nil {
		output.StdErrLogger().Msgf("%s", err)
		os.Exit(1)
	}
}

func run() error {
	app := commands.NewApp(build.Version, build.CommitSHA)
	if err := app.Execute(); err != nil {
		return err
	}
	return nil
}
