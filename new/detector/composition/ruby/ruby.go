package ruby

import (
	"fmt"
	"os"

	"golang.org/x/exp/slices"

	"github.com/bearer/bearer/new/detector/composition/types"
	"github.com/bearer/bearer/new/detector/evaluator"
	"github.com/bearer/bearer/new/detector/implementation/custom"
	"github.com/bearer/bearer/new/detector/implementation/generic/datatype"
	"github.com/bearer/bearer/new/detector/implementation/generic/insecureurl"
	"github.com/bearer/bearer/new/detector/implementation/ruby/object"
	"github.com/bearer/bearer/new/detector/implementation/ruby/property"
	"github.com/bearer/bearer/new/language"
	"github.com/rs/zerolog/log"

	"github.com/bearer/bearer/pkg/classification"
	"github.com/bearer/bearer/pkg/commands/process/settings"
	"github.com/bearer/bearer/pkg/util/file"

	stringdetector "github.com/bearer/bearer/new/detector/implementation/ruby/string"
	detectorset "github.com/bearer/bearer/new/detector/set"
	detectortypes "github.com/bearer/bearer/new/detector/types"
	languagetypes "github.com/bearer/bearer/new/language/types"
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

	detectorIterator := 0

	for _, detectorCreator := range staticDetectors {
		localIterator := detectorIterator
		creator := detectorCreator
		detectorIterator++
		go func() {
			log.Debug().Msgf("adding ruby rule %d", localIterator)
			detector, err := creator.constructor(lang)
			receiver <- types.DetectorInitResult{
				Error:        err,
				Detector:     detector,
				DetectorName: creator.name,
				Order:        localIterator,
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
		localIterator := detectorIterator

		composition.customDetectorTypes = append(composition.customDetectorTypes, ruleName)
		detectorIterator++
		log.Debug().Msgf("adding ruby custom rule %d", localIterator)

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
				Order:        localIterator,
			}
		}()
	}

	constructedDetectors := map[int]types.DetectorInitResult{}
	for i := 0; i < detectorsLen; i++ {
		response := <-receiver
		if response.Error != nil {
			composition.Close()
			return nil, fmt.Errorf("failed to create detector %s: %s", response.DetectorName, response.Error)
		}
		constructedDetectors[response.Order] = response
	}

	for _, constructed := range constructedDetectors {
		detectors = append(detectors, constructed.Detector)
		composition.closers = append(composition.closers, constructed.Detector.Close)
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
