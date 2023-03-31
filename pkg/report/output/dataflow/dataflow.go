package dataflow

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/bearer/bearer/pkg/commands/process/settings"
	"github.com/bearer/bearer/pkg/report/customdetectors"
	"github.com/bearer/bearer/pkg/report/detections"
	reportdetectors "github.com/bearer/bearer/pkg/report/detectors"
	"github.com/bearer/bearer/pkg/report/output/dataflow/components"
	"github.com/bearer/bearer/pkg/report/output/dataflow/datatypes"
	"github.com/bearer/bearer/pkg/report/output/dataflow/detectiondecoder"
	"github.com/bearer/bearer/pkg/report/output/dataflow/risks"
	"github.com/bearer/bearer/pkg/util/file"
	"github.com/bearer/bearer/pkg/util/output"
	"github.com/hhatto/gocloc"

	"github.com/bearer/bearer/pkg/report/output/dataflow/types"
)

type DataFlow struct {
	Datatypes  []types.Datatype     `json:"data_types,omitempty" yaml:"data_types,omitempty"`
	Risks      []types.RiskDetector `json:"risks,omitempty" yaml:"risks,omitempty"`
	Components []types.Component    `json:"components" yaml:"components"`
}

var allowedDetections []detections.DetectionType = []detections.DetectionType{
	detections.TypeSchemaClassified,
	detections.TypeCustomClassified,
	detections.TypeDependencyClassified,
	detections.TypeInterfaceClassified,
	detections.TypeFrameworkClassified,
	detections.TypeCustomRisk,
	detections.TypeSecretleak,
}

func contains(detections []detections.DetectionType, detection detections.DetectionType) bool {
	for _, v := range detections {
		if v == detection {
			return true
		}
	}

	return false
}

func GetOutput(input []interface{}, config settings.Config, isInternal bool) (*DataFlow, *gocloc.Result, *DataFlow, error) {
	dataTypesHolder := datatypes.New(config, isInternal)
	risksHolder := risks.New(config, isInternal)
	componentsHolder := components.New(isInternal)

	extras, err := datatypes.NewExtras(input, config)
	if err != nil {
		return nil, nil, nil, err
	}

	customExtras, err := datatypes.NewCustomExtras(input, config)
	if err != nil {
		return nil, nil, nil, err
	}

	for _, detection := range input {
		detectionMap, ok := detection.(map[string]interface{})
		if !ok {
			return nil, nil, nil, fmt.Errorf("found detection in report which is not object")
		}

		detectionTypeS, ok := detectionMap["type"].(string)

		if !ok {
			continue
		}

		detectionType := detections.DetectionType(detectionTypeS)

		isDataflow := contains(allowedDetections, detectionType)
		if !isDataflow {
			continue
		}

		var castDetection detections.Detection
		buf := bytes.NewBuffer(nil)
		err := json.NewEncoder(buf).Encode(detection)
		if err != nil {
			return nil, nil, nil, err
		}
		err = json.NewDecoder(buf).Decode(&castDetection)
		if err != nil {
			return nil, nil, nil, err
		}

		// add full path to filename
		fullFilename := file.GetFullFilename(config.Target, castDetection.Source.Filename)
		castDetection.Source.Filename = fullFilename

		switch detectionType {
		case detections.TypeSchemaClassified:
			var detectionExtras *datatypes.ExtraFields
			if castDetection.DetectorType == reportdetectors.DetectorSchemaRb {
				detectionExtras = customExtras.Get(detection)
			}

			err = dataTypesHolder.AddSchema(castDetection, detectionExtras)
			if err != nil {
				return nil, nil, nil, err
			}
		case detections.TypeCustomRisk:
			risksHolder.AddRiskPresence(castDetection)
		case detections.TypeCustomClassified:
			ruleName := string(castDetection.DetectorType)
			customDetector, ok := config.Rules[ruleName]
			if !ok {
				customDetector, ok = config.BuiltInRules[ruleName]
				if !ok {
					return nil, nil, nil, fmt.Errorf("custom detector not in config %s", ruleName)
				}
			}

			switch customDetector.Type {
			case customdetectors.TypeVerifier:
				continue
			case customdetectors.TypeRisk:
				err := risksHolder.AddSchema(castDetection)
				if err != nil {
					return nil, nil, nil, err
				}
			case customdetectors.TypeDatatype:
				var detectionExtras *datatypes.ExtraFields
				detectionExtras = extras.Get(detection)

				err = dataTypesHolder.AddSchema(castDetection, detectionExtras)
				if err != nil {
					return nil, nil, nil, err
				}
			}

		case detections.TypeSecretleak:
			risksHolder.AddRiskPresence(castDetection)

		case detections.TypeDependencyClassified:
			classifiedDetection, err := detectiondecoder.GetClassifiedDependency(detection)
			if err != nil {
				return nil, nil, nil, err
			}

			classifiedDetection.Source.Filename = fullFilename
			err = componentsHolder.AddDependency(classifiedDetection)
			if err != nil {
				return nil, nil, nil, err
			}
		case detections.TypeInterfaceClassified:
			classifiedDetection, err := detectiondecoder.GetClassifiedInterface(detection)
			if err != nil {
				return nil, nil, nil, err
			}

			classifiedDetection.Source.Filename = fullFilename
			err = componentsHolder.AddInterface(classifiedDetection)
			if err != nil {
				return nil, nil, nil, err
			}
		case detections.TypeFrameworkClassified:
			classifiedDetection, err := detectiondecoder.GetClassifiedFramework(detection)
			if err != nil {
				return nil, nil, nil, err
			}
			classifiedDetection.Source.Filename = fullFilename
			err = componentsHolder.AddFramework(classifiedDetection)
			if err != nil {
				return nil, nil, nil, err
			}
		}
	}

	if !config.Scan.Quiet {
		output.StdErrLogger().Msgf("Generating dataflow")
	}
	dataflow := &DataFlow{
		Datatypes:  dataTypesHolder.ToDataFlow(),
		Risks:      risksHolder.ToDataFlow(),
		Components: componentsHolder.ToDataFlow(),
	}
	return dataflow, nil, nil, nil
}
