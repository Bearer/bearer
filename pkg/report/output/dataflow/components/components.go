package components

import (
	"strings"

	"github.com/bearer/curio/pkg/report/output/dataflow/types"

	dependenciesclassification "github.com/bearer/curio/pkg/classification/dependencies"
	frameworkclassification "github.com/bearer/curio/pkg/classification/frameworks"
	interfaceclassification "github.com/bearer/curio/pkg/classification/interfaces"
	"github.com/bearer/curio/pkg/util/classify"
	"github.com/bearer/curio/pkg/util/maputil"
)

type Holder struct {
	components map[string]*component // group components by name
	isInternal bool
}

type component struct {
	name           string
	component_type string
	uuid           string
	detectors      map[string]*detector // group detectors by detectorName
}

type detector struct {
	name  string
	files map[string]*fileHolder // group files by filename
}

type fileHolder struct {
	name        string
	lineNumbers map[int]int //group lines by linenumber
}

func New(isInternal bool) *Holder {
	return &Holder{
		components: make(map[string]*component),
		isInternal: isInternal,
	}
}

func getComponentType(classificationReason string) string {
	if strings.HasPrefix(classificationReason, "internal_domain_") {
		return "internal_service"
	}

	if strings.HasPrefix(classificationReason, "recipe_match") {
		return "external_service"
	}

	return ""
}

func (holder *Holder) AddInterface(classifiedDetection interfaceclassification.ClassifiedInterface) error {
	if classifiedDetection.Classification == nil {
		return nil
	}

	componentType := getComponentType(classifiedDetection.Classification.Decision.Reason)

	componentUUID := classifiedDetection.Classification.RecipeUUID
	if componentUUID == "" {
		componentUUID = classifiedDetection.Classification.Name()
	}

	if classifiedDetection.Classification.Decision.State == classify.Valid {
		holder.addComponent(
			classifiedDetection.Classification.Name(),
			componentType,
			componentUUID,
			string(classifiedDetection.DetectorType),
			classifiedDetection.Source.Filename,
			*classifiedDetection.Source.LineNumber,
		)
	}

	return nil
}

func (holder *Holder) AddDependency(classifiedDetection dependenciesclassification.ClassifiedDependency) error {
	if classifiedDetection.Classification == nil {
		return nil
	}

	componentType := getComponentType(classifiedDetection.Classification.Decision.Reason)

	if classifiedDetection.Classification.Decision.State == classify.Valid {
		holder.addComponent(
			classifiedDetection.Classification.RecipeName,
			componentType,
			classifiedDetection.Classification.RecipeUUID,
			string(classifiedDetection.DetectorType),
			classifiedDetection.Source.Filename,
			*classifiedDetection.Source.LineNumber,
		)
	}

	return nil
}

func (holder *Holder) AddFramework(classifiedDetection frameworkclassification.ClassifiedFramework) error {
	if classifiedDetection.Classification == nil {
		return nil
	}

	componentType := getComponentType(classifiedDetection.Classification.Decision.Reason)

	if classifiedDetection.Classification.Decision.State == classify.Valid {
		holder.addComponent(
			classifiedDetection.Classification.RecipeName,
			componentType,
			classifiedDetection.Classification.RecipeUUID,
			string(classifiedDetection.DetectorType),
			classifiedDetection.Source.Filename,
			*classifiedDetection.Source.LineNumber,
		)
	}

	return nil
}

// addComponent adds component to hash list and at the same time blocks duplicates
func (holder *Holder) addComponent(componentName string, componentType string, componentUUID string, detectorName string, fileName string, lineNumber int) {
	// create component entry if it doesn't exist
	if _, exists := holder.components[componentUUID]; !exists {
		var uuid string
		if holder.isInternal {
			uuid = componentUUID
		}
		holder.components[componentUUID] = &component{
			name:           componentName,
			component_type: componentType,
			uuid:           uuid,
			detectors:      make(map[string]*detector),
		}
	}

	targetComponent := holder.components[componentUUID]
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
			Type:      targetComponent.component_type,
			UUID:      targetComponent.uuid,
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
