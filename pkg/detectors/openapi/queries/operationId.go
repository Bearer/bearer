package queries

import (
	"github.com/bearer/curio/pkg/parser"
	"github.com/bearer/curio/pkg/parser/nodeid"
	"github.com/bearer/curio/pkg/report/schema/schemahelper"
	"github.com/bearer/curio/pkg/util/stringutil"
	sitter "github.com/smacker/go-tree-sitter"
)

type OperationIdRequest struct {
	Tree        *parser.Tree
	FoundValues map[parser.Node]*schemahelper.Schema
	Query       *sitter.Query
	ChildMatch  ChildMatch
	NodeIDMap   *nodeid.Map
}

func AnnotateOperationId(request OperationIdRequest) error {
	captures := request.Tree.QueryMustPass(request.Query)

	for _, capture := range captures {

		if capture["param_operation_id"] == nil || capture["param_parameters"] == nil {
			continue
		}

		if stringutil.StripQuotes(capture["helperOperationId"].Content()) != "operationId" ||
			stringutil.StripQuotes(capture["helperParameters"].Content()) != "parameters" {
			continue
		}

		operationNode := capture["param_operation_id"]
		paramsNode := capture["param_parameters"]

		for i := 0; i < paramsNode.ChildCount(); i++ {

			paramsObject := request.ChildMatch.Match(paramsNode.Child(i))

			if paramsObject == nil {
				continue
			}

			for j := 0; j < paramsObject.ChildCount(); j++ {

				keyValue := paramsObject.Child(j)
				key := keyValue.ChildByFieldName("key")

				if key != nil && stringutil.StripQuotes(key.Content()) == "name" {
					value := keyValue.ChildByFieldName("value")

					if value == nil {
						continue
					}

					if existingValue, ok := request.FoundValues[*value]; ok {
						existingValue.Value.ObjectName = operationNode.Content()
						existingValue.Value.ObjectUUID = request.NodeIDMap.ValueForNode(operationNode.ID())
					}
				}
			}
		}

	}

	return nil
}
