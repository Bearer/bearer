package stats

import (
	"fmt"
	"strings"

	"github.com/bearer/curio/pkg/classification/db"
	"github.com/bearer/curio/pkg/commands/process/settings"

	"github.com/bearer/curio/pkg/report/output/dataflow"
	"github.com/bearer/curio/pkg/util/maputil"
	"github.com/hhatto/gocloc"
	"github.com/fatih/color"
)

type DataType struct {
	Name         string `json:"name" yaml:"name"`
	CategoryUUID string `json:"-" yaml:"-"`
	Encrypted    bool   `json:"-" yaml:"-"`
	Occurrences  int    `json:"occurrences" yaml:"occurrences"`
}

type DataStore struct {
	Name                       string `json:"name" yaml:"name"`
	NumberOfDataTypes          int `json:"number_of_data_types" yaml:"number_of_data_types"`
	NumberOfEncryptedDataTypes int `json:"number_of_encrypted_data_types" yaml:"number_of_encrypted_data_types"`
}

type Stats struct {
	NumberOfLines        int32            `json:"number_of_lines" yaml:"number_of_lines"`
	NumberOfDataTypes    int              `json:"number_of_data_types" yaml:"number_of_data_types"`
	DataTypes            []DataType       `json:"data_types" yaml:"data_types"`
	DataStores           []DataStore      `json:"-" yaml:"-"`
	NumberOfExternalAPIs int              `json:"-" yaml:"-"`
	NumberOfInternalAPIs int              `json:"-" yaml:"-"`
	Languages            map[string]int32 `json:"-" yaml:"-"`
	DataGroups           []string         `json:"-" yaml:"-"`
}

func GetOutput(inputgocloc *gocloc.Result, inputDataflow *dataflow.DataFlow, config settings.Config) (*Stats, error) {
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

	dataStores := []DataStore{}
	numberOfExternalAPIs := 0
	numberOfInternalAPIs := 0
	for _, component := range inputDataflow.Components {
		if component.Type == "internal_service" {
			numberOfInternalAPIs++
		}
		if component.Type == "external_service" {
			numberOfExternalAPIs++
		}

		// @todo FIXME: Collect statistics for data stores

		// detectors := []string{}
		// for _, location := range component.Locations {
		//	detectors = append(detectors, location.Detector)
		// }

		// for _, detector := range detectors {
		//	if detector == string(reportdetectors.DetectorSQL) {
		//		dataStores = append(dataStores, DataStore{
		//			Name: "",
		//			NumberOfDataTypes: 0,
		//			NumberOfEncryptedDataTypes: 0,
		//		})
		//	}
		// }
	}

	languages := map[string]int32{}
	for _, language := range inputgocloc.Languages {
		languages[language.Name] = language.Code
	}

	return &Stats{
		NumberOfLines:        inputgocloc.Total.Code,
		NumberOfDataTypes:    numberOfDataTypesFound,
		DataTypes:            data_types,
		DataStores:           dataStores,
		NumberOfExternalAPIs: numberOfExternalAPIs,
		NumberOfInternalAPIs: numberOfInternalAPIs,
		Languages:            languages,
		DataGroups:           dataGroupNames,
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

func GetPlaceholderOutput(inputgocloc *gocloc.Result, inputDataflow *dataflow.DataFlow, config settings.Config) (outputStr *strings.Builder, err error) {
	outputStr = &strings.Builder{}
	statistics, err := GetOutput(inputgocloc, inputDataflow, config)

	totalDataTypeOccurrences := 0
	for _, dataType := range statistics.DataTypes {
		totalDataTypeOccurrences += dataType.Occurrences
	}

	supportURL := "https://curio.sh/explanations/reports/"
	outputStr.WriteString(fmt.Sprintf(`
The policy report is not yet available for your stack. Learn more at %s

Though this doesn’t mean the curious bear comes empty-handed, it found:

- %d unique data type(s), representing %d occurrences, including %s.`,
		supportURL,
		statistics.NumberOfDataTypes,
		totalDataTypeOccurrences,
		strings.Join(statistics.DataGroups, ", ")))

	if len(statistics.DataStores) != 0 {
		totalDataStoreDataTypes := 0
		totalDataStoreEncryptedDataTypes := 0
		for _, dataStore := range statistics.DataStores {
			totalDataStoreDataTypes += dataStore.NumberOfDataTypes
			totalDataStoreEncryptedDataTypes += dataStore.NumberOfEncryptedDataTypes
		}

		outputStr.WriteString(fmt.Sprintf(
			`
- %d database(s) storing %d data type(s) including %d encrypted data type(s).`,
			len(statistics.DataStores),
			totalDataStoreDataTypes,
			totalDataStoreEncryptedDataTypes))
	}

	if statistics.NumberOfExternalAPIs != 0 {
		outputStr.WriteString(fmt.Sprintf(
			`
- %d external service(s).`,
			statistics.NumberOfExternalAPIs))
	}

	if statistics.NumberOfInternalAPIs != 0 {
		outputStr.WriteString(fmt.Sprintf(
			`
- %d internal URL(s).`,
			statistics.NumberOfInternalAPIs))
	}

	suggestedCommand := color.New(color.Italic).Sprintf("curio scan --report dataflow")
	outputStr.WriteString(fmt.Sprintf(`

Run the data flow report if you want the full output using: %s`, suggestedCommand))

	return
}
