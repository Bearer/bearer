package datatypes

import (
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
	datatypes  map[string]datatypeHolder // group datatypeHolders by name
	config     settings.Config
	isInternal bool
}

type datatypeHolder struct {
	name         string
	uuid         string
	categoryUUID string
	categoryName string
	detectors    map[string]*detectorHolder // group detectors by detectorName
}

type detectorHolder struct {
	name  string
	files map[string]*fileHolder // group files by filename
}

type fileHolder struct {
	name        string
	lineNumbers map[int]*lineNumberHolder
}

type lineNumberHolder struct {
	startLineNumber   int
	startColumnNumber int
	endColumnNumber   int
	encrypted         *bool
	verifiedBy        []types.DatatypeVerifiedBy
	stored            *bool
	parent            *schema.Parent
	fieldName         string
	objectName        string
	subjectName       *string
}

func New(config settings.Config, isInternal bool) *Holder {
	return &Holder{
		datatypes:  make(map[string]datatypeHolder),
		config:     config,
		isInternal: isInternal,
	}
}

func (holder *Holder) AddSchema(detection detections.Detection, extras *ExtraFields) error {
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
			classification.DataType,
			string(detection.DetectorType),
			detection.Source.Filename,
			*detection.Source.StartLineNumber,
			*detection.Source.StartColumnNumber,
			*detection.Source.EndColumnNumber,
			classification.SubjectName,
			extras,
			schema,
		)
	}

	return nil
}

// addDatatype adds datatype to hash list and at the same time blocks duplicates
func (holder *Holder) addDatatype(
	classification *db.DataType,
	detectorName string,
	fileName string,
	lineNumber int,
	startColumnNumber int,
	endColumnNumber int,
	subjectName *string,
	extras *ExtraFields,
	schema schema.Schema,
) {
	// create datatype entry if it doesn't exist
	if _, exists := holder.datatypes[classification.Name]; !exists {
		datatype := datatypeHolder{
			name:         classification.Name,
			categoryName: classification.Category.Name,
			detectors:    make(map[string]*detectorHolder),
		}

		if holder.isInternal {
			datatype.categoryUUID = classification.CategoryUUID
			datatype.uuid = classification.UUID
		}

		holder.datatypes[classification.Name] = datatype
	}

	datatype := holder.datatypes[classification.Name]
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
			lineNumbers: make(map[int]*lineNumberHolder),
		}
	}

	file := datatype.detectors[detectorName].files[fileName]
	// create line number entry if it doesn't exist
	if _, exists := file.lineNumbers[lineNumber]; !exists {
		file.lineNumbers[lineNumber] = &lineNumberHolder{
			startLineNumber:   lineNumber,
			startColumnNumber: startColumnNumber,
			endColumnNumber:   endColumnNumber,
			fieldName:         schema.FieldName,
			objectName:        schema.ObjectName,
			subjectName:       subjectName,
			parent:            schema.Parent,
		}
	}

	lineEntry := file.lineNumbers[lineNumber]

	if extras != nil {
		lineEntry.encrypted = extras.encrypted
		lineEntry.verifiedBy = extras.verifiedBy
	}

	if detectorName == string(detectors.DetectorSchemaRb) {
		storedFlag := true
		lineEntry.stored = &storedFlag
	} else if customDetector, isCustomDetector := holder.config.Rules[detectorName]; isCustomDetector {
		if customDetector.Stored {
			storedFlag := true
			lineEntry.stored = &storedFlag
		}
	} else if customDetector, isCustomDetector := holder.config.BuiltInRules[detectorName]; isCustomDetector {
		if customDetector.Stored {
			storedFlag := true
			lineEntry.stored = &storedFlag
		}
	}
}

func (holder *Holder) ToDataFlow() []types.Datatype {
	data := make([]types.Datatype, 0)

	datatypes := maputil.ToSortedSlice(holder.datatypes)

	for _, datatype := range datatypes {
		constructedDatatype := types.Datatype{
			Name:         datatype.name,
			UUID:         datatype.uuid,
			CategoryUUID: datatype.categoryUUID,
			CategoryName: datatype.categoryName,
		}

		detectors := maputil.ToSortedSlice(datatype.detectors)

		for _, detectorHolder := range detectors {
			constructedDetector := types.DatatypeDetector{
				Name:      detectorHolder.name,
				Locations: make([]types.DatatypeLocation, 0),
			}

			for _, fileHolder := range maputil.ToSortedSlice(detectorHolder.files) {
				for _, lineNumber := range maputil.ToSortedSlice(fileHolder.lineNumbers) {
					location := types.DatatypeLocation{
						Filename:          fileHolder.name,
						StartLineNumber:   lineNumber.startLineNumber,
						StartColumnNumber: lineNumber.startColumnNumber,
						EndColumnNumber:   lineNumber.endColumnNumber,
						Encrypted:         lineNumber.encrypted,
						VerifiedBy:        lineNumber.verifiedBy,
						Stored:            lineNumber.stored,
						Parent:            lineNumber.parent,
						FieldName:         lineNumber.fieldName,
						ObjectName:        lineNumber.objectName,
						SubjectName:       lineNumber.subjectName,
					}
					constructedDetector.Locations = append(constructedDetector.Locations, location)
				}
			}
			constructedDatatype.Detectors = append(constructedDatatype.Detectors, constructedDetector)
		}

		data = append(data, constructedDatatype)
	}

	return data
}
