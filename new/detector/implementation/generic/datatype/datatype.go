package datatype

import (
	"github.com/bearer/bearer/new/detector/types"
	"github.com/bearer/bearer/new/language/tree"
	"github.com/bearer/bearer/pkg/report/detectors"
	"github.com/bearer/bearer/pkg/report/schema"

	generictypes "github.com/bearer/bearer/new/detector/implementation/generic/types"
	languagetypes "github.com/bearer/bearer/new/language/types"
	classificationschema "github.com/bearer/bearer/pkg/classification/schema"
)

type Data struct {
	Name           string
	Classification classificationschema.Classification
	Properties     []Property
}

func (data *Data) toClassifcationRequestDetection() *classificationschema.ClassificationRequestDetection {
	req := &classificationschema.ClassificationRequestDetection{
		Name:       data.Name,
		SimpleType: schema.SimpleTypeUnknown,
	}
	for _, property := range data.Properties {
		req.Properties = append(req.Properties, &classificationschema.ClassificationRequestDetection{
			Name:       property.Name,
			SimpleType: schema.SimpleTypeUnknown,
			Properties: []*classificationschema.ClassificationRequestDetection{},
		})
	}
	return req
}

type Property struct {
	Name           string
	Detection      *types.Detection
	Classification classificationschema.Classification
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
		var properties []Property

		objectData := object.Data.(generictypes.Object)

		for _, property := range objectData.Properties {
			propertyData := property.Data.(generictypes.Property)
			properties = append(properties, Property{
				Detection: property,
				Name:      propertyData.Name,
			})
		}

		data := Data{
			Name:       objectData.Name,
			Properties: properties,
		}

		classificationReqDetection := data.toClassifcationRequestDetection()

		classification := detector.classifier.Classify(classificationschema.ClassificationRequest{
			Value:        classificationReqDetection,
			DetectorType: detectors.DetectorRuby,
			Filename:     evaluator.FileName(),
		})

		mergeClassification(&data, classification)

		result = append(result, data)
	}

	return result, nil
}

func (detector *datatypeDetector) Close() {}

// NOTE: presumption for mergeClassification is that classification will have all properties that detection has in same order
func mergeClassification(detection *Data, classification *classificationschema.ClassifiedDatatype) {
	detection.Classification = classification.Classification
	for i := range detection.Properties {
		detection.Properties[i].Classification = classification.Properties[i].Classification
	}
}
