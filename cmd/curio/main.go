package main

import (
	"log"

	"github.com/Bearer/curio/pkg/commands"
)

var (
	version = "dev"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	app := commands.NewApp(version)
	if err := app.Execute(); err != nil {
		return err
	}
	return nil
}
