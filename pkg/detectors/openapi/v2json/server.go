package v2json

import (
	"strings"

	"github.com/bearer/bearer/pkg/parser"
	"github.com/bearer/bearer/pkg/report/operations"
	"github.com/bearer/bearer/pkg/util/stringutil"
	"github.com/smacker/go-tree-sitter/javascript"
)

var queryServersHost = parser.QueryMustCompile(javascript.GetLanguage(), `
(expression_statement
	(object
		(pair
			key: (string) @helper_host
			(#match? @helper_host "^\"host\"$")
			value: (string) @param_host
		)
	)
)
`)

var queryServersBasePath = parser.QueryMustCompile(javascript.GetLanguage(), `
(expression_statement
	(object
		(pair
			key: (string) @helper_basePath
			(#match? @helper_basePath "^\"basePath\"$")
			value: (string) @param_base_path
		)
	)
)
`)

var queryServerSchemes = parser.QueryMustCompile(javascript.GetLanguage(), `
(expression_statement
	(object
		(pair
			key: (string) @helper_schemes
			(#match? @helper_schemes "^\"schemes\"$")
			value: (array
				(string) @param_schemes
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
		schemes = append(schemes, stringutil.StripQuotes(capture["param_schemes"].Content()))
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
