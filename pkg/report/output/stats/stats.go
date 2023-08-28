package stats

import (
	"fmt"
	"strings"

	"github.com/bearer/bearer/pkg/classification/db"
	"github.com/bearer/bearer/pkg/commands/process/settings"

	"github.com/bearer/bearer/pkg/report/output/types"
	"github.com/bearer/bearer/pkg/util/maputil"
	"github.com/hhatto/gocloc"
)

type DataType struct {
	Name         string `json:"name" yaml:"name"`
	CategoryUUID string `json:"-" yaml:"-"`
	Encrypted    bool   `json:"-" yaml:"-"`
	Occurrences  int    `json:"occurrences" yaml:"occurrences"`
}

type Stats struct {
	NumberOfLines        int32            `json:"number_of_lines" yaml:"number_of_lines"`
	NumberOfDataTypes    int              `json:"number_of_data_types" yaml:"number_of_data_types"`
	DataTypes            []DataType       `json:"data_types" yaml:"data_types"`
	NumberOfDatabases    int              `json:"-" yaml:"-"`
	NumberOfExternalAPIs int              `json:"-" yaml:"-"`
	NumberOfInternalAPIs int              `json:"-" yaml:"-"`
	Languages            map[string]int32 `json:"-" yaml:"-"`
	DataGroups           []string         `json:"-" yaml:"-"`
}

func GetOutput(
	inputgocloc *gocloc.Result,
	inputDataflow *types.DataFlow,
	config settings.Config,
) (*types.Output[Stats], error) {
	numberOfDataTypesFound := len(inputDataflow.Datatypes)
	data_types := []DataType{}

	for _, data_type := range inputDataflow.Datatypes {
		occurrences := 0
		for _, detector := range data_type.Detectors {
			occurrences += len(detector.Locations)
		}

		encrypted := false
	outer:
		for _, detector := range data_type.Detectors {
			for _, location := range detector.Locations {
				if location.Encrypted != nil && *location.Encrypted {
					encrypted = true
					break outer
				}
			}
		}

		data_types = append(data_types, DataType{
			Name:         data_type.Name,
			CategoryUUID: data_type.CategoryUUID,
			Encrypted:    encrypted,
			Occurrences:  occurrences,
		})
	}

	dataGroupNames := getDataGroupNames(config, data_types)

	numberOfDatabases := 0
	numberOfExternalAPIs := 0
	numberOfInternalAPIs := 0
	for _, component := range inputDataflow.Components {
		if component.Type == "internal_service" {
			numberOfInternalAPIs++
		}
		if component.Type == "external_service" {
			numberOfExternalAPIs++
		}

		if component.SubType == "database" {
			numberOfDatabases++
		}
	}

	languages := map[string]int32{}
	for _, language := range inputgocloc.Languages {
		languages[language.Name] = language.Code
	}

	return &types.Output[Stats]{
		Data: Stats{
			NumberOfLines:        inputgocloc.Total.Code,
			NumberOfDataTypes:    numberOfDataTypesFound,
			DataTypes:            data_types,
			NumberOfDatabases:    numberOfDatabases,
			NumberOfExternalAPIs: numberOfExternalAPIs,
			NumberOfInternalAPIs: numberOfInternalAPIs,
			Languages:            languages,
			DataGroups:           dataGroupNames,
		},
		Dataflow: inputDataflow,
	}, nil
}

func getDataGroupNames(config settings.Config, dataTypes []DataType) []string {
	dataCategories := db.DefaultWithContext(config.Scan.Context).DataCategories
	dataGroups := make(map[string]bool)
	for _, dataType := range dataTypes {
		for _, category := range dataCategories {
			if category.UUID == dataType.CategoryUUID {
				for _, group := range category.Groups {
					dataGroups[group.Name] = true
				}
				break
			}
		}
	}

	return maputil.SortedStringKeys(dataGroups)
}

func AnythingFoundFor(statistics *Stats) bool {
	return statistics.NumberOfDataTypes != 0 ||
		statistics.NumberOfDatabases != 0 ||
		statistics.NumberOfExternalAPIs != 0 ||
		statistics.NumberOfInternalAPIs != 0
}

func WriteStatsToString(outputStr *strings.Builder, statistics *Stats) {
	totalDataTypeOccurrences := 0
	for _, dataType := range statistics.DataTypes {
		totalDataTypeOccurrences += dataType.Occurrences
	}

	if statistics.NumberOfDataTypes != 0 {
		outputStr.WriteString(fmt.Sprintf(`- %d unique data type(s), representing %d occurrences, including %s.`,
			statistics.NumberOfDataTypes,
			totalDataTypeOccurrences,
			strings.Join(statistics.DataGroups, ", ")))
	}

	if statistics.NumberOfDatabases != 0 {
		numberOfEncryptedDataTypes := 0
		for _, dataType := range statistics.DataTypes {
			if dataType.Encrypted {
				numberOfEncryptedDataTypes += 1
			}
		}

		outputStr.WriteString(fmt.Sprintf(`
- %d database(s) storing %d data type(s) including %d encrypted data type(s).`,
			statistics.NumberOfDatabases,
			statistics.NumberOfDataTypes,
			numberOfEncryptedDataTypes))
	}

	if statistics.NumberOfExternalAPIs != 0 {
		outputStr.WriteString(fmt.Sprintf(`
- %d external service(s).`,
			statistics.NumberOfExternalAPIs))
	}

	if statistics.NumberOfInternalAPIs != 0 {
		outputStr.WriteString(fmt.Sprintf(`
- %d internal URL(s).`,
			statistics.NumberOfInternalAPIs))
	}
}

func GetPlaceholderOutput(
	inputgocloc *gocloc.Result,
	inputDataflow *types.DataFlow,
	config settings.Config,
) (outputStr *strings.Builder, err error) {
	outputStr = &strings.Builder{}
	statisticsOutput, err := GetOutput(inputgocloc, inputDataflow, config)

	outputStr.WriteString(`
The security report is not yet available for your application.
Learn more about language support at https://docs.bearer.com/reference/supported-languages/`)

	if AnythingFoundFor(&statisticsOutput.Data) {
		outputStr.WriteString(`

Though this doesn’t mean the curious bear comes empty-handed, it found:

`)
	}

	WriteStatsToString(outputStr, &statisticsOutput.Data)

	return
}
