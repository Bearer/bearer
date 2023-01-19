package datatype

import (
	objectdetector "github.com/bearer/curio/new/detector/implementation/ruby/object"
	propertydetector "github.com/bearer/curio/new/detector/implementation/ruby/property"
	"github.com/bearer/curio/new/detector/types"
	"github.com/bearer/curio/new/language/tree"
	languagetypes "github.com/bearer/curio/new/language/types"
	classificationschema "github.com/bearer/curio/pkg/classification/schema"
	"github.com/bearer/curio/pkg/report/detectors"
	"github.com/bearer/curio/pkg/report/schema"
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

func (detector *datatypeDetector) DetectAt(
	node *tree.Node,
	evaluator types.Evaluator,
) ([]*types.Detection, error) {
	objectDetections, err := evaluator.ForNode(node, "object")
	if err != nil {
		return nil, err
	}

	var result []*types.Detection

	for _, object := range objectDetections {
		var properties []Property

		objectData := object.Data.(objectdetector.Data)

		for _, property := range objectData.Properties {
			propertyData := property.Data.(propertydetector.Data)
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

		result = append(result, &types.Detection{
			ContextNode: node,
			MatchNode:   object.MatchNode,
			Data:        data,
		})

	}

	return result, nil
}

func (detector *datatypeDetector) Close() {}

// NOTE: presumption for mergeClassification is that classification will have all properties that detection has in same order
func mergeClassification(detection *Data, clasffication *classificationschema.ClassifiedDatatype) {
	detection.Classification = clasffication.Classification
	for i := range detection.Properties {
		detection.Properties[i].Classification = clasffication.Properties[i].Classification
	}
}
