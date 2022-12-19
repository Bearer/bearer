package ffi

import (
	"errors"
	"log"

	"github.com/open-policy-agent/opa/ast"
	"github.com/open-policy-agent/opa/rego"
	"github.com/open-policy-agent/opa/types"

	detectiontypes "github.com/bearer/curio/new/detection/types"
	"github.com/bearer/curio/new/language"
	treeevaluatortypes "github.com/bearer/curio/new/treeevaluator/types"
)

type Data struct {
	evaluator treeevaluatortypes.Evaluator
	nodes     opaqueTypeMap[language.Node]
}

var (
	currentData *Data

	opaqueType = types.NewObject([]*types.StaticProperty{
		{Key: "type", Value: types.S},
		{Key: "id", Value: types.N},
	}, nil)

	detectionType = types.NewObject([]*types.StaticProperty{{Key: "node", Value: opaqueType}}, nil)
)

func SetData(data *Data) {
	currentData = data
}

func NewData(evaluator treeevaluatortypes.Evaluator) *Data {
	return &Data{
		evaluator: evaluator,
		nodes:     newOpaqueTypeMap[language.Node]("node"),
	}
}

func GetData() (*Data, error) {
	if currentData == nil {
		return nil, errors.New("ffi data not set")
	}

	return currentData, nil
}

func (data *Data) NodeToRegoInput(node *language.Node) interface{} {
	return data.nodes.CastToRegoInput(node)
}

func (data *Data) DetectionsToRego(detections []*detectiontypes.Detection) (*ast.Term, error) {
	detectionTerms := make([]*ast.Term, len(detections))

	for i, detection := range detections {
		dataValue, err := interfaceToValue(detection.Data)
		if err != nil {
			return nil, err
		}

		detectionTerms[i] = ast.ObjectTerm(
			ast.Item(ast.StringTerm("match_node"), data.nodes.CastToRego(detection.MatchNode)),
			ast.Item(ast.StringTerm("context_node"), data.nodes.CastToRego(detection.ContextNode)),
			ast.Item(ast.StringTerm("data"), ast.NewTerm(dataValue)),
		)
	}

	return ast.ArrayTerm(detectionTerms...), nil
}

func EvaluatorNodeDetections() func(*rego.Rego) {
	return rego.Function2(
		&rego.Function{
			Name: "curio.evaluator.node_detections",
			Decl: types.NewFunction(
				[]types.Type{opaqueType, types.S},
				types.NewArray([]types.Type{detectionType}, nil),
			),
			Memoize: true,
		},
		func(bctx rego.BuiltinContext, nodeTerm, detectorTypeTerm *ast.Term) (*ast.Term, error) {
			log.Printf("REGO! start")
			data, err := GetData()
			if err != nil {
				return nil, err
			}

			node, err := data.nodes.CastToGo(nodeTerm)
			if err != nil {
				log.Printf("REGO! err: %s", err)
				return nil, err
			}

			detectorType := string(detectorTypeTerm.Value.(ast.String))

			log.Printf("REGO! detector type: %s\n%s\n", detectorType, node.Content())

			detections, err := data.evaluator.NodeDetections(node, detectorType)
			if err != nil {
				return nil, err
			}

			return data.DetectionsToRego(detections)
		},
	)
}
