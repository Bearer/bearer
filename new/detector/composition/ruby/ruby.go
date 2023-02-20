package ruby

import (
	"fmt"
	"os"

	"golang.org/x/exp/slices"

	"github.com/bearer/curio/new/detector/composition/types"
	"github.com/bearer/curio/new/detector/evaluator"
	"github.com/bearer/curio/new/detector/implementation/custom"
	"github.com/bearer/curio/new/detector/implementation/generic/datatype"
	"github.com/bearer/curio/new/detector/implementation/generic/insecureurl"
	"github.com/bearer/curio/new/detector/implementation/ruby/object"
	"github.com/bearer/curio/new/detector/implementation/ruby/property"
	"github.com/bearer/curio/new/language"

	"github.com/bearer/curio/pkg/classification"
	"github.com/bearer/curio/pkg/commands/process/settings"
	"github.com/bearer/curio/pkg/util/file"

	stringdetector "github.com/bearer/curio/new/detector/implementation/ruby/string"
	detectorset "github.com/bearer/curio/new/detector/set"
	detectortypes "github.com/bearer/curio/new/detector/types"
	languagetypes "github.com/bearer/curio/new/language/types"
)

type Composition struct {
	customDetectorTypes []string
	detectorSet         detectortypes.DetectorSet
	lang                languagetypes.Language
	closers             []func()
}

func New(rules map[string]*settings.Rule, classifier *classification.Classifier) (detectortypes.Composition, error) {
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
			constructor: object.New,
			name:        "object detector",
		},
		{
			constructor: property.New,
			name:        "property detector",
		},
		{
			constructor: stringdetector.New,
			name:        "string detector",
		},
		{
			constructor: insecureurl.New,
			name:        "insecure url detector",
		},
	}

	// instantiate custom ruby detectors
	rubyRules := map[string]*settings.Rule{}
	for ruleName, rule := range rules {
		if !slices.Contains(rule.Languages, "ruby") {
			continue
		}
		rubyRules[ruleName] = rule
	}

	detectorsLen := len(rubyRules) + len(staticDetectors)
	receiver := make(chan types.DetectorInitResult, detectorsLen)

	var detectors []detectortypes.Detector

	for _, detectorCreator := range staticDetectors {
		creator := detectorCreator
		go func() {
			detector, err := creator.constructor(lang)
			receiver <- types.DetectorInitResult{
				Error:        err,
				Detector:     detector,
				DetectorName: creator.name,
			}
		}()
	}

	detector, err := datatype.New(lang, classifier.Schema)
	if err != nil {
		composition.Close()
		return nil, fmt.Errorf("failed to create datatype detector: %s", err)
	}
	detectors = append(detectors, detector)
	composition.closers = append(composition.closers, detector.Close)

	for ruleName, rule := range rubyRules {
		patterns := rule.Patterns
		localRuleName := ruleName

		composition.customDetectorTypes = append(composition.customDetectorTypes, ruleName)
		go func() {
			customDetector, err := custom.New(
				lang,
				localRuleName,
				patterns,
			)

			receiver <- types.DetectorInitResult{
				Error:        err,
				Detector:     customDetector,
				DetectorName: "customDetector:" + localRuleName,
			}
		}()
	}

	for i := 0; i < detectorsLen; i++ {
		response := <-receiver
		detectors = append(detectors, response.Detector)
		composition.closers = append(composition.closers, response.Detector.Close)
		if response.Error != nil {
			composition.Close()
			return nil, fmt.Errorf("failed to create detector %s: %s", response.DetectorName, err)
		}
	}

	detectorSet, err := detectorset.New(detectors)
	if err != nil {
		composition.Close()
		return nil, fmt.Errorf("failed to create detector set: %s", err)
	}
	composition.detectorSet = detectorSet

	return composition, nil
}

func (composition *Composition) Close() {
	for _, closeFunc := range composition.closers {
		closeFunc()
	}
}

func (composition *Composition) DetectFromFile(file *file.FileInfo) ([]*detectortypes.Detection, error) {
	return composition.DetectFromFileWithTypes(file, composition.customDetectorTypes)
}

func (composition *Composition) DetectFromFileWithTypes(file *file.FileInfo, detectorTypes []string) ([]*detectortypes.Detection, error) {
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

	evaluator := evaluator.New(composition.lang, composition.detectorSet, tree, file.FileInfo.Name())

	var result []*detectortypes.Detection
	for _, detectorType := range detectorTypes {
		detections, err := evaluator.ForTree(tree.RootNode(), detectorType, false)
		if err != nil {
			return nil, err
		}

		result = append(result, detections...)
	}

	return result, nil
}
