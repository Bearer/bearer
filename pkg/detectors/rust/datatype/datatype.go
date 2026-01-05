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
	"github.com/smacker/go-tree-sitter/rust"
)

var structsQuery = parser.QueryMustCompile(rust.GetLanguage(),
	`(struct_item
		name: (type_identifier) @param_name
	) @param_struct`)

var structFieldsQuery = parser.QueryMustCompile(rust.GetLanguage(),
	`(struct_item
		name: (type_identifier)
		body: (field_declaration_list
			(field_declaration
				name: (field_identifier) @param_field_name
				type: (_) @param_field_type
			) @param_field
		)
	) @param_struct`)

var enumQuery = parser.QueryMustCompile(rust.GetLanguage(),
	`(enum_item
		name: (type_identifier) @param_name
	) @param_enum`)

var implMethodsQuery = parser.QueryMustCompile(rust.GetLanguage(),
	`(impl_item
		type: (_) @param_type_name
		body: (declaration_list
			(function_item
				name: (identifier) @param_method_name
			)
		)
	)`)

var functionParametersQuery = parser.QueryMustCompile(rust.GetLanguage(),
	`(function_item
		parameters: (parameters
			(parameter
				pattern: (identifier) @param_parameter_name
				type: (_) @param_parameter_type
			) @param_parameter
		)
	)`)

func Discover(report report.Report, tree *parser.Tree, idGenerator nodeid.Generator) {
	datatypes := make(map[parser.NodeID]*schemadatatype.DataType)

	// Add structs
	captures := tree.QueryConventional(structsQuery)
	for _, capture := range captures {
		name := capture["param_name"].Content()
		structNode := capture["param_struct"]
		datatypes[structNode.ID()] = &schemadatatype.DataType{
			Node:       structNode,
			Name:       name,
			Type:       schema.SimpleTypeObject,
			TextType:   "struct",
			Properties: make(map[string]schemadatatype.DataTypable),
		}
	}

	// Add enums
	captures = tree.QueryConventional(enumQuery)
	for _, capture := range captures {
		name := capture["param_name"].Content()
		enumNode := capture["param_enum"]
		datatypes[enumNode.ID()] = &schemadatatype.DataType{
			Node:       enumNode,
			Name:       name,
			Type:       schema.SimpleTypeObject,
			TextType:   "enum",
			Properties: make(map[string]schemadatatype.DataTypable),
		}
	}

	discoverStructFields(tree, datatypes)
	discoverImplMethods(tree, datatypes)
	discoverFunctionParameters(tree, datatypes)

	datatype.PruneMap(datatypes)

	report.AddDataType(detections.TypeSchema, detectors.DetectorRust, idGenerator, datatypes, nil)
}

func discoverStructFields(tree *parser.Tree, datatypes map[parser.NodeID]*schemadatatype.DataType) {
	captures := tree.QueryConventional(structFieldsQuery)
	for _, capture := range captures {
		structNode := capture["param_struct"]

		if datatypes[structNode.ID()] == nil {
			continue
		}

		fieldNameNode := capture["param_field_name"]
		if fieldNameNode == nil {
			continue
		}

		fieldName := fieldNameNode.Content()
		fieldTypeNode := capture["param_field_type"]

		fieldTextType := getTypeFromNode(fieldTypeNode)
		fieldType := standardizeDataType(fieldTextType)

		datatypes[structNode.ID()].Properties[fieldName] = &schemadatatype.DataType{
			Node:       fieldNameNode,
			Name:       fieldName,
			Type:       fieldType,
			TextType:   fieldTextType,
			Properties: make(map[string]schemadatatype.DataTypable),
		}
	}
}

func discoverImplMethods(tree *parser.Tree, datatypes map[parser.NodeID]*schemadatatype.DataType) {
	captures := tree.QueryConventional(implMethodsQuery)
	for _, capture := range captures {
		methodNameNode := capture["param_method_name"]
		methodName := methodNameNode.Content()
		typeNameNode := capture["param_type_name"]

		var typeName string
		if typeNameNode.Type() == "type_identifier" {
			typeName = typeNameNode.Content()
		} else if typeNameNode.Type() == "generic_type" {
			// For generic types like Vec<T>, get the base type
			typeNode := typeNameNode.ChildByFieldName("type")
			if typeNode != nil {
				typeName = typeNode.Content()
			}
		}

		// Find the corresponding datatype
		var targetDataType *schemadatatype.DataType
		for _, dt := range datatypes {
			if dt.Name == typeName {
				targetDataType = dt
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

func discoverFunctionParameters(tree *parser.Tree, datatypes map[parser.NodeID]*schemadatatype.DataType) {
	captures := tree.QueryConventional(functionParametersQuery)
	for _, capture := range captures {
		parameterNode := capture["param_parameter"]
		parameterName := capture["param_parameter_name"].Content()
		parameterTypeNode := capture["param_parameter_type"]

		paramTextType := getTypeFromNode(parameterTypeNode)
		paramType := standardizeDataType(paramTextType)

		datatypes[parameterNode.ID()] = &schemadatatype.DataType{
			Node:     parameterNode,
			Name:     parameterName,
			Type:     paramType,
			TextType: paramTextType,
		}
	}
}

func getTypeFromNode(node *parser.Node) string {
	if node == nil {
		return ""
	}

	switch node.Type() {
	case "type_identifier":
		return node.Content()
	case "primitive_type":
		return node.Content()
	case "reference_type":
		// &str, &String, etc.
		typeNode := node.ChildByFieldName("type")
		if typeNode != nil {
			return getTypeFromNode(typeNode)
		}
	case "generic_type":
		// Vec<T>, Option<T>, etc.
		typeNode := node.ChildByFieldName("type")
		if typeNode != nil {
			return typeNode.Content()
		}
	}

	return ""
}

func standardizeDataType(propertyType string) string {
	standardizedType := strings.TrimPrefix(propertyType, "&")
	standardizedType = strings.TrimPrefix(standardizedType, "mut ")

	switch standardizedType {
	case "i8", "i16", "i32", "i64", "i128", "isize",
		"u8", "u16", "u32", "u64", "u128", "usize",
		"f32", "f64":
		return schema.SimpleTypeNumber
	case "bool":
		return schema.SimpleTypeBool
	case "str", "String", "char":
		return schema.SimpleTypeString
	case "Vec", "Box", "Rc", "Arc":
		return schema.SimpleTypeObject
	}

	return schema.SimpleTypeObject
}

