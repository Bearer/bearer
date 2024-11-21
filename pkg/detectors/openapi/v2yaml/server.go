package v2yaml

import (
	"strings"

	"github.com/bearer/bearer/pkg/parser"
	"github.com/bearer/bearer/pkg/report/operations"
	"github.com/bearer/bearer/pkg/util/stringutil"
	"github.com/smacker/go-tree-sitter/yaml"
)

var queryServersHost = parser.QueryMustCompile(yaml.GetLanguage(), `
(document
	(block_node
    	(block_mapping
        	(block_mapping_pair
              key: (flow_node) @helper_host
			  value: (flow_node) @param_host
            )
        )
    )
)
`)

var queryServersBasePath = parser.QueryMustCompile(yaml.GetLanguage(), `
(document
	(block_node
    	(block_mapping
        	(block_mapping_pair
              key: (flow_node) @helper_basePath
			  value: (flow_node) @param_base_path
            )
        )
    )
)
`)

var queryServerSchemes = parser.QueryMustCompile(yaml.GetLanguage(), `
(document
	(block_node
    	(block_mapping
        	(block_mapping_pair
              key: (flow_node) @helper_schemes
              value:
              	(block_node
                	(block_sequence
                    	(block_sequence_item
                        	(flow_node) @param_scheme
                        )
                    )
                )
            )
        )
    )
)
`)

func findServers(tree *parser.Tree) (result []operations.Url) {
	host := ""
	captures := tree.QueryConventional(queryServersHost)
	for _, capture := range captures {
		host = stringutil.StripQuotes(capture["param_host"].Content())
	}

	basePath := ""
	captures = tree.QueryConventional(queryServersBasePath)
	for _, capture := range captures {
		basePath = stringutil.StripQuotes(capture["param_base_path"].Content())
	}

	schemes := make([]string, 0)
	captures = tree.QueryConventional(queryServerSchemes)
	for _, capture := range captures {
		schemes = append(schemes, stringutil.StripQuotes(capture["param_scheme"].Content()))
	}

	if len(schemes) == 0 {
		schemes = append(schemes, "http://", "https://")
	}

	for _, scheme := range schemes {
		if !strings.HasSuffix(scheme, "://") {
			scheme = scheme + "://"
		}

		if !strings.HasPrefix(basePath, "/") {
			basePath = "/" + basePath
		}
		result = append(result, operations.Url{
			Url: scheme + host + basePath,
		})
	}

	return result
}
