package datatype

import (
	"strings"

	"github.com/bearer/bearer/pkg/parser"
	"github.com/bearer/bearer/pkg/parser/datatype"
	"github.com/bearer/bearer/pkg/parser/nodeid"
	"github.com/bearer/bearer/pkg/report"
	"github.com/bearer/bearer/pkg/report/detections"
	"github.com/bearer/bearer/pkg/report/detectors"
	"github.com/bearer/bearer/pkg/report/schema"
	schemadatatype "github.com/bearer/bearer/pkg/report/schema/datatype"
	"github.com/smacker/go-tree-sitter/golang"
)

var structsQuery = parser.QueryMustCompile(golang.GetLanguage(),
	`(type_declaration
		(type_spec
			name: (type_identifier) @param_name
			type: (struct_type)
		)
	)@param_struct`)

var structPropertiesQuery = parser.QueryMustCompile(golang.GetLanguage(),
	`(type_declaration
		(type_spec
			name: (type_identifier)
			type: (struct_type
				(field_declaration_list
					(field_declaration) @param_property
				)
			)
		)
	)@param_struct`)

var methodsQuery = parser.QueryMustCompile(
	golang.GetLanguage(),
	`(method_declaration
		receiver: (parameter_list
			(parameter_declaration) @param_struct
		)
		name: (field_identifier) @param_method_name
	)`)

var queryMethodParameters = parser.QueryMustCompile(
	golang.GetLanguage(),
	`(method_declaration
		parameters: (parameter_list
			(parameter_declaration
				name: (identifier) @param_parameter_name
			) @param_parameter
		)
	)`,
)

var queryFunctionParameters = parser.QueryMustCompile(
	golang.GetLanguage(),
	`(function_declaration
		parameters: (parameter_list
			(parameter_declaration
				name: (identifier) @param_parameter_name
			) @param_parameter
		)
	)`,
)

func Discover(report report.Report, tree *parser.Tree, idGenerator nodeid.Generator) {
	datatypes := make(map[parser.NodeID]*schemadatatype.DataType)

	// add structs
	captures := tree.QueryConventional(structsQuery)
	for _, capture := range captures {
		id := capture["param_name"].Content()
		structNode := capture["param_struct"]
		datatypes[structNode.ID()] = &schemadatatype.DataType{
			Node:       structNode,
			Name:       id,
			Type:       schema.SimpleTypeObject,
			TextType:   "struct",
			Properties: make(map[string]schemadatatype.DataTypable),
		}
	}

	discoverProperties(tree, datatypes)
	discoverFunctions(tree, datatypes)
	discoverParameters(tree, datatypes)

	datatype.PruneMap(datatypes)

	report.AddDataType(detections.TypeSchema, detectors.DetectorGo, idGenerator, datatypes, nil)
}

func discoverProperties(tree *parser.Tree, datatypes map[parser.NodeID]*schemadatatype.DataType) {
	captures := tree.QueryConventional(structPropertiesQuery)
	for _, capture := range captures {
		structNode := capture["param_struct"]

		if datatypes[structNode.ID()] == nil {
			continue
		}

		propertyNode := capture["param_property"]
		propertyNameNode := propertyNode.ChildByFieldName("name")
		// name is undefined for embeds
		if propertyNameNode == nil {
			continue
		}

		propertyName := propertyNameNode.Content()

		propertyTypeNode := propertyNode.ChildByFieldName("type")

		propertyTextType := getTypeFromNode(propertyTypeNode)
		propertyType := standardizeDataType(propertyTypeNode, propertyTextType)

		datatypes[structNode.ID()].Properties[propertyName] = &schemadatatype.DataType{
			Node:       propertyNameNode,
			Name:       propertyName,
			Type:       propertyType,
			TextType:   propertyTextType,
			Properties: make(map[string]schemadatatype.DataTypable),
		}
	}
}

