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
		json.Unmarshal(bytes, &verifiedBy)

		return verifiedBy, nil
	}

	return nil, nil
}
