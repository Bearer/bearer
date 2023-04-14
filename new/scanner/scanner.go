package scanner

import (
	"fmt"

	"github.com/bearer/bearer/new/detector/composition"
	"github.com/bearer/bearer/new/detector/composition/java"
	"github.com/bearer/bearer/new/detector/composition/javascript"
	"github.com/bearer/bearer/new/detector/composition/ruby"
	"github.com/bearer/bearer/new/detector/types"
	"github.com/bearer/bearer/pkg/classification"
	"github.com/bearer/bearer/pkg/commands/process/settings"
	"github.com/bearer/bearer/pkg/report"
	"github.com/bearer/bearer/pkg/util/file"
)

type language struct {
	name        string
	composition types.Composition
}

type scannerType []language

var scanner scannerType

func (scanner scannerType) Close() {
	for _, language := range scanner {
		language.composition.Close()
	}
}

func Setup(config *settings.Config, classifier *classification.Classifier) (err error) {
	var toInstantiate = []struct {
		constructor func(map[string]*settings.Rule, *classification.Classifier) (types.Composition, error)
		name        string
	}{
		{
			constructor: ruby.New,
			name:        "ruby",
		},
		{
			constructor: javascript.New,
			name:        "javascript",
		},
		{
			constructor: java.New,
			name:        "java",
		},
	}

	for _, instantiatior := range toInstantiate {
		composition, err := instantiatior.constructor(config.Rules, classifier)
		if err != nil {
			return fmt.Errorf("failed to instantiate composition %s:%s", instantiatior.name, err)
		}

		scanner = append(scanner, language{
			name:        instantiatior.name,
			composition: composition,
		})
	}

	return err
}

func Detect(report report.Report, file *file.FileInfo) (err error) {
	for _, language := range scanner {
		detections, err := language.composition.DetectFromFile(file)
		if err != nil {
			return fmt.Errorf("%s failed to detect in file %s: %s", language.name, file.AbsolutePath, err)
		}

		composition.ReportDetections(report, file, detections)
	}

	return nil
}
