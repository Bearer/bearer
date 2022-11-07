package risks

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/bearer/curio/pkg/report/output/dataflow/types"

	"github.com/bearer/curio/pkg/report/detections"
	"github.com/bearer/curio/pkg/report/schema"
	"github.com/bearer/curio/pkg/util/maputil"
)

type Holder struct {
	detectors map[string]detectorHolder // group datatypeHolders by name
}

type detectorHolder struct {
	id        string
	datatypes map[string]*datatypeHolder // group detectors by detectorName
}
type datatypeHolder struct {
	name  string
	files map[string]*fileHolder // group files by filename
}
type fileHolder struct {
	name       string
	lineNumber map[int]int
}

func New() *Holder {
	return &Holder{
		detectors: make(map[string]detectorHolder),
	}
}

func (holder *Holder) AddSchema(detection detections.Detection) error {
	var value schema.Schema
	buf := bytes.NewBuffer(nil)
	err := json.NewEncoder(buf).Encode(detection.Value)
	if err != nil {
		return fmt.Errorf("expect detection to have value of type schema %#v", detection.Value)
	}
	err = json.NewDecoder(buf).Decode(&value)
	if err != nil {
		return fmt.Errorf("expect detection to have value of type schema %#v", detection.Value)
	}

	if value.FieldName != "" {
		holder.addDatatype(string(detection.DetectorType), value.FieldName, detection.Source.Filename, *detection.Source.LineNumber)
	}

	return nil
}

// addDatatype adds datatype to hash list and at the same time blocks duplicates
func (holder *Holder) addDatatype(ruleName string, datatypeName string, fileName string, lineNumber int) {
	// create detector entry if it doesn't exist
	if _, exists := holder.detectors[ruleName]; !exists {
		holder.detectors[ruleName] = detectorHolder{
			id:        ruleName,
			datatypes: make(map[string]*datatypeHolder),
		}
	}

	detector := holder.detectors[ruleName]
	// create datatype entry if it doesn't exist
	if _, exists := detector.datatypes[datatypeName]; !exists {
		detector.datatypes[datatypeName] = &datatypeHolder{
			name:  datatypeName,
			files: make(map[string]*fileHolder),
		}
	}

	datatype := detector.datatypes[datatypeName]
	// create file entry if it doesn't exist
	if _, exists := datatype.files[fileName]; !exists {
		datatype.files[fileName] = &fileHolder{
			name:       fileName,
			lineNumber: make(map[int]int),
		}
	}

	datatype.files[fileName].lineNumber[lineNumber] = lineNumber

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
			constructedDatatype := types.RiskDatatype{
				Name:      datatype.name,
				Locations: make([]types.RiskLocation, 0),
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
