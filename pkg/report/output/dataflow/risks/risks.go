package risks

import (
	"bytes"
	"encoding/json"
	"sort"

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
	files        map[string]map[int]*fileHolder // group files by filename
}

type fileHolder struct {
	name        string
	lineNumber  int
	parent      *schema.Parent
	fieldName   string
	objectName  string
	subjectName *string
}

func New(config settings.Config, isInternal bool) *Holder {
	return &Holder{
		detectors:  make(map[string]detectorHolder),
		config:     config,
		isInternal: isInternal,
	}
}

func (holder *Holder) AddRiskPresence(detection detections.Detection) {
	// create entry if it doesn't exist
	ruleName := string(detection.DetectorType)
	if _, exists := holder.detectors[ruleName]; !exists {
		holder.detectors[ruleName] = detectorHolder{
			id:        ruleName,
			datatypes: make(map[string]*datatypeHolder),
		}
	}


	detector := holder.detectors[ruleName]

	if _, exists := detector[]

	riskLocation := &types.RiskLocation{
		Filename:   detection.Source.Filename,
		LineNumber: *detection.Source.LineNumber,
	}

	// add parent information if possible
	parent := extractCustomRiskParent(detection.Value)
	if parent != nil {
		riskLocation.Parent = parent
	}

	if detection.DetectorType == detectors.DetectorGitleaks {
		value := detection.Value.(map[string]interface{})["description"]

		holder.presentRisks[ruleName].Locations = append(holder.presentRisks[ruleName].Locations, types.RiskDetectionLocation{
			RiskLocation: riskLocation,
			Content:      value.(string),
		})
	} else {
		holder.presentRisks[ruleName].Locations = append(holder.presentRisks[ruleName].Locations, types.RiskDetectionLocation{
			RiskLocation: riskLocation,
			Content:      *detection.Source.Text,
		})
	}
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
		)
	}

	return nil
}

// addDatatype adds detector to hash list and at the same time blocks duplicates
func (holder *Holder) addDatatype(ruleName string, datatype *db.DataType, subjectName *string, fileName string, lineNumber int, schema schema.Schema) {
	if datatype == nil {
		// FIXME: we end up with empty field Name and no datatype with the new code
		// Might be related to the bug with the Unique Identifier classification
		return
	}

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
				files:        make(map[string]map[int]*fileHolder),
			}
		} else {
			detector.datatypes[datatype.Name] = &datatypeHolder{
				name:  datatype.Name,
				files: make(map[string]map[int]*fileHolder),
			}
		}
	}

	detectorDatatype := detector.datatypes[datatype.Name]
	// create file entry if it doesn't exist
	if _, exists := detectorDatatype.files[fileName]; !exists {
		detectorDatatype.files[fileName] = make(map[int]*fileHolder, 0)
	}

	detectorDatatype.files[fileName][lineNumber] = &fileHolder{
		name:        fileName,
		lineNumber:  lineNumber,
		parent:      schema.Parent,
		fieldName:   schema.FieldName,
		objectName:  schema.ObjectName,
		subjectName: subjectName,
	}
}

func (holder *Holder) ToDataFlow() []interface{} {
	data := make([]interface{}, 0)

	detectors := maputil.ToSortedSlice(holder.detectors)

	for _, detector := range detectors {
		constructedDetector := types.RiskDetector{
			DetectorID: detector.id,
			DataTypes:  make([]types.RiskDatatype, 0),
		}

		datatypes := maputil.ToSortedSlice(detector.datatypes)

		for _, datatype := range datatypes {
			stored := false
			if customDetector, isCustomDetector := holder.config.Rules[detector.id]; isCustomDetector {
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
			for _, locations := range files {
				sortedLocations := maputil.ToSortedSlice(locations)
				for _, location := range sortedLocations {
					constructedDatatype.Locations = append(constructedDatatype.Locations, types.RiskLocation{
						Filename:    location.name,
						LineNumber:  location.lineNumber,
						Parent:      location.parent,
						FieldName:   location.fieldName,
						ObjectName:  location.objectName,
						SubjectName: location.subjectName,
					})
				}
			}
			constructedDetector.DataTypes = append(constructedDetector.DataTypes, constructedDatatype)
		}

		data = append(data, constructedDetector)
	}

	sortedRisks := maputil.ToSortedSlice(holder.presentRisks)
	for _, presentRisk := range sortedRisks {
		sortLocations(presentRisk)
		data = append(data, presentRisk)
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

func sortLocations(risk *types.RiskDetection) {
	sort.Slice(risk.Locations, func(i, j int) bool {
		locationA := risk.Locations[i]
		locationB := risk.Locations[j]

		if locationA.Filename < locationB.Filename {
			return true
		}
		if locationA.Filename > locationB.Filename {
			return false
		}

		if locationA.LineNumber < locationB.LineNumber {
			return true
		}
		if locationA.LineNumber > locationB.LineNumber {
			return false
		}

		if locationA.Parent == nil && locationB.Parent != nil {
			return true
		}
		if locationA.Parent != nil && locationB.Parent == nil {
			return false
		}
		if locationA.Parent != nil {
			if locationA.Parent.LineNumber < locationB.Parent.LineNumber {
				return true
			}
			if locationA.Parent.LineNumber > locationB.Parent.LineNumber {
				return false
			}

			if locationA.Parent.Content < locationB.Parent.Content {
				return true
			}
			if locationA.Parent.Content > locationB.Parent.Content {
				return false
			}
		}

		return locationA.Content < locationB.Content
	})
}
