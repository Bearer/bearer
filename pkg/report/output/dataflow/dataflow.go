package dataflow

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/bearer/curio/pkg/commands/process/settings"
	"github.com/bearer/curio/pkg/report/customdetectors"
	"github.com/bearer/curio/pkg/report/detections"
	reportdetectors "github.com/bearer/curio/pkg/report/detectors"
	"github.com/bearer/curio/pkg/report/output/dataflow/components"
	"github.com/bearer/curio/pkg/report/output/dataflow/datatypes"
	"github.com/bearer/curio/pkg/report/output/dataflow/detectiondecoder"
	"github.com/bearer/curio/pkg/report/output/dataflow/risks"
	"github.com/bearer/curio/pkg/util/output"

	"github.com/bearer/curio/pkg/report/output/dataflow/types"
)

type DataFlow struct {
	Datatypes  []types.Datatype  `json:"data_types,omitempty" yaml:"data_types,omitempty"`
	Risks      []interface{}     `json:"risks,omitempty" yaml:"risks,omitempty"`
	Components []types.Component `json:"components" yaml:"components"`
}

var allowedDetections []detections.DetectionType = []detections.DetectionType{detections.TypeSchemaClassified, detections.TypeCustomClassified, detections.TypeDependencyClassified, detections.TypeInterfaceClassified, detections.TypeFrameworkClassified, detections.TypeCustomRisk}

func GetOutput(input []interface{}, config settings.Config, isInternal bool) (*DataFlow, error) {
	dataTypesHolder := datatypes.New(config, isInternal)
	risksHolder := risks.New(config, isInternal)
	componentsHolder := components.New(isInternal)

	extras, err := datatypes.NewExtras(input)
	if err != nil {
		return nil, err
	}
	railsExtras, err := datatypes.NewRailsExtras(input)
	if err != nil {
		return nil, err
	}

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

		// add full path to filename
		fullFilename := getFullFilename(config.Target, castDetection.Source.Filename)
		castDetection.Source.Filename = fullFilename

		switch detectionType {
		case detections.TypeSchemaClassified:
			var detectionExtras *datatypes.ExtraFields
			if castDetection.DetectorType == reportdetectors.DetectorSchemaRb {
				detectionExtras = railsExtras.Get(detection)
			}

			err = dataTypesHolder.AddSchema(castDetection, detectionExtras)
			if err != nil {
				return nil, err
			}
		case detections.TypeCustomRisk:
			risksHolder.AddRiskPresence(castDetection)
		case detections.TypeCustomClassified:
			ruleName := string(castDetection.DetectorType)
			customDetector, ok := config.CustomDetector[ruleName]
			if !ok {
				return nil, fmt.Errorf("there is a custom detector in report that is not in the config %s", ruleName)
			}

			switch customDetector.Type {
			case customdetectors.TypeVerfifier:
				continue
			case customdetectors.TypeRisk:
				err := risksHolder.AddSchema(castDetection)
				if err != nil {
					return nil, err
				}
			case customdetectors.TypeDatatype:
				var detectionExtras *datatypes.ExtraFields
				if castDetection.DetectorType == "detect_sql_create_public_table" {
					detectionExtras = extras.Get(detection)
				}

				err = dataTypesHolder.AddSchema(castDetection, detectionExtras)
				if err != nil {
					return nil, err
				}
			}

		case detections.TypeDependencyClassified:
			classifiedDetection, err := detectiondecoder.GetClassifiedDependency(detection)
			if err != nil {
				return nil, err
			}

			classifiedDetection.Source.Filename = fullFilename
			err = componentsHolder.AddDependency(classifiedDetection)
			if err != nil {
				return nil, err
			}
		case detections.TypeInterfaceClassified:
			classifiedDetection, err := detectiondecoder.GetClassifiedInterface(detection)
			if err != nil {
				return nil, err
			}

			classifiedDetection.Source.Filename = fullFilename
			err = componentsHolder.AddInterface(classifiedDetection)
			if err != nil {
				return nil, err
			}
		case detections.TypeFrameworkClassified:
			classifiedDetection, err := detectiondecoder.GetClassifiedFramework(detection)
			if err != nil {
				return nil, err
			}

			classifiedDetection.Source.Filename = fullFilename
			err = componentsHolder.AddFramework(classifiedDetection)
			if err != nil {
				return nil, err
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
	return dataflow, nil
}

func getFullFilename(path string, filename string) string {
	path = strings.TrimSuffix(path, "/")
	filename = strings.TrimPrefix(filename, "/")

	if filename == "." {
		return path
	}

	if path == "" || path == "." {
		return filename
	}

	return path + "/" + filename
}
