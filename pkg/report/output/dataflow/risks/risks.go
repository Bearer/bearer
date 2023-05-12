package risks

import (
	"bytes"
	"encoding/json"

	"github.com/bearer/bearer/pkg/classification/db"
	"github.com/bearer/bearer/pkg/commands/process/settings"
	"github.com/bearer/bearer/pkg/report/detectors"
	"github.com/bearer/bearer/pkg/report/output/dataflow/detectiondecoder"
	"github.com/bearer/bearer/pkg/report/output/dataflow/types"
	"github.com/bearer/bearer/pkg/report/schema"

	"github.com/bearer/bearer/pkg/report/detections"
	"github.com/bearer/bearer/pkg/util/classify"
	"github.com/bearer/bearer/pkg/util/maputil"
)

type Holder struct {
	detectors map[string]detectorHolder // group detections by detector name
	config    settings.Config
}

type detectorHolder struct {
	id    string
	files map[string]fileHolder // group detectors by file name
}

type fileHolder struct {
	name            string
	startLineNumber map[int]lineHolder // group detections by line number
}

type lineHolder struct {
	startLineNumber int
	parent          map[string]parentHolder // group detections by parent
}

type parentHolder struct {
	name    string
	parent  *schema.Parent
	matches map[string]matchCategoryHolder // group detections by datatype category
}

type matchCategoryHolder struct {
	name         string
	categoryUUID *string
	category     string
	dataType     map[string]dataTypeHolder
}

type dataTypeHolder struct {
	content     *string
	fieldName   *string
	objectName  *string
	subjectName *string
}

var categoryPresence = "presence"
var categoryDatatype = "datatype"

func New(config settings.Config, isInternal bool) *Holder {
	return &Holder{
		detectors: make(map[string]detectorHolder),
		config:    config,
	}
}

func (holder *Holder) AddRiskPresence(detection detections.Detection) {
	// create entry if it doesn't exist
	ruleName := string(detection.DetectorType)
	fileName := detection.Source.Filename
	startLineNumber := *detection.Source.StartLineNumber

	var parent *schema.Parent
	var content string

	if detection.DetectorType == detectors.DetectorGitleaks {
		value := detection.Value.(map[string]interface{})["description"]
		content = value.(string)
		parent = &schema.Parent{
			StartLineNumber:   *detection.Source.StartLineNumber,
			StartColumnNumber: *detection.Source.StartColumnNumber,
			EndLineNumber:     *detection.Source.EndLineNumber,
			EndColumnNumber:   *detection.Source.EndColumnNumber,
			Content:           *detection.Source.Text,
		}
	} else {
		// parent can be nil
		parent = extractCustomRiskParent(detection.Value)
		content = *detection.Source.Text
	}

	holder.addDatatype(ruleName, &db.DataType{Name: content}, nil, fileName, startLineNumber, schema.Schema{Parent: parent}, categoryPresence)
}

func (holder *Holder) AddSchema(detection detections.Detection) error {
	schema, err := detectiondecoder.GetSchema(detection)
	if err != nil {
		return err
	}

	classification, err := detectiondecoder.GetSchemaClassification(schema)
	if err != nil {
		return err
	}

	if classification.Decision.State == classify.Valid {
		holder.addDatatype(
			string(detection.DetectorType),
			classification.DataType,
			classification.SubjectName,
			detection.Source.Filename,
			*detection.Source.StartLineNumber,
			schema,
			categoryDatatype,
		)
	}

	return nil
}

// addDatatype adds detector to hash list and at the same time blocks duplicates
func (holder *Holder) addDatatype(ruleName string, datatype *db.DataType, subjectName *string, fileName string, startLineNumber int, schema schema.Schema, category string) {
	if datatype == nil {
		// FIXME: we end up with empty field Name and no datatype with the new code
		// Might be related to the bug with the Unique Identifier classification
		return
	}

	// create detector entry if it doesn't exist
	if _, exists := holder.detectors[ruleName]; !exists {
		holder.detectors[ruleName] = detectorHolder{
			id:    ruleName,
			files: make(map[string]fileHolder),
		}
	}

	detector := holder.detectors[ruleName]
	// create file entry if it doesn't exist
	if _, exists := detector.files[fileName]; !exists {
		detector.files[fileName] = fileHolder{
			name:            fileName,
			startLineNumber: make(map[int]lineHolder),
		}
	}

	file := detector.files[fileName]
	// create line number entry if it doesn't exist
	if _, exists := file.startLineNumber[startLineNumber]; !exists {
		file.startLineNumber[startLineNumber] = lineHolder{
			startLineNumber: startLineNumber,
			parent:          make(map[string]parentHolder),
		}
	}

	line := file.startLineNumber[startLineNumber]
	// create datatype parent entry if it doesn't exist
	parentKey := "undefined_parent"
	if schema.Parent != nil {
		parentKey = schema.Parent.Content
	}

	if _, exists := line.parent[parentKey]; !exists {
		line.parent[parentKey] = parentHolder{
			name:    parentKey,
			parent:  schema.Parent,
			matches: make(map[string]matchCategoryHolder),
		}
	}

	parent := line.parent[parentKey]
	// create datatype category if it doesn't exist
	if _, exists := parent.matches[datatype.Name]; !exists {
		categoryToAdd := matchCategoryHolder{
			name:         datatype.Name,
			category:     category,
			categoryUUID: &datatype.CategoryUUID,
			dataType:     make(map[string]dataTypeHolder),
		}

		parent.matches[datatype.Name] = categoryToAdd
	}

	if category == "datatype" {
		datatypeCategory := parent.matches[datatype.Name]
		datatypeKey := schema.FieldName + schema.ObjectName
		// create datatype if it doesn't exists
		if _, exists := datatypeCategory.dataType[datatypeKey]; !exists {
			datatypeCategory.dataType[datatypeKey] = dataTypeHolder{
				fieldName:   &schema.FieldName,
				objectName:  &schema.ObjectName,
				subjectName: subjectName,
				content:     nil,
			}
		}
	}
}

