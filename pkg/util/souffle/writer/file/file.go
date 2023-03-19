package file

import (
	"io"
	"strings"

	"github.com/bearer/bearer/pkg/util/souffle/writer/base"
)

type Writer struct {
	base.Base
	output io.StringWriter
}

func New(output io.StringWriter) *Writer {
	return &Writer{output: output}
}

func (writer *Writer) Any() base.Any {
	return base.Any{}
}

func (writer *Writer) Identifier(value string) base.Identifier {
	return base.Identifier(value)
}

func (writer *Writer) Conjunction(literals ...base.Literal) base.Conjunction {
	return base.Conjunction(literals)
}

func (writer *Writer) Disjunction(literals ...base.Literal) base.Disjunction {
	return base.Disjunction(literals)
}

func (writer *Writer) Constraint(left base.LiteralElement, operator string, right base.LiteralElement) base.Constraint {
	return base.Constraint{Left: left, Operator: operator, Right: right}
}

func (writer *Writer) Predicate(name string, elements ...base.LiteralElement) base.Predicate {
	return base.Predicate{Name: name, Elements: elements}
}

func (writer *Writer) NegativePredicate(name string, elements ...base.LiteralElement) base.NegativePredicate {
	return base.NegativePredicate(base.Predicate{Name: name, Elements: elements})
}

func (writer *Writer) WriteComment(text string) error {
	return writer.write("// " + text + "\n")
}

// FIXME: use proper attribute type
func (writer *Writer) WriteRelation(name string, typ base.RelationType, attributes ...string) error {
	builder := strings.Builder{}
	builder.WriteString(".decl ")
	builder.WriteString(name)
	builder.WriteString("(")

	for i, attribute := range attributes {
		if i != 0 {
			builder.WriteString(", ")
		}

		builder.WriteString(attribute)
	}

	builder.WriteString(")\n")

	if typ != base.Intermediate {
		builder.WriteString(".")
		builder.WriteString(string(typ))
		builder.WriteString(" ")
		builder.WriteString(name)
		builder.WriteString("\n")
	}

	return writer.write(builder.String())
}

func (writer *Writer) WriteRule(heads []base.Predicate, body []base.Literal) error {
	builder := strings.Builder{}

	for i, head := range heads {
		if i != 0 {
			builder.WriteString(", ")
		}

		builder.WriteString(head.String())
	}

	builder.WriteString(" :- ")
	builder.WriteString(writer.Conjunction(body...).String())
	builder.WriteString(".\n")

	return writer.write(builder.String())
}

func (writer *Writer) WriteFact(relation string, elements ...base.Element) error {
	builder := strings.Builder{}

	builder.WriteString(relation)
	builder.WriteString("(")

	for i, element := range elements {
		if i != 0 {
			builder.WriteString(", ")
		}

		builder.WriteString(base.ElementString(element))
	}

	builder.WriteString(").\n")

	return writer.write(builder.String())
}

func (writer *Writer) write(value string) error {
	_, err := writer.output.WriteString(value)

	return err
}
