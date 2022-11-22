package risks

import (
	"github.com/bearer/curio/pkg/classification/db"
	"github.com/bearer/curio/pkg/commands/process/settings"
	"github.com/bearer/curio/pkg/report/output/dataflow/detectiondecoder"
	"github.com/bearer/curio/pkg/report/output/dataflow/types"

	"github.com/bearer/curio/pkg/report/detections"
	"github.com/bearer/curio/pkg/util/classify"
	"github.com/bearer/curio/pkg/util/maputil"
)

type Holder struct {
	detectors  map[string]detectorHolder // group datatypeHolders by name
	config     settings.Config
	isInternal bool
}

type detectorHolder struct {
	id        string
	datatypes map[string]*datatypeHolder // group detectors by detectorName
}
type datatypeHolder struct {
	name         string
	uuid         string
	categoryUUID string
	files        map[string]*fileHolder // group files by filename
}
type fileHolder struct {
	name       string
	lineNumber map[int]int
}

func New(config settings.Config, isInternal bool) *Holder {
	return &Holder{
		detectors:  make(map[string]detectorHolder),
		config:     config,
		isInternal: isInternal,
	}
}

func (holder *Holder) AddSchema(detection detections.Detection) error {
	classification, err := detectiondecoder.GetSchemaClassification(detection)
	if err != nil {
		return err
	}

	if classification.Decision.State == classify.Valid {
		holder.addDatatype(string(detection.DetectorType), classification.DataType, detection.Source.Filename, *detection.Source.LineNumber)
	}

	return nil
}

func (holder *Holder) AddComputedRisk(risk types.RiskDetector) {
	for _, datatype := range risk.DataTypes {
		for _, location := range datatype.Locations {
			holder.addDatatype(string(risk.DetectorID), datatype.Name, location.Filename, location.LineNumber)
		}
	}
}

// addDatatype adds detector to hash list and at the same time blocks duplicates
func (holder *Holder) addDatatype(ruleName string, datatype *db.DataType, fileName string, lineNumber int) {
	// create detector entry if it doesn't exist
	if _, exists := holder.detectors[ruleName]; !exists {
		holder.detectors[ruleName] = detectorHolder{
			id:        ruleName,
			datatypes: make(map[string]*datatypeHolder),
		}
	}

	detector := holder.detectors[ruleName]
	// create datatype entry if it doesn't exist
	if _, exists := detector.datatypes[datatype.Name]; !exists {
		if holder.isInternal {
			detector.datatypes[datatype.Name] = &datatypeHolder{
				name:         datatype.Name,
				uuid:         datatype.UUID,
				categoryUUID: datatype.CategoryUUID,
				files:        make(map[string]*fileHolder),
			}
		} else {
			detector.datatypes[datatype.Name] = &datatypeHolder{
				name:  datatype.Name,
				files: make(map[string]*fileHolder),
			}
		}
	}

	detectorDatatype := detector.datatypes[datatype.Name]
	// create file entry if it doesn't exist
	if _, exists := detectorDatatype.files[fileName]; !exists {
		detectorDatatype.files[fileName] = &fileHolder{
			name:       fileName,
			lineNumber: make(map[int]int),
		}
	}

	detectorDatatype.files[fileName].lineNumber[lineNumber] = lineNumber

}

func (holder *Holder) ToDataFlow() []types.RiskDetector {
	data := make([]types.RiskDetector, 0)

	detectors := maputil.ToSortedSlice(holder.detectors)

	for _, detector := range detectors {
		constructedDetector := types.RiskDetector{
			DetectorID: detector.id,
			DataTypes:  make([]types.RiskDatatype, 0),
		}

		datatypes := maputil.ToSortedSlice(detector.datatypes)

		for _, datatype := range datatypes {

			stored := false
			if customDetector, isCustomDetector := holder.config.CustomDetector[detector.id]; isCustomDetector {
				stored = customDetector.Stored
			}

			constructedDatatype := types.RiskDatatype{
				Name:         datatype.name,
				UUID:         datatype.uuid,
				CategoryUUID: datatype.categoryUUID,
				Stored:       stored,
				Locations:    make([]types.RiskLocation, 0),
			}

			files := maputil.ToSortedSlice(datatype.files)
			for _, file := range files {

				lineNumbers := maputil.ToSortedSlice(file.lineNumber)
				for _, lineNumber := range lineNumbers {
					constructedDatatype.Locations = append(constructedDatatype.Locations, types.RiskLocation{
						Filename:   file.name,
						LineNumber: lineNumber,
					})
				}
			}
			constructedDetector.DataTypes = append(constructedDetector.DataTypes, constructedDatatype)
		}

		data = append(data, constructedDetector)
	}

	return data
}
