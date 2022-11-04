package dataflow

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/bearer/curio/pkg/report/detections"
	"github.com/bearer/curio/pkg/report/output/dataflow/datatypes"
	"github.com/bearer/curio/pkg/report/output/dataflow/risks"

	"github.com/bearer/curio/pkg/report/output/dataflow/types"
)

type DataFlow struct {
	Datatypes []types.Datatype     `json:"data_types,omitempty"`
	Risks     []types.RiskDetector `json:"risks,omitempty"`
}

var allowedDetections []detections.DetectionType = []detections.DetectionType{detections.TypeSchema, detections.TypeSchemaClassified, detections.TypeCustom, detections.TypeCustomClassified}

func GetOuput(input []interface{}) (*DataFlow, error) {
	dataTypesHolder := datatypes.New()
	risksHolder := risks.New()

	for _, detection := range input {
		detection, ok := detection.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("found detection in report which is not object")
		}

		detectionTypeS, ok := detection["type"].(string)
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

		var castedDetection detections.Detection
		buf := bytes.NewBuffer(nil)
		err := json.NewEncoder(buf).Encode(detection)
		if err != nil {
			return nil, err
		}
		err = json.NewDecoder(buf).Decode(&castedDetection)
		if err != nil {
			return nil, err
		}

		if detectionType == detections.TypeSchema || detectionType == (detections.TypeSchemaClassified) {
			err = dataTypesHolder.AddSchema(castedDetection)
			if err != nil {
				return nil, err
			}
		}

		if detectionType == detections.TypeCustom || detectionType == (detections.TypeCustomClassified) {
			err := risksHolder.AddSchema(castedDetection)
			if err != nil {
				return nil, err
			}
		}

	}

	dataflow := &DataFlow{
		Datatypes: dataTypesHolder.ToDataFlow(),
		Risks:     risksHolder.ToDataFlow(),
	}

	return dataflow, nil
}
