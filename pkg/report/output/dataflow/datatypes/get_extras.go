package datatypes

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/bearer/curio/pkg/commands/process/settings"
	"github.com/bearer/curio/pkg/report/detections"
	"github.com/bearer/curio/pkg/report/detectors"
	"github.com/bearer/curio/pkg/report/output/dataflow/types"
	regohelper "github.com/bearer/curio/pkg/util/rego"
	"github.com/open-policy-agent/opa/rego"
)

type processorInput struct {
	AllDetections    []interface{} `json:"all_detections"`
	TargetDetections []interface{} `json:"target_detections"`
}

type ExtraFields struct {
	encrypted  *bool
	verifiedBy []types.DatatypeVerifiedBy
}

func getRailsTargetDetections(allDetections []interface{}) ([]interface{}, error) {
	var result []interface{}

	for _, detection := range allDetections {
		detectionMap, ok := detection.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("found detection in report which is not object")
		}

		detectionType, ok := detectionMap["type"].(string)
		if !ok {
			continue
		}

		if detections.DetectionType(detectionType) != detections.TypeSchemaClassified {
			continue
		}

		detectorType, ok := detectionMap["detector_type"].(string)
		if !ok {
			continue
		}

		if detectors.Type(detectorType) != detectors.DetectorSchemaRb {
			continue
		}

		result = append(result, detection)
	}

	return result, nil
}

func getEncryptedField(result rego.Vars, detection interface{}) (bool, error) {
	rawEncryptedFields, ok := result["encrypted"]
	if !ok {
		return false, errors.New("no 'encrypted' value in output")
	}

	encryptedFields, ok := rawEncryptedFields.([]interface{})
	if !ok {
		return false, errors.New("invalid type for 'encrypted' value")
	}

	detectionMap := detection.(map[string]interface{})
	detectionID := detectionMap["id"].(string)

	for _, rawResultDetection := range encryptedFields {
		resultDetection, ok := rawResultDetection.(map[string]interface{})
		if !ok {
			return false, errors.New("invalid type for 'encrypted' detection")
		}

		rawResultDetectionID, ok := resultDetection["id"]
		if !ok {
			return false, errors.New("missing id for 'encrypted' detection")
		}

		resultDetectionID, ok := rawResultDetectionID.(string)
		if !ok {
			return false, errors.New("invalid type for 'encrypted' detection id")
		}

		if resultDetectionID == detectionID {
			return true, nil
		}
	}

	return false, nil
}

func getVerifiedBy(result rego.Vars, detection interface{}) ([]types.DatatypeVerifiedBy, error) {
	rawVerifiedBy, ok := result["verified_by"]
	if !ok {
		return nil, errors.New("no 'verified_by' value in output")
	}

	verifiedBy, ok := rawVerifiedBy.([]interface{})
	if !ok {
		return nil, errors.New("invalid type for 'verified_by' value")
	}

	detectionMap := detection.(map[string]interface{})
	detectionID := detectionMap["id"].(string)

	for _, rawItem := range verifiedBy {
		item, ok := rawItem.([]interface{})
		if !ok {
			return nil, errors.New("invalid type for 'verified_by' item")
		}

		if len(item) != 2 {
			return nil, errors.New("invalid length for 'verified_by' item")
		}

		rawItemDetection := item[0]
		rawItemVerifiedBy := item[1]

		itemDetection, ok := rawItemDetection.(map[string]interface{})
		if !ok {
			return nil, errors.New("invalid type for 'verified_by' item detection")
		}

		rawItemDetectionID, ok := itemDetection["id"]
		if !ok {
			return nil, errors.New("missing id for 'verified_by' item detection")
		}

		itemDetectionID, ok := rawItemDetectionID.(string)
		if !ok {
			return nil, errors.New("invalid type for 'verified_by' item detection id")
		}

		if itemDetectionID != detectionID {
			continue
		}

		var verifiedBy []types.DatatypeVerifiedBy
		bytes, err := json.Marshal(rawItemVerifiedBy)
		if err != nil {
			return nil, fmt.Errorf("failed to serialize 'verified_by' item: %s", err)
		}
		err = json.Unmarshal(bytes, &verifiedBy)
		if err != nil {
			return nil, fmt.Errorf("invalid format for 'verified_by' item: %s", err)
		}

		return verifiedBy, nil
	}

	return nil, nil
}

type extrasObj struct {
	data map[string]*ExtraFields
}

func NewRailsExtras(detections []interface{}) (*extrasObj, error) {
	return newExtrasObj(detections, getRailsTargetDetections)
}

func NewExtras(detections []interface{}) (*extrasObj, error) {
	return newExtrasObj(detections, getTargetDetections)
}

func newExtrasObj(
	detections []interface{},
	targetDetectionsFunc func(detections []interface{}) ([]interface{}, error),
) (*extrasObj, error) {
	targetDetections, err := targetDetectionsFunc(detections)
	if err != nil {
		return nil, err
	}

	module, err := settings.EncryptedVerifiedRegoModuleText()
	if err != nil {
		return nil, err
	}

	data, err := runExtrasQuery(
		`
			verified_by = data.bearer.encrypted_verified.verified_by
			encrypted = data.bearer.encrypted_verified.encrypted
		`,
		[]regohelper.Module{{
			Name:    "bearer.encrypted_verified",
			Content: module,
		}},
		detections,
		targetDetections,
	)
	if err != nil {
		return nil, err
	}

	return &extrasObj{data: data}, nil
}

func runExtrasQuery(
	query string,
	modules []regohelper.Module,
	detections, targetDetections []interface{},
) (map[string]*ExtraFields, error) {
	data := make(map[string]*ExtraFields)

	result, err := regohelper.RunQuery(query, processorInput{
		AllDetections:    detections,
		TargetDetections: targetDetections,
	}, modules)
	if err != nil {
		return nil, err
	}

	for _, detection := range targetDetections {
		extras := &ExtraFields{}
		encrypted, err := getEncryptedField(result, detection)
		if err != nil {
			return nil, err
		}

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

		detectionMap := detection.(map[string]interface{})
		detectionID := detectionMap["id"].(string)
		data[detectionID] = extras
	}

	return data, nil
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

func (extras *extrasObj) Get(detection interface{}) *ExtraFields {
	detectionMap := detection.(map[string]interface{})
	detectionID := detectionMap["id"].(string)

	return extras.data[detectionID]
}
