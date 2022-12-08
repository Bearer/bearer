package customdetector

import (
	"strings"

	"github.com/bearer/curio/pkg/detectors/ruby/datatype"
	"github.com/bearer/curio/pkg/parser"
	parserdatatype "github.com/bearer/curio/pkg/parser/datatype"
	"github.com/bearer/curio/pkg/parser/nodeid"
	"github.com/bearer/curio/pkg/report/schema"
	"github.com/bearer/curio/pkg/util/stringutil"

	schemadatatype "github.com/bearer/curio/pkg/report/schema/datatype"
)

func (detector *Detector) ExtractArguments(node *parser.Node, idGenerator nodeid.Generator, variableReconciliation *parserdatatype.ReconciliationRequest) (map[parser.NodeID]*schemadatatype.DataType, error) {
	extractedDatatypes, err := detector.extractArguments(node, idGenerator)
	if err != nil {
		return nil, err
	}

	if variableReconciliation != nil {
		parserdatatype.VariableReconciliation(extractedDatatypes, variableReconciliation)
	}

	return extractedDatatypes, nil
}

func (detector *Detector) extractArguments(node *parser.Node, idGenerator nodeid.Generator) (map[parser.NodeID]*schemadatatype.DataType, error) {
	if node == nil {
		return nil, nil
	}

	joinedDatatypes := make(map[parser.NodeID]*schemadatatype.DataType)

	// handle class name
	if node.Type() == "constant" || node.Type() == "identifier" || node.Type() == "simple_symbol" || node.Type() == "bare_symbol" {
		datatype := &schemadatatype.DataType{
			Node:       node,
			Name:       node.Content(),
			Type:       schema.SimpleTypeObject,
			Properties: make(map[string]schemadatatype.DataTypable),
		}
		joinedDatatypes[datatype.Node.ID()] = datatype
		return joinedDatatypes, nil
	}

	if node.Type() == "hash" {
		extractHash(joinedDatatypes, node)
		return joinedDatatypes, nil
	}

	if node.Type() == "argument_list" {
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
					Type:       schema.SimpleTypeUnknown,
					Properties: make(map[string]schemadatatype.DataTypable),
				}
				joinedDatatypes[datatype.Node.ID()] = datatype
				continue
			}

			if singleArgument.Type() == "hash" {
				extractHash(joinedDatatypes, singleArgument)
				continue
			}
		}

		complexDatatypes := datatype.Discover(node, idGenerator)
		for nodeID, target := range complexDatatypes {
			_, exists := joinedDatatypes[nodeID]
			if exists {
				continue
			}

			joinedDatatypes[nodeID] = target
		}

		return joinedDatatypes, nil
	}

	complexDatatypes := datatype.Discover(node.Parent(), idGenerator)
	for nodeID, target := range complexDatatypes {
		_, exists := joinedDatatypes[nodeID]
		if exists {
			continue
		}

		if parserdatatype.IsParentedByNodeID(node.ID(), target.Node) {
			joinedDatatypes[nodeID] = target
		}
	}

	return joinedDatatypes, nil

}

func extractHash(datatypes map[parser.NodeID]*schemadatatype.DataType, node *parser.Node) {
	var parentDatatype *schemadatatype.DataType
	for i := 0; i < node.ChildCount(); i++ {
		pair := node.Child(i)

		if pair.Type() != "pair" {
			continue
		}

		key := pair.ChildByFieldName("key")

		if key == nil {
			continue
		}

		if parentDatatype == nil {
			parentDatatype = &schemadatatype.DataType{
				Node:       node,
				Name:       "",
				Type:       schema.SimpleTypeObject,
				Properties: make(map[string]schemadatatype.DataTypable),
			}
		}

		keyName := stringutil.StripQuotes(key.Content())

		parentDatatype.Properties[keyName] = &schemadatatype.DataType{
			Node:       pair,
			Name:       keyName,
			Type:       schema.SimpleTypeUnknown,
			Properties: make(map[string]schemadatatype.DataTypable),
		}
	}

	if parentDatatype != nil {
		datatypes[parentDatatype.Node.ID()] = parentDatatype
	}
}
