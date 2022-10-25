package customdetector

import (
	"github.com/bearer/curio/pkg/detectors/sql/util"
	"github.com/bearer/curio/pkg/parser"
	"github.com/bearer/curio/pkg/parser/nodeid"
	"github.com/bearer/curio/pkg/report/schema"
	schemadatatype "github.com/bearer/curio/pkg/report/schema/datatype"
)

func (detector *Detector) ExtractArguments(node *parser.Node, idGenerator nodeid.Generator) (map[parser.NodeID]schemadatatype.DataTypable, error) {
	if node == nil {
		return nil, nil
	}

	joinedDatatypes := make(map[parser.NodeID]schemadatatype.DataTypable)

	if node.Type() == "identifier" && node.Parent() != nil && node.Parent().Type() == "table_column" {
		parent := node.Parent()
		typeNode := parent.ChildByFieldName("type")
		typeIdentifierNode := typeNode.Child(0)

		simpleType := util.ConvertToSimpleType(typeIdentifierNode.Content())

		datatype := &schemadatatype.DataType{
			Node:       node,
			Name:       node.Content(),
			Type:       simpleType,
			TextType:   typeIdentifierNode.Content(),
			Properties: make(map[string]schemadatatype.DataTypable),
		}

		joinedDatatypes[datatype.Node.ID()] = datatype
		return joinedDatatypes, nil
	}

	if node.Type() == "identifier" && node.Parent() != nil && node.Parent().Type() == "create_function_statement" {
		datatype := &schemadatatype.DataType{
			Node:       node,
			Name:       node.Content(),
			Type:       schema.SimpleTypeObject,
			TextType:   "",
			Properties: make(map[string]schemadatatype.DataTypable),
		}

		joinedDatatypes[datatype.Node.ID()] = datatype
		return joinedDatatypes, nil
	}

	// handle generic datatype extraction
	if node.Type() == "identifier" {
		tableNameNode := node

		tableNameDatatype := &schemadatatype.DataType{
			Node:       tableNameNode,
			Name:       tableNameNode.Content(),
			Type:       schema.SimpleTypeObject,
			Properties: make(map[string]schemadatatype.DataTypable),
		}

		joinedDatatypes[tableNameDatatype.Node.ID()] = tableNameDatatype

		return joinedDatatypes, nil
	}

	return joinedDatatypes, nil
}
