package composition

import (
	"strings"

	"github.com/gertd/go-pluralize"

	"github.com/bearer/bearer/new/detector/implementation/custom"
	"github.com/bearer/bearer/new/detector/implementation/generic/datatype"
	detectortypes "github.com/bearer/bearer/new/detector/types"
	reportdetections "github.com/bearer/bearer/pkg/report/detections"
	"github.com/bearer/bearer/pkg/report/detectors"
	"github.com/bearer/bearer/pkg/report/schema"
	"github.com/bearer/bearer/pkg/report/source"
	"github.com/bearer/bearer/pkg/util/file"
)

func ReportDetections(report reportdetections.ReportDetection, file *file.FileInfo, detections []*detectortypes.Detection) {
	pluralizer := pluralize.NewClient()

	for _, detection := range detections {
		detectorType := detectors.Type(detection.DetectorType)
		data := detection.Data.(custom.Data)

		if len(data.Datatypes) == 0 {
			report.AddDetection(reportdetections.TypeCustomRisk,
				detectorType,
				source.New(
					file,
					file.Path,
					detection.MatchNode.StartLineNumber(),
					detection.MatchNode.StartColumnNumber(),
					detection.MatchNode.EndLineNumber(),
					detection.MatchNode.EndColumnNumber(),
					data.Pattern,
				),
				schema.Source{
					StartLineNumber:   detection.MatchNode.StartLineNumber(),
					EndLineNumber:     detection.MatchNode.EndLineNumber(),
					StartColumnNumber: detection.MatchNode.StartColumnNumber(),
					EndColumnNumber:   detection.MatchNode.EndColumnNumber(),
					Content:           detection.MatchNode.Content(),
				})
		}

		for _, datatypeDetection := range data.Datatypes {
			reportDatatypeDetection(
				report,
				pluralizer,
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
	pluralizer *pluralize.Client,
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
				property.Node.StartLineNumber(),
				property.Node.StartColumnNumber(),
				property.Node.EndLineNumber(),
				property.Node.EndColumnNumber(),
				"",
			),
			schema.Schema{
				ObjectName:           objectName,
				NormalizedObjectName: pluralizer.Singular(strings.ToLower(objectName)),
				FieldName:            property.Name,
				NormalizedFieldName:  pluralizer.Singular(strings.ToLower(property.Name)),
				Classification:       property.Classification,
				Source: &schema.Source{
					StartLineNumber:   detection.MatchNode.StartLineNumber(),
					EndLineNumber:     detection.MatchNode.EndLineNumber(),
					StartColumnNumber: detection.MatchNode.StartColumnNumber(),
					EndColumnNumber:   detection.MatchNode.EndColumnNumber(),
					Content:           detection.MatchNode.Content(),
				},
			},
		)

		if property.Datatype != nil {
			reportDatatypeDetection(
				report,
				pluralizer,
				file,
				detectorType,
				detection,
				property.Datatype,
				property.Name,
			)
		}
	}
}