func (holder *Holder) ToDataFlow() []types.RiskDetector {
	data := make([]types.RiskDetector, 0)

	for _, detector := range maputil.ToSortedSlice(holder.detectors) {
		constructedDetector := types.RiskDetector{
			DetectorID: detector.id,
		}
		locations := []types.RiskLocation{}

		for _, file := range maputil.ToSortedSlice(detector.files) {

			for _, line := range maputil.ToSortedSlice(file.startLineNumber) {

				for _, parent := range maputil.ToSortedSlice(line.parent) {

					location := types.RiskLocation{
						Filename:        file.name,
						StartLineNumber: line.startLineNumber,
						Parent:          parent.parent,
					}

					hasDatatype := false
					matches := maputil.ToSortedSlice(parent.matches)
					for _, dataType := range matches {
						if dataType.category == categoryDatatype {
							hasDatatype = true
						}
					}

					for _, dataTypeCategory := range matches {
						if dataTypeCategory.category == categoryPresence {
							if hasDatatype {
								continue
							}

							location.PresenceMatches = append(location.PresenceMatches, types.RiskPresence{
								Name: dataTypeCategory.name,
							})
							continue
						}

						match := types.RiskDatatype{
							Name:         dataTypeCategory.name,
							CategoryUUID: *dataTypeCategory.categoryUUID,
						}

						for _, dataType := range maputil.ToSortedSlice(dataTypeCategory.dataType) {
							riskSchema := types.RiskSchema{
								SubjectName: dataType.subjectName,
							}

							if dataType.fieldName != nil {
								riskSchema.FieldName = *dataType.fieldName
							}
							if dataType.objectName != nil {
								riskSchema.ObjectName = *dataType.objectName
							}

							match.Schemas = append(match.Schemas, riskSchema)
						}

						location.DataTypes = append(location.DataTypes, match)
					}

					locations = append(locations, location)
				}

			}

		}
		constructedDetector.Locations = locations

		data = append(data, constructedDetector)
	}

	data = removeParentBasedDuplicates(data)

	return data
}

// removeParentBasedDuplicates checks if there are 2 risk locations one with presence and one with datatype which have same parent line number and parentContent and if it finds such case it discards the presence one
func removeParentBasedDuplicates(data []types.RiskDetector) []types.RiskDetector {
	filteredData := []types.RiskDetector{}
	for _, detector := range data {
		newDetector := types.RiskDetector{
			DetectorID: detector.DetectorID,
		}
		for _, location := range detector.Locations {
			// presence matches are always alone per location
			if len(location.PresenceMatches) > 0 && location.Parent != nil {
				hasSameParentLocation := false

				for _, otherLocation := range detector.Locations {
					if len(otherLocation.DataTypes) > 0 &&
						otherLocation.Filename == location.Filename &&
						otherLocation.Parent != nil &&
						otherLocation.Parent.Content == location.Parent.Content &&
						otherLocation.Parent.StartLineNumber == location.Parent.StartLineNumber {

						hasSameParentLocation = true
					}
				}

				if hasSameParentLocation {
					continue
				}
			}

			newDetector.Locations = append(newDetector.Locations, location)
		}
		filteredData = append(filteredData, newDetector)
	}

	return filteredData
}

func extractCustomRiskParent(value interface{}) *schema.Parent {
	if value == nil {
		return nil
	}

	var parent schema.Parent
	buf := bytes.NewBuffer(nil)
	err := json.NewEncoder(buf).Encode(value)
	if err != nil {
		return nil
	}

	err = json.NewDecoder(buf).Decode(&parent)
	if err != nil {
		return nil
	}

	return &parent
}
