package knex

import (
	"github.com/bearer/curio/pkg/parser"
	"github.com/bearer/curio/pkg/report"
	"github.com/bearer/curio/pkg/report/detectors"
	reportknex "github.com/bearer/curio/pkg/report/frameworks/knex"
	sitter "github.com/smacker/go-tree-sitter"
)

var queryModule = `(ambient_declaration
	(module
    	name: (string) @param_module_name
    ) @param_module
)`

const knexTablesModule = "knex/types/tables"

var queryTable = `(interface_declaration
	name: (type_identifier) @helper_Tables
	body: (object_type) @param_table_body
)`

func detectTableDeclarationModule(report report.Report, tree *parser.Tree, language *sitter.Language) {
	compiledQueryModule := parser.QueryMustCompile(language, queryModule)

	captures := tree.QueryConventional(compiledQueryModule)

	var knexModule *parser.Node

	for _, capture := range captures {
		moduleName := capture["param_module_name"]
		if moduleName.Content() != knexTablesModule {
			knexModule = capture["param_module"]
			break
		}
	}

	if knexModule == nil {
		return
	}

	compiledQueryTable := parser.QueryMustCompile(language, queryTable)

	captures = tree.QueryConventional(compiledQueryTable)

	var tableNode *parser.Node

	for _, capture := range captures {
		table := capture["param_table_body"]

		isDescendant := parser.IsDescendant(table, knexModule)
		if !isDescendant {
			continue
		}

		tableNode = table
		break
	}

	if tableNode == nil {
		return
	}

	for i := 0; i < tableNode.ChildCount(); i++ {
		tableChild := tableNode.Child(i)
		if tableChild.Type() == "property_signature" {
			propertyName := tableChild.ChildByFieldName("name")
			propertyType := tableChild.ChildByFieldName("type")

			if propertyName == nil || propertyType == nil {
				continue
			}

			if propertyType.Type() != "type_annotation" {
				continue
			}

			report.AddFramework(detectors.DetectorTypescript, reportknex.TypeSchema, reportknex.Schema{
				DataType:     propertyType.Content(),
				PropertyName: propertyName.Content(),
			}, tableChild.Source(false))
		}
	}
}
