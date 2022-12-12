package datatypes

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/bearer/curio/pkg/commands/process/settings"
	"github.com/bearer/curio/pkg/report/customdetectors"
	"github.com/bearer/curio/pkg/report/detections"
	"github.com/bearer/curio/pkg/report/output/dataflow/types"
	regohelper "github.com/bearer/curio/pkg/util/rego"
	"github.com/open-policy-agent/opa/rego"
)

type processorInput struct {
	AllDetections    []interface{} `json:"all_detections"`
	TargetDetections []interface{} `json:"target_detections"`
}

type extraFields struct {
	encrypted  *bool
	verifiedBy []types.DatatypeVerifiedBy
}

// func GetRailsExtras(input []interface{}, detection interface{}) (*extraFields, error) {
// 	extras := &extraFields{}

// 	processorContent := `
// package bearer.rails_encrypted_verified

// import future.keywords

// default encrypted := false

// ruby_encrypted[location] {
// 		some detection in input.all_detections
// 		detection.detector_type == "detect_encrypted_ruby_class_properties"
// 		detection.value.classification.decision.state == "valid"
// 		location = detection
// }

// encrypted = true {
// 		some detection in ruby_encrypted
// 		detection.value.transformed_object_name == input.target.value.transformed_object_name
// 		detection.value.field_name == input.target.value.field_name
// 		input.target.value.field_name != ""
// 		input.target.value.object_name != ""
// }

// verified_by[verification] {
// 		some detection in ruby_encrypted
// 		detection.value.transformed_object_name == input.target.value.transformed_object_name
// 		detection.value.field_name == input.target.value.field_name

// 		verification = {
// 				"detector": "detect_encrypted_ruby_class_properties",
// 				"filename": detection.source.filename,
// 				"line_number": detection.source.line_number
// 		}
// }
// `

// 	query := `
// verified_by = data.bearer.rails_encrypted_verified.verified_by
// encrypted = data.bearer.rails_encrypted_verified.encrypted
// `

// 	module := regohelper.Module{
// 		Name:    "bearer.rails_encrypted_verified",
// 		Content: processorContent,
// 	}

// 	result, err := regohelper.RunQuery(query, processorInput{
// 		AllDetections: input,
// 		Target:        detection,
// 	}, []regohelper.Module{module})
// 	if err != nil {
// 		return nil, err
// 	}

// 	encrypted := getEncryptedField(result)

// 	if encrypted {
// 		extras.encrypted = &encrypted

// 		verified, err := getVerifiedBy(result)
// 		if err != nil {
// 			return nil, err
// 		}

// 		if verified != nil {
// 			extras.verifiedBy = append(extras.verifiedBy, verified...)
// 		}
// 	}

// 	return extras, nil
// }

// func GetExtras(customDetector settings.Rule, input []interface{}, detection interface{}) (*extraFields, error) {
// 	extras := &extraFields{}

// 	for _, processor := range customDetector.Processors {
// 		result, err := regohelper.RunQuery(processor.Query, processorInput{
// 			AllDetections: input,
// 			Target:        detection,
// 		}, processor.Modules.ToRegoModules())
// 		if err != nil {
// 			return nil, err
// 		}

// 		encrypted := getEncryptedField(result)

// 		if encrypted {
// 			extras.encrypted = &encrypted

// 			verified, err := getVerifiedBy(result)
// 			if err != nil {
// 				return nil, err
// 			}

// 			if verified != nil {
// 				extras.verifiedBy = append(extras.verifiedBy, verified...)
// 			}
// 		}
// 	}

// 	return extras, nil
// }

func getEncryptedField(result rego.Vars, detection interface{}) bool {
	rawEncryptedFields, hasEncrypted := result["encrypted"]

	if !hasEncrypted {
		return false
	}

	encryptedFields, ok := rawEncryptedFields.(map[interface{}]interface{})
	if !ok {
		return false
	}

	encryptedField, hasEncryptedField := encryptedFields[detection]

	if hasEncryptedField {
		encrypted, ok := encryptedField.(bool)

		if ok && encrypted {
			return true
		}
	}

	return false
}

func getVerifiedBy(result rego.Vars, detection interface{}) ([]types.DatatypeVerifiedBy, error) {
	rawVerifiedByFields, hasVerifiedBy := result["verified_by"]

	if !hasVerifiedBy {
		return nil, nil
	}

	verifiedByFields, ok := rawVerifiedByFields.(map[interface{}]interface{})
	if !ok {
		return nil, nil
	}

	verifiedByField, hasVerifiedByField := verifiedByFields[detection]

	if hasVerifiedByField {
		var verifiedBy []types.DatatypeVerifiedBy
		bytes, err := json.Marshal(verifiedByField)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(bytes, &verifiedBy)
		if err != nil {
			return nil, err
		}

		return verifiedBy, nil
	}

	return nil, nil
}

type extrasObj struct {
	data map[string][]extrasResult
}

type extrasResult struct {
	detection interface{}
	extras    *extraFields
}

func NewExtras(customRules map[string]settings.Rule, detections []interface{}) (*extrasObj, error) {
	targetDetections, err := getTargetDetections(detections)
	if err != nil {
		return nil, err
	}

	data := make(map[string][]extrasResult)

	for ruleName, rule := range customRules {
		if rule.Type != customdetectors.TypeDatatype {
			continue
		}

		if len(rule.Processors) == 0 {
			continue
		}

		processor := rule.Processors[0]

		result, err := regohelper.RunQuery(processor.Query, processorInput{
			AllDetections:    detections,
			TargetDetections: targetDetections,
		}, processor.Modules.ToRegoModules())
		if err != nil {
			return nil, err
		}

		for _, detection := range targetDetections {
			extras := &extraFields{}
			encrypted := getEncryptedField(result, detection)

			if encrypted {
				extras.encrypted = &encrypted

				verified, err := getVerifiedBy(result, detection)
				if err != nil {
					return nil, err
				}

				if verified != nil {
					extras.verifiedBy = append(extras.verifiedBy, verified...)
				}
			}

			data[ruleName] = append(data[ruleName], extrasResult{
				detection: detection,
				extras:    extras,
			})
		}
	}

	return &extrasObj{data: data}, nil
}

func getTargetDetections(allDetections []interface{}) ([]interface{}, error) {
	var result []interface{}

	for _, detection := range allDetections {
		detectionMap, ok := detection.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("found detection in report which is not object")
		}

		detectionTypeS, ok := detectionMap["type"].(string)
		if !ok {
			continue
		}

		detectionType := detections.DetectionType(detectionTypeS)
		if detectionType != detections.TypeCustomClassified {
			continue
		}

		result = append(result, detection)
	}

	return result, nil
}

func (extras *extrasObj) Get(ruleName string, detection interface{}) *extraFields {
	ruleExtras, ok := extras.data[ruleName]
	if !ok {
		return nil
	}

	for _, result := range ruleExtras {
		if reflect.DeepEqual(result.detection, detection) {
			return result.extras
		}
	}

	return nil
}
