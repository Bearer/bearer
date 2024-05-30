package yaml

import (
	"github.com/bearer/bearer/pkg/detectors/openapi/queries"
	"github.com/bearer/bearer/pkg/parser"
	"github.com/bearer/bearer/pkg/parser/nodeid"
	"github.com/bearer/bearer/pkg/report/schema/schemahelper"
	"github.com/smacker/go-tree-sitter/yaml"
)

var queryObjects = parser.QueryMustCompile(yaml.GetLanguage(), `
(_
	(
      block_mapping_pair
        key:
            (flow_node) @param_object_name
         value:
         	(block_node
            	(block_mapping
                	(block_mapping_pair
                        key:
                            (flow_node) @helperProperties
                            (#match? @helperProperties "^properties$")
                         value:
                            (block_node (block_mapping) @param_object_properties)
                    )
                )
            )
	)
)
`)

type ObjectChildMatcher struct {
}

func (childMatches ObjectChildMatcher) Match(input *parser.Node) *parser.Node {
	if input == nil || input.ChildCount() == 0 {
		return nil
	}

	return input.Child(0)
}

func AnnotateObjects(nodeIDMap *nodeid.Map, tree *parser.Tree, foundValues map[parser.Node]*schemahelper.Schema) error {
	return queries.AnnotateObjects(queries.ObjectsRequest{
		Tree:        tree,
		Query:       queryObjects,
		FoundValues: foundValues,
		ChildMatch:  ObjectChildMatcher{},
		NodeIDMap:   nodeIDMap,
	})
}
