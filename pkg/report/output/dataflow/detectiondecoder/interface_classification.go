package detectiondecoder

import (
	"bytes"
	"encoding/json"
	"fmt"

	interfaceclassification "github.com/bearer/curio/pkg/classification/interfaces"
)

func GetClassifiedInterface(detection interface{}) (interfaceclassification.ClassifiedInterface, error) {
	var value interfaceclassification.ClassifiedInterface
	buf := bytes.NewBuffer(nil)
	err := json.NewEncoder(buf).Encode(detection)
	if err != nil {
		return interfaceclassification.ClassifiedInterface{}, fmt.Errorf("expect detection to have value of type schema %#v", detection)
	}
	err = json.NewDecoder(buf).Decode(&value)
	if err != nil {
		return interfaceclassification.ClassifiedInterface{}, fmt.Errorf("expect detection to have value of type schema %#v", detection)
	}

	return value, nil
}
