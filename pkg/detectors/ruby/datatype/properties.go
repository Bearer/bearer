package datatype

import (
	"strings"

	"github.com/bearer/curio/pkg/parser"
	parserdatatype "github.com/bearer/curio/pkg/parser/datatype"
	"github.com/bearer/curio/pkg/parser/nodeid"
	"github.com/bearer/curio/pkg/report/schema"
	schemadatatype "github.com/bearer/curio/pkg/report/schema/datatype"
	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/ruby"
)

// person[:city][:number]()
var elementSimpleSymbolQuery = parser.QueryMustCompile(ruby.GetLanguage(),
	`(element_reference
		object: (_) @param_object
		(simple_symbol) @param_id
	) @param_parent`)

var elementIdentifierQuery = parser.QueryMustCompile(ruby.GetLanguage(),
	`(element_reference
		object: (_) @param_object
		(identifier) @param_id
	) @param_parent`)

// person.city.number()
var callsQuery = parser.QueryMustCompile(ruby.GetLanguage(),
	`(call
		receiver: (_) @param_receiver
		method: (identifier) @param_id
	) @param_parent`)

func addProperties(tree *parser.Tree, helperDatatypes map[parser.NodeID]*schemadatatype.DataType) {
	// add element references
	var doElementsQuery = func(query *sitter.Query) {
		captures := tree.QueryConventional(query)
		for _, capture := range captures {
			objectNode := capture["param_object"]
			if objectNode.Type() == "identifier" || objectNode.Type() == "instance_variable" {
				id := strings.TrimLeft(objectNode.Content(), "@")
				helperDatatypes[objectNode.ID()] = &schemadatatype.DataType{
					Node:       objectNode,
					Name:       id,
					Type:       schema.SimpleTypeUknown,
					TextType:   "",
					Properties: make(map[string]schemadatatype.DataTypable),
					UUID:       "",
				}
			}

			elementNode := capture["param_parent"]
			idNode := capture["param_id"]

			id := strings.TrimLeft(idNode.Content(), ":")

			helperDatatypes[elementNode.ID()] = &schemadatatype.DataType{
				Node:       idNode,
				Name:       id,
				Type:       schema.SimpleTypeUknown,
				TextType:   "",
				Properties: make(map[string]schemadatatype.DataTypable),
				UUID:       "",
			}
		}
	}

	doElementsQuery(elementIdentifierQuery)
	doElementsQuery(elementSimpleSymbolQuery)

	// add calls
	captures := tree.QueryConventional(callsQuery)
	for _, capture := range captures {
		receiverNode := capture["param_receiver"]
		if receiverNode.Type() == "identifier" || receiverNode.Type() == "instance_variable" {
			id := strings.TrimLeft(receiverNode.Content(), "@")
			helperDatatypes[receiverNode.ID()] = &schemadatatype.DataType{
				Node:       receiverNode,
				Name:       id,
				Type:       schema.SimpleTypeUknown,
				TextType:   "",
				Properties: make(map[string]schemadatatype.DataTypable),
				UUID:       "",
			}
		}

		elementNode := capture["param_parent"]
		idNode := capture["param_id"]

		id := idNode.Content()

		helperDatatypes[elementNode.ID()] = &schemadatatype.DataType{
			Node:       idNode,
			Name:       id,
			Type:       schema.SimpleTypeUknown,
			TextType:   "",
			Properties: make(map[string]schemadatatype.DataTypable),
			UUID:       "",
		}
	}
}

func linkProperties(tree *parser.Tree, datatypes, helperDatatypes map[parser.NodeID]*schemadatatype.DataType) {
	for _, helperType := range helperDatatypes {
		node := helperType.Node
		parent := node.Parent()
		// add root node

		if parent.Type() == "call" {
			receiver := parent.ChildByFieldName("receiver")

			// add root calls
			if receiver != nil && receiver.ID() == node.ID() {
				datatypes[node.ID()] = helperType
				continue
			}

			// make sure it is not a function call
			isAllowedCall := false
			if receiver.Type() == "call" {
				isAllowedCall = true
				if receiver.ChildByFieldName("arguments") != nil {
					isAllowedCall = false
				}
			}

			// link chains
			if isAllowedCall || receiver.Type() == "element_reference" || receiver.Type() == "identifier" || receiver.Type() == "instance_variable" {
				// there are wierd cases like [-2].to_sym where there is no id
				_, hasID := helperDatatypes[receiver.ID()]
				if hasID {
					helperDatatypes[receiver.ID()].Properties[helperType.Name] = helperType
					continue
				}
			}
		}

		if parent.Type() == "element_reference" {
			object := parent.ChildByFieldName("object")

			// add root element references
			if object != nil && object.ID() == node.ID() {
				datatypes[node.ID()] = helperType
				continue
			}

			// make sure it is not a function call
			isAllowedCall := false
			if object.Type() == "call" {
				isAllowedCall = true
				if object.ChildByFieldName("arguments") != nil {
					isAllowedCall = false
				}
			}

			// link chains
			if isAllowedCall || object.Type() == "element_reference" || object.Type() == "identifier" || object.Type() == "instance_variable" {
				// there are wierd cases like [-2].to_sym where there is no id
				_, hasID := helperDatatypes[object.ID()]
				if hasID {
					helperDatatypes[object.ID()].Properties[helperType.Name] = helperType
					continue
				}
			}
		}

		// link to root document
		datatypes[node.ID()] = helperType

	}
}

func scopeAndMergeProperties(propertiesDatatypes, classDataTypes map[parser.NodeID]*schemadatatype.DataType, idGenerator nodeid.Generator) {
	// replace root instance variables with class names for classes
	for key, datatype := range propertiesDatatypes {
		if datatype.Node.Type() == "instance_variable" {
			class, err := datatype.Node.FindParent("class")
			if err != nil {
				continue
			}

			nameNode := class.ChildByFieldName("name")
			if nameNode == nil {
				continue
			}

			// add a parent class data type
			classDataTypes[datatype.Node.ID()] = &schemadatatype.DataType{
				Name:       nameNode.Content(),
				Node:       datatype.Node,
				Type:       schema.SimpleTypeUknown,
				TextType:   "",
				Properties: make(map[string]schemadatatype.DataTypable),
			}

			classDataTypes[datatype.Node.ID()].Properties[datatype.Name] = datatype
			delete(propertiesDatatypes, key)
		}
	}
	// replace root instance variables with class names for classes assigmnent
	for key, datatype := range propertiesDatatypes {
		if datatype.Node.Type() == "instance_variable" {
			classNode, err := datatype.Node.FindParent("do_block")
			if err != nil {
				continue
			}

			classNodeDatatype, exists := classDataTypes[classNode.ID()]
			if !exists {
				continue
			}

			className := classNodeDatatype.Name

			// add a parent class data type
			classDataTypes[datatype.Node.ID()] = &schemadatatype.DataType{
				Name:       className,
				Node:       datatype.Node,
				Type:       schema.SimpleTypeUknown,
				TextType:   "",
				Properties: make(map[string]schemadatatype.DataTypable),
			}

			classDataTypes[datatype.Node.ID()].Properties[datatype.Name] = datatype
			delete(propertiesDatatypes, key)
		}
	}

	// scope replaced properties
	terminatorKeywords := []string{"program"}
	parserdatatype.ScopeDatatypes(classDataTypes, idGenerator, terminatorKeywords)

	// // scoped data
	terminatorKeywords = []string{"program", "method", "block", "lambda", "singleton_method"}
	// pull all scope terminator nodes
	parserdatatype.ScopeDatatypes(propertiesDatatypes, idGenerator, terminatorKeywords)
}
