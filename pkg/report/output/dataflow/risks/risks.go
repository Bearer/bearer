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
	name       string
	lineNumber map[int]lineHolder // group detections by line number
}

type lineHolder struct {
	lineNumber       int
	dataTypeCategory map[string]dataTypeCategoryHolder // group detections by datatype category
}

type dataTypeCategoryHolder struct {
	name     string
	category string
	parent   *schema.Parent
	dataType map[string]dataTypeHolder
}

type dataTypeHolder struct {
	parent      *schema.Parent
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
	lineNumber := *detection.Source.LineNumber
	// can be nil
	parent := extractCustomRiskParent(detection.Value)
	var content string

	if detection.DetectorType == detectors.DetectorGitleaks {
		value := detection.Value.(map[string]interface{})["description"]
		content = value.(string)
	} else {
		content = *detection.Source.Text
	}

	holder.addDatatype(ruleName, &db.DataType{Name: content}, nil, fileName, lineNumber, schema.Schema{Parent: parent}, categoryPresence)
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
			*detection.Source.LineNumber,
			schema,
			categoryDatatype,
		)
	}

	return nil
}

// addDatatype adds detector to hash list and at the same time blocks duplicates
func (holder *Holder) addDatatype(ruleName string, datatype *db.DataType, subjectName *string, fileName string, lineNumber int, schema schema.Schema, category string) {
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
			name:       fileName,
			lineNumber: make(map[int]lineHolder),
		}
	}

	file := detector.files[fileName]
	// create line number entry if it doesn't exist
	if _, exists := file.lineNumber[lineNumber]; !exists {
		file.lineNumber[lineNumber] = lineHolder{
			lineNumber:       lineNumber,
			dataTypeCategory: make(map[string]dataTypeCategoryHolder),
		}
	}

	line := file.lineNumber[lineNumber]
	// create datatype category entry if it doesn't exist
	if _, exists := line.dataTypeCategory[datatype.Name]; !exists {
		categoryToAdd := dataTypeCategoryHolder{
			name:     datatype.Name,
			category: category,
			dataType: make(map[string]dataTypeHolder),
		}

		if category == "presence" {
			categoryToAdd.parent = schema.Parent
		}

		line.dataTypeCategory[datatype.Name] = categoryToAdd
	}

	if category == "datatype" {
		datatypeCategory := line.dataTypeCategory[datatype.Name]
		datatypeKey := schema.FieldName + schema.ObjectName
		// create datatype if it doesn't exists
		if _, exists := datatypeCategory.dataType[datatypeKey]; !exists {
			datatypeCategory.dataType[datatypeKey] = dataTypeHolder{
				parent:      schema.Parent,
				fieldName:   &schema.FieldName,
				objectName:  &schema.ObjectName,
				subjectName: subjectName,
				content:     nil,
			}
		}
	}
}

func (holder *Holder) ToDataFlow() []interface{} {
	data := make([]interface{}, 0)

	for _, detector := range maputil.ToSortedSlice(holder.detectors) {
		stored := false
		if customDetector, isCustomDetector := holder.config.Rules[detector.id]; isCustomDetector {
			stored = customDetector.Stored
		}

		constructedDetector := types.RiskDetector{
			DetectorID: detector.id,
		}
		locations := []types.RiskLocation{}

		for _, file := range maputil.ToSortedSlice(detector.files) {

			for _, line := range maputil.ToSortedSlice(file.lineNumber) {
				location := types.RiskLocation{
					Filename:   file.name,
					LineNumber: line.lineNumber,
				}

				for _, dataTypeCategory := range maputil.ToSortedSlice(line.dataTypeCategory) {
					category := types.RiskDatatypeCategory{
						Name:     dataTypeCategory.name,
						Category: dataTypeCategory.category,
					}

					if category.Category == categoryPresence {
						category.Parent = dataTypeCategory.parent
						category.Stored = &stored
					}

					for _, dataType := range maputil.ToSortedSlice(dataTypeCategory.dataType) {
						riskDatatype := types.RiskDatatype{
							Parent:      dataType.parent,
							SubjectName: dataType.subjectName,
							Stored:      stored,
						}

						if dataType.fieldName != nil {
							riskDatatype.FieldName = *dataType.fieldName
						}
						if dataType.objectName != nil {
							riskDatatype.ObjectName = *dataType.objectName
						}

						category.DataTypes = append(category.DataTypes, riskDatatype)
					}

					location.DataTypeCategories = append(location.DataTypeCategories, category)
				}

				locations = append(locations, location)
			}

		}
		constructedDetector.Locations = locations

		data = append(data, constructedDetector)
	}

	return data
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
