package queries

import (
	"github.com/bearer/curio/pkg/parser"
	"github.com/bearer/curio/pkg/report/operations"
	"github.com/bearer/curio/pkg/report/operations/operationshelper"
	"github.com/bearer/curio/pkg/util/stringutil"
	sitter "github.com/smacker/go-tree-sitter"
)

type PathsRequest struct {
	Tree        *parser.Tree
	Query       *sitter.Query
	FoundValues map[parser.Node]*operationshelper.Operation
}

func AnnotatePaths(request PathsRequest) error {
	captures := request.Tree.QueryMustPass(request.Query)

	for _, capture := range captures {

		if stringutil.StripQuotes(capture["helper_paths"].Content()) != "paths" {
			continue
		}

		path := capture["param_path"]
		requestType := capture["param_request_type"]

		request.FoundValues[*requestType] = &operationshelper.Operation{
			Source: requestType.Source(true),
			Value: operations.Operation{
				Path: stringutil.StripQuotes(path.Content()),
				Type: stringutil.StripQuotes(requestType.Content()),
			},
		}
	}

	return nil
}
