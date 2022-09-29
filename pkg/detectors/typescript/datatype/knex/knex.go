package knex

import (
	"github.com/bearer/curio/pkg/detectors/javascript/util"
	"github.com/bearer/curio/pkg/parser"
	"github.com/bearer/curio/pkg/report"
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
