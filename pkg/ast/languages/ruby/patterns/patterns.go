package patterns

import (
	"fmt"
	"log"

	sitter "github.com/smacker/go-tree-sitter"

	builderinput "github.com/bearer/bearer/new/language/patternquery/builder/input"
	querytypes "github.com/bearer/bearer/new/language/patternquery/types"
	"github.com/bearer/bearer/pkg/ast/idgenerator"
	"github.com/bearer/bearer/pkg/ast/languages/ruby/common"
	"github.com/bearer/bearer/pkg/ast/walker"
	"github.com/bearer/bearer/pkg/util/set"
	writerbase "github.com/bearer/bearer/pkg/util/souffle/writer/base"
	filewriter "github.com/bearer/bearer/pkg/util/souffle/writer/file"
)

type nodeVariableGenerator struct {
	ids         map[*sitter.Node]uint32
	idGenerator *idgenerator.Generator
}

func newNodeVariableGenerator() *nodeVariableGenerator {
	return &nodeVariableGenerator{
		ids:         make(map[*sitter.Node]uint32),
		idGenerator: idgenerator.NewGenerator(),
	}
}

func (generator *nodeVariableGenerator) Get(node *sitter.Node) string {
	id := generator.getId(node)
	return fmt.Sprintf("node%d", id)
}

func (generator *nodeVariableGenerator) getId(node *sitter.Node) uint32 {
	if id, cached := generator.ids[node]; cached {
		return id
	}

	id := generator.idGenerator.Get()
	generator.ids[node] = id
	return id
}

type patternWriter struct {
	*filewriter.Writer
	inputParams           *builderinput.InputParams
	input                 []byte
	literals              []writerbase.Literal
	childIndex            uint32
	rootElement           writerbase.LiteralElement
	parentElement         writerbase.LiteralElement
	nodeVariableGenerator *nodeVariableGenerator
	tempIdGenerator       *idgenerator.Generator
	handled               set.Set[*sitter.Node]
	variableNodes         map[string][]writerbase.Identifier
}

func CompileRule(
	walker *walker.Walker,
	inputParams *builderinput.InputParams,
	ruleRelation,
	variableRelation string,
	patternIndex int,
	input []byte,
	rootNode *sitter.Node,
	writer *filewriter.Writer,
) error {
	w := &patternWriter{
		Writer:                writer,
		inputParams:           inputParams,
		input:                 input,
		nodeVariableGenerator: newNodeVariableGenerator(),
		tempIdGenerator:       idgenerator.NewGenerator(),
		handled:               set.New[*sitter.Node](),
		variableNodes:         make(map[string][]writerbase.Identifier),
	}

	err := walker.Walk(rootNode, w.visitNode)
	if err != nil {
		return err
	}

	var variablePredicates []writerbase.Predicate
	for name, variables := range w.variableNodes {
		for _, variable := range variables {
			variablePredicates = append(variablePredicates, writer.Predicate(
				variableRelation,
				w.rootElement,
				writer.Symbol(name),
				variable,
			))
		}
	}

	if len(w.literals) > 20 {
		log.Printf("rule too large, skipping")
		return nil
	}
	log.Printf("#literals: %d", len(w.literals))

	if err := writer.WriteRule(
		// append(
		// 	[]writerbase.Predicate(nil), //variablePredicates,
		// 	writer.Predicate(ruleRelation, w.rootElement),
		// ),
		[]writerbase.Predicate{writer.Predicate(ruleRelation, writer.Unsigned(uint32(patternIndex)), w.rootElement)},
		w.literals,
	); err != nil {
		return err
	}

	return nil
}

