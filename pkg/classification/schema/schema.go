package schema

import (
	"fmt"
	"regexp"

	"github.com/bearer/curio/pkg/report/schema/datatype"

	"github.com/bearer/curio/pkg/classification/db"
	"github.com/bearer/curio/pkg/report/detectors"
	"github.com/bearer/curio/pkg/util/classify"
	"github.com/bearer/curio/pkg/util/normalize_key"
)

var regexpIdentifierMatcher = regexp.MustCompile(`(uu)?id\z`)
var objectStopWords = map[string]struct{}{
	"this":        {},
	"props":       {},
	"prop types":  {},
	"exports":     {},
	"export":      {},
	"env":         {},
	"argv":        {},
	"arguments":   {},
	"errors":      {},
	"args":        {},
	"state":       {},
	"filter":      {},
	"memberships": {},
}
var databaseDetectorTypes = map[string]struct{}{
	"sql":   {},
	"rails": {},
}
var expectedIdentifierDataTypeIds = map[string]struct{}{
	"132": {}, // Unique Identifier
	"13":  {}, // Device Identifier
}

type ClassifiedDatatype struct {
	*datatype.DataType
	Classification Classification `json:"classification"`
}

type Classification struct {
	Name     string
	Decision classify.ClassificationDecision `json:"decision"`
}

type Classifier struct {
	config Config
}

type Config struct {
	DataTypes                      []db.DataType
	DataTypeClassificationPatterns []db.DataTypeClassificationPattern
	KnownPersonObjectPatterns      []db.KnownPersonObjectPattern
}

func New(config Config) *Classifier {
	return &Classifier{config: config}
}

type DataTypeDetection struct {
	Value        datatype.DataTypable
	Filename     string
	DetectorType detectors.Type
}

func extractDataType(value datatype.DataTypable) *datatype.DataType {
	return &datatype.DataType{
		Node:       value.GetNode(),
		Name:       value.GetName(),
		Type:       value.GetType(),
		TextType:   value.GetTextType(),
		Properties: value.GetProperties(),
		IsHelper:   value.GetIsHelper(),
		UUID:       value.GetUUID(),
	}
}

func (classifier *Classifier) Classify(data DataTypeDetection) (*ClassifiedDatatype, error) {
	detectedDataType := extractDataType(data.Value)
	normalizedObjectName := normalize_key.Normalize(detectedDataType.Name)
	classifiedDataType := &ClassifiedDatatype{
		DataType:       detectedDataType,
		Classification: Classification{Name: normalizedObjectName},
	}

	if classify.IsVendored(data.Filename) {
		classifiedDataType.Classification.Decision = classify.ClassificationDecision{
			State:  classify.Invalid,
			Reason: classify.IncludedInVendorFolderReason,
		}
	}

	if classify.IsPotentialDetector(data.DetectorType) {
		classifiedDataType.Classification.Decision = classify.ClassificationDecision{
			State:  classify.Invalid,
			Reason: classify.PotentialDetectorReason,
		}
	}

	// stop_word
	if objectStopWordDetected(normalizedObjectName) {
		classifiedDataType.Classification.Decision = classify.ClassificationDecision{
			State:  classify.Invalid,
			Reason: "stop_word",
		}
	}

	// todo: do we need these checks too? app/services/detection_engine/classify_detections/classify_schema_object_detection.rb

	if classifiedDataType.Classification.Decision.State == classify.Invalid {
		// schema object did not pass initial checks
		// mark all first level children as invalid

		// todo: handle children that are themselves schema objects
		for _, property := range detectedDataType.Properties {
			detectedDataType.Properties[property.GetName()] = classifyAsInvalid(property)
		}

		return classifiedDataType, nil
	}

	if classifier.isKnownPersonObject(normalizedObjectName) {
		hasKnownObjectProperties := false
		hasKnownDBIdentifierProperties := false

		// todo: handle children that are themselves schema objects
		for _, property := range detectedDataType.Properties {
			propertyDataType := extractDataType(property)
			normalizedPropertyName := normalize_key.Normalize(propertyDataType.Name)

			if classifier.isKnownObjectPattern(normalizedPropertyName, propertyDataType.Type) {
				hasKnownObjectProperties = true
				detectedDataType.Properties[propertyDataType.Name] = ClassifiedDatatype{
					DataType: propertyDataType,
					Classification: Classification{
						Name: normalizedPropertyName,
						Decision: classify.ClassificationDecision{
							State:  classify.Valid,
							Reason: "known_classification_pattern",
						},
					},
				}

				continue
			}

			if classifier.isKnownDBIdentifierPattern(normalizedPropertyName) {
				hasKnownDBIdentifierProperties = true
				detectedDataType.Properties[propertyDataType.Name] = ClassifiedDatatype{
					DataType: propertyDataType,
					Classification: Classification{
						Name: normalizedPropertyName,
						Decision: classify.ClassificationDecision{
							State:  classify.Valid,
							Reason: "known_database_identifier",
						},
					},
				}

				continue
			}

			// todo: check for field stop word?

			detectedDataType.Properties[propertyDataType.Name] = ClassifiedDatatype{
				DataType: propertyDataType,
				Classification: Classification{
					Name: normalizedPropertyName,
					Decision: classify.ClassificationDecision{
						State:  classify.Invalid,
						Reason: "invalid_property",
					},
				},
			}
		}

		if hasKnownObjectProperties || hasKnownDBIdentifierProperties {
			classifiedDataType.Classification.Decision = classify.ClassificationDecision{
				State:  classify.Valid,
				Reason: "valid_object_with_valid_properties",
			}
			return classifiedDataType, nil
		}

		objectState := classify.Invalid
		if isDatabase(data.DetectorType) {
			objectState = classify.Potential
		}
		classifiedDataType.Classification.Decision = classify.ClassificationDecision{
			State:  objectState,
			Reason: "valid_object_with_invalid_properties",
		}

		return classifiedDataType, nil
	}

	return classifiedDataType, nil
}

