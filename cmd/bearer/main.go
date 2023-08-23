package main

import (
	"github.com/bearer/bearer/cmd/bearer/build"

	"github.com/bearer/bearer/internal/commands"
	"github.com/bearer/bearer/internal/util/output"
)

func main() {
	if err := run(); err != nil {
		output.Fatal(err.Error())
	}
}

func run() error {
	app := commands.NewApp(build.Version, build.CommitSHA)
	if err := app.Execute(); err != nil {
		return err
	}
	return nil
}
