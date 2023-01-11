package testhelper

import (
	"sort"

	"github.com/bearer/curio/pkg/classification/schema"
	"github.com/bearer/curio/pkg/report/schema/datatype"
	"github.com/bearer/curio/pkg/util/classify"

	"github.com/stretchr/testify/assert"
)

type KPI struct {
	DetectionsCount                    int
	ValidDetectionsCount               int
	ValidObjectDetectionsCount         int
	ValidFieldDetectionsCount          int
	ExpectedValidDetectionsCount       int // TODO: remove expected counts from KPIs
	ExpectedValidObjectDetectionsCount int
	ExpectedValidFieldDetectionsCount  int
	ExpectedTruePositivesCount         int
	ExpectedFalsePositivesCount        int
}

type ClassificationResult struct {
	Name       string
	Decision   classify.ClassificationDecision
	Properties map[string]ClassificationResult
}

type Output struct {
	KPI                     KPI
	ValidClassifications    []ClassificationResult
	ExpectedClassifications []ClassificationResult
}

type InputProperties struct {
	Name          string `json:"name" yaml:"name"`
	Type          string `json:"type" yaml:"type"`
	State         string `json:"state" yaml:"state"`
	Reason        string `json:"reason" yaml:"reason"`
	FalsePositive bool   `json:"false_positive" yaml:"false_positive"`
}

type Input struct {
	Name          string            `json:"name" yaml:"name"`
	Filename      string            `json:"filename" yaml:"filename"`
	DetectorType  string            `json:"detector_type" yaml:"detector_type"`
	Properties    []InputProperties `json:"properties" yaml:"properties"`
	State         string            `json:"state" yaml:"state"`
	Reason        string            `json:"reason" yaml:"reason"`
	FalsePositive bool              `json:"false_positive" yaml:"false_positive"`
}

func ExtractExpectedOutput(
	t assert.TestingT,
	lang string,
	classifier *schema.Classifier,
) Output {

	result := Output{
		KPI: KPI{
			DetectionsCount: 0,
		},
		ValidClassifications:    []ClassificationResult{},
		ExpectedClassifications: []ClassificationResult{},
	}

	// val, err := ioutil.ReadFile("././fixtures/" + lang + ".json")
	// if err != nil {
	// 	t.Errorf("error opening file %e", err)
	// }

	// var input []Input
	// rawBytes := []byte(val)
	// err = json.Unmarshal(rawBytes, &input)
	// if err != nil {
	// 	t.Errorf("error unmarshalling JSON %e", err)
	// }

	// for _, inputItem := range input {
	// 	includeResult := false
	// 	result.KPI.DetectionsCount += 1

	// 	expectedClassification := ClassificationResult{
	// 		Name: inputItem.Name,
	// 		Decision: classify.ClassificationDecision{
	// 			State:  classify.ValidationState(inputItem.State),
	// 			Reason: inputItem.Reason,
	// 		},
	// 	}

	// 	expectedProperties := map[string]ClassificationResult{}
	// 	detectionProperties := map[string]datatype.DataTypable{}
	// 	for _, inputItemProperty := range inputItem.Properties {
	// 		result.KPI.DetectionsCount += 1

	// 		if inputItemProperty.State == "valid" {
	// 			result.KPI.ExpectedValidDetectionsCount += 1
	// 			result.KPI.ExpectedValidFieldDetectionsCount += 1
	// 			if inputItemProperty.FalsePositive {
	// 				result.KPI.ExpectedFalsePositivesCount += 1
	// 			} else {
	// 				result.KPI.ExpectedTruePositivesCount += 1
	// 			}

	// 			expectedProperties[inputItemProperty.Name] = ClassificationResult{
	// 				Name: inputItemProperty.Name,
	// 				Decision: classify.ClassificationDecision{
	// 					State:  classify.Valid,
	// 					Reason: inputItemProperty.Reason,
	// 				},
	// 			}
	// 		} else {
	// 			expectedProperties[inputItemProperty.Name] = ClassificationResult{
	// 				Name: inputItemProperty.Name,
	// 				Decision: classify.ClassificationDecision{
	// 					State:  classify.Invalid,
	// 					Reason: inputItemProperty.Reason,
	// 				},
	// 			}
	// 		}

	// 		detectionProperties[inputItemProperty.Name] = &datatype.DataType{
	// 			Name:       inputItemProperty.Name,
	// 			Type:       inputItemProperty.Type,
	// 			Properties: map[string]datatype.DataTypable{},
	// 		}
	// 	}

	// 	expectedClassification.Properties = expectedProperties

	// 	if inputItem.State == "valid" {
	// 		result.KPI.ExpectedValidDetectionsCount += 1
	// 		result.KPI.ExpectedValidObjectDetectionsCount += 1
	// 		// include as expected classification
	// 		result.ExpectedClassifications = append(result.ExpectedClassifications, expectedClassification)

	// 		if inputItem.FalsePositive {
	// 			result.KPI.ExpectedFalsePositivesCount += 1
	// 		} else {
	// 			result.KPI.ExpectedTruePositivesCount += 1
	// 		}
	// 	}

	// 	detection := schema.DataTypeDetection{
	// 		Filename:     inputItem.Filename,
	// 		DetectorType: detectors.Type(inputItem.DetectorType),
	// 		Value: &datatype.DataType{
	// 			Name:       inputItem.Name,
	// 			Properties: detectionProperties,
	// 		},
	// 	}

	// 	classification := classifier.Classify(detection)

	// 	if classification.Classification.Decision.State == classify.Valid {
	// 		includeResult = true
	// 		result.KPI.ValidObjectDetectionsCount += 1
	// 	}
	// 	classificationResult := ClassificationResult{
	// 		Name: classification.GetName(),
	// 		Decision: classify.ClassificationDecision{
	// 			State:  classification.Classification.Decision.State,
	// 			Reason: classification.Classification.Decision.Reason,
	// 		},
	// 	}

	// 	classifiedProperties := map[string]ClassificationResult{}
	// 	// sort properties to ensure consistency for snapshot
	// 	fields := classification.DataTypable.GetProperties()
	// 	for _, key := range sortKeys(fields) {
	// 		field := fields[key]
	// 		fieldClassification, _ := field.GetClassification().(schema.Classification)
	// 		// TODO: casting does not work if classification is empty
	// 		if fieldClassification.Decision.State == "valid" {
	// 			includeResult = true
	// 			result.KPI.ValidFieldDetectionsCount += 1
	// 		}

	// 		classifiedProperties[field.GetName()] = ClassificationResult{
	// 			Name: field.GetName(),
	// 			Decision: classify.ClassificationDecision{
	// 				State:  fieldClassification.Decision.State,
	// 				Reason: fieldClassification.Decision.Reason,
	// 			},
	// 		}
	// 	}

	// 	classificationResult.Properties = classifiedProperties

	// 	if includeResult {
	// 		result.ValidClassifications = append(result.ValidClassifications, classificationResult)
	// 	}
	// }

	// result.KPI.ValidDetectionsCount = result.KPI.ValidObjectDetectionsCount + result.KPI.ValidFieldDetectionsCount

	return result
}

func sortKeys(fields map[string]datatype.DataTypable) []string {
	keys := make([]string, 0, len(fields))
	for k := range fields {
		keys = append(keys, k)
	}

	sort.Strings(keys)
	return keys
}
