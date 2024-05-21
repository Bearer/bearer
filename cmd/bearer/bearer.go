package main

import (
	"os"

	"github.com/bearer/bearer/cmd/bearer/build"
	"github.com/bearer/bearer/internal/commands"
	"github.com/bearer/bearer/internal/languages/golang"
	"github.com/bearer/bearer/internal/languages/java"
	"github.com/bearer/bearer/internal/languages/javascript"
	"github.com/bearer/bearer/internal/languages/php"
	"github.com/bearer/bearer/internal/languages/python"
	"github.com/bearer/bearer/internal/languages/ruby"
	"github.com/bearer/bearer/internal/scanner/language"
)

func main() {
	languageset.Register([]language.Language{
		java.Get(),
		javascript.Get(),
		ruby.Get(),
		php.Get(),
		golang.Get(),
		python.Get(),
	})

	app := commands.NewApp(build.Version, build.CommitSHA)
	if err := app.Execute(); err != nil {
		// error messages are printed by the framework
		os.Exit(1)
	}
}
