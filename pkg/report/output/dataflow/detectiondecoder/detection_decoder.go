package detectiondecoder

import (
	"bytes"
	"encoding/json"
	"fmt"

	schemaclassification "github.com/bearer/curio/pkg/classification/schema"

	"github.com/bearer/curio/pkg/report/detections"
	"github.com/bearer/curio/pkg/report/schema"
)

func GetClassification(detection detections.Detection) (schemaclassification.Classification, error) {
	// decode schema
	var value schema.Schema
	buf := bytes.NewBuffer(nil)
	err := json.NewEncoder(buf).Encode(detection.Value)
	if err != nil {
		return schemaclassification.Classification{}, fmt.Errorf("expect detection to have value of type schema %#v", detection.Value)
	}
	err = json.NewDecoder(buf).Decode(&value)
	if err != nil {
		return schemaclassification.Classification{}, fmt.Errorf("expect detection to have value of type schema %#v", detection.Value)
	}

	// decode classification
	var classification schemaclassification.Classification
	buf = bytes.NewBuffer(nil)
	err = json.NewEncoder(buf).Encode(value.Classification)
	if err != nil {
		return schemaclassification.Classification{}, fmt.Errorf("expecting classification got %#v", value.Classification)
	}
	err = json.NewDecoder(buf).Decode(&classification)
	if err != nil {
		return schemaclassification.Classification{}, fmt.Errorf("expecting classification got %#v", value.Classification)
	}

	return classification, nil
}
