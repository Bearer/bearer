package components

import (
	"strings"

	"github.com/bearer/curio/pkg/report/output/dataflow/detectiondecoder"
	"github.com/bearer/curio/pkg/report/output/dataflow/types"

	"github.com/bearer/curio/pkg/util/classify"
	"github.com/bearer/curio/pkg/util/maputil"
)

type Holder struct {
	components map[string]*component // group components by name
}

type component struct {
	name      string
	detectors map[string]*detector // group detectors by detectorName
}
type detector struct {
	name  string
	files map[string]*fileHolder // group files by filename
}
type fileHolder struct {
	name        string
	lineNumbers map[int]int //group lines by linenumber
}

func New() *Holder {
	return &Holder{
		components: make(map[string]*component),
	}
}

func (holder *Holder) AddInterface(detection interface{}) error {
	value, err := detectiondecoder.GetClassifiedInterface(detection)
	if err != nil {
		return err
	}

	if value.Classification == nil {
		return nil
	}

	if value.Classification.Decision.State == classify.Valid {
		holder.addComponent(strings.ToLower(value.Classification.RecipeName), string(value.DetectorType), value.Source.Filename, *value.Source.LineNumber)
	}

	return nil
}

func (holder *Holder) AddDependency(detection interface{}) error {
	value, err := detectiondecoder.GetClassifiedDependency(detection)
	if err != nil {
		return err
	}

	if value.Classification == nil {
		return nil
	}

	if value.Classification.Decision.State == classify.Valid {
		holder.addComponent(strings.ToLower(value.Classification.RecipeName), string(value.DetectorType), value.Source.Filename, *value.Source.LineNumber)
	}

	return nil
}

// addComponent adds component to hash list and at the same time blocks duplicates
func (holder *Holder) addComponent(componentName string, detectorName string, fileName string, lineNumber int) {
	// create component entry if it doesn't exist
	if _, exists := holder.components[componentName]; !exists {
		holder.components[componentName] = &component{
			name:      componentName,
			detectors: make(map[string]*detector),
		}
	}

	targetComponent := holder.components[componentName]
	// create detector entry if it doesn't exist
	if _, exists := targetComponent.detectors[detectorName]; !exists {
		targetComponent.detectors[detectorName] = &detector{
			name:  detectorName,
			files: make(map[string]*fileHolder),
		}
	}

	targetDetector := targetComponent.detectors[detectorName]
	// create file entry if it doesn't exist
	if _, exists := targetDetector.files[fileName]; !exists {
		targetDetector.files[fileName] = &fileHolder{
			name:        fileName,
			lineNumbers: make(map[int]int),
		}
	}

	targetDetector.files[fileName].lineNumbers[lineNumber] = lineNumber

}

func (holder *Holder) ToDataFlow() []types.Component {
	data := make([]types.Component, 0)

	availableComponents := maputil.ToSortedSlice(holder.components)

	for _, targetComponent := range availableComponents {
		constructedComponent := types.Component{
			Name:      targetComponent.name,
			Locations: make([]types.ComponentLocation, 0),
		}

		for _, targetDetector := range maputil.ToSortedSlice(targetComponent.detectors) {
			for _, targetFile := range maputil.ToSortedSlice(targetDetector.files) {
				for _, targetLineNumber := range maputil.ToSortedSlice(targetFile.lineNumbers) {
					constructedComponent.Locations = append(constructedComponent.Locations, types.ComponentLocation{
						Filename:   targetFile.name,
						LineNumber: targetLineNumber,
						Detector:   targetDetector.name,
					})
				}
			}
		}

		data = append(data, constructedComponent)
	}

	return data
}
