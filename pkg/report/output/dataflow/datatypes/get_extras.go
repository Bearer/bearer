package datatypes

import (
	"encoding/json"

	"github.com/bearer/curio/pkg/commands/process/settings"
	"github.com/bearer/curio/pkg/report/output/dataflow/types"
	regohelper "github.com/bearer/curio/pkg/util/rego"
	"github.com/open-policy-agent/opa/rego"
)

type processorInput struct {
	AllDetections []interface{} `json:"all_detections"`
	Target        interface{}   `json:"target"`
}

type extraFields struct {
	encrypted  *bool
	verifiedBy []types.DatatypeVerifiedBy
}

func GetRailsExtras(input []interface{}, detection interface{}) (*extraFields, error) {
	extras := &extraFields{}

	processorContent := `
package bearer.rails_encrypted_verified

import future.keywords

default encrypted := false

ruby_encrypted[location] {
		some detection in input.all_detections
		detection.detector_type == "detect_encrypted_ruby_class_properties"
		detection.value.classification.decision.state == "valid"
		location = detection
}

encrypted = true {
		some detection in ruby_encrypted
		detection.value.transformed_object_name == input.target.value.transformed_object_name
		detection.value.field_name == input.target.value.field_name
		input.target.value.field_name != ""
		input.target.value.object_name != ""
}

verified_by[verification] {
		some detection in ruby_encrypted
		detection.value.transformed_object_name == input.target.value.transformed_object_name
		detection.value.field_name == input.target.value.field_name

		verification = {
				"detector": "detect_encrypted_ruby_class_properties",
				"filename": detection.source.filename,
				"line_number": detection.source.line_number
		}
}
`

	query := `
verified_by = data.bearer.rails_encrypted_verified.verified_by
encrypted = data.bearer.rails_encrypted_verified.encrypted
`

	module := regohelper.Module{
		Name: "bearer.rails_encrypted_verified",
		Content: processorContent,
	}

	result, err := regohelper.RunQuery(query, processorInput{
		AllDetections: input,
		Target:        detection,
	}, []regohelper.Module{module})
	if err != nil {
		return nil, err
	}

	encrypted := getEncryptedField(result)

	if encrypted {
		extras.encrypted = &encrypted

		verified, err := getVerifiedBy(result)
		if err != nil {
			return nil, err
		}

		if verified != nil {
			extras.verifiedBy = append(extras.verifiedBy, verified...)
		}
	}

	return extras, nil
}

func GetExtras(customDetector settings.Rule, input []interface{}, detection interface{}) (*extraFields, error) {
	extras := &extraFields{}

	for _, processor := range customDetector.Processors {
		result, err := regohelper.RunQuery(processor.Query, processorInput{
			AllDetections: input,
			Target:        detection,
		}, processor.Modules.ToRegoModules())
		if err != nil {
			return nil, err
		}

		encrypted := getEncryptedField(result)

		if encrypted {
			extras.encrypted = &encrypted

			verified, err := getVerifiedBy(result)
			if err != nil {
				return nil, err
			}

			if verified != nil {
				extras.verifiedBy = append(extras.verifiedBy, verified...)
			}
		}
	}

	return extras, nil
}

func getEncryptedField(result rego.Vars) bool {
	encryptedField, hasEncryptedField := result["encrypted"]

	if hasEncryptedField {
		encrypted, ok := encryptedField.(bool)

		if ok && encrypted {
			return true
		}
	}

	return false
}

func getVerifiedBy(result rego.Vars) ([]types.DatatypeVerifiedBy, error) {
	verifiedByField, hasverifiedByField := result["verified_by"]

	if hasverifiedByField {
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
