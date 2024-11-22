package detectiondecoder

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/bearer/bearer/pkg/report/operations/operationshelper"
)

func GetOperation(detection interface{}) (operationshelper.Operation, error) {
	var value operationshelper.Operation
	buf := bytes.NewBuffer(nil)
	err := json.NewEncoder(buf).Encode(detection)
	if err != nil {
		return operationshelper.Operation{}, fmt.Errorf("expect detection to have value of type operation %#v", detection)
	}
	err = json.NewDecoder(buf).Decode(&value)
	if err != nil {
		return operationshelper.Operation{}, fmt.Errorf("expect detection to have value of type operation %#v", detection)
	}

	return value, nil
}
