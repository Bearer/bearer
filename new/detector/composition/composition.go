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
)

func ReportDetections(report reportdetections.ReportDetection, compositionDetections []*detectortypes.CompositionDetection) {
	pluralizer := pluralize.NewClient()

	for _, compositionDetection := range compositionDetections {

		for _, detection := range compositionDetection.Detections {

			data := detection.Data.(custom.Data)

			if len(data.Datatypes) == 0 {
				report.AddDetection(reportdetections.TypeCustomRisk,
					detectors.Type(compositionDetection.RuleName),
					source.New(
						compositionDetection.File,
						compositionDetection.File.Path,
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
					detectors.Type(compositionDetection.RuleName),
					source.New(
						compositionDetection.File,
						compositionDetection.File.Path,
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
						detectors.Type(compositionDetection.RuleName),
						source.New(
							compositionDetection.File,
							compositionDetection.File.Path,
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

}
