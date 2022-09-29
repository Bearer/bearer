package datatype

import (
	"strings"

	"github.com/bearer/curio/pkg/parser"
	parserdatatype "github.com/bearer/curio/pkg/parser/datatype"
	"github.com/bearer/curio/pkg/report/schema"
	"github.com/smacker/go-tree-sitter/ruby"
)

var classAssignmentQuery = parser.QueryMustCompile(ruby.GetLanguage(),
	`(assignment
		left: (constant) @param_id
		right: 
			(call
				receiver: (constant) @helper_Class
				method: (identifier) @helper_new
				block: (do_block) @param_class
			)
	)`)

var classAssignmentPropertiesQuery = parser.QueryMustCompile(ruby.GetLanguage(),
	`(do_block
		( call
			arguments: (argument_list
				(simple_symbol) @param_id
			)
		)
	) @param_class`)

var classAssignmentFunctionsQuery = parser.QueryMustCompile(ruby.GetLanguage(),
	`(do_block
		( method
			name: (identifier) @param_id
		)
	) @param_class`)

func discoverClassAssignment(tree *parser.Tree, datatypes map[parser.NodeID]*parserdatatype.DataType) {
	captures := tree.QueryConventional(classAssignmentQuery)
	for _, capture := range captures {
		id := capture["param_id"].Content()
		classNode := capture["param_class"]
		datatypes[classNode.ID()] = &parserdatatype.DataType{
			Node:       classNode,
			Name:       id,
			Type:       schema.SimpleTypeObject,
			TextType:   "class",
			Properties: make(map[string]*parserdatatype.DataType),
		}
	}

	discoverClassAssignmentProperties(tree, datatypes)
	discoverClassAssignmentFunctions(tree, datatypes)
}

func discoverClassAssignmentProperties(tree *parser.Tree, datatypes map[parser.NodeID]*parserdatatype.DataType) {
	// add class assigment properties
	captures := tree.QueryConventional(classAssignmentPropertiesQuery)
	for _, capture := range captures {
		classNode := capture["param_class"]
		if datatypes[classNode.ID()] == nil {
			continue
		}

		// get property name
		propertyNode := capture["param_id"]
		propertyName := strings.TrimLeft(propertyNode.Content(), ":")

		datatypes[classNode.ID()].Properties[propertyName] = &parserdatatype.DataType{
			Node:       propertyNode,
			Name:       propertyName,
			Type:       schema.SimpleTypeUknown,
			Properties: make(map[string]*parserdatatype.DataType),
			TextType:   "",
		}
	}
}

func discoverClassAssignmentFunctions(tree *parser.Tree, datatypes map[parser.NodeID]*parserdatatype.DataType) {
	captures := tree.QueryConventional(classAssignmentFunctionsQuery)
	for _, capture := range captures {
		classNode := capture["param_class"]
		if datatypes[classNode.ID()] == nil {
			continue
		}

		// get method name
		functionNameNode := capture["param_id"]
		functionName := functionNameNode.Content()

		datatypes[classNode.ID()].Properties[functionName] = &parserdatatype.DataType{
			Node:     functionNameNode,
			Name:     functionName,
			Type:     schema.SimpleTypeUknown,
			TextType: "",
		}
	}
}
