package dataflow

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/bearer/curio/pkg/report/detections"
	"github.com/bearer/curio/pkg/report/output/dataflow/datatypes"
	"github.com/bearer/curio/pkg/report/output/dataflow/types"
)

type DataFlow struct {
	Datatypes []types.Datatype `json:"data_types,omitempty"`
	Risks     []types.Datatype `json:"risks,omitempty"`
}

var allowedDetections []detections.DetectionType = []detections.DetectionType{detections.TypeSchema, detections.TypeCustom, detections.TypeSchemaClassified}

func GetOuput(input []interface{}) (*DataFlow, error) {
	dataTypesHolder := datatypes.New()

	for _, detection := range input {
		detection, ok := detection.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("found detection in report which is not object")
		}

		detectionType, ok := detection["type"].(string)

		isDataflow := false
		for _, allowedDetection := range allowedDetections {
			if detections.DetectionType(detectionType) == allowedDetection {
				isDataflow = true
			}
		}

		if !ok || !isDataflow {
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

		err = dataTypesHolder.AddSchema(castedDetection)
		if err != nil {
			return nil, err
		}
	}

	dataflow := &DataFlow{
		Datatypes: dataTypesHolder.ToDataFlow(),
	}

	return dataflow, nil
}
