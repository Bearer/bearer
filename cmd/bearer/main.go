package main

import (
	"os"

	"github.com/bearer/bearer/cmd/bearer/build"
	"github.com/bearer/bearer/internal/commands"
)

func main() {
	app := commands.NewApp(build.Version, build.CommitSHA)
	if err := app.Execute(); err != nil {
		// error messages are printed by the framework
		os.Exit(1)
	}
}
