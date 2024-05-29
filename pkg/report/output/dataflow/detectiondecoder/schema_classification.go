package detectiondecoder

import (
	"bytes"
	"encoding/json"
	"fmt"

	schemaclassification "github.com/bearer/bearer/pkg/classification/schema"

	"github.com/bearer/bearer/pkg/report/detections"
	"github.com/bearer/bearer/pkg/report/schema"
)

func GetSchemaClassification(schema schema.Schema) (schemaclassification.Classification, error) {
	// decode classification
	var classification schemaclassification.Classification
	buf := bytes.NewBuffer(nil)
	err := json.NewEncoder(buf).Encode(schema.Classification)
	if err != nil {
		return schemaclassification.Classification{}, fmt.Errorf("expecting classification got %#v", schema.Classification)
	}
	err = json.NewDecoder(buf).Decode(&classification)
	if err != nil {
		return schemaclassification.Classification{}, fmt.Errorf("expecting classification got %#v", schema.Classification)
	}

	return classification, nil
}

func GetSchema(detection detections.Detection) (schema.Schema, error) {
	var value schema.Schema
	buf := bytes.NewBuffer(nil)
	err := json.NewEncoder(buf).Encode(detection.Value)
	if err != nil {
		return schema.Schema{}, fmt.Errorf("expect detection to have value of type schema %#v", detection.Value)
	}
	err = json.NewDecoder(buf).Decode(&value)
	if err != nil {
		return schema.Schema{}, fmt.Errorf("expect detection to have value of type schema %#v", detection.Value)
	}

	return value, nil
}
