package scanner

import (
	"context"
	"fmt"
	"strings"

	"github.com/bearer/bearer/new/detector/detection"
	"github.com/bearer/bearer/new/detector/evaluator"
	"github.com/bearer/bearer/new/detector/evaluator/stats"
	"github.com/bearer/bearer/new/detector/implementation/custom"
	"github.com/bearer/bearer/new/detector/implementation/generic/datatype"
	"github.com/bearer/bearer/new/language/implementation"
	"github.com/bearer/bearer/new/language/implementation/java"
	"github.com/bearer/bearer/new/language/implementation/javascript"
	"github.com/bearer/bearer/new/language/implementation/ruby"
	schemaclassifier "github.com/bearer/bearer/pkg/classification/schema"
	"github.com/bearer/bearer/pkg/commands/process/settings"
	"github.com/bearer/bearer/pkg/report"
	reportdetections "github.com/bearer/bearer/pkg/report/detections"
	"github.com/bearer/bearer/pkg/report/detectors"
	reportschema "github.com/bearer/bearer/pkg/report/schema"
	"github.com/bearer/bearer/pkg/report/source"
	"github.com/bearer/bearer/pkg/util/file"
	"github.com/bearer/bearer/pkg/util/pluralize"
)

type Scanner struct {
	evaluators []*evaluator.Evaluator
}

func New(schemaClassifier *schemaclassifier.Classifier, rules map[string]*settings.Rule) (*Scanner, error) {
	langImplementations := []implementation.Implementation{
		java.Get(),
		javascript.Get(),
		ruby.Get(),
	}

	evaluators := make([]*evaluator.Evaluator, len(langImplementations))

	for i, langImplementation := range langImplementations {
		evaluator, err := evaluator.New(langImplementation, schemaClassifier, rules)
		if err != nil {
			return nil, fmt.Errorf("error creating %s evaluator: %w", langImplementation.Name(), err)
		}

		evaluators[i] = evaluator
	}

	return &Scanner{evaluators: evaluators}, nil
}

func (scanner *Scanner) Scan(
	ctx context.Context,
	report report.Report,
	fileStats *stats.FileStats,
	file *file.FileInfo,
) error {
	if scanner == nil {
		return nil
	}

	for _, evaluator := range scanner.evaluators {
		detections, err := evaluator.DetectFromFile(ctx, fileStats, file)
		if err != nil {
			return fmt.Errorf("%s failed to detect in file %s: %w", evaluator.LanguageName(), file.AbsolutePath, err)
		}

		for _, detection := range detections {
			detectorType := detectors.Type(detection.RuleID)
			data := detection.Data.(custom.Data)

			if len(data.Datatypes) == 0 {
				report.AddDetection(reportdetections.TypeCustomRisk,
					detectorType,
					source.New(
						file,
						file.Path,
						detection.MatchNode.ContentStart.Line,
						detection.MatchNode.ContentStart.Column,
						detection.MatchNode.ContentEnd.Line,
						detection.MatchNode.ContentEnd.Column,
						data.Pattern,
					),
					reportschema.Source{
						StartLineNumber:   detection.MatchNode.ContentStart.Line,
						EndLineNumber:     detection.MatchNode.ContentEnd.Line,
						StartColumnNumber: detection.MatchNode.ContentStart.Column,
						EndColumnNumber:   detection.MatchNode.ContentEnd.Column,
						Content:           detection.MatchNode.Content(),
					})
			}

			for _, datatypeDetection := range data.Datatypes {
				reportDatatypeDetection(
					report,
					file,
					detectorType,
					detection,
					datatypeDetection,
					"",
				)
			}
		}
	}

	return nil
}

func reportDatatypeDetection(
	report reportdetections.ReportDetection,
	file *file.FileInfo,
	detectorType detectors.Type,
	detection,
	datatypeDetection *detection.Detection,
	objectName string,
) {
	data := datatypeDetection.Data.(datatype.Data)

	for _, property := range data.Properties {
		report.AddDetection(
			reportdetections.TypeCustomClassified,
			detectorType,
			source.New(
				file,
				file.Path,
				property.Node.ContentStart.Line,
				property.Node.ContentStart.Column,
				property.Node.ContentEnd.Line,
				property.Node.ContentEnd.Column,
				"",
			),
			reportschema.Schema{
				ObjectName:           objectName,
				NormalizedObjectName: pluralize.Singular(strings.ToLower(objectName)),
				FieldName:            property.Name,
				NormalizedFieldName:  pluralize.Singular(strings.ToLower(property.Name)),
				Classification:       property.Classification,
				Source: &reportschema.Source{
					StartLineNumber:   detection.MatchNode.ContentStart.Line,
					EndLineNumber:     detection.MatchNode.ContentEnd.Line,
					StartColumnNumber: detection.MatchNode.ContentStart.Column,
					EndColumnNumber:   detection.MatchNode.ContentEnd.Column,
					Content:           detection.MatchNode.Content(),
				},
			},
		)

		if property.Datatype != nil {
			reportDatatypeDetection(
				report,
				file,
				detectorType,
				detection,
				property.Datatype,
				property.Name,
			)
		}
	}
}

func (scanner *Scanner) Close() {
	for _, evaluator := range scanner.evaluators {
		evaluator.Close()
	}
}
