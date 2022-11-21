package datatypes

import (
	"github.com/bearer/curio/pkg/report/output/dataflow/detectiondecoder"
	"github.com/bearer/curio/pkg/report/output/dataflow/types"

	"github.com/bearer/curio/pkg/report/detections"
	"github.com/bearer/curio/pkg/util/classify"
	"github.com/bearer/curio/pkg/util/maputil"
)

type Holder struct {
	datatypes map[string]datatypeHolder // group datatypeHolders by name
}

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

func New() *Holder {
	return &Holder{
		datatypes: make(map[string]datatypeHolder),
	}
}

func (holder *Holder) AddSchema(detection detections.Detection) error {
	classification, err := detectiondecoder.GetSchemaClassification(detection)
	if err != nil {
		return err
	}

	if classification.Decision.State == classify.Valid {
		holder.addDatatype(classification.DataType.Name, string(detection.DetectorType), detection.Source.Filename, *detection.Source.LineNumber)
	}

	return nil
}

// addDatatype adds datatype to hash list and at the same time blocks duplicates
func (holder *Holder) addDatatype(datatypeName string, detectorName string, fileName string, lineNumber int) {
	// create datatype entry if it doesn't exist
	if _, exists := holder.datatypes[datatypeName]; !exists {
		holder.datatypes[datatypeName] = datatypeHolder{
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

func (holder *Holder) ToDataFlow() []types.Datatype {
	data := make([]types.Datatype, 0)

	datatypes := maputil.ToSortedSlice(holder.datatypes)

	for _, datatype := range datatypes {
		constructedDatatype := types.Datatype{
			Name: datatype.name,
		}

		detectors := maputil.ToSortedSlice(datatype.detectors)

		for _, detectorHolder := range detectors {
			constructedDetector := types.DatatypeDetector{
				Name:      detectorHolder.name,
				Locations: make([]types.DatatypeLocation, 0),
			}

			for _, fileHolder := range maputil.ToSortedSlice(detectorHolder.files) {
				for _, lineNumber := range maputil.ToSortedSlice(fileHolder.lineNumbers) {
					constructedDetector.Locations = append(constructedDetector.Locations, types.DatatypeLocation{
						Filename:   fileHolder.name,
						LineNumber: lineNumber,
					})
				}
			}
			constructedDatatype.Detectors = append(constructedDatatype.Detectors, constructedDetector)
		}

		data = append(data, constructedDatatype)
	}

	return data
}
