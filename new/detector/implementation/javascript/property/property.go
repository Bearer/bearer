package property

import (
	"fmt"

	"github.com/bearer/curio/new/detector/types"
	"github.com/bearer/curio/new/language/tree"

	generictypes "github.com/bearer/curio/new/detector/implementation/generic/types"
	languagetypes "github.com/bearer/curio/new/language/types"
)

type propertyDetector struct {
	pairQuery         *tree.Query
	functionNameQuery *tree.Query
	methodNameQuery   *tree.Query
}

func New(lang languagetypes.Language) (types.Detector, error) {
	// { user: "admin_user" }
	pairQuery, err := lang.CompileQuery(`(pair key: (_) @key value: (_) @value) @root`)
	if err != nil {
		return nil, fmt.Errorf("error compiling pair query: %s", err)
	}

	// function getName(){}
	functionNameQuery, err := lang.CompileQuery(`(function_declaration name: (identifier) @name) @root`)
	if err != nil {
		return nil, fmt.Errorf("error compiling class method query: %s", err)
	}

	// class User {
	//	constructor(name, surname)
	//	GetName()
	// }
	methodNameQuery, err := lang.CompileQuery(`(method_definition name: (property_identifier) @name (formal_parameters) @params ) @root`)
	if err != nil {
		return nil, fmt.Errorf("error compiling class method query: %s", err)
	}

	return &propertyDetector{
		pairQuery:         pairQuery,
		functionNameQuery: functionNameQuery,
		methodNameQuery:   methodNameQuery,
	}, nil
}

func (detector *propertyDetector) Name() string {
	return "property"
}

func (detector *propertyDetector) DetectAt(
	node *tree.Node,
	evaluator types.Evaluator,
) ([]interface{}, error) {
	// run pair query
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

	// run function name query
	result, err = detector.functionNameQuery.MatchOnceAt(node)
	if err != nil {
		return nil, err
	}
	if len(result) != 0 {
		return []interface{}{generictypes.Property{Name: result["name"].Content()}}, nil
	}

	return detector.getMethod(node, evaluator)
}

func (detector *propertyDetector) getMethod(
	node *tree.Node,
	evaluator types.Evaluator,
) ([]interface{}, error) {
	// run method query
	result, err := detector.methodNameQuery.MatchOnceAt(node)
	if err != nil || len(result) == 0 {
		return nil, err
	}

	// fetch all arguments from constructor
	if result["name"].Type() == "constructor" {
		properties := []interface{}{}

		params := result["params"]
		for i := 0; i < params.ChildCount(); i++ {
			param := params.Child(i)
			if param.Type() != "required_parameter" {
				continue
			}

			identifier := param.ChildByFieldName("pattern")
			if identifier == nil || identifier.Type() != "identifier" {
				continue
			}

			properties = append(properties, generictypes.Property{Name: result["name"].Content()})
		}

		return properties, nil
	}

	return []interface{}{generictypes.Property{Name: result["name"].Content()}}, nil

}

func (detector *propertyDetector) Close() {
	detector.pairQuery.Close()
}
