package schema

import (
	"regexp"

	"github.com/bearer/curio/pkg/report/schema/datatype"

	"github.com/bearer/curio/pkg/classification/db"
	"github.com/bearer/curio/pkg/report/detectors"
	"github.com/bearer/curio/pkg/util/classify"
)

var regexpIdentifierMatcher = regexp.MustCompile(`(uu)?id\z`)
var regexpTimestampsMatcher = regexp.MustCompile(`\A(created|updated)\sat\z`)

type ClassifiedDatatype struct {
	datatype.DataTypable
	Classification Classification `json:"classification"`
}

func (datatype ClassifiedDatatype) GetClassification() interface{} {
	return datatype.Classification
}

type Classification struct {
	Name     string                          `json:"name"`
	DataType *db.DataType                    `json:"data_type,omitempty"`
	Decision classify.ClassificationDecision `json:"decision"`
}

type Classifier struct {
	config Config
}

type Config struct {
	DataTypes                      []db.DataType
	DataTypeClassificationPatterns []db.DataTypeClassificationPattern
	KnownPersonObjectPatterns      []db.KnownPersonObjectPattern
	Context                        string
}

func New(config Config) *Classifier {
	return &Classifier{config: config}
}

type DataTypeDetection struct {
	Value        datatype.DataTypable
	Filename     string
	DetectorType detectors.Type
}