func (writer *patternWriter) visitNode(node *sitter.Node, visitChildren func() error) error {
	if writer.handled.Has(node) {
		return nil
	}

	nodeElement := writer.Identifier(writer.nodeVariableGenerator.Get(node))

	if node.Type() == "program" {
		if node.ChildCount() != 1 {
			return fmt.Errorf("expected 1 root node in pattern but got %d", node.ChildCount())
		}

		return visitChildren()
	}

	if writer.rootElement == nil {
		writer.rootElement = nodeElement
	}

	if writer.parentElement != nil {
		if fname := common.FieldName(node); fname != "" {
			writer.literals = append(
				writer.literals,
				writer.Predicate("AST_NodeField", writer.parentElement, nodeElement, writer.Symbol(fname)),
			)
		} else {
			writer.literals = append(
				writer.literals,
				writer.Predicate("AST_ParentChild", writer.parentElement, writer.Unsigned(writer.childIndex), nodeElement),
			)

			writer.childIndex++
		}
	}

	if variable := writer.getVariableFor(node); variable != nil {
		if variable.Name != "_" {
			writer.variableNodes[variable.Name] = append(writer.variableNodes[variable.Name], nodeElement)
		}

		writer.literals = append(
			writer.literals,
			writer.Predicate("AST_NodeType", nodeElement, writer.Any()),
		)

		return nil
	}

	writer.literals = append(
		writer.literals,
		writer.Predicate("AST_NodeType", nodeElement, writer.Symbol(node.Type())),
	)

	if common.MatchContent(node) {
		writer.literals = append(
			writer.literals,
			writer.Predicate("AST_NodeContent", nodeElement, writer.Symbol(node.Content(writer.input))),
		)
	}

	for fieldName := range common.GetMissingFields(node) {
		writer.literals = append(
			writer.literals,
			writer.NegativePredicate("AST_NodeField", nodeElement, writer.Any(), writer.Symbol(fieldName)),
		)
	}

	// Handle `call` vs `call()`
	if node.Type() == "call" {
		if arguments := node.ChildByFieldName("arguments"); arguments == nil || arguments.NamedChildCount() == 0 {
			writer.handled.Add(arguments)
			writer.literals = append(writer.literals, writer.optionalEmpty(node, "arguments"))
		}
	}

	// Handle { a: 1 } vs { :a => 1 }
	if node.Type() == "pair" {
		key := node.ChildByFieldName("key")
		if key.Type() == "simple_symbol" || key.Type() == "hash_key_symbol" {
			writer.handled.Add(key)

			symbolName := key.Content(writer.input)
			if key.Type() == "simple_symbol" {
				symbolName = symbolName[1:]
			}

			writer.literals = append(writer.literals, writer.eitherNodeLiteral(
				node,
				"key",
				literalNode{nodeType: "simple_symbol", content: ":" + symbolName},
				literalNode{nodeType: "hash_key_symbol", content: symbolName},
			))
		}
	}

	oldParentElement := writer.parentElement
	oldChildIndex := writer.childIndex
	writer.childIndex = 0
	writer.parentElement = nodeElement
	err := visitChildren()
	writer.childIndex = oldChildIndex
	writer.parentElement = oldParentElement

	return err
}

func (writer *patternWriter) optionalEmpty(node *sitter.Node, fieldName string) writerbase.Literal {
	nodeElement := writer.Identifier(writer.nodeVariableGenerator.Get(node))
	tempVariable := writer.tempVariable()

	return writer.Disjunction(
		writer.NegativePredicate("AST_NodeField", nodeElement, writer.Any(), writer.Symbol(fieldName)),
		writer.Conjunction(
			writer.Predicate("AST_NodeField", nodeElement, tempVariable, writer.Symbol(fieldName)),
			writer.NegativePredicate(
				"AST_ParentChild",
				tempVariable,
				writer.Any(),
				writer.Any(),
			),
		),
	)
}

type literalNode struct {
	nodeType string
	content  string
}

func (writer *patternWriter) eitherNodeLiteral(
	node *sitter.Node,
	fieldName string,
	literalNodes ...literalNode,
) writerbase.Literal {
	nodeElement := writer.Identifier(writer.nodeVariableGenerator.Get(node))
	tempVariable := writer.tempVariable()

	literals := make([]writerbase.Literal, len(literalNodes))
	for i, literalNode := range literalNodes {
		literals[i] = writer.Conjunction(
			writer.Predicate("AST_NodeType", tempVariable, writer.Symbol(literalNode.nodeType)),
			writer.Predicate("AST_NodeContent", tempVariable, writer.Symbol(literalNode.content)),
		)
	}

	return writer.Conjunction(
		writer.Predicate("AST_NodeField", nodeElement, tempVariable, writer.Symbol(fieldName)),
		writer.Disjunction(literals...),
	)
}

func (writer *patternWriter) tempVariable() writerbase.LiteralElement {
	return writer.Identifier(fmt.Sprintf("tmp%d", writer.tempIdGenerator.Get()))
}

func (writer *patternWriter) getVariableFor(node *sitter.Node) *querytypes.Variable {
	for _, variable := range writer.inputParams.Variables {
		if node.Content(writer.input) == variable.DummyValue {
			return &variable
		}
	}

	return nil
}
