package stats

import (
	"github.com/bearer/curio/pkg/commands/process/settings"
	"github.com/bearer/curio/pkg/report/output/dataflow"
	"github.com/hhatto/gocloc"
)

type DataType struct {
	Name        string `json:"name"`
	Occurrences int    `json:"occurrences"`
}

type Stats struct {
	NumberOfLines     int32      `json:"number_of_lines"`
	NumberOfDataTypes int        `json:"number_of_data_types"`
	DataTypes         []DataType `json:"data_types"`
}

func GetOutput(inputgocloc *gocloc.Result, inputDataflow *dataflow.DataFlow, config settings.Config) (*Stats, error) {
	numberOfDataTypesFound := len(inputDataflow.Datatypes)
	data_types := []DataType{}

	for _, data_type := range inputDataflow.Datatypes {
		occurrences := 0
		for _, detector := range data_type.Detectors {
			occurrences += len(detector.Locations)
		}

		data_types = append(data_types, DataType{
			Name:        data_type.Name,
			Occurrences: occurrences,
		})
	}

	return &Stats{
		NumberOfLines:     inputgocloc.Total.Code,
		NumberOfDataTypes: numberOfDataTypesFound,
		DataTypes:         data_types,
	}, nil
}
