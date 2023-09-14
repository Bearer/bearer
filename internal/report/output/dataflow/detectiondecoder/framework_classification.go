package detectiondecoder

import (
	"bytes"
	"encoding/json"
	"fmt"

	frameworkclassification "github.com/bearer/bearer/internal/classification/frameworks"
)

func GetClassifiedFramework(detection interface{}) (frameworkclassification.ClassifiedFramework, error) {
	var value frameworkclassification.ClassifiedFramework
	buf := bytes.NewBuffer(nil)
	err := json.NewEncoder(buf).Encode(detection)
	if err != nil {
		return frameworkclassification.ClassifiedFramework{}, fmt.Errorf("expect detection to have value of type framework %#v", detection)
	}
	err = json.NewDecoder(buf).Decode(&value)
	if err != nil {
		return frameworkclassification.ClassifiedFramework{}, fmt.Errorf("expect detection to have value of type framework %#v", detection)
	}

	return value, nil
}
