package stats

import (
	"fmt"
	"strings"

	"github.com/hhatto/gocloc"

	"github.com/bearer/bearer/pkg/classification/db"
	"github.com/bearer/bearer/pkg/commands/process/settings"
	"github.com/bearer/bearer/pkg/util/maputil"

	"github.com/bearer/bearer/pkg/report/output/stats/types"
	outputtypes "github.com/bearer/bearer/pkg/report/output/types"
)

func AddReportData(
	reportData *outputtypes.ReportData,
	inputgocloc *gocloc.Result,
	config settings.Config,
) error {
	numberOfDataTypesFound := len(reportData.Dataflow.Datatypes)
	datatypes := []types.DataType{}

	for _, datatype := range reportData.Dataflow.Datatypes {
		occurrences := 0
		for _, detector := range datatype.Detectors {
			occurrences += len(detector.Locations)
		}

		encrypted := false
	outer:
		for _, detector := range datatype.Detectors {
			for _, location := range detector.Locations {
				if location.Encrypted != nil && *location.Encrypted {
					encrypted = true
					break outer
				}
			}
		}

		datatypes = append(datatypes, types.DataType{
			Name:         datatype.Name,
			CategoryUUID: datatype.CategoryUUID,
			Encrypted:    encrypted,
			Occurrences:  occurrences,
		})
	}

	dataGroupNames := getDataGroupNames(config, datatypes)

	numberOfDatabases := 0
	numberOfExternalAPIs := 0
	numberOfInternalAPIs := 0
	for _, component := range reportData.Dataflow.Components {
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

	reportData.Stats = &types.Stats{
		NumberOfLines:        inputgocloc.Total.Code,
		NumberOfDataTypes:    numberOfDataTypesFound,
		DataTypes:            datatypes,
		NumberOfDatabases:    numberOfDatabases,
		NumberOfExternalAPIs: numberOfExternalAPIs,
		NumberOfInternalAPIs: numberOfInternalAPIs,
		Languages:            languages,
		DataGroups:           dataGroupNames,
	}

	return nil
}

func getDataGroupNames(config settings.Config, dataTypes []types.DataType) []string {
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

func AnythingFoundFor(statistics *types.Stats) bool {
	return statistics.NumberOfDataTypes != 0 ||
		statistics.NumberOfDatabases != 0 ||
		statistics.NumberOfExternalAPIs != 0 ||
		statistics.NumberOfInternalAPIs != 0
}

func WriteStatsToString(outputStr *strings.Builder, statistics *types.Stats) {
	totalDataTypeOccurrences := 0
	for _, dataType := range statistics.DataTypes {
		totalDataTypeOccurrences += dataType.Occurrences
	}

	if statistics.NumberOfDataTypes != 0 {
		fmt.Fprintf(outputStr, `- %d unique data type(s), representing %d occurrences, including %s.`,
			statistics.NumberOfDataTypes,
			totalDataTypeOccurrences,
			strings.Join(statistics.DataGroups, ", "))
	}

	if statistics.NumberOfDatabases != 0 {
		numberOfEncryptedDataTypes := 0
		for _, dataType := range statistics.DataTypes {
			if dataType.Encrypted {
				numberOfEncryptedDataTypes += 1
			}
		}

		fmt.Fprintf(outputStr, `
- %d database(s) storing %d data type(s) including %d encrypted data type(s).`,
			statistics.NumberOfDatabases,
			statistics.NumberOfDataTypes,
			numberOfEncryptedDataTypes)
	}

	if statistics.NumberOfExternalAPIs != 0 {
		fmt.Fprintf(outputStr, `
- %d external service(s).`,
			statistics.NumberOfExternalAPIs)
	}

	if statistics.NumberOfInternalAPIs != 0 {
		fmt.Fprintf(outputStr, `
- %d internal URL(s).`,
			statistics.NumberOfInternalAPIs)
	}
}

func GetPlaceholderOutput(
	reportData *outputtypes.ReportData,
	inputgocloc *gocloc.Result,
	config settings.Config,
) (outputStr *strings.Builder, err error) {
	outputStr = &strings.Builder{}
	if err := AddReportData(reportData, inputgocloc, config); err != nil {
		return nil, err
	}

	outputStr.WriteString(`
The security report is not yet available for your application.
Learn more about language support at https://docs.bearer.com/reference/supported-languages/`)

	if AnythingFoundFor(reportData.Stats) {
		outputStr.WriteString(`

Though this doesn’t mean the curious bear comes empty-handed, it found:

`)
	}

	WriteStatsToString(outputStr, reportData.Stats)

	return
}
