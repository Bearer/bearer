package scanner

import (
	"context"
	"fmt"
	"strings"

	schemaclassifier "github.com/bearer/bearer/internal/classification/schema"
	"github.com/bearer/bearer/internal/commands/process/settings"
	"github.com/bearer/bearer/internal/languages/golang"
	"github.com/bearer/bearer/internal/languages/java"
	"github.com/bearer/bearer/internal/languages/javascript"
	"github.com/bearer/bearer/internal/languages/php"
	"github.com/bearer/bearer/internal/languages/ruby"
	"github.com/bearer/bearer/internal/report"
	reportdetections "github.com/bearer/bearer/internal/report/detections"
	"github.com/bearer/bearer/internal/report/detectors"
	reportschema "github.com/bearer/bearer/internal/report/schema"
	"github.com/bearer/bearer/internal/report/source"
	customruletypes "github.com/bearer/bearer/internal/scanner/detectors/customrule/types"
	"github.com/bearer/bearer/internal/scanner/detectors/datatype"
	detectortypes "github.com/bearer/bearer/internal/scanner/detectors/types"
	"github.com/bearer/bearer/internal/scanner/language"
	"github.com/bearer/bearer/internal/util/file"
	"github.com/bearer/bearer/internal/util/pluralize"

	"github.com/bearer/bearer/internal/scanner/languagescanner"
	"github.com/bearer/bearer/internal/scanner/stats"
)

type Scanner struct {
	languageScanners []*languagescanner.Scanner
}

func New(schemaClassifier *schemaclassifier.Classifier, rules map[string]*settings.Rule) (*Scanner, error) {
	languages := []language.Language{
		java.Get(),
		javascript.Get(),
		ruby.Get(),
		php.Get(),
		golang.Get(),
	}

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
		detections, err := languageScanner.Scan(ctx, fileStats, file)
		if err != nil {
			return fmt.Errorf("%s scan failed: %w", languageScanner.LanguageID(), err)
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
	datatypeDetection *detectortypes.Detection,
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
	for _, languageScanner := range scanner.languageScanners {
		languageScanner.Close()
	}
}
