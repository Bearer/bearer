package property

import (
	"fmt"

	"github.com/bearer/bearer/new/detector/types"
	"github.com/bearer/bearer/new/language/tree"
	soufflequery "github.com/bearer/bearer/pkg/souffle/query"

	generictypes "github.com/bearer/bearer/new/detector/implementation/generic/types"
	"github.com/bearer/bearer/new/detector/implementation/ruby/common"
	languagetypes "github.com/bearer/bearer/new/language/types"
)

type propertyDetector struct {
	types.DetectorBase
	pairQuery       *tree.Query
	accessorQuery   *tree.Query
	methodNameQuery *tree.Query
}

func New(lang languagetypes.Language) (types.Detector, error) {
	// { user: "admin_user" }
	pairQuery, err := lang.CompileQuery(`(pair key: (_) @key value: (_) @value) @root`)
	if err != nil {
		return nil, fmt.Errorf("error compiling pair query: %s", err)
	}

	// 	attr_accessor :name
	accessorQuery, err := lang.CompileQuery(`(call arguments: (argument_list (simple_symbol) @name ))@root`)
	if err != nil {
		return nil, fmt.Errorf("error compiling class accessor query: %s", err)
	}

	// def get_first_name()
	// end
	methodNameQuery, err := lang.CompileQuery(`(method name: (identifier) @name) @root`)
	if err != nil {
		return nil, fmt.Errorf("error compiling class method query: %s", err)
	}

	return &propertyDetector{
		pairQuery:       pairQuery,
		accessorQuery:   accessorQuery,
		methodNameQuery: methodNameQuery,
	}, nil
}

func (detector *propertyDetector) Name() string {
	return "property"
}

func (detector *propertyDetector) DetectAt(
	node *tree.Node,
	evaluator types.Evaluator,
	queryContext *soufflequery.QueryContext,
) ([]interface{}, error) {
	// run hash pair query
	result, err := detector.pairQuery.MatchOnceAt(node)
	if err != nil {
		return nil, err
	}

	if result != nil {
		if key := common.GetLiteralKey(result["key"]); key != "" {
			return []interface{}{
				generictypes.Property{Name: key},
			}, nil
		}
	}

	// run accessor query
	results, err := detector.accessorQuery.MatchAt(node)
	if err != nil {
		return nil, err
	}
	if len(results) != 0 {
		properties := make([]interface{}, len(results))

		for i, result := range results {
			properties[i] = generictypes.Property{Name: result["name"].Content()[1:]}
		}

		return properties, nil
	}

	// run method name query
	result, err = detector.methodNameQuery.MatchOnceAt(node)
	if err != nil {
		return nil, err
	}
	if result != nil {
		return []interface{}{
			generictypes.Property{Name: result["name"].Content()},
		}, nil
	}

	return nil, nil
}

func (detector *propertyDetector) Close() {
	detector.pairQuery.Close()
	detector.accessorQuery.Close()
	detector.methodNameQuery.Close()
}
