package ruby

import (
	"fmt"
	"os"
	"strings"

	"github.com/gertd/go-pluralize"
	"github.com/rs/zerolog/log"
	"golang.org/x/exp/slices"

	"github.com/bearer/curio/new/detector/evaluator"
	"github.com/bearer/curio/new/detector/implementation/custom"
	"github.com/bearer/curio/new/detector/implementation/generic/datatype"
	"github.com/bearer/curio/new/detector/implementation/generic/insecureurl"
	"github.com/bearer/curio/new/detector/implementation/ruby/object"
	"github.com/bearer/curio/new/detector/implementation/ruby/property"
	"github.com/bearer/curio/new/language"
	"github.com/bearer/curio/new/language/tree"

	"github.com/bearer/curio/pkg/classification"
	"github.com/bearer/curio/pkg/commands/process/settings"
	"github.com/bearer/curio/pkg/report/detectors"
	"github.com/bearer/curio/pkg/report/schema"
	"github.com/bearer/curio/pkg/report/source"
	"github.com/bearer/curio/pkg/util/file"

	compositiontypes "github.com/bearer/curio/new/detector/composition/types"
	stringdetector "github.com/bearer/curio/new/detector/implementation/ruby/string"
	detectorset "github.com/bearer/curio/new/detector/set"
	detectortypes "github.com/bearer/curio/new/detector/types"
	languagetypes "github.com/bearer/curio/new/language/types"
	reportdetections "github.com/bearer/curio/pkg/report/detections"
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

	detector, err := datatype.New(lang, classifier.Schema)
	if err != nil {
		composition.Close()
		return nil, fmt.Errorf("failed to create datatype detector: %s", err)
	}
	detectors = append(detectors, detector)
	composition.closers = append(composition.closers, detector.Close)

	// instantiate custom ruby detectors
	for ruleName, rule := range rules {
		if !slices.Contains(rule.Languages, "ruby") {
			continue
		}

		composition.customDetectorTypes = append(composition.customDetectorTypes, ruleName)

		customDetector, err := custom.New(
			lang,
			ruleName,
			rule.Patterns,
		)
		if err != nil {
			composition.Close()
			return nil, fmt.Errorf("failed to create rule %s: %s", ruleName, err)
		}
		detectors = append(detectors, customDetector)
		composition.closers = append(composition.closers, customDetector.Close)
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

func (composition *Composition) DetectFromFile(file *file.FileInfo) ([]compositiontypes.Detection, error) {
	if file.Language != "Ruby" {
		return nil, nil
	}

	log.Debug().Msgf("reading file %s", file.AbsolutePath)
	fileContent, err := os.ReadFile(file.AbsolutePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s", err)
	}

	log.Debug().Msgf("file content is %s", fileContent)

	tree, err := composition.lang.Parse(string(fileContent))
	if err != nil {
		return nil, fmt.Errorf("failed to parse file %s", err)
	}

	evaluator := evaluator.New(composition.lang, composition.detectorSet, tree, file.FileInfo.Name())

	return composition.extractCustomDetectors(evaluator, tree, file)

}

func (composition *Composition) extractCustomDetectors(evaluator detectortypes.Evaluator, tree *tree.Tree, file *file.FileInfo) ([]compositiontypes.Detection, error) {
	customDetections := []compositiontypes.Detection{}
	pluralizer := pluralize.NewClient()

	for _, detectorType := range composition.customDetectorTypes {
		detections, err := evaluator.ForTree(tree.RootNode(), detectorType)

		if err != nil {
			return nil, err
		}

		for _, detection := range detections {
			data := detection.Data.(custom.Data)

			if len(data.Datatypes) == 0 {
				customDetections = append(customDetections, compositiontypes.Detection{
					CustomDetector: detectors.Type(detectorType),
					DetectionType:  reportdetections.TypeCustomRisk,
					Source: source.New(
						file,
						file.Path,
						detection.MatchNode.LineNumber(),
						detection.MatchNode.ColumnNumber(),
						data.Pattern,
					),
					Value: schema.Parent{
						LineNumber: detection.MatchNode.LineNumber(),
						Content:    detection.MatchNode.Content(),
					},
				})

				continue
			}

			for _, datatypeDetection := range data.Datatypes {
				data := datatypeDetection.Data.(datatype.Data)

				customDetections = append(customDetections, compositiontypes.Detection{
					CustomDetector: detectors.Type(detectorType),
					DetectionType:  reportdetections.TypeCustomClassified,
					Source: source.New(
						file,
						file.Path,
						datatypeDetection.MatchNode.LineNumber(),
						datatypeDetection.MatchNode.ColumnNumber(),
						"",
					),
					Value: schema.Schema{
						ObjectName:           data.Name,
						NormalizedObjectName: pluralizer.Singular(strings.ToLower(data.Name)),
						Classification:       data.Classification,
						Parent: &schema.Parent{
							LineNumber: datatypeDetection.MatchNode.LineNumber(),
							Content:    datatypeDetection.MatchNode.Content(),
						},
					},
				})

				for _, property := range data.Properties {

					customDetections = append(customDetections, compositiontypes.Detection{
						CustomDetector: detectors.Type(detectorType),
						DetectionType:  reportdetections.TypeCustomClassified,
						Source: source.New(
							file,
							file.Path,
							property.Detection.MatchNode.LineNumber(),
							property.Detection.MatchNode.ColumnNumber(),
							"",
						),
						Value: schema.Schema{
							ObjectName:           data.Name,
							NormalizedObjectName: pluralizer.Singular(strings.ToLower(data.Name)),
							FieldName:            property.Name,
							NormalizedFieldName:  pluralizer.Singular(strings.ToLower(property.Name)),
							Classification:       property.Classification,
							Parent: &schema.Parent{
								LineNumber: property.Detection.MatchNode.LineNumber(),
								Content:    property.Detection.MatchNode.Content(),
							},
						},
					})
				}
			}
		}
	}

	return customDetections, nil
}