func (classifier *Classifier) isKnownPersonObject(name string) bool {
	result := false
	for _, pattern := range classifier.config.KnownPersonObjectPatterns {
		if pattern.IncludeRegexpMatcher.MatchString(name) {
			result = true
			break
		}

		if pattern.ExcludeRegexpMatcher != nil && !pattern.ExcludeRegexpMatcher.MatchString(name) {
			result = true
			break
		}
	}

	return result
}

func (classifier *Classifier) isKnownObjectPattern(name string, simpleType string) bool {
	result := false
	for _, pattern := range classifier.config.DataTypeClassificationPatterns {
		if _, isKnown := pattern.ObjectTypeMapping["known"]; !isKnown {
			continue
		}
		if pattern.DataTypeUUID == nil {
			continue
		}
		if !pattern.IncludeRegexpMatcher.MatchString(name) {
			continue
		}
		if !isExpectedIdentifierDataTypeId(pattern.Id) && regexpIdentifierMatcher.MatchString(name) {
			continue
		}
		if pattern.ExcludeRegexpMatcher != nil && pattern.ExcludeRegexpMatcher.MatchString(name) {
			continue
		}
		if _, isExcludedType := pattern.ExcludeTypesMapping[simpleType]; isExcludedType {
			continue
		}
		if !pattern.MatchColumn {
			continue
		}

		result = true
		break
	}

	return result
}

func (classifier *Classifier) isKnownDBIdentifierPattern(name string) bool {
	// todo: support health context
	result := false
	for _, pattern := range classifier.config.KnownPersonObjectPatterns {
		if !pattern.ActAsIdentifier {
			continue
		}
		if !pattern.IncludeRegexpMatcher.MatchString(name) {
			continue
		}
		if pattern.ExcludeRegexpMatcher != nil && pattern.ExcludeRegexpMatcher.MatchString(name) {
			continue
		}
		if pattern.IdentifierRegexpMatcher != nil && !pattern.IdentifierRegexpMatcher.MatchString(name) {
			continue
		}

		result = true
		break
	}

	return result
}

func isDatabase(detectorType detectors.Type) bool {
	_, ok := databaseDetectorTypes[string(detectorType)]
	return ok
}

func isExpectedIdentifierDataTypeId(id int) bool {
	_, ok := expectedIdentifierDataTypeIds[fmt.Sprint(id)]
	return ok
}

func objectStopWordDetected(name string) bool {
	_, ok := objectStopWords[name]
	return ok
}

func classifyAsInvalid(property datatype.DataTypable) ClassifiedDatatype {
	propertyDataType := extractDataType(property)
	normalizedPropertyName := normalize_key.Normalize(propertyDataType.Name)

	return ClassifiedDatatype{
		DataType: propertyDataType,
		Classification: Classification{
			Name: normalizedPropertyName,
			Decision: classify.ClassificationDecision{
				State:  classify.Invalid,
				Reason: "belongs_to_invalid_object",
			},
		},
	}
}
