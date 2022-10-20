package main

import (
	"log"

	"github.com/bearer/curio/cmd/curio/build"

	"github.com/bearer/curio/pkg/commands"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	app := commands.NewApp(build.Version)
	if err := app.Execute(); err != nil {
		return err
	}
	return nil
}
