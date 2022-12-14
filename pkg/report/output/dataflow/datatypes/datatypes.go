package datatypes

import (
	"github.com/bearer/curio/pkg/classification/db"
	"github.com/bearer/curio/pkg/commands/process/settings"
	"github.com/bearer/curio/pkg/report/detectors"
	"github.com/bearer/curio/pkg/report/output/dataflow/detectiondecoder"
	"github.com/bearer/curio/pkg/report/output/dataflow/types"
	"github.com/bearer/curio/pkg/report/schema"

	"github.com/bearer/curio/pkg/report/detections"
	"github.com/bearer/curio/pkg/util/classify"
	"github.com/bearer/curio/pkg/util/maputil"
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
	parent       *schema.Parent
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
	lineNumber int
	encrypted  *bool
	verifiedBy []types.DatatypeVerifiedBy
	stored     *bool
	parent     *schema.Parent
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

	classification, err := detectiondecoder.GetSchemaClassification(schema, detection)
	if err != nil {
		return err
	}

	if classification.Decision.State == classify.Valid {
		holder.addDatatype(classification.DataType, string(detection.DetectorType), detection.Source.Filename, *detection.Source.LineNumber, extras, schema.Parent)
	}

	return nil
}

// addDatatype adds datatype to hash list and at the same time blocks duplicates
func (holder *Holder) addDatatype(classification *db.DataType, detectorName string, fileName string, lineNumber int, extras *ExtraFields, parent *schema.Parent) {
	// create datatype entry if it doesn't exist
	if _, exists := holder.datatypes[classification.Name]; !exists {
		datatype := datatypeHolder{
			name:      classification.Name,
			detectors: make(map[string]*detectorHolder),
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
			lineNumber: lineNumber,
			parent:     parent,
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
	} else if customDetector, isCustomDetector := holder.config.CustomDetector[detectorName]; isCustomDetector {
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
						Filename:   fileHolder.name,
						LineNumber: lineNumber.lineNumber,
						Encrypted:  lineNumber.encrypted,
						VerifiedBy: lineNumber.verifiedBy,
						Stored:     lineNumber.stored,
						Parent:     lineNumber.parent,
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
