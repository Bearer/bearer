package datatype

import (
	generictypes "github.com/bearer/bearer/new/detector/implementation/generic/types"
	"github.com/bearer/bearer/new/detector/types"
	"github.com/bearer/bearer/new/language/tree"
	languagetypes "github.com/bearer/bearer/new/language/types"
	classificationschema "github.com/bearer/bearer/pkg/classification/schema"
	"github.com/bearer/bearer/pkg/report/detectors"
	"github.com/bearer/bearer/pkg/report/schema"
	"github.com/bearer/bearer/pkg/util/classify"
)

type Data struct {
	Properties []Property
}

type Property struct {
	Name           string
	Classification classificationschema.Classification
	Datatype       *types.Detection
}

type datatypeDetector struct {
	types.DetectorBase
	classifier *classificationschema.Classifier
}

func New(lang languagetypes.Language, classifier *classificationschema.Classifier) (types.Detector, error) {
	return &datatypeDetector{
		classifier: classifier,
	}, nil
}

func (detector *datatypeDetector) Name() string {
	return "datatype"
}

func (detector *datatypeDetector) NestedDetections() bool {
	return false
}

func (detector *datatypeDetector) DetectAt(
	node *tree.Node,
	evaluator types.Evaluator,
) ([]interface{}, error) {
	objectDetections, err := evaluator.ForNode(node, "object", false)
	if err != nil {
		return nil, err
	}

	var result []interface{}

	for _, object := range objectDetections {
		data, _ := detector.classifyObject(evaluator.FileName(), "", object)
		result = append(result, data)
	}

	return result, nil
}

func (detector *datatypeDetector) Close() {}

func (detector *datatypeDetector) classifyObject(
	filename,
	name string,
	detection *types.Detection,
) (Data, classificationschema.Classification) {
	objectData := detection.Data.(generictypes.Object)

	classification := detector.classifier.Classify(buildClassificationRequest(filename, name, objectData))

	var properties []Property

	// NOTE: assumption is that classification will have all properties that detection has in same order
	for i, property := range objectData.Properties {
		propertyDetection, propertyClassification := detector.classifyProperty(
			filename,
			property.Name,
			property.Object,
			classification.Properties[i].Classification,
		)

		properties = append(properties, Property{
			Datatype:       propertyDetection,
			Name:           property.Name,
			Classification: propertyClassification,
		})
	}

	return Data{Properties: properties}, classification.Classification
}

func (detector *datatypeDetector) classifyProperty(
	filename,
	name string,
	detection *types.Detection,
	parentClassification classificationschema.Classification,
) (*types.Detection, classificationschema.Classification) {
	if detection == nil {
		return nil, parentClassification
	}

	data, propertyClassification := detector.classifyObject(filename, name, detection)

	propertyDetection := &types.Detection{
		DetectorType: "datatype",
		MatchNode:    detection.MatchNode,
		Data:         data,
	}

	if parentClassification.Decision.State == classify.Valid ||
		(parentClassification.Decision.State == classify.Potential && propertyClassification.Decision.State == classify.Invalid) ||
		(parentClassification.Decision.State == classify.Invalid && propertyClassification.Decision.State == classify.Invalid) {
		return propertyDetection, parentClassification
	}

	return propertyDetection, propertyClassification
}

func buildClassificationRequest(filename, name string, data generictypes.Object) classificationschema.ClassificationRequest {
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
		DetectorType: detectors.DetectorRuby,
		Filename:     filename,
	}
}
