package risks

import (
	"github.com/bearer/curio/pkg/classification/db"
	"github.com/bearer/curio/pkg/commands/process/settings"
	"github.com/bearer/curio/pkg/report/output/dataflow/detectiondecoder"
	"github.com/bearer/curio/pkg/report/output/dataflow/types"
	"github.com/bearer/curio/pkg/report/schema"

	"github.com/bearer/curio/pkg/report/detections"
	"github.com/bearer/curio/pkg/util/classify"
	"github.com/bearer/curio/pkg/util/maputil"
)

type Holder struct {
	detectors    map[string]detectorHolder // group datatypeHolders by name
	config       settings.Config
	isInternal   bool
	presentRisks map[string]*types.RiskDetection
}

type detectorHolder struct {
	id        string
	datatypes map[string]*datatypeHolder // group detectors by detectorName
}

type datatypeHolder struct {
	name         string
	uuid         string
	categoryUUID string
	files        map[string]*fileHolder // group files by filename
}

type fileHolder struct {
	name       string
	lineNumber map[int]int
	parent     *schema.Parent
}

func New(config settings.Config, isInternal bool) *Holder {
	return &Holder{
		detectors:    make(map[string]detectorHolder),
		config:       config,
		isInternal:   isInternal,
		presentRisks: make(map[string]*types.RiskDetection),
	}
}

func (holder *Holder) AddRiskPresence(detection detections.Detection) {
	// create entry if it doesn't exist
	ruleName := string(detection.DetectorType)
	if _, exists := holder.presentRisks[ruleName]; !exists {
		holder.presentRisks[ruleName] = &types.RiskDetection{
			DetectorID: ruleName,
		}
	}

	holder.presentRisks[ruleName].Locations = append(holder.presentRisks[ruleName].Locations, types.RiskDetectionLocation{
		RiskLocation: &types.RiskLocation{
			Filename:   detection.Source.Filename,
			LineNumber: *detection.Source.LineNumber,
		},
		Content: *detection.Source.Text,
	})
}

func (holder *Holder) AddSchema(detection detections.Detection) error {
	schema, err := detectiondecoder.GetSchema(detection)
	if err != nil {
		return err
	}

	classification, err := detectiondecoder.GetSchemaClassification(schema, detection)
	if err != nil {
		return err
	}

	if classification.Decision.State == classify.Valid {
		holder.addDatatype(string(detection.DetectorType), classification.DataType, detection.Source.Filename, *detection.Source.LineNumber, schema.Parent)
	}

	return nil
}

// addDatatype adds detector to hash list and at the same time blocks duplicates
func (holder *Holder) addDatatype(ruleName string, datatype *db.DataType, fileName string, lineNumber int, parent *schema.Parent) {
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
				files:        make(map[string]*fileHolder),
			}
		} else {
			detector.datatypes[datatype.Name] = &datatypeHolder{
				name:  datatype.Name,
				files: make(map[string]*fileHolder),
				// parent: parent,
			}
		}
	}

	detectorDatatype := detector.datatypes[datatype.Name]
	// create file entry if it doesn't exist
	if _, exists := detectorDatatype.files[fileName]; !exists {
		detectorDatatype.files[fileName] = &fileHolder{
			name:       fileName,
			lineNumber: make(map[int]int),
			parent:     parent,
		}
	}

	detectorDatatype.files[fileName].lineNumber[lineNumber] = lineNumber

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
			if customDetector, isCustomDetector := holder.config.CustomDetector[detector.id]; isCustomDetector {
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
			for _, file := range files {

				lineNumbers := maputil.ToSortedSlice(file.lineNumber)
				for _, lineNumber := range lineNumbers {
					constructedDatatype.Locations = append(constructedDatatype.Locations, types.RiskLocation{
						Filename:   file.name,
						LineNumber: lineNumber,
						Parent:     file.parent,
					})
				}
			}
			constructedDetector.DataTypes = append(constructedDetector.DataTypes, constructedDatatype)
		}

		data = append(data, constructedDetector)
	}

	for _, presentRisk := range holder.presentRisks {
		data = append(data, presentRisk)
	}

	return data
}