func discoverFunctions(tree *parser.Tree, datatypes map[parser.NodeID]*schemadatatype.DataType) {
	captures := tree.QueryConventional(methodsQuery)
	for _, capture := range captures {
		methodNameNode := capture["param_method_name"]
		methodName := methodNameNode.Content()
		targetDatatypeTypeNode := capture["param_struct"].ChildByFieldName("type")
		var targetDatatypeName string

		if targetDatatypeTypeNode.Type() == "type_identifier" {
			targetDatatypeName = targetDatatypeTypeNode.Content()
		} else if targetDatatypeTypeNode.Type() == "pointer_type" {
			targetDatatypeName = targetDatatypeTypeNode.Child(0).Content()
		}

		var targetDataType *schemadatatype.DataType
		for _, datatype := range datatypes {
			if datatype.Name == targetDatatypeName {
				targetDataType = datatype
				break
			}
		}

		if targetDataType == nil {
			continue
		}

		targetDataType.Properties[methodName] = &schemadatatype.DataType{
			Node:       methodNameNode,
			Name:       methodName,
			Type:       schema.SimpleTypeFunction,
			TextType:   schema.SimpleTypeFunction,
			Properties: make(map[string]schemadatatype.DataTypable),
		}
	}
}

func discoverParameters(tree *parser.Tree, datatypes map[parser.NodeID]*schemadatatype.DataType) {
	addDatatype := func(capture parser.Captures) {
		parameterNode := capture["param_parameter"]
		parameterName := capture["param_parameter_name"].Content()

		parameterTypeNode := parameterNode.ChildByFieldName("type")

		propertyTextType := getTypeFromNode(parameterTypeNode)
		propertyType := standardizeDataType(parameterTypeNode, propertyTextType)

		datatypes[parameterNode.ID()] = &schemadatatype.DataType{
			Node:     parameterNode,
			Name:     parameterName,
			Type:     propertyType,
			TextType: propertyTextType,
		}
	}

	captures := tree.QueryConventional(queryFunctionParameters)
	for _, capture := range captures {
		addDatatype(capture)
	}

	captures = tree.QueryConventional(queryMethodParameters)
	for _, capture := range captures {
		addDatatype(capture)
	}
}

func getTypeFromNode(node *parser.Node) string {
	// primitive types
	if node.Type() == "type_identifier" {
		return node.Content()
	}

	// same package pointers
	if node.Type() == "pointer_type" && node.ChildCount() == 1 {
		if node.Child(0).Type() == "type_identifier" {
			return node.Child(0).Content()
		}
	}

	// imported packages
	if node.Type() == "qualified_type" {
		child := node.ChildByFieldName("name")
		if child.Type() == "type_identifier" {
			return child.Content()
		}
	}

	// imported package pointer
	if node.Type() == "pointer_type" && node.ChildCount() == 1 {
		child := node.Child(0)
		if child.Type() == "qualified_type" {
			grandchild := child.ChildByFieldName("name")
			if grandchild.Type() == "type_identifier" {
				return grandchild.Content()
			}
		}
	}

	return ""
}

func standardizeDataType(node *parser.Node, propertyType string) string {
	standardizedType := strings.Trim(propertyType, "*")

	if standardizedType == "int" || standardizedType == "int8" || standardizedType == "int16" || standardizedType == "int32" || standardizedType == "int64" ||
		standardizedType == "uint" || standardizedType == "uint8" || standardizedType == "uint16" || standardizedType == "uint32" || standardizedType == "uint64" ||
		standardizedType == "float32" || standardizedType == "float64" {
		return schema.SimpleTypeNumber
	}

	if standardizedType == "byte" {
		return schema.SimpleTypeBinary
	}

	if standardizedType == "bool" {
		return schema.SimpleTypeBool
	}

	if standardizedType == "Time" {
		return schema.SimpleTypeDate
	}

	if standardizedType == "string" || standardizedType == "rune" {
		return schema.SimpleTypeString
	}

	return schema.SimpleTypeObject
}
