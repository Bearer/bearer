package components

import (
	"regexp"
	"strings"

	"github.com/bearer/bearer/pkg/report/output/dataflow/types"

	dependenciesclassification "github.com/bearer/bearer/pkg/classification/dependencies"
	frameworkclassification "github.com/bearer/bearer/pkg/classification/frameworks"
	interfaceclassification "github.com/bearer/bearer/pkg/classification/interfaces"
	"github.com/bearer/bearer/pkg/util/classify"
	"github.com/bearer/bearer/pkg/util/maputil"
)

type Holder struct {
	dependencies map[string][]*dependency // group dependencies by detector name
	components   map[string]*component    // group components by name
	isInternal   bool
}

type dependency struct {
	name     string
	filename string
	// version  types.Version
	version string
}

type component struct {
	name               string
	component_type     string
	component_sub_type string
	uuid               string
	detectors          map[string]*detector // group detectors by detectorName
}

type detector struct {
	name  string
	files map[string]*fileHolder // group files by filename
}

type fileHolder struct {
	name        string
	fullName    string
	lineNumbers map[int]int //group lines by linenumber
}

var (
	unwantedVersionCharRegex = regexp.MustCompile(`[^0-9.]+`)
)

func New(isInternal bool) *Holder {
	return &Holder{
		dependencies: make(map[string][]*dependency),
		components:   make(map[string]*component),
		isInternal:   isInternal,
	}
}

func getComponentType(recipeType string, reason string) string {
	if recipeType != "" {
		return recipeType
	}

	if strings.HasPrefix(reason, "internal") {
		return "internal_service"
	}
	return "external_service"
}

func (holder *Holder) AddInterface(classifiedDetection interfaceclassification.ClassifiedInterface) error {
	if classifiedDetection.Classification == nil {
		return nil
	}
	componentType := getComponentType(classifiedDetection.Classification.RecipeType, classifiedDetection.Classification.Decision.Reason)
	componentSubType := classifiedDetection.Classification.RecipeSubType

	componentUUID := classifiedDetection.Classification.RecipeUUID
	if componentUUID == "" {
		componentUUID = classifiedDetection.Classification.Name()
	}

	if classifiedDetection.Classification.Decision.State == classify.Valid {
		holder.addComponent(
			classifiedDetection.Classification.Name(),
			componentType,
			componentSubType,
			componentUUID,
			string(classifiedDetection.DetectorType),
			classifiedDetection.Source.Filename,
			classifiedDetection.Source.FullFilename,
			*classifiedDetection.Source.StartLineNumber,
		)
	}

	return nil
}

func (holder *Holder) AddDependency(classifiedDetection dependenciesclassification.ClassifiedDependency) error {
	if classifiedDetection.Value != nil {
		value := classifiedDetection.Value.(map[string]interface{})
		version := convertVersion(value["version"].(string))
		name := value["name"].(string)

		holder.addDependency(
			string(classifiedDetection.DetectorType),
			classifiedDetection.Source.Filename,
			name,
			version,
		)
	}

	if classifiedDetection.Classification == nil {
		return nil
	}

	componentType := getComponentType(classifiedDetection.Classification.RecipeType, classifiedDetection.Classification.Decision.Reason)
	componentSubType := classifiedDetection.Classification.RecipeSubType

	if classifiedDetection.Classification.Decision.State == classify.Valid {
		holder.addComponent(
			classifiedDetection.Classification.RecipeName,
			componentType,
			componentSubType,
			classifiedDetection.Classification.RecipeUUID,
			string(classifiedDetection.DetectorType),
			classifiedDetection.Source.Filename,
			classifiedDetection.Source.FullFilename,
			*classifiedDetection.Source.StartLineNumber,
		)
	}

	return nil
}

func convertVersion(version string) string {
	return unwantedVersionCharRegex.ReplaceAllString(version, "")
}

func (holder *Holder) AddFramework(classifiedDetection frameworkclassification.ClassifiedFramework) error {
	if classifiedDetection.Classification == nil {
		return nil
	}

	componentType := getComponentType(classifiedDetection.Classification.Decision.Reason, classifiedDetection.Classification.Decision.Reason)
	componentSubType := classifiedDetection.Classification.RecipeSubType

	if classifiedDetection.Classification.Decision.State == classify.Valid {
		holder.addComponent(
			classifiedDetection.Classification.RecipeName,
			componentType,
			componentSubType,
			classifiedDetection.Classification.RecipeUUID,
			string(classifiedDetection.DetectorType),
			classifiedDetection.Source.Filename,
			classifiedDetection.Source.FullFilename,
			*classifiedDetection.Source.StartLineNumber,
		)
	}

	return nil
}

// addComponent adds component to hash list and at the same time blocks duplicates
func (holder *Holder) addDependency(
	detectorName string,
	fileName string,
	name string,
	version string,
) {
	if _, exists := holder.dependencies[detectorName]; !exists {
		holder.dependencies[detectorName] = make([]*dependency, 0)
	}

	holder.dependencies[detectorName] = append(
		holder.dependencies[detectorName],
		&dependency{
			name:     name,
			version:  version,
			filename: fileName,
		},
	)
}

// addComponent adds component to hash list and at the same time blocks duplicates
func (holder *Holder) addComponent(
	componentName string,
	componentType string,
	componentSubType string,
	componentUUID string,
	detectorName string,
	fileName string,
	fullFilename string,
	lineNumber int,
) {
	// create component entry if it doesn't exist
	if _, exists := holder.components[componentUUID]; !exists {
		var uuid string
		if holder.isInternal {
			uuid = componentUUID
		}
		holder.components[componentUUID] = &component{
			name:               componentName,
			component_type:     componentType,
			component_sub_type: componentSubType,
			uuid:               uuid,
			detectors:          make(map[string]*detector),
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
			fullName:    fullFilename,
			lineNumbers: make(map[int]int),
		}
	}

	targetDetector.files[fileName].lineNumbers[lineNumber] = lineNumber
}

func (holder *Holder) ToDataFlowForDependencies() []types.Dependency {
	data := make([]types.Dependency, 0)

	for detectorName, dependencies := range holder.dependencies {
		for _, dependency := range dependencies {
			data = append(data, types.Dependency{
				Name:     dependency.name,
				Version:  dependency.version,
				Filename: dependency.filename,
				Detector: detectorName,
			})
		}
	}

	return data
}

func (holder *Holder) ToDataFlow() []types.Component {
	data := make([]types.Component, 0)

	availableComponents := maputil.ToSortedSlice(holder.components)

	for _, targetComponent := range availableComponents {
		constructedComponent := types.Component{
			Name:      targetComponent.name,
			Type:      targetComponent.component_type,
			SubType:   targetComponent.component_sub_type,
			UUID:      targetComponent.uuid,
			Locations: make([]types.ComponentLocation, 0),
		}

		for _, targetDetector := range maputil.ToSortedSlice(targetComponent.detectors) {
			for _, targetFile := range maputil.ToSortedSlice(targetDetector.files) {
				for _, targetLineNumber := range maputil.ToSortedSlice(targetFile.lineNumbers) {
					constructedComponent.Locations = append(constructedComponent.Locations, types.ComponentLocation{
						Filename:     targetFile.name,
						FullFilename: targetFile.fullName,
						LineNumber:   targetLineNumber,
						Detector:     targetDetector.name,
					})
				}
			}
		}

		data = append(data, constructedComponent)
	}

	return data
}
