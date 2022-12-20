package property

import (
	"fmt"

	detectiontypes "github.com/bearer/curio/new/detection/types"
	"github.com/bearer/curio/new/detector"
	"github.com/bearer/curio/new/language/tree"
	languagetypes "github.com/bearer/curio/new/language/types"
	treeevaluatortypes "github.com/bearer/curio/new/treeevaluator/types"
)

type Data struct {
	Name string
}

type propertyDetector struct {
	pairQuery *tree.Query
}

func New(lang languagetypes.Language) (detector.Detector, error) {
	pairQuery, err := lang.CompileQuery(`(pair key: (hash_key_symbol) @key value: (_) @value) @root`)
	if err != nil {
		return nil, fmt.Errorf("error compiling pair query: %s", err)
	}

	return &propertyDetector{pairQuery: pairQuery}, nil
}

func (detector *propertyDetector) Name() string {
	return "property"
}

func (detector *propertyDetector) DetectAt(
	node *tree.Node,
	evaluator treeevaluatortypes.Evaluator,
) ([]*detectiontypes.Detection, error) {
	result, err := detector.pairQuery.MatchOnceAt(node)
	if result == nil || err != nil {
		return nil, err
	}

	return []*detectiontypes.Detection{{
		MatchNode: node,
		Data:      Data{Name: result["key"].Content()},
	}}, nil
}

func (detector *propertyDetector) Close() {
	detector.pairQuery.Close()
}
