package ruby

import (
	"fmt"

	"github.com/bearer/curio/new/detector/implementation/custom"
	"github.com/bearer/curio/new/detector/implementation/generic/datatype"
	"github.com/bearer/curio/new/detector/implementation/ruby/object"
	"github.com/bearer/curio/new/detector/implementation/ruby/property"
	detectorset "github.com/bearer/curio/new/detector/set"
	detectortypes "github.com/bearer/curio/new/detector/types"
	languagetypes "github.com/bearer/curio/new/language/types"

	"github.com/bearer/curio/new/language"
	"github.com/bearer/curio/pkg/util/file"
)

type Composition struct {
	detectorsSet detectortypes.DetectorSet
	closers      []func()
}

func New() (*Composition, error) {
	composition := &Composition{}

	lang, err := language.Get("ruby")
	if err != nil {
		return nil, fmt.Errorf("failed to lookup language: %s", err)
	}

	staticDetectors := []struct {
		constructor func(languagetypes.Language) (detectortypes.Detector, error)
		name        string
	}{
		{
			constructor: property.New,
			name:        "property detector",
		},
		{
			constructor: object.New,
			name:        "object detector",
		},
		{
			constructor: datatype.New,
			name:        "datatype detector",
		},
	}

	var detectors []detectortypes.Detector

	for _, detectorCreator := range staticDetectors {
		detector, err := detectorCreator.constructor(lang)
		if err != nil {
			composition.Close()
			return nil, fmt.Errorf("failed to create %s: %s", detectorCreator.name, err)
		}
		detectors = append(detectors, detector)
		composition.closers = append(composition.closers, detector.Close)
	}

	rubyFileDetector, err := custom.New(
		lang,
		"ruby_file_detection",
		custom.Rule{
			Pattern: `
				Sentry.init do |$<CONFIG:identifier>|
					$<CONFIG>.before_breadcrumb = lambda do |$<BREADCRUMB:identifier>, hint|
						$<BREADCRUMB>.message = $<MESSAGE>
					end
				end`,
			Filters: []custom.Filter{
				{
					Variable:  "MESSAGE",
					Detection: "datatype",
				},
			},
		},
	)

	if err != nil {
		composition.Close()
		return nil, fmt.Errorf("failed to create ruby file detector: %s", err)
	}
	composition.closers = append(composition.closers, rubyFileDetector.Close)
	detectors = append(detectors, rubyFileDetector)

	detectorSet, err := detectorset.New(detectors)
	if err != nil {
		composition.Close()
		return nil, fmt.Errorf("failed to create detector set: %s", err)
	}
	composition.detectorsSet = detectorSet

	return composition, nil
}

func (composition *Composition) Close() {
	for _, closeFunc := range composition.closers {
		closeFunc()
	}
}

func (composition *Composition) ParseFile(file *file.FileInfo) {
	if file.Language != "Ruby" {
		return
	}
}
