package schema

import (
	"regexp"

	"github.com/bearer/curio/pkg/flag"

	"github.com/bearer/curio/pkg/classification/db"
	"github.com/bearer/curio/pkg/report/detectors"
	"github.com/bearer/curio/pkg/util/classify"
	"github.com/bearer/curio/pkg/util/normalize_key"
)

var regexpIdentifierMatcher = regexp.MustCompile(`(uu)?id\z`)
var regexpTimestampsMatcher = regexp.MustCompile(`\A(created|updated)\sat\z`)

type ClassifiedDatatype struct {
	Name           string
	Properties     []*ClassifiedDatatype
	Classification Classification
}

func (datatype ClassifiedDatatype) GetClassification() interface{} {
	return datatype.Classification
}

type Classification struct {
	Name     string                          `json:"name" yaml:"name"`
	DataType *db.DataType                    `json:"data_type,omitempty"`
	Decision classify.ClassificationDecision `json:"decision" yaml:"decision"`
}

type Classifier struct {
	config Config
}

type Config struct {
	DataTypes                      []db.DataType
	DataTypeClassificationPatterns []db.DataTypeClassificationPattern
	KnownPersonObjectPatterns      []db.KnownPersonObjectPattern
	Context                        flag.Context
}

func New(config Config) *Classifier {
	return &Classifier{config: config}
}

type ClassificationRequestDetection struct {
	Name       string
	SimpleType string
	Properties []*ClassificationRequestDetection
}

type ClassificationRequest struct {
	Value        *ClassificationRequestDetection
	Filename     string
	DetectorType detectors.Type
}

func (classifier *Classifier) Classify(data ClassificationRequest) *ClassifiedDatatype {
	var classifiedDatatype *ClassifiedDatatype
	var normalizedName = normalize_key.Normalize(data.Value.Name)

	// general checks
	if classify.IsVendored(data.Filename) {
		classifiedDatatype = classifyObjectAsInvalid(data.Value, classify.IncludedInVendorFolderReason)
	}
	if classify.IsPotentialDetector(data.DetectorType) {
		classifiedDatatype = classifyObjectAsInvalid(data.Value, classify.PotentialDetectorReason)
	}
	if classify.ObjectStopWordDetected(normalizedName) {
		classifiedDatatype = classifyObjectAsInvalid(data.Value, "stop_word")
	}
	if data.Value.Name == "" {
		classifiedDatatype = classifyObjectAsInvalid(data.Value, "blank_object_name")
	}

	if classifiedDatatype != nil && classifiedDatatype.Classification.Decision.State == classify.Invalid {
		return classifiedDatatype
	}

	// schema-specific checks
	var properties []*ClassifiedDatatype
	for _, v := range data.Value.Properties {
		properties = append(properties, &ClassifiedDatatype{
			Name: v.Name,
			Classification: Classification{
				Name: normalize_key.Normalize(v.Name),
			},
		})
	}

	classifiedDatatype = &ClassifiedDatatype{
		Name:           data.Value.Name,
		Classification: Classification{Name: normalize_key.Normalize(normalizedName)},
		Properties:     properties,
	}

	matchedKnownPersonObject := classifier.matchKnownPersonObjectPatterns(normalizedName, false)
	if matchedKnownPersonObject != nil {
		// add data type to object
		classifiedDatatype.Classification.DataType = &matchedKnownPersonObject.DataType
		return classifier.classifyKnownObject(classifiedDatatype, data.Value, data.DetectorType)
	}

	// do we have an object with unknown or unknown extended properties?
	isJSDetection := classify.IsJSDetection(data.DetectorType)
	if classifier.hasUnknownObjectProperties(data.Value.Properties, isJSDetection) {
		return classifier.classifyObjectWithUnknownProperties(classifiedDatatype, data.Value, isJSDetection)
	}

	hasIdentifierProperties := classifier.hasIdentifierProperties(data.Value.Properties, isJSDetection)
	if hasIdentifierProperties {
		// object is somehow linked with a "person" e.g. an invoice with a user_id property
		return classifier.classifyObjectWithIdentifierProperties(classifiedDatatype, data.Value, isJSDetection)
	}

	// object and properties are unknown
	objectState := classify.Invalid
	if classify.IsDatabase(data.DetectorType) {
		objectState = classify.Potential
	}
	classifiedDatatype.Classification.Decision = classify.ClassificationDecision{
		State:  objectState,
		Reason: "unknown_data_object",
	}

	return classifiedDatatype
}

func (classifier *Classifier) hasIdentifierProperties(objectProperties []*ClassificationRequestDetection, isJSDetection bool) bool {
	for _, property := range objectProperties {
		if isJSDetection && classify.PropertyStopWordDetected(normalize_key.Normalize(property.Name)) {
			continue
		}

		matchedIdentifier := classifier.matchKnownPersonObjectPatterns(normalize_key.Normalize(property.Name), true)
		if matchedIdentifier != nil {
			return true
		}
	}

	return false
}

