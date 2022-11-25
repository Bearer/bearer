package customdetector

import (
	"strings"

	"github.com/bearer/curio/pkg/detectors/ruby/datatype"
	"github.com/bearer/curio/pkg/parser"
	"github.com/bearer/curio/pkg/parser/nodeid"
	"github.com/bearer/curio/pkg/report/schema"
	schemadatatype "github.com/bearer/curio/pkg/report/schema/datatype"
	"github.com/bearer/curio/pkg/util/file"
	"github.com/smacker/go-tree-sitter/ruby"
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

		content := singleArgument.Content()
		tree, err := parser.ParseBytes(fileinfo, filepath, []byte(content), ruby.GetLanguage(), singleArgument.LineNumber()-1)
		if err != nil {
			return nil, err
		}

		singleArgumentDatatypes := datatype.Discover(tree, idGenerator)
		allDatatypes := datatype.Discover(node.Tree(), idGenerator)

		for nodeID, target := range singleArgumentDatatypes {
			joinedDatatypes[nodeID] = target
		}

	}

	return joinedDatatypes, nil
}

func variableReconciliation(singleArgumentDatatypes map[parser.NodeID]*schemadatatype.DataType, allDataTypes map[parser.NodeID]*schemadatatype.DataType) map[parser.NodeID]*schemadatatype.DataType {
	result := make(map[parser.NodeID]*schemadatatype.DataType)

	for nodeID, argumentDatatype := range singleArgumentDatatypes {
		currentNode := argumentDatatype.Node
		for {
			parent := currentNode.Parent()
			if parent == nil {
				break
			}

			isTerminating := false

			for _, terminator := range datatype.ScopeTerminators {
				if parent.Type() == terminator {
					isTerminating = true
					break
				}
			}

			if !isTerminating {
				continue
			}

			for globalNodeID, globalDatatype := range allDataTypes {
				// not in the same scope
				if globalNodeID != parent.ID() {
					continue
				}
				// in the same scope but it doesn't interest us
				if globalDatatype.Name != argumentDatatype.Name {

				}
				// merge properties of argumentDatatype and globalDatatype
			}
		}

	}
}
