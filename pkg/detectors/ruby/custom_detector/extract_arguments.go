package customdetector

import (
	"strings"

	"github.com/bearer/curio/pkg/detectors/ruby/datatype"
	"github.com/bearer/curio/pkg/parser"
	parserdatatype "github.com/bearer/curio/pkg/parser/datatype"
	"github.com/bearer/curio/pkg/parser/nodeid"
	"github.com/bearer/curio/pkg/report/schema"

	schemadatatype "github.com/bearer/curio/pkg/report/schema/datatype"
	"github.com/bearer/curio/pkg/util/file"
)

func (detector *Detector) ExtractArguments(node *parser.Node, idGenerator nodeid.Generator, fileinfo *file.FileInfo, filepath *file.Path) (map[parser.NodeID]*schemadatatype.DataType, error) {
	if node == nil {
		return nil, nil
	}

	joinedDatatypes := make(map[parser.NodeID]*schemadatatype.DataType)

	// handle classs name
	if node.Type() == "constant" {
		datatype := &schemadatatype.DataType{
			Node:       node,
			Name:       node.Content(),
			Type:       schema.SimpleTypeObject,
			Properties: make(map[string]schemadatatype.DataTypable),
		}
		joinedDatatypes[datatype.Node.ID()] = datatype
		return joinedDatatypes, nil
	}

	if node.ChildCount() == 0 {
		return joinedDatatypes, nil
	}

	allDatatypes := datatype.Discover(node.Tree().RootNode(), idGenerator)

	for i := 0; i < node.ChildCount(); i++ {
		singleArgument := node.Child(i)

		if singleArgument.Type() == "identifier" || singleArgument.Type() == "simple_symbol" || singleArgument.Type() == "bare_symbol" {

			content := singleArgument.Content()

			if singleArgument.Type() == "simple_symbol" {
				content = strings.TrimLeft(content, ":")
			}

			datatype := &schemadatatype.DataType{
				Node:       singleArgument,
				Name:       content,
				Type:       schema.SimpleTypeUknown,
				Properties: make(map[string]schemadatatype.DataTypable),
			}
			joinedDatatypes[datatype.Node.ID()] = datatype
			continue
		}

		singleArgumentDatatypes := datatype.Discover(singleArgument, idGenerator)
		parserdatatype.VariableReconciliation(singleArgumentDatatypes, allDatatypes, datatype.ScopeTerminators)

		for nodeID, target := range singleArgumentDatatypes {
			joinedDatatypes[nodeID] = target
		}

	}

	return joinedDatatypes, nil
}
