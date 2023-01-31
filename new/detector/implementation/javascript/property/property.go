package property

import (
	"fmt"

	"github.com/bearer/curio/new/detector/types"
	"github.com/bearer/curio/new/language/tree"

	generictypes "github.com/bearer/curio/new/detector/implementation/generic/types"
	languagetypes "github.com/bearer/curio/new/language/types"
)

type propertyDetector struct {
	pairQuery       *tree.Query
	methodNameQuery *tree.Query
}

func New(lang languagetypes.Language) (types.Detector, error) {
	// { user: "admin_user" }
	pairQuery, err := lang.CompileQuery(`(pair key: (_) @key value: (_) @value) @root`)
	if err != nil {
		return nil, fmt.Errorf("error compiling pair query: %s", err)
	}

	// function getName(){}
	methodNameQuery, err := lang.CompileQuery(`(function_declaration name: (identifier) @name) @root`)
	if err != nil {
		return nil, fmt.Errorf("error compiling class method query: %s", err)
	}

	return &propertyDetector{
		pairQuery:       pairQuery,
		methodNameQuery: methodNameQuery,
	}, nil
}

func (detector *propertyDetector) Name() string {
	return "property"
}

func (detector *propertyDetector) DetectAt(
	node *tree.Node,
	evaluator types.Evaluator,
) ([]interface{}, error) {
	// run hash pair query
	result, err := detector.pairQuery.MatchOnceAt(node)
	if err != nil {
		return nil, err
	}

	if len(result) != 0 {
		key := result["key"]
		// { user: "admin_user"} || {"user": "admin_user"}
		if key.Type() == "property_identifier" || key.Type() == "string" {
			return []interface{}{generictypes.Property{Name: result["key"].Content()}}, nil
		}
	}

	// run method name query
	result, err = detector.methodNameQuery.MatchOnceAt(node)
	if err != nil {
		return nil, err
	}
	if len(result) != 0 {
		return []interface{}{generictypes.Property{Name: result["name"].Content()}}, nil
	}

	return nil, nil
}

func (detector *propertyDetector) Close() {
	detector.pairQuery.Close()
}
