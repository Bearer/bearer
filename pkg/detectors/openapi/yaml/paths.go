package yaml

import (
	"github.com/bearer/bearer/pkg/detectors/openapi/queries"
	"github.com/bearer/bearer/pkg/parser"
	"github.com/bearer/bearer/pkg/report/operations/operationshelper"
	"github.com/smacker/go-tree-sitter/yaml"
)

var queryPaths = parser.QueryMustCompile(yaml.GetLanguage(), `
(block_mapping
	(block_mapping_pair
    	key: (flow_node) @helper_paths
        (#match? @helper_paths "^paths$")
        value:
        	(block_node
            	(block_mapping
                	(block_mapping_pair
                    	key: (flow_node) @param_path
                        value: (
                        	block_node
                        		(block_mapping
                                	(block_mapping_pair
                                    	key: (flow_node) @param_request_type
                                    )
                                )
                        )
                    )
                )
            )
    )
)`)

func AnnotatePaths(tree *parser.Tree, foundValues map[parser.Node]*operationshelper.Operation) error {
	return queries.AnnotatePaths(queries.PathsRequest{
		Tree:        tree,
		Query:       queryPaths,
		FoundValues: foundValues,
	})
}
