package ruby

import (
	"fmt"
	"os"

	"github.com/bearer/curio/new/detector/evaluator"
	"github.com/bearer/curio/new/detector/implementation/custom"
	"github.com/bearer/curio/new/detector/implementation/generic/datatype"
	"github.com/bearer/curio/new/detector/implementation/ruby/object"
	"github.com/bearer/curio/new/detector/implementation/ruby/property"
	detectorset "github.com/bearer/curio/new/detector/set"
	"github.com/bearer/curio/new/detector/types"
	detectortypes "github.com/bearer/curio/new/detector/types"
	languagetypes "github.com/bearer/curio/new/language/types"

	"github.com/bearer/curio/new/language"
	"github.com/bearer/curio/pkg/commands/process/settings"
	"github.com/bearer/curio/pkg/util/file"
)

type Composition struct {
	detectorsSet detectortypes.DetectorSet
	lang         languagetypes.Language
	closers      []func()
}

func New(rules map[string]settings.Rule) (types.Composition, error) {
	lang, err := language.Get("ruby")
	if err != nil {
		return nil, fmt.Errorf("failed to lookup language: %s", err)
	}

	composition := &Composition{
		lang: lang,
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

	// instantiate custom ruby detectors
	for ruleName, rule := range rules {
		for i, pattern := range rule.Patterns {
			// todo: Figure out how to have multiple patterns for same ruleName, or support in dataflow and policies multiple rule_names

			customDetector, err := custom.New(
				lang,
				fmt.Sprintf("%s_%d", ruleName, i),
				custom.Rule{
					Pattern: pattern.Pattern,
					Filters: pattern.Filters.ToCustomFilters(),
				},
			)
			if err != nil {
				composition.Close()
				return nil, fmt.Errorf("failed to create custom detector %s for pattern count %d: %s", ruleName, i, err)
			}
			detectors = append(detectors, customDetector)
			composition.closers = append(composition.closers, customDetector.Close)
		}
	}

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

func (composition *Composition) DetectFromFile(file *file.FileInfo) ([]*detectortypes.Detection, error) {
	if file.Language != "Ruby" {
		return nil, nil
	}

	fileContent, err := os.ReadFile(file.AbsolutePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s", err)
	}

	tree, err := composition.lang.Parse(string(fileContent))
	if err != nil {
		return nil, fmt.Errorf("failed to parse file %s", err)
	}

	evaluator := evaluator.New(composition.detectorsSet, tree)

	return evaluator.ForTree(tree.RootNode(), "ruby_file_detection")
}
