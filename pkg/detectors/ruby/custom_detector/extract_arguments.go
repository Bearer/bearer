package customdetector

import (
	"strings"

	"github.com/bearer/curio/pkg/detectors/ruby/datatype"
	"github.com/bearer/curio/pkg/parser"
	parserdatatype "github.com/bearer/curio/pkg/parser/datatype"
	"github.com/bearer/curio/pkg/parser/nodeid"
	"github.com/bearer/curio/pkg/report/schema"
	"github.com/bearer/curio/pkg/util/file"
	"github.com/smacker/go-tree-sitter/ruby"
)

func (detector *Detector) ExtractArguments(node *parser.Node, idGenerator nodeid.Generator) (map[parser.NodeID]*parserdatatype.DataType, error) {
	if node == nil {
		return nil, nil
	}

	joinedDatatypes := make(map[parser.NodeID]*parserdatatype.DataType)

	// handle classs name
	if node.Type() == "constant" {
		datatype := &parserdatatype.DataType{
			Node:       node,
			Name:       node.Content(),
			Type:       schema.SimpleTypeObject,
			Properties: make(map[string]*parserdatatype.DataType),
		}
		joinedDatatypes[datatype.Node.ID()] = datatype
		return joinedDatatypes, nil
	}

	for i := 0; i < node.ChildCount(); i++ {
		singleArgument := node.Child(i)

		if singleArgument.Type() == "identifier" || singleArgument.Type() == "simple_symbol" || singleArgument.Type() == "bare_symbol" {

			content := singleArgument.Content()

			if singleArgument.Type() == "simple_symbol" {
				content = strings.TrimLeft(content, ":")
			}

			datatype := &parserdatatype.DataType{
				Node:       singleArgument,
				Name:       content,
				Type:       schema.SimpleTypeUknown,
				Properties: make(map[string]*parserdatatype.DataType),
			}
			joinedDatatypes[datatype.Node.ID()] = datatype
			continue
		}

		content := singleArgument.Content()
		tree, err := parser.ParseBytes(&file.FileInfo{}, &file.Path{}, []byte(content), ruby.GetLanguage(), singleArgument.LineNumber()-1)
		if err != nil {
			return nil, err
		}

		singleArgumentDatatypes := datatype.Discover(tree, idGenerator)

		for nodeID, datatype := range singleArgumentDatatypes {
			joinedDatatypes[nodeID] = datatype
		}
	}

	return joinedDatatypes, nil
}