func (classifier *Classifier) hasUnknownObjectProperties(objectProperties []*ClassificationRequestDetection, isJSDetection bool) bool {
	for _, property := range objectProperties {
		if isJSDetection && classify.PropertyStopWordDetected(normalize_key.Normalize(property.Name)) {
			continue
		}

		matchedUnknownObject := classifier.matchObjectPatterns(normalize_key.Normalize(property.Name), property.SimpleType, db.UnknownObject)
		if matchedUnknownObject != nil {
			return true
		}
	}

	return false
}

func (classifier *Classifier) matchObjectPatterns(name string, simpleType string, objectType db.ObjectType) *db.DataTypeClassificationPattern {
	matchObject := objectType == db.KnownDataObject // we're matching on a schema object, not a property

	var matchedPattern *db.DataTypeClassificationPattern
	for _, pattern := range classifier.config.DataTypeClassificationPatterns {
		if _, correctType := pattern.ObjectTypeMapping[string(objectType)]; !correctType {
			continue
		}
		if !classifier.healthContext() && pattern.DataTypeUUID == "" {
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
		if matchObject && !pattern.MatchObject {
			continue
		}
		if _, isExcludedType := pattern.ExcludeTypesMapping[simpleType]; isExcludedType {
			continue
		}
		if !matchObject && !pattern.MatchColumn {
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

func (classifier *Classifier) classifyKnownObject(classifiedDatatype *ClassifiedDatatype, detection *ClassificationRequestDetection, detectorType detectors.Type) *ClassifiedDatatype {
	isJSDetection := classify.IsJSDetection(detectorType)

	validProperties := false
	for i, property := range classifiedDatatype.Properties {
		if isJSDetection && classify.PropertyStopWordDetected(property.Classification.Name) {
			classifiedDatatype.Properties[i] = classifyAsInvalid(detection.Properties[i], "stop_word")
			continue
		}

		matchedKnownObject := classifier.matchObjectPatterns(property.Classification.Name, detection.Properties[i].SimpleType, db.KnownObject)
		if matchedKnownObject != nil {
			validProperties = true
			classifiedDatatype.Properties[i] = classifyAsValid(detection.Properties[i], classifier.datatypeFromPattern(matchedKnownObject), "known_pattern")
			continue
		}

		matchedKnownIdentifier := classifier.matchKnownPersonObjectPatterns(normalize_key.Normalize(property.Name), true)
		if matchedKnownIdentifier != nil {
			validProperties = true
			classifiedDatatype.Properties[i] = classifyAsValid(detection.Properties[i], matchedKnownIdentifier.DataType, "known_database_identifier")
			continue
		}

		classifiedDatatype.Properties[i] = classifyAsInvalid(detection.Properties[i], "invalid_property")

	}

	if validProperties {
		classifiedDatatype.Classification.Decision = classify.ClassificationDecision{
			State:  classify.Valid,
			Reason: "valid_object_with_valid_properties",
		}
		return classifiedDatatype
	}

	objectState := classify.Invalid
	if classify.IsDatabase(detectorType) {
		objectState = classify.Potential
	}
	classifiedDatatype.Classification.Decision = classify.ClassificationDecision{
		State:  objectState,
		Reason: "valid_object_with_invalid_properties",
	}

	return classifiedDatatype
}

func (classifier *Classifier) classifyObjectWithUnknownProperties(classifiedDatatype *ClassifiedDatatype, detection *ClassificationRequestDetection, isJSDetection bool) *ClassifiedDatatype {
	for i, property := range classifiedDatatype.Properties {
		if isJSDetection && classify.PropertyStopWordDetected(normalize_key.Normalize(property.Name)) {
			classifiedDatatype.Properties[i] = classifyAsInvalid(detection.Properties[i], "stop_word")
			continue
		}

		// check unknown object patterns
		unknownObject := classifier.matchObjectPatterns(normalize_key.Normalize(property.Name), detection.Properties[i].SimpleType, db.UnknownObject)
		if unknownObject != nil {
			classifiedDatatype.Properties[i] = classifyAsValid(detection.Properties[i], classifier.datatypeFromPattern(unknownObject), "valid_unknown_pattern")
			continue
		}

		// check extended patterns
		extendedUnknownObject := classifier.matchObjectPatterns(normalize_key.Normalize(property.Name), detection.Properties[i].SimpleType, db.ExtendedUnknownObject)
		if extendedUnknownObject != nil {
			classifiedDatatype.Properties[i] = classifyAsValid(detection.Properties[i], classifier.datatypeFromPattern(extendedUnknownObject), "valid_extended_pattern")
			continue
		}

		// check identifier patterns
		matchedKnownIdentifier := classifier.matchKnownPersonObjectPatterns(normalize_key.Normalize(property.Name), true)
		if matchedKnownIdentifier != nil {
			classifiedDatatype.Properties[i] = classifyAsValid(detection.Properties[i], matchedKnownIdentifier.DataType, "known_database_identifier")
			continue
		}

		classifiedDatatype.Properties[i] = classifyAsInvalid(detection.Properties[i], "invalid_property")
	}

	classifiedDatatype.Classification.Decision = classify.ClassificationDecision{
		State:  classify.Valid,
		Reason: "invalid_object_with_valid_properties",
	}

	return classifiedDatatype
}

func (classifier *Classifier) classifyObjectWithIdentifierProperties(classifiedDatatype *ClassifiedDatatype, detection *ClassificationRequestDetection, isJSDetection bool) *ClassifiedDatatype {
	associatedObjectProperties := false
	for i, property := range classifiedDatatype.Properties {
		if isJSDetection && classify.PropertyStopWordDetected(normalize_key.Normalize(property.Name)) {
			classifiedDatatype.Properties[i] = classifyAsInvalid(detection.Properties[i], "stop_word")
			continue
		}

		matchedDBIdentifier := classifier.matchKnownPersonObjectPatterns(normalize_key.Normalize(property.Name), true)
		if matchedDBIdentifier != nil {
			classifiedDatatype.Properties[i] = classifyAsValid(detection.Properties[i], matchedDBIdentifier.DataType, "known_database_identifier")
			continue
		}

		matchedAssociatedObjectPattern := classifier.matchObjectPatterns(normalize_key.Normalize(property.Name), detection.Properties[i].SimpleType, db.AssociatedObject)
		if matchedAssociatedObjectPattern != nil {
			associatedObjectProperties = true
			classifiedDatatype.Properties[i] = classifyAsValid(detection.Properties[i], classifier.datatypeFromPattern(matchedAssociatedObjectPattern), "valid_associated_object_pattern")
			continue
		}

		classifiedDatatype.Properties[i] = classifyAsInvalid(detection.Properties[i], "invalid_property")

	}

	if associatedObjectProperties {
		classifiedDatatype.Classification.Decision = classify.ClassificationDecision{
			State:  classify.Valid,
			Reason: "invalid_object_with_valid_properties",
		}
		return classifiedDatatype
	}

	// object composed only of DB identifiers ; check for known data object patterns
	return classifier.classifySchemaObject(classifiedDatatype)
}

func (classifier *Classifier) classifySchemaObject(classifiedDatatype *ClassifiedDatatype) *ClassifiedDatatype {
	if identifiersOnly(classifiedDatatype) {
		classifiedDatatype.Classification.Decision = classify.ClassificationDecision{
			State:  classify.Invalid,
			Reason: "only_db_identifiers",
		}
		return classifiedDatatype
	}

	matchedObjectPattern := classifier.matchObjectPatterns(classifiedDatatype.Classification.Name, "", db.KnownDataObject)
	if matchedObjectPattern != nil {
		classifiedDatatype.Classification.Decision = classify.ClassificationDecision{
			State:  classify.Valid,
			Reason: "known_data_object",
		}
		return classifiedDatatype
	}

	classifiedDatatype.Classification.Decision = classify.ClassificationDecision{
		State:  classify.Invalid,
		Reason: "unknown_data_object",
	}
	return classifiedDatatype
}

func (classifier *Classifier) healthContext() bool {
	return classifier.config.Context == "health"
}

func (classifier *Classifier) datatypeFromPattern(pattern *db.DataTypeClassificationPattern) db.DataType {
	if classifier.healthContext() && pattern.HealthContextDataTypeUUID != "" {
		return pattern.HealthContextDataType
	}

	return pattern.DataType
}

func identifiersOnly(classifiedDatatype *ClassifiedDatatype) bool {
	identifiersOnly := true
	for _, property := range classifiedDatatype.Properties {
		normalizedName := property.Classification.Name
		if !regexpIdentifierMatcher.MatchString(normalizedName) && !regexpTimestampsMatcher.MatchString(normalizedName) {
			identifiersOnly = false
			break
		}
	}

	return identifiersOnly
}

func classifyObjectAsInvalid(D *ClassificationRequestDetection, reason string) *ClassifiedDatatype {
	classifiedDatatype := &ClassifiedDatatype{
		Classification: Classification{
			Name: normalize_key.Normalize(D.Name),
			Decision: classify.ClassificationDecision{
				State:  classify.Invalid,
				Reason: reason,
			},
		},
	}

	// schema object did not pass initial checks ; mark all fields as invalid
	for _, property := range D.Properties {
		classifiedDatatype.Properties = append(classifiedDatatype.Properties, classifyAsInvalid(property, "belongs_to_invalid_object"))
	}

	return classifiedDatatype
}

func classifyAsValid(D *ClassificationRequestDetection, datatype db.DataType, reason string) *ClassifiedDatatype {
	return &ClassifiedDatatype{
		Name: D.Name,
		Classification: Classification{
			Name:     normalize_key.Normalize(D.Name),
			DataType: &datatype,
			Decision: classify.ClassificationDecision{
				State:  classify.Valid,
				Reason: reason,
			},
		},
	}
}

func classifyAsInvalid(D *ClassificationRequestDetection, reason string) *ClassifiedDatatype {
	return &ClassifiedDatatype{
		Name: D.Name,
		Classification: Classification{
			Name: normalize_key.Normalize(D.Name),
			Decision: classify.ClassificationDecision{
				State:  classify.Invalid,
				Reason: reason,
			},
		},
	}
}
