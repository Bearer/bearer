package main

import (
	"log"

	"github.com/bearer/curio/pkg/commands"
)

var (
	Version = "dev"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	app := commands.NewApp(Version)
	if err := app.Execute(); err != nil {
		return err
	}
	return nil
}
