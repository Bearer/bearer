package dataflow

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/bearer/curio/pkg/report/detections"
)

type DataFlow struct {
	Datatypes []Datatype `json:"data_types,omitempty"`
	Risks     []Datatype `json:"risks,omitempty"`
}

type Datatype struct {
	Name      string     `json:"name"`
	Detectors []Detector `json:"detectors"`
}

type Detector struct {
	Name      string     `json:"name"`
	Stored    bool       `json:"stored"`
	Locations []Location `json:"locations"`
}

type Location struct {
	Filename   string `json:"filename"`
	LineNumber int    `json:"line_number"`
}

var allowedDetections []detections.DetectionType = []detections.DetectionType{detections.TypeSchema, detections.TypeCustom, detections.TypeSchemaClassified}

func GetOuput(input []interface{}) (*DataFlow, error) {
	holder := dataFlowHolder{
		datatypes: make(map[string]*datatypeHolder),
	}

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

		err = holder.addSchema(castedDetection)
		if err != nil {
			return nil, err
		}
	}

	return holder.toDataFlow(), nil
}

type dataFlowHolder struct {
	datatypes map[string]*datatypeHolder // group datatypeHolders by name
}
