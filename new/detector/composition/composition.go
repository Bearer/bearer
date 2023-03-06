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
					detection.MatchNode.LineNumber(),
					detection.MatchNode.ColumnNumber(),
					data.Pattern,
				),
				schema.Parent{
					LineNumber: detection.MatchNode.LineNumber(),
					Content:    detection.MatchNode.Content(),
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
				datatypeDetection.MatchNode.LineNumber(),
				datatypeDetection.MatchNode.ColumnNumber(),
				"",
			),
			schema.Schema{
				ObjectName:           objectName,
				NormalizedObjectName: pluralizer.Singular(strings.ToLower(objectName)),
				FieldName:            property.Name,
				NormalizedFieldName:  pluralizer.Singular(strings.ToLower(property.Name)),
				Classification:       property.Classification,
				Parent: &schema.Parent{
					LineNumber: detection.MatchNode.LineNumber(),
					Content:    detection.MatchNode.Content(),
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
