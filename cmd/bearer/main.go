package main

import (
	"os"

	"github.com/bearer/bearer/cmd/bearer/build"

	"github.com/bearer/bearer/pkg/commands"
	"github.com/bearer/bearer/pkg/util/output"
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
