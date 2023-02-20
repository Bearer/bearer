package knex

import (
	"github.com/bearer/bearer/pkg/detectors/javascript/util"
	"github.com/bearer/bearer/pkg/parser"
	"github.com/bearer/bearer/pkg/report"
	sitter "github.com/smacker/go-tree-sitter"
)

func Discover(report report.Report, tree *parser.Tree, language *sitter.Language) {
	knexImports := util.GetImports(tree, language, []string{"knex"})

	if len(knexImports) == 0 {
		return
	}

	detectFunctionTypes(report, tree, language, knexImports)
	detectTableDeclarationModule(report, tree, language)
}
