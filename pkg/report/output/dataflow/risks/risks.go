package risks

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"

	"github.com/bearer/curio/pkg/report/output/dataflow/types"

	"github.com/bearer/curio/pkg/report/detections"
	"github.com/bearer/curio/pkg/report/schema"
	"github.com/bearer/curio/pkg/util/maputil"
)

// - detector_id: "rails_logger_detector"
// data_types:
//   - name: "email"
// 	locations:
// 	  - filename: "app/models/user.rb"
// 		line_number: 5
// 	  - filename: "app/models/employee.rb"
// 		line_number: 5

type Holder struct {
	detectors map[string]detectorHolder // group datatypeHolders by name
}

type detectorHolder struct {
	id        string
	datatypes map[string]*datatypeHolder // group detectors by detectorName
}
type datatypeHolder struct {
	name  string
	files []*fileHolder // group files by filename
}
type fileHolder struct {
	name       string
	lineNumber int
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
	// create detector entry if it doesn't exist
	if _, exists := detector.datatypes[datatypeName]; !exists {
		detector.datatypes[datatypeName] = &datatypeHolder{
			name:  datatypeName,
			files: make([]*fileHolder, 0),
		}
	}

	detector.datatypes[datatypeName].files = append(detector.datatypes[datatypeName].files, &fileHolder{
		name:       fileName,
		lineNumber: lineNumber,
	})

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

			// sort by name asc and then by line number asc
			sort.Slice(datatype.files, func(i, j int) bool {
				a := datatype.files[i]
				b := datatype.files[j]

				if a.name < b.name {
					return true
				}

				return a.lineNumber < b.lineNumber
			})
			for _, file := range datatype.files {
				constructedDatatype.Locations = append(constructedDatatype.Locations, types.RiskLocation{
					Filename:   file.name,
					LineNumber: file.lineNumber,
				})
			}
			constructedDetector.DataTypes = append(constructedDetector.DataTypes, constructedDatatype)
		}

		data = append(data, constructedDetector)
	}

	return data
}
