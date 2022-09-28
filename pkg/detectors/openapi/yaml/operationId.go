package yaml

import (
	"github.com/bearer/curio/pkg/detectors/openapi/queries"
	"github.com/bearer/curio/pkg/parser"
	"github.com/bearer/curio/pkg/parser/nodeid"
	"github.com/bearer/curio/pkg/report/schema/schemahelper"
	"github.com/smacker/go-tree-sitter/yaml"
)

var queryOperationId = parser.QueryMustCompile(yaml.GetLanguage(), `
(_
	(
      block_mapping_pair
        key:
            (flow_node) @helperOperationId
            (#match? @helperOperationId "^operationId$")
         value:
            (flow_node) @param_operation_id
	)
	(
      block_mapping_pair
      
        key:
            (flow_node) @helperParameters
            (#match? @helperParameters "^parameters$")
         value:
          (block_node (block_sequence)  @param_parameters)
	)
)
`)

type OperationIdChildMatcher struct {
}

func (childMatch OperationIdChildMatcher) Match(input *parser.Node) *parser.Node {
	if input == nil || input.ChildCount() == 0 {
		return nil
	}

	firstChild := input.Child(0)

	if firstChild == nil {
		return nil
	}

	return firstChild.Child(0)
}

func AnnotateOperationId(nodeIDMap *nodeid.Map, tree *parser.Tree, foundValues map[parser.Node]*schemahelper.Schema) error {
	return queries.AnnotateOperationId(queries.OperationIdRequest{
		Tree:        tree,
		FoundValues: foundValues,
		Query:       queryOperationId,
		ChildMatch:  OperationIdChildMatcher{},
		NodeIDMap:   nodeIDMap,
	})
}
