package queries

import (
	"github.com/bearer/bearer/pkg/parser"
	"github.com/bearer/bearer/pkg/parser/nodeid"
	"github.com/bearer/bearer/pkg/report/schema"
	"github.com/bearer/bearer/pkg/report/schema/schemahelper"
	"github.com/bearer/bearer/pkg/util/stringutil"
	sitter "github.com/smacker/go-tree-sitter"
)

func AnnotateV3Paramaters(
	nodeIDMap *nodeid.Map,
	tree *parser.Tree,
	foundValues map[parser.Node]*schemahelper.Schema,
	queryParameters *sitter.Query,
) error {
	captures := tree.QueryMustPass(queryParameters)

	for _, capture := range captures {
		if capture["param_name"] == nil || capture["param_schema"] == nil {
			continue
		}

		if stringutil.StripQuotes(capture["helperName"].Content()) != "name" ||
			stringutil.StripQuotes(capture["helperSchema"].Content()) != "schema" {
			continue
		}

		nameNode := capture["param_name"]
		schemaNode := capture["param_schema"]

		fieldType := ""

		var typeNode *parser.Node

		for i := 0; i < schemaNode.ChildCount(); i++ {
			objectMapping := schemaNode.Child(i)
			key := objectMapping.ChildByFieldName("key")

			if key != nil && stringutil.StripQuotes(key.Content()) == "type" {
				typeNode = objectMapping.ChildByFieldName("value")
			}
		}

		if typeNode != nil {
			fieldType = typeNode.Content()
		}

		foundValues[*nameNode] = &schemahelper.Schema{
			Source: nameNode.Source(true),
			Value: schema.Schema{
				FieldName: nameNode.Content(),
				FieldType: fieldType,
				FieldUUID: nodeIDMap.ValueForNode(nameNode.ID()),
			},
		}
	}

	return nil
}
