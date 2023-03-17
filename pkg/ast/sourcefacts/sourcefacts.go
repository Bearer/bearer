package sourcefacts

import (
	sitter "github.com/smacker/go-tree-sitter"

	"github.com/bearer/bearer/pkg/ast/idgenerator"
	"github.com/bearer/bearer/pkg/ast/languages/ruby/common"
	"github.com/bearer/bearer/pkg/ast/walker"
	writer "github.com/bearer/bearer/pkg/util/souffle/writer"
	writerbase "github.com/bearer/bearer/pkg/util/souffle/writer/base"
)

type astWriter struct {
	writer.FactWriter
	input           []byte
	nodeIdGenerator *idgenerator.NodeIdGenerator
	childIndex      uint32
}

func WriteFacts(
	walker *walker.Walker,
	input []byte,
	rootNode *sitter.Node,
	nodeIdGenerator *idgenerator.NodeIdGenerator,
	writer writer.FactWriter,
) error {
	w := &astWriter{
		FactWriter:      writer,
		input:           input,
		nodeIdGenerator: nodeIdGenerator,
	}

	return walker.Walk(rootNode, w.visitNode)
}

func (writer *astWriter) visitNode(node *sitter.Node, visitChildren func() error) error {
	nodeElement := writer.node(node)

	if node.Parent() != nil {
		if fname := common.FieldName(node); fname != "" {
			if err := writer.WriteFact("AST_NodeField", writer.node(node.Parent()), nodeElement, writer.Symbol(fname)); err != nil {
				return err
			}
		} else {
			if err := writer.WriteFact(
				"AST_ParentChild",
				writer.node(node.Parent()),
				writer.Unsigned(writer.childIndex),
				nodeElement,
			); err != nil {
				return err
			}

			writer.childIndex++
		}
	}

	if err := writer.WriteFact("AST_NodeType", nodeElement, writer.Symbol(node.Type())); err != nil {
		return err
	}

	if common.MatchContent(node) {
		if err := writer.WriteFact("AST_NodeContent", nodeElement, writer.Symbol(node.Content(writer.input))); err != nil {
			return err
		}
	}

	if err := writer.WriteFact(
		"AST_NodeLocation",
		nodeElement,
		writer.Record(
			writer.Unsigned(node.StartByte()),
			writer.Unsigned(node.StartPoint().Row+1),
			writer.Unsigned(node.StartPoint().Column+1),
			writer.Unsigned(node.EndPoint().Row+1),
			writer.Unsigned(node.EndPoint().Column+1),
		),
	); err != nil {
		return err
	}

	oldChildIndex := writer.childIndex
	writer.childIndex = 0
	err := visitChildren()
	writer.childIndex = oldChildIndex

	return err
}

func (writer *astWriter) node(node *sitter.Node) writerbase.Element {
	return writer.Unsigned(writer.nodeIdGenerator.Get(node))
}
