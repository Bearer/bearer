package stats

import (
	"github.com/hhatto/gocloc"
)

func GoclocDetectorOutput(path string) (*gocloc.Result, error) {
	clocOpts := gocloc.NewClocOptions()
	clocOpts.SkipDuplicated = true

	languages := gocloc.NewDefinedLanguages()
	processor := gocloc.NewProcessor(languages, clocOpts)

	return processor.Analyze([]string{path})
}
