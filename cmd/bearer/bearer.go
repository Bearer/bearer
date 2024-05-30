package main

import (
	"github.com/bearer/bearer/cmd/bearer/build"
	"github.com/bearer/bearer/external/run"
)

func main() {
	run.Run(build.Version, build.CommitSHA, run.NewEngine(run.DefaultLanguages()))
}
