package datatype

import (
	objectdetector "github.com/bearer/curio/new/detector/implementation/ruby/object"
	propertydetector "github.com/bearer/curio/new/detector/implementation/ruby/property"
	"github.com/bearer/curio/new/detector/types"
	"github.com/bearer/curio/new/language/tree"
	languagetypes "github.com/bearer/curio/new/language/types"
	"github.com/bearer/curio/pkg/classification/db"
	"github.com/bearer/curio/pkg/classification/schema"
	"github.com/bearer/curio/pkg/parser/nodeid"
	"github.com/bearer/curio/pkg/util/classify"
)

type Data struct {
	Classification schema.Classification
	Properties     []Property
}

type Property struct {
	Detection      *types.Detection
	Classification schema.Classification
}

type datatypeDetector struct {
	classifier  *schema.Classifier
	idGenerator nodeid.Generator
}

func New(lang languagetypes.Language, classifier *schema.Classifier, idGenerator nodeid.Generator) (types.Detector, error) {
	return &datatypeDetector{
		classifier:  classifier,
		idGenerator: idGenerator,
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

		data := object.Data.(objectdetector.Data)

		for _, property := range data.Properties {
			propertyData := property.Data.(propertydetector.Data)
			var classification schema.Classification
			if propertyData.Name == "first_name" && data.Name == "user" {
				classification = schema.Classification{
					DataType: &db.DataType{
						Name:         "Firstname",
						UUID:         "380c8cde-ca2e-44ed-82db-2ab1e7c255c7",
						CategoryUUID: "14124881-6b92-4fc5-8005-ea7c1c09592e",
					},
					Name: "first_name",
					Decision: classify.ClassificationDecision{
						State:  classify.Valid,
						Reason: "matches piid",
					},
				}
			} else {
				classification = schema.Classification{
					DataType: nil,
					Name:     propertyData.Name,
					Decision: classify.ClassificationDecision{
						State:  classify.Invalid,
						Reason: "coudln't find datatype",
					},
				}
			}

			properties = append(properties, Property{
				Detection:      property,
				Classification: classification,
			})
		}
		result = append(result, &types.Detection{
			ContextNode: node,
			MatchNode:   object.MatchNode,
			Data: Data{
				Classification: schema.Classification{
					DataType: nil,
					Name:     data.Name,
					Decision: classify.ClassificationDecision{
						State:  classify.Invalid,
						Reason: "coudln't find datatype",
					},
				},
				Properties: properties,
			},
		})
	}

	return result, nil
}

func (detector *datatypeDetector) Close() {}
