package dataflow

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/bearer/curio/pkg/commands/process/settings"
	"github.com/bearer/curio/pkg/report/customdetectors"
	"github.com/bearer/curio/pkg/report/detections"
	"github.com/bearer/curio/pkg/report/output/dataflow/components"
	"github.com/bearer/curio/pkg/report/output/dataflow/datatypes"
	"github.com/bearer/curio/pkg/report/output/dataflow/risks"

	"github.com/bearer/curio/pkg/report/output/dataflow/types"
)

type DataFlow struct {
	Datatypes  []types.Datatype     `json:"data_types,omitempty" yaml:"data_types,omitempty"`
	Risks      []types.RiskDetector `json:"risks,omitempty" yaml:"risks,omitempty"`
	Components []types.Component    `json:"components" yaml:"components"`
}

var allowedDetections []detections.DetectionType = []detections.DetectionType{detections.TypeSchemaClassified, detections.TypeCustomClassified, detections.TypeDependencyClassified, detections.TypeInterfaceClassified, detections.TypeFrameworkClassified}

func GetOutput(input []interface{}, config settings.Config, isInternal bool) (*DataFlow, error) {
	dataTypesHolder := datatypes.New()
	risksHolder := risks.New(config, isInternal)
	componentsHolder := components.New(isInternal)

	for _, detection := range input {
		detectionMap, ok := detection.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("found detection in report which is not object")
		}

		detectionTypeS, ok := detectionMap["type"].(string)
		if !ok {
			continue
		}

		detectionType := detections.DetectionType(detectionTypeS)

		isDataflow := false
		for _, allowedDetection := range allowedDetections {
			if (detectionType) == allowedDetection {
				isDataflow = true
			}
		}

		if !isDataflow {
			continue
		}

		var castDetection detections.Detection
		buf := bytes.NewBuffer(nil)
		err := json.NewEncoder(buf).Encode(detection)
		if err != nil {
			return nil, err
		}
		err = json.NewDecoder(buf).Decode(&castDetection)
		if err != nil {
			return nil, err
		}

		switch detectionType {
		case detections.TypeSchemaClassified:
			err = dataTypesHolder.AddSchema(castDetection, nil)
			if err != nil {
				return nil, err
			}
		case detections.TypeCustomClassified:
			ruleName := string(castDetection.DetectorType)
			customDetector, ok := config.CustomDetector[ruleName]
			if !ok {
				return nil, fmt.Errorf("there is a custom detector in report that is not in the config %s", ruleName)
			}

			if customDetector.Verifier {
				continue
			}

			switch customDetector.Type {
			case customdetectors.TypeRisk:
				err := risksHolder.AddSchema(castDetection)
				if err != nil {
					return nil, err
				}
			case customdetectors.TypeDatatype:
				extras, err := datatypes.GetExtras(customDetector, input, detection)
				if err != nil {
					return nil, err
				}

				err = dataTypesHolder.AddSchema(castDetection, extras)
				if err != nil {
					return nil, err
				}
			}

		case detections.TypeDependencyClassified:
			err := componentsHolder.AddDependency(detection)
			if err != nil {
				return nil, err
			}
		case detections.TypeInterfaceClassified:
			err := componentsHolder.AddInterface(detection)
			if err != nil {
				return nil, err
			}
		case detections.TypeFrameworkClassified:
			err := componentsHolder.AddFramework(detection)
			if err != nil {
				return nil, err
			}
		}
	}

	dataflow := &DataFlow{
		Datatypes:  dataTypesHolder.ToDataFlow(),
		Risks:      risksHolder.ToDataFlow(),
		Components: componentsHolder.ToDataFlow(),
	}

	return dataflow, nil
}
