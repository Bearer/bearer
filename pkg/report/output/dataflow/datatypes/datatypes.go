package datatypes

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/bearer/curio/pkg/report/detections"
	"github.com/bearer/curio/pkg/report/schema"
	"github.com/bearer/curio/pkg/util/maputil"
)

type datatypeHolder struct {
	name      string
	detectors map[string]*detectorHolder // group detectors by detectorName
}
type detectorHolder struct {
	name  string
	files map[string]*fileHolder // group files by filename
}
type fileHolder struct {
	name        string
	lineNumbers map[int]int // group occurences by line number
}

func (holder *dataFlowHolder) toDataFlow() *DataFlow {
	data := &DataFlow{}

	datatypes := maputil.ToSortedSlice(holder.datatypes)

	for _, datatype := range datatypes {
		constructedDatatype := Datatype{
			Name: datatype.name,
		}

		detectors := maputil.ToSortedSlice(datatype.detectors)

		for _, detectorHolder := range detectors {
			constructedDetector := Detector{
				Name:      detectorHolder.name,
				Stored:    true,
				Locations: make([]Location, 0),
			}

			for _, fileHolder := range maputil.ToSortedSlice(detectorHolder.files) {
				for _, lineNumber := range maputil.ToSortedSlice(fileHolder.lineNumbers) {
					constructedDetector.Locations = append(constructedDetector.Locations, Location{
						Filename:   fileHolder.name,
						LineNumber: lineNumber,
					})
				}
			}
			constructedDatatype.Detectors = append(constructedDatatype.Detectors, constructedDetector)
		}

		data.Datatypes = append(data.Datatypes, constructedDatatype)
	}

	return data
}

func (holder *dataFlowHolder) addSchema(detection detections.Detection) error {
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
		holder.addDatatype(value.FieldName, string(detection.DetectorType), detection.Source.Filename, *detection.Source.LineNumber)
	}

	return nil
}

// addDatatype adds datatype to hash list and at the same time blocks duplicates
func (holder *dataFlowHolder) addDatatype(datatypeName string, detectorName string, fileName string, lineNumber int) {
	// create datatype entry if it doesn't exist
	if _, exists := holder.datatypes[datatypeName]; !exists {
		holder.datatypes[datatypeName] = &datatypeHolder{
			name:      datatypeName,
			detectors: make(map[string]*detectorHolder),
		}
	}

	datatype := holder.datatypes[datatypeName]
	// create detector entry if it doesn't exist
	if _, exists := datatype.detectors[detectorName]; !exists {
		datatype.detectors[detectorName] = &detectorHolder{
			name:  detectorName,
			files: make(map[string]*fileHolder),
		}
	}

	detector := datatype.detectors[detectorName]
	// create file entry if it doesn't exist
	if _, exists := detector.files[fileName]; !exists {
		detector.files[fileName] = &fileHolder{
			name:        fileName,
			lineNumbers: make(map[int]int),
		}
	}

	file := datatype.detectors[detectorName].files[fileName]
	// create line number entry if it doesn't exist
	if _, exists := detector.files[fileName]; !exists {
		detector.files[fileName] = &fileHolder{
			name:        fileName,
			lineNumbers: make(map[int]int),
		}
	}

	file.lineNumbers[lineNumber] = lineNumber

}
