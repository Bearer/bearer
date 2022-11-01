package schema

import (
	"regexp"

	"github.com/bearer/curio/pkg/report/schema/datatype"

	"github.com/bearer/curio/pkg/classification/db"
	"github.com/bearer/curio/pkg/report/detectors"
	"github.com/bearer/curio/pkg/util/classify"
	"github.com/bearer/curio/pkg/util/normalize_key"
)

var regexpIdentifierMatcher = regexp.MustCompile(`(uu)?id\z`)

type ClassifiedDatatype struct {
	datatype.DataTypable
	Classification Classification `json:"classification"`
}

type Classification struct {
	Name     string                          `json:"name"`
	DataType db.DataType                     `json:"data_type,omitempty"`
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

func (classifier *Classifier) Classify(data DataTypeDetection) (*ClassifiedDatatype, error) {
	dataTypeable := data.Value
	normalizedObjectName := normalize_key.Normalize(dataTypeable.GetName())
	objectProperties := dataTypeable.GetProperties()
	classifiedDataType := &ClassifiedDatatype{
		DataTypable:    dataTypeable,
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

	if classify.ObjectStopWordDetected(normalizedObjectName) {
		classifiedDataType.Classification.Decision = classify.ClassificationDecision{
			State:  classify.Invalid,
			Reason: "stop_word",
		}
	}

	if classifiedDataType.Classification.Decision.State == classify.Invalid {
		// schema object did not pass initial checks
		// mark all first level children as invalid

		// todo: handle children that are themselves schema objects
		for _, property := range objectProperties {
			objectProperties[property.GetName()] = classifyAsInvalid(property)
		}

		return classifiedDataType, nil
	}

	var isJSDetection = data.DetectorType == detectors.DetectorJavascript || data.DetectorType == detectors.DetectorTypescript

	matchedKnownPersonObject := classifier.matchKnownPersonObjectPatterns(normalizedObjectName, false)
	if matchedKnownPersonObject != nil {
		// add data type to object
		classifiedDataType.Classification.DataType = matchedKnownPersonObject.DataType

		hasKnownObjectProperties := false
		hasKnownDBIdentifierProperties := false

		// todo: handle children that are themselves schema objects
		for _, property := range objectProperties {
			normalizedPropertyName := normalize_key.Normalize(property.GetName())

			if isJSDetection && classify.PropertyStopWordDetected(normalizedPropertyName) {
				objectProperties[property.GetName()] = ClassifiedDatatype{
					DataTypable: property,
					Classification: Classification{
						Name: normalizedPropertyName,
						Decision: classify.ClassificationDecision{
							State:  classify.Invalid,
							Reason: "stop_word",
						},
					},
				}

				continue
			}

			matchedKnownObject := classifier.matchKnownObjectPatterns(normalizedPropertyName, property.GetType())
			if matchedKnownObject != nil {
				hasKnownObjectProperties = true
				objectProperties[property.GetName()] = ClassifiedDatatype{
					DataTypable: property,
					Classification: Classification{
						Name:     normalizedPropertyName,
						DataType: matchedKnownObject.DataType, // todo: check for health context
						Decision: classify.ClassificationDecision{
							State:  classify.Valid,
							Reason: "known_classification_pattern",
						},
					},
				}

				continue
			}

			matchedKnownIdentifier := classifier.matchKnownPersonObjectPatterns(normalizedPropertyName, true)
			if matchedKnownIdentifier != nil {
				hasKnownDBIdentifierProperties = true
				objectProperties[property.GetName()] = ClassifiedDatatype{
					DataTypable: property,
					Classification: Classification{
						Name:     normalizedPropertyName,
						DataType: matchedKnownIdentifier.DataType, // always "Unique Identifier"
						Decision: classify.ClassificationDecision{
							State:  classify.Valid,
							Reason: "known_database_identifier",
						},
					},
				}

				continue
			}

			objectProperties[property.GetName()] = ClassifiedDatatype{
				DataTypable: property,
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
		if classify.IsDatabase(data.DetectorType) {
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

func (classifier *Classifier) matchKnownObjectPatterns(name string, simpleType string) *db.DataTypeClassificationPattern {
	var matchedPattern *db.DataTypeClassificationPattern
	for _, pattern := range classifier.config.DataTypeClassificationPatterns {
		if _, isKnown := pattern.ObjectTypeMapping["known"]; !isKnown {
			continue
		}
		if pattern.DataTypeUUID == "" {
			continue
		}
		if !pattern.IncludeRegexpMatcher.MatchString(name) {
			continue
		}
		if !classify.IsExpectedIdentifierDataTypeId(pattern.Id) && regexpIdentifierMatcher.MatchString(name) {
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

		matchedPattern = &pattern
		break
	}

	return matchedPattern
}

func (classifier *Classifier) matchKnownPersonObjectPatterns(name string, matchAsIdentifier bool) *db.KnownPersonObjectPattern {
	var matchedPattern *db.KnownPersonObjectPattern

	// todo: support health context
	for _, pattern := range classifier.config.KnownPersonObjectPatterns {
		if matchAsIdentifier && !pattern.ActAsIdentifier {
			continue
		}
		if !pattern.IncludeRegexpMatcher.MatchString(name) {
			continue
		}
		if pattern.ExcludeRegexpMatcher != nil && pattern.ExcludeRegexpMatcher.MatchString(name) {
			continue
		}
		if matchAsIdentifier && pattern.IdentifierRegexpMatcher != nil && !pattern.IdentifierRegexpMatcher.MatchString(name) {
			continue
		}

		matchedPattern = &pattern
		break
	}

	return matchedPattern
}

func classifyAsInvalid(property datatype.DataTypable) ClassifiedDatatype {
	normalizedPropertyName := normalize_key.Normalize(property.GetName())

	return ClassifiedDatatype{
		DataTypable: property,
		Classification: Classification{
			Name: normalizedPropertyName,
			Decision: classify.ClassificationDecision{
				State:  classify.Invalid,
				Reason: "belongs_to_invalid_object",
			},
		},
	}
}
