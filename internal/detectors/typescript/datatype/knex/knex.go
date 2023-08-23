package knex

import (
	"github.com/bearer/bearer/internal/detectors/javascript/util"
	"github.com/bearer/bearer/internal/parser"
	"github.com/bearer/bearer/internal/report"
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
