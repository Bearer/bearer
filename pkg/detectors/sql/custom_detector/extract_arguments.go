package customdetector

import (
	"github.com/bearer/curio/pkg/detectors/sql/util"
	"github.com/bearer/curio/pkg/parser"
	parserdatatype "github.com/bearer/curio/pkg/parser/datatype"
	"github.com/bearer/curio/pkg/parser/nodeid"
	"github.com/bearer/curio/pkg/report/schema"
)

func (detector *Detector) ExtractArguments(node *parser.Node, idGenerator nodeid.Generator) (map[parser.NodeID]*parserdatatype.DataType, error) {
	if node == nil {
		return nil, nil
	}

	joinedDatatypes := make(map[parser.NodeID]*parserdatatype.DataType)

	if node.Type() == "identifier" && node.Parent() != nil && node.Parent().Type() == "table_column" {
		parent := node.Parent()
		typeNode := parent.ChildByFieldName("type")
		typeIdentifierNode := typeNode.Child(0)

		simpleType := util.ConvertToSimpleType(typeIdentifierNode.Content())

		datatype := &parserdatatype.DataType{
			Node:       node,
			Name:       node.Content(),
			Type:       simpleType,
			TextType:   typeIdentifierNode.Content(),
			Properties: make(map[string]*parserdatatype.DataType),
		}

		joinedDatatypes[datatype.Node.ID()] = datatype
		return joinedDatatypes, nil
	}

	if node.Type() == "identifier" && node.Parent() != nil && node.Parent().Type() == "create_function_statement" {
		datatype := &parserdatatype.DataType{
			Node:       node,
			Name:       node.Content(),
			Type:       schema.SimpleTypeObject,
			TextType:   "",
			Properties: make(map[string]*parserdatatype.DataType),
		}

		joinedDatatypes[datatype.Node.ID()] = datatype
		return joinedDatatypes, nil
	}

	// handle generic datatype extraction
	if node.Type() == "identifier" {
		tableNameNode := node

		tableNameDatatype := &parserdatatype.DataType{
			Node:       tableNameNode,
			Name:       tableNameNode.Content(),
			Type:       schema.SimpleTypeObject,
			Properties: make(map[string]*parserdatatype.DataType),
		}

		joinedDatatypes[tableNameDatatype.Node.ID()] = tableNameDatatype

		return joinedDatatypes, nil
	}

	return joinedDatatypes, nil
}