func (classifier *Classifier) Classify(data DataTypeDetection) *ClassifiedDatatype {
	var classifiedDatatype *ClassifiedDatatype

	// general checks
	if classify.IsVendored(data.Filename) {
		classifiedDatatype = classifyObjectAsInvalid(data.Value, classify.IncludedInVendorFolderReason)
	}
	if classify.IsPotentialDetector(data.DetectorType) {
		classifiedDatatype = classifyObjectAsInvalid(data.Value, classify.PotentialDetectorReason)
	}
	if classify.ObjectStopWordDetected(data.Value.GetNormalizedName()) {
		classifiedDatatype = classifyObjectAsInvalid(data.Value, "stop_word")
	}

	if classifiedDatatype != nil && classifiedDatatype.Classification.Decision.State == classify.Invalid {
		return classifiedDatatype
	}

	// schema-specific checks
	classifiedDatatype = &ClassifiedDatatype{
		DataTypable:    data.Value,
		Classification: Classification{Name: data.Value.GetNormalizedName()},
	}

	matchedKnownPersonObject := classifier.matchKnownPersonObjectPatterns(data.Value.GetNormalizedName(), false)
	if matchedKnownPersonObject != nil {
		// add data type to object
		classifiedDatatype.Classification.DataType = &matchedKnownPersonObject.DataType
		return classifier.classifyKnownObject(classifiedDatatype, data.DetectorType)
	}

	// do we have an object with unknown or unknown extended properties?
	isJSDetection := classify.IsJSDetection(data.DetectorType)
	if classifier.hasUnknownObjectProperties(classifiedDatatype.DataTypable.GetProperties(), isJSDetection) {
		return classifier.classifyObjectWithUnknownProperties(classifiedDatatype, isJSDetection)
	}

	hasIdentifierProperties := classifier.hasIdentifierProperties(classifiedDatatype.DataTypable.GetProperties(), isJSDetection)
	if hasIdentifierProperties {
		// object is somehow linked with a "person" e.g. an invoice with a user_id property
		return classifier.classifyObjectWithIdentifierProperties(classifiedDatatype, isJSDetection)
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

func (classifier *Classifier) hasIdentifierProperties(objectProperties map[string]datatype.DataTypable, isJSDetection bool) bool {
	for _, property := range objectProperties {
		if isJSDetection && classify.PropertyStopWordDetected(property.GetNormalizedName()) {
			continue
		}

		matchedIdentifier := classifier.matchKnownPersonObjectPatterns(property.GetNormalizedName(), true)
		if matchedIdentifier != nil {
			return true
		}
	}

	return false
}

func (classifier *Classifier) hasUnknownObjectProperties(objectProperties map[string]datatype.DataTypable, isJSDetection bool) bool {
	for _, property := range objectProperties {
		if isJSDetection && classify.PropertyStopWordDetected(property.GetNormalizedName()) {
			continue
		}

		matchedUnknownObject := classifier.matchObjectPatterns(property.GetNormalizedName(), property.GetType(), db.UnknownObject)
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

func (classifier *Classifier) classifyKnownObject(classifiedDatatype *ClassifiedDatatype, detectorType detectors.Type) *ClassifiedDatatype {
	isJSDetection := classify.IsJSDetection(detectorType)

	validProperties := false
	for _, property := range classifiedDatatype.DataTypable.GetProperties() {
		if isJSDetection && classify.PropertyStopWordDetected(property.GetNormalizedName()) {
			classifiedDatatype.DataTypable.SetProperty(
				property.GetName(),
				classifyAsInvalid(property, "stop_word"),
			)

			continue
		}

		matchedKnownObject := classifier.matchObjectPatterns(property.GetNormalizedName(), property.GetType(), db.KnownObject)
		if matchedKnownObject != nil {
			validProperties = true
			classifiedDatatype.DataTypable.SetProperty(
				property.GetName(),
				classifyAsValid(property, classifier.datatypeFromPattern(matchedKnownObject), "known_classification_pattern"),
			)

			continue
		}

		matchedKnownIdentifier := classifier.matchKnownPersonObjectPatterns(property.GetNormalizedName(), true)
		if matchedKnownIdentifier != nil {
			validProperties = true
			classifiedDatatype.DataTypable.SetProperty(
				property.GetName(),
				classifyAsValid(property, matchedKnownIdentifier.DataType, "known_database_identifier"),
			)

			continue
		}

		classifiedDatatype.DataTypable.SetProperty(
			property.GetName(),
			classifyAsInvalid(property, "invalid_property"),
		)
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

func (classifier *Classifier) classifyObjectWithUnknownProperties(classifiedDatatype *ClassifiedDatatype, isJSDetection bool) *ClassifiedDatatype {
	for _, property := range classifiedDatatype.GetProperties() {
		if isJSDetection && classify.PropertyStopWordDetected(property.GetNormalizedName()) {
			classifiedDatatype.DataTypable.SetProperty(
				property.GetName(),
				classifyAsInvalid(property, "stop_word"),
			)

			continue
		}

		// check unknown object patterns
		unknownObject := classifier.matchObjectPatterns(property.GetNormalizedName(), property.GetType(), db.UnknownObject)
		if unknownObject != nil {
			classifiedDatatype.DataTypable.SetProperty(
				property.GetName(),
				classifyAsValid(property, classifier.datatypeFromPattern(unknownObject), "valid_unknown_pattern"),
			)

			continue
		}

		// check extended patterns
		extendedUnknownObject := classifier.matchObjectPatterns(property.GetNormalizedName(), property.GetType(), db.ExtendedUnknownObject)
		if extendedUnknownObject != nil {
			classifiedDatatype.DataTypable.SetProperty(
				property.GetName(),
				classifyAsValid(property, classifier.datatypeFromPattern(extendedUnknownObject), "valid_extended_pattern"),
			)

			continue
		}

		classifiedDatatype.DataTypable.SetProperty(
			property.GetName(),
			classifyAsInvalid(property, "invalid_property"),
		)
	}

	classifiedDatatype.Classification.Decision = classify.ClassificationDecision{
		State:  classify.Valid,
		Reason: "invalid_object_with_valid_properties",
	}

	return classifiedDatatype
}

func (classifier *Classifier) classifyObjectWithIdentifierProperties(classifiedDatatype *ClassifiedDatatype, isJSDetection bool) *ClassifiedDatatype {
	associatedObjectProperties := false
	for _, property := range classifiedDatatype.GetProperties() {
		if isJSDetection && classify.PropertyStopWordDetected(property.GetNormalizedName()) {
			classifiedDatatype.DataTypable.SetProperty(
				property.GetName(),
				classifyAsInvalid(property, "stop_word"),
			)

			continue
		}

		matchedDBIdentifier := classifier.matchKnownPersonObjectPatterns(property.GetNormalizedName(), true)
		if matchedDBIdentifier != nil {
			classifiedDatatype.DataTypable.SetProperty(
				property.GetName(),
				classifyAsValid(property, matchedDBIdentifier.DataType, "known_database_identifier"),
			)

			continue
		}

		matchedAssociatedObjectPattern := classifier.matchObjectPatterns(property.GetNormalizedName(), property.GetType(), db.AssociatedObject)
		if matchedAssociatedObjectPattern != nil {
			associatedObjectProperties = true
			classifiedDatatype.DataTypable.SetProperty(
				property.GetName(),
				classifyAsValid(property, classifier.datatypeFromPattern(matchedAssociatedObjectPattern), "valid_associated_object_pattern"),
			)

			continue
		}

		classifiedDatatype.DataTypable.SetProperty(
			property.GetName(),
			classifyAsInvalid(property, "invalid_property"),
		)
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

	matchedObjectPattern := classifier.matchObjectPatterns(classifiedDatatype.GetNormalizedName(), "", db.KnownDataObject)
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
	if classifier.healthContext() {
		return pattern.HealthContextDataType
	}

	return pattern.DataType
}

func identifiersOnly(classifiedDatatype *ClassifiedDatatype) bool {
	identifiersOnly := true
	for _, property := range classifiedDatatype.GetProperties() {
		normalizedName := property.GetNormalizedName()
		if !regexpIdentifierMatcher.MatchString(normalizedName) && !regexpTimestampsMatcher.MatchString(normalizedName) {
			identifiersOnly = false
			break
		}
	}

	return identifiersOnly
}

func classifyObjectAsInvalid(D datatype.DataTypable, reason string) *ClassifiedDatatype {
	classifiedDatatype := &ClassifiedDatatype{
		DataTypable: D,
		Classification: Classification{
			Name: D.GetNormalizedName(),
			Decision: classify.ClassificationDecision{
				State:  classify.Invalid,
				Reason: reason,
			},
		},
	}

	// schema object did not pass initial checks ; mark all fields as invalid
	for _, property := range D.GetProperties() {
		classifiedDatatype.DataTypable.SetProperty(
			property.GetName(),
			classifyAsInvalid(property, "belongs_to_invalid_object"),
		)
	}

	return classifiedDatatype
}

func classifyAsValid(D datatype.DataTypable, datatype db.DataType, reason string) ClassifiedDatatype {
	return ClassifiedDatatype{
		DataTypable: D,
		Classification: Classification{
			Name:     D.GetNormalizedName(),
			DataType: &datatype,
			Decision: classify.ClassificationDecision{
				State:  classify.Valid,
				Reason: reason,
			},
		},
	}
}

func classifyAsInvalid(D datatype.DataTypable, reason string) ClassifiedDatatype {
	return ClassifiedDatatype{
		DataTypable: D,
		Classification: Classification{
			Name: D.GetNormalizedName(),
			Decision: classify.ClassificationDecision{
				State:  classify.Invalid,
				Reason: reason,
			},
		},
	}
}
