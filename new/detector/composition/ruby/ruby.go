package ruby

import (
	"fmt"
	"os"

	"golang.org/x/exp/slices"

	"github.com/bearer/bearer/new/detector/composition/types"
	"github.com/bearer/bearer/new/detector/detection"
	"github.com/bearer/bearer/new/detector/evaluator"
	cachepkg "github.com/bearer/bearer/new/detector/evaluator/cache"
	"github.com/bearer/bearer/new/detector/evaluator/stats"
	"github.com/bearer/bearer/new/detector/implementation/custom"
	"github.com/bearer/bearer/new/detector/implementation/generic/datatype"
	"github.com/bearer/bearer/new/detector/implementation/generic/insecureurl"
	"github.com/bearer/bearer/new/detector/implementation/generic/stringliteral"
	"github.com/bearer/bearer/new/detector/implementation/ruby/object"
	"github.com/bearer/bearer/new/language"
	"github.com/rs/zerolog/log"

	"github.com/bearer/bearer/pkg/classification"
	"github.com/bearer/bearer/pkg/commands/process/settings"
	"github.com/bearer/bearer/pkg/report/customdetectors"
	"github.com/bearer/bearer/pkg/util/file"

	stringdetector "github.com/bearer/bearer/new/detector/implementation/ruby/string"
	detectorset "github.com/bearer/bearer/new/detector/set"
	detectortypes "github.com/bearer/bearer/new/detector/types"
	"github.com/bearer/bearer/new/language/implementation"
	"github.com/bearer/bearer/new/language/implementation/ruby"
	languagetypes "github.com/bearer/bearer/new/language/types"
)

type Composition struct {
	customDetectorTypes []string
	sharedDetectorTypes []string
	detectorSet         detectortypes.DetectorSet
	langImplementation  implementation.Implementation
	lang                languagetypes.Language
	closers             []func()
	rules               map[string]*settings.Rule
	stats               *stats.Stats
}

func New(
	debugProfile bool,
	rules map[string]*settings.Rule,
	classifier *classification.Classifier,
) (detectortypes.Composition, error) {
	lang, err := language.Get("ruby")
	if err != nil {
		return nil, fmt.Errorf("failed to lookup language: %s", err)
	}

	composition := &Composition{
		langImplementation: ruby.Get(),
		lang:               lang,
	}

	if debugProfile {
		composition.stats = stats.New()
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
			constructor: stringdetector.New,
			name:        "string detector",
		},
		{
			constructor: insecureurl.New,
			name:        "insecure url detector",
		},
		{
			constructor: stringliteral.New,
			name:        "string literal detector",
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

	composition.rules = rubyRules

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

	presenceRules := map[string]bool{}
	for _, rule := range rubyRules {
		if rule.Trigger.RequiredDetection != nil {
			presenceRules[*rule.Trigger.RequiredDetection] = true
		}
	}

	for ruleName, rule := range rubyRules {
		patterns := rule.Patterns
		localRuleName := ruleName

		if rule.Type == customdetectors.TypeShared {
			composition.sharedDetectorTypes = append(composition.sharedDetectorTypes, ruleName)
		} else {
			if !rule.IsAuxilary || presenceRules[ruleName] {
				composition.customDetectorTypes = append(composition.customDetectorTypes, ruleName)
			}
		}

		go func() {
			customDetector, err := custom.New(
				lang,
				localRuleName,
				patterns,
				rubyRules,
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
		if response.Error != nil {
			composition.Close()
			return nil, fmt.Errorf("failed to create detector %s: %s", response.DetectorName, response.Error)
		}
		detectors = append(detectors, response.Detector)
		composition.closers = append(composition.closers, response.Detector.Close)
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
	if composition.stats != nil {
		log.Debug().Msgf("ruby stats:\n%s", composition.stats)
	}

	for _, closeFunc := range composition.closers {
		closeFunc()
	}
}

func (composition *Composition) DetectFromFile(file *file.FileInfo) ([]*detection.Detection, error) {
	return composition.DetectFromFileWithTypes(
		file,
		composition.customDetectorTypes,
		composition.sharedDetectorTypes,
	)
}

func (composition *Composition) DetectFromFileWithTypes(
	file *file.FileInfo,
	detectorTypes, sharedDetectorTypes []string,
) ([]*detection.Detection, error) {
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

	evaluator := evaluator.New(
		composition.langImplementation,
		composition.lang,
		composition.detectorSet,
		tree,
		file.FileInfo.Name(),
		composition.stats,
	)

	sharedCache := cachepkg.NewShared(sharedDetectorTypes)

	var result []*detection.Detection
	for _, detectorType := range detectorTypes {
		cache := cachepkg.NewCache(sharedCache)
		rule := composition.rules[detectorType]
		sanitizerRuleID := ""
		if rule != nil {
			sanitizerRuleID = rule.SanitizerRuleID
		}
		detections, err := evaluator.Evaluate(
			tree.RootNode(),
			detectorType,
			sanitizerRuleID,
			cache,
			settings.DefaultScope,
			false,
		)
		if err != nil {
			return nil, err
		}

		result = append(result, detections...)
	}

	return result, nil
}
