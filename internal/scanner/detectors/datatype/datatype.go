package datatype

import (
	classificationschema "github.com/bearer/bearer/internal/classification/schema"
	"github.com/bearer/bearer/internal/commands/process/settings"
	"github.com/bearer/bearer/internal/report/detectors"
	"github.com/bearer/bearer/internal/report/schema"
	"github.com/bearer/bearer/internal/scanner/ast/tree"
	"github.com/bearer/bearer/internal/scanner/detectors/common"
	"github.com/bearer/bearer/internal/scanner/detectors/types"
	"github.com/bearer/bearer/internal/scanner/ruleset"
	"github.com/bearer/bearer/internal/util/classify"
)

type Data struct {
	Properties []Property
}

type Property struct {
	Name           string
	Node           *tree.Node
	Classification classificationschema.Classification
	Datatype       *types.Detection
}

type datatypeDetector struct {
	detectorType detectors.Type
	types.DetectorBase
	classifier *classificationschema.Classifier
}

func New(detectorType detectors.Type, classifier *classificationschema.Classifier) types.Detector {
	return &datatypeDetector{
		detectorType: detectorType,
		classifier:   classifier,
	}
}

func (detector *datatypeDetector) Rule() *ruleset.Rule {
	return ruleset.BuiltinDatatypeRule
}

func (detector *datatypeDetector) DetectAt(
	node *tree.Node,
	detectorContext types.Context,
) ([]interface{}, error) {
	objectDetections, err := detectorContext.Scan(node, ruleset.BuiltinObjectRule, settings.CURSOR_STRICT_SCOPE)
	if err != nil {
		return nil, err
	}

	var result []interface{}

	for _, object := range objectDetections {
		data, _, containsValidClassification := detector.classifyObject(detectorContext.Filename(), "", object)
		if containsValidClassification {
			result = append(result, data)
		}
	}

	return result, nil
}

func (detector *datatypeDetector) classifyObject(
	filename,
	name string,
	detection *types.Detection,
) (Data, classificationschema.Classification, bool) {
	objectData := detection.Data.(common.Object)

	classification := detector.classifier.Classify(buildClassificationRequest(detector.detectorType, filename, name, objectData))
	containsValidClassification := classification.Classification.Decision.State == classify.Valid

	properties := make([]Property, len(objectData.Properties))

	// NOTE: assumption is that classification will have all properties that detection has in same order
	for i, property := range objectData.Properties {
		propertyDetection, propertyClassification, containsValidPropertyClassification := detector.classifyProperty(
			filename,
			property.Name,
			property.Object,
			classification.Properties[i].Classification,
		)

		if !containsValidClassification && containsValidPropertyClassification {
			containsValidClassification = true
		}

		node := property.Node
		if node == nil {
			node = detection.MatchNode
		}

		properties[i] = Property{
			Datatype:       propertyDetection,
			Node:           node,
			Name:           property.Name,
			Classification: propertyClassification,
		}
	}

	return Data{Properties: properties}, classification.Classification, containsValidClassification
}

func (detector *datatypeDetector) classifyProperty(
	filename,
	name string,
	detection *types.Detection,
	parentClassification classificationschema.Classification,
) (*types.Detection, classificationschema.Classification, bool) {
	if detection == nil {
		return nil, parentClassification, false
	}

	data, propertyClassification, containsValidClassification := detector.classifyObject(filename, name, detection)

	propertyDetection := &types.Detection{
		RuleID:    "datatype",
		MatchNode: detection.MatchNode,
		Data:      data,
	}

	if parentClassification.Decision.State == classify.Valid {
		return propertyDetection, parentClassification, true
	}

	if (parentClassification.Decision.State == classify.Potential && propertyClassification.Decision.State == classify.Invalid) ||
		(parentClassification.Decision.State == classify.Invalid && propertyClassification.Decision.State == classify.Invalid) {

		return propertyDetection, parentClassification, containsValidClassification
	}

	return propertyDetection,
		propertyClassification,
		containsValidClassification || propertyClassification.Decision.State == classify.Valid
}

func buildClassificationRequest(
	detectorType detectors.Type,
	filename,
	name string,
	data common.Object,
) classificationschema.ClassificationRequest {
	var properties []*classificationschema.ClassificationRequestDetection

	for _, property := range data.Properties {
		properties = append(properties, &classificationschema.ClassificationRequestDetection{
			Name:       property.Name,
			SimpleType: schema.SimpleTypeUnknown,
		})
	}

	return classificationschema.ClassificationRequest{
		Value: &classificationschema.ClassificationRequestDetection{
			Name:       name,
			SimpleType: schema.SimpleTypeUnknown,
			Properties: properties,
		},
		DetectorType: detectorType,
		Filename:     filename,
	}
}
