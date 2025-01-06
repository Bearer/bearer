package scanner

import (
	"context"
	"fmt"
	"strings"

	schemaclassifier "github.com/bearer/bearer/pkg/classification/schema"
	"github.com/bearer/bearer/pkg/commands/process/settings"
	"github.com/bearer/bearer/pkg/engine"
	"github.com/bearer/bearer/pkg/report"
	reportdetections "github.com/bearer/bearer/pkg/report/detections"
	"github.com/bearer/bearer/pkg/report/detectors"
	reportschema "github.com/bearer/bearer/pkg/report/schema"
	"github.com/bearer/bearer/pkg/report/source"
	customruletypes "github.com/bearer/bearer/pkg/scanner/detectors/customrule/types"
	"github.com/bearer/bearer/pkg/scanner/detectors/datatype"
	detectortypes "github.com/bearer/bearer/pkg/scanner/detectors/types"
	"github.com/bearer/bearer/pkg/util/file"
	"github.com/bearer/bearer/pkg/util/pluralize"

	"github.com/bearer/bearer/pkg/scanner/languagescanner"
	"github.com/bearer/bearer/pkg/scanner/stats"
)

type Scanner struct {
	languageScanners []*languagescanner.Scanner
}

func New(
	engine engine.Engine,
	schemaClassifier *schemaclassifier.Classifier,
	rules map[string]*settings.Rule,
) (*Scanner, error) {
	languages := engine.GetLanguages()

	languageScanners := make([]*languagescanner.Scanner, len(languages))

	for i, language := range languages {
		languageScanner, err := languagescanner.New(language, schemaClassifier, rules)
		if err != nil {
			return nil, fmt.Errorf("error creating %s language scanner: %w", language.ID(), err)
		}

		languageScanners[i] = languageScanner
	}

	return &Scanner{languageScanners: languageScanners}, nil
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

	for _, languageScanner := range scanner.languageScanners {
		detections, expectedDetections, err := languageScanner.Scan(ctx, fileStats, file)
		if err != nil {
			return fmt.Errorf("%s scan failed: %w", languageScanner.LanguageID(), err)
		}

		for _, detection := range expectedDetections {
			value := ""
			detectorType := detectors.Type(detection.RuleID)
			report.AddDetection(reportdetections.TypeExpectedDetection,
				detectorType,
				source.New(
					file,
					file.Path,
					detection.MatchNode.ContentStart.Line,
					detection.MatchNode.ContentStart.Column,
					detection.MatchNode.ContentEnd.Line,
					detection.MatchNode.ContentEnd.Column,
					"",
				),
				reportschema.Source{
					StartLineNumber:   detection.MatchNode.ContentStart.Line,
					EndLineNumber:     detection.MatchNode.ContentEnd.Line,
					StartColumnNumber: detection.MatchNode.ContentStart.Column,
					EndColumnNumber:   detection.MatchNode.ContentEnd.Column,
					Content:           value,
				})
		}

		for _, detection := range detections {
			detectorType := detectors.Type(detection.RuleID)
			data := detection.Data.(customruletypes.Data)

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
						data.Value,
					),
					reportschema.Source{
						StartLineNumber:   detection.MatchNode.ContentStart.Line,
						EndLineNumber:     detection.MatchNode.ContentEnd.Line,
						StartColumnNumber: detection.MatchNode.ContentStart.Column,
						EndColumnNumber:   detection.MatchNode.ContentEnd.Column,
						Content:           data.Value,
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
	datatypeDetection *detectortypes.Detection,
	objectName string,
) {
	data := datatypeDetection.Data.(datatype.Data)

	for _, property := range data.Properties {
		detectionContent := detection.MatchNode.Content()

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
					Content:           detectionContent,
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
	for _, languageScanner := range scanner.languageScanners {
		languageScanner.Close()
	}
}
