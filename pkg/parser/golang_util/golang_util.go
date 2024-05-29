package golang_util

import (
	"regexp"
	"strings"

	"github.com/smacker/go-tree-sitter/golang"

	"github.com/bearer/bearer/pkg/parser"
)

var (
	importsQuery = parser.QueryMustCompile(golang.GetLanguage(), `
		(import_spec
			name: (_)? @name
			path: (_) @path)
	`)

	/*
	 * github.com/example/mypkg => mypkg
	 * github.com/example/mypkg.v2 => mypkg
	 * github.com/example/mypkg/v3 => mypkg
	 */
	nameFromPathRegex = regexp.MustCompile(`([^/]+)((\.|/)v\d+)?$`)
)

func GetImports(tree *parser.Tree) (map[string]string, error) {
	imports := make(map[string]string)

	err := tree.Query(importsQuery, func(captures parser.Captures) error {
		nameNode := captures["name"]
		path := stripQuotes(captures["path"].Content())

		name := ""
		if nameNode != nil {
			name = nameNode.Content()
		}

		if name == "_" {
			return nil
		}

		if name == "" {
			name = nameFromPathRegex.FindStringSubmatch(path)[1]
		}

		imports[name] = path

		return nil
	})

	return imports, err
}

func AliasesFor(imports map[string]string, packageName string) []string {
	var result []string

	for alias, candidatePackage := range imports {
		if candidatePackage == packageName {
			result = append(result, alias)
		}
	}

	return result
}

func stripQuotes(value string) string {
	return strings.Trim(value, "\"`")
}
