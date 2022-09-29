package util

import (
	"github.com/bearer/curio/pkg/parser"
	"github.com/bearer/curio/pkg/util/stringutil"
	sitter "github.com/smacker/go-tree-sitter"
)

var (
	importsQuery = `
	(import_statement
		source: (string) @param_source
	)
`
)

type Import struct {
	Name     string
	Alias    string
	IsMaster bool
	IsRoot   bool
}

// handled typescript imports
// import { knexMember as knexAlias, knexMember2 as knexAlias2 } from "knex";
// import { knexMember3 as knexAlias3 } from "knex";
// import { knexMember3 } from "knex";
// import * as knex from 'knex';
// import knex from 'knex'; // master import
// import 'knex'; // root import
func GetImports(tree *parser.Tree, language *sitter.Language, matchingImports []string) []Import {
	query := parser.QueryMustCompile(language, importsQuery)

	var imports []Import

	captures := tree.QueryConventional(query)

	for _, capture := range captures {
		importSource := stringutil.StripQuotes(capture["param_source"].Content())
		matches := false
		for _, importMatch := range matchingImports {
			if importSource == importMatch {
				matches = true
				break
			}
		}

		if !matches {
			continue
		}

		if capture["param_source"].Parent().ChildCount() == 1 {
			imports = append(imports, Import{
				Name:     "",
				IsMaster: false,
				IsRoot:   true,
			})
			continue
		}

		importClauseNode := capture["param_source"].Parent().Child(0)

		if importClauseNode.Type() != "import_clause" {
			continue
		}

		for i := 0; i < importClauseNode.ChildCount(); i++ {
			child := importClauseNode.Child(i)
			switch child.Type() {
			case "identifier":
				imports = append(imports, Import{
					Name:     child.Content(),
					Alias:    child.Content(),
					IsMaster: true,
				})
			case "namespace_import":
				imports = append(imports, Import{
					Name:     child.Child(0).Content(),
					Alias:    child.Child(0).Content(),
					IsMaster: true,
				})
			case "named_imports":
				for i := 0; i < child.ChildCount(); i++ {
					importSpecifier := child.Child(i)
					if importSpecifier.Type() != "import_specifier" {
						continue
					}

					newImport := Import{}

					name := importSpecifier.ChildByFieldName("name")
					if name != nil {
						newImport.Name = name.Content()
						newImport.Alias = name.Content()
					}

					alias := importSpecifier.ChildByFieldName("alias")
					if alias != nil {
						newImport.Alias = alias.Content()
					}

					imports = append(imports, newImport)
				}
			}
		}

	}

	return imports
}
