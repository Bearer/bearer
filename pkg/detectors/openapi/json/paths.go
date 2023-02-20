package json

import (
	"github.com/bearer/bearer/pkg/detectors/openapi/queries"
	"github.com/bearer/bearer/pkg/parser"
	"github.com/bearer/bearer/pkg/report/operations/operationshelper"
	"github.com/smacker/go-tree-sitter/javascript"
)

var queryPaths = parser.QueryMustCompile(javascript.GetLanguage(), `
(object
	(pair
    	key: (string) @helper_paths
        (#match? @helper_paths "^\"paths\"$")
        value:
        	(object
            	(pair
                	key: (string) @param_path
					value:
                    	(object
                        	(pair
                            	key:
                                	(string) @param_request_type
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
