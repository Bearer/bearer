package datatypes

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/bearer/bearer/internal/commands/process/settings"
	"github.com/bearer/bearer/internal/report/detections"
	"github.com/bearer/bearer/internal/report/detectors"
	"github.com/bearer/bearer/internal/report/output/dataflow/types"
	regohelper "github.com/bearer/bearer/internal/util/rego"
	"github.com/open-policy-agent/opa/rego"
)

type processorInput struct {
	Rule             *settings.Rule `json:"rule"`
	AllDetections    []interface{}  `json:"all_detections"`
	TargetDetections []interface{}  `json:"target_detections"`
}

type ExtraFields struct {
	encrypted  *bool
	verifiedBy []types.DatatypeVerifiedBy
}

func getCustomTargetDetections(allDetections []interface{}) ([]interface{}, error) {
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

func NewCustomExtras(detections []interface{}, config settings.Config) (*extrasObj, error) {
	return newExtrasObj(detections, getCustomTargetDetections, config)
}

func NewExtras(detections []interface{}, config settings.Config) (*extrasObj, error) {
	return newExtrasObj(detections, getTargetDetections, config)
}

func newExtrasObj(
	detections []interface{},
	targetDetectionsFunc func(detections []interface{}) ([]interface{}, error),
	config settings.Config,
) (*extrasObj, error) {
	targetDetections, err := targetDetectionsFunc(detections)
	if err != nil {
		return nil, err
	}

	data := make(map[string]*ExtraFields)

	for _, rule := range config.Rules {
		for _, processor := range rule.Processors {
			dataForProcessor, err := runProcessor(
				processor,
				detections,
				targetDetections,
				rule,
			)
			if err != nil {
				return nil, err
			}
			for k, v := range dataForProcessor {
				existingExtraFields, keyPresent := data[k]
				if keyPresent {
					// Merge in the new processor data
					if existingExtraFields.encrypted == nil {
						data[k].encrypted = v.encrypted
					}
					data[k].verifiedBy = append(data[k].verifiedBy, v.verifiedBy...)
				} else {
					data[k] = v
				}
			}
		}
	}
	return &extrasObj{data: data}, nil
}

func runExtrasQuery(
	query string,
	modules []regohelper.Module,
	detections, targetDetections []interface{},
	rule *settings.Rule,
) (map[string]*ExtraFields, error) {
	data := make(map[string]*ExtraFields)

	result, err := regohelper.RunQuery(query, processorInput{
		Rule:             rule,
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

func processorModules(processorName string) (modules []regohelper.Module, err error) {
	moduleText, err := settings.ProcessorRegoModuleText(processorName)
	if err != nil {
		return
	}

	fullModuleName := fmt.Sprintf("bearer.%s", processorName)
	modules = []regohelper.Module{{
		Name:    fullModuleName,
		Content: moduleText,
	}}

	return
}

func runProcessor(
	processorName string,
	detections []any,
	targetDetections []any,
	rule *settings.Rule,
) (data map[string]*ExtraFields, err error) {
	modules, err := processorModules(processorName)
	if err != nil {
		return
	}

	query := fmt.Sprintf(`
			verified_by = data.bearer.%s.verified_by
			encrypted = data.bearer.%s.encrypted
		`, processorName, processorName)

	data, err = runExtrasQuery(
		query,
		modules,
		detections,
		targetDetections,
		rule,
	)

	return
}
