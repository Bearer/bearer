package composition

import (
	"strings"

	"github.com/gertd/go-pluralize"

	"github.com/bearer/curio/new/detector/implementation/custom"
	"github.com/bearer/curio/new/detector/implementation/generic/datatype"
	detectortypes "github.com/bearer/curio/new/detector/types"
	reportdetections "github.com/bearer/curio/pkg/report/detections"
	"github.com/bearer/curio/pkg/report/detectors"
	"github.com/bearer/curio/pkg/report/schema"
	"github.com/bearer/curio/pkg/report/source"
	"github.com/bearer/curio/pkg/util/file"
)

func ReportDetections(report reportdetections.ReportDetection, file *file.FileInfo, detections []*detectortypes.Detection) {
	pluralizer := pluralize.NewClient()

	for _, detection := range detections {
		data := detection.Data.(custom.Data)

		if len(data.Datatypes) == 0 {
			report.AddDetection(reportdetections.TypeCustomRisk,
				detectors.Type(detection.DetectorType),
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
			datatypeData := datatypeDetection.Data.(datatype.Data)

			report.AddDetection(
				reportdetections.TypeCustomClassified,
				detectors.Type(detection.DetectorType),
				source.New(
					file,
					file.Path,
					datatypeDetection.MatchNode.LineNumber(),
					datatypeDetection.MatchNode.ColumnNumber(),
					"",
				),
				schema.Schema{
					ObjectName:           datatypeData.Name,
					NormalizedObjectName: pluralizer.Singular(strings.ToLower(datatypeData.Name)),
					Classification:       datatypeData.Classification,
					Parent: &schema.Parent{
						LineNumber: detection.MatchNode.LineNumber(),
						Content:    detection.MatchNode.Content(),
					},
				},
			)

			for _, property := range datatypeData.Properties {

				report.AddDetection(
					reportdetections.TypeCustomClassified,
					detectors.Type(detection.DetectorType),
					source.New(
						file,
						file.Path,
						property.Detection.MatchNode.LineNumber(),
						property.Detection.MatchNode.ColumnNumber(),
						"",
					),
					schema.Schema{
						ObjectName:           datatypeData.Name,
						NormalizedObjectName: pluralizer.Singular(strings.ToLower(property.Name)),
						FieldName:            property.Name,
						NormalizedFieldName:  pluralizer.Singular(strings.ToLower(property.Name)),
						Classification:       property.Classification,
						Parent: &schema.Parent{
							LineNumber: detection.MatchNode.LineNumber(),
							Content:    detection.MatchNode.Content(),
						},
					},
				)
			}
		}
	}

}
