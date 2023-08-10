package composition

import (
	"strings"

	"github.com/bearer/bearer/new/detector/detection"
	"github.com/bearer/bearer/new/detector/implementation/custom"
	"github.com/bearer/bearer/new/detector/implementation/generic/datatype"
	reportdetections "github.com/bearer/bearer/pkg/report/detections"
	"github.com/bearer/bearer/pkg/report/detectors"
	"github.com/bearer/bearer/pkg/report/schema"
	"github.com/bearer/bearer/pkg/report/source"
	"github.com/bearer/bearer/pkg/util/file"
	"github.com/bearer/bearer/pkg/util/pluralize"
)

func ReportDetections(report reportdetections.ReportDetection, file *file.FileInfo, detections []*detection.Detection) {
	for _, detection := range detections {
		detectorType := detectors.Type(detection.DetectorType)
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
				schema.Source{
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
			schema.Schema{
				ObjectName:           objectName,
				NormalizedObjectName: pluralize.Singular(strings.ToLower(objectName)),
				FieldName:            property.Name,
				NormalizedFieldName:  pluralize.Singular(strings.ToLower(property.Name)),
				Classification:       property.Classification,
				Source: &schema.Source{
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
