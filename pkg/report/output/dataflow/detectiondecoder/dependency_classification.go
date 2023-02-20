package detectiondecoder

import (
	"bytes"
	"encoding/json"
	"fmt"

	dependenciesclassification "github.com/bearer/bearer/pkg/classification/dependencies"
)

func GetClassifiedDependency(detection interface{}) (dependenciesclassification.ClassifiedDependency, error) {
	var value dependenciesclassification.ClassifiedDependency
	buf := bytes.NewBuffer(nil)
	err := json.NewEncoder(buf).Encode(detection)
	if err != nil {
		return dependenciesclassification.ClassifiedDependency{}, fmt.Errorf("expect detection to have value of type dependency %#v", detection)
	}
	err = json.NewDecoder(buf).Decode(&value)
	if err != nil {
		return dependenciesclassification.ClassifiedDependency{}, fmt.Errorf("expect detection to have value of type dependency %#v", detection)
	}

	return value, nil
}
