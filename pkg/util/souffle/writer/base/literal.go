package base

import (
	"fmt"
	"strings"
)

type Any struct{}
type Identifier string

type Conjunction []Literal
type Disjunction []Literal

type Predicate struct {
	Name     string
	Elements []LiteralElement
}

type NegativePredicate Predicate

type Constraint struct {
	Left     LiteralElement
	Operator string
	Right    LiteralElement
}

type Literal interface {
	String() string
	sealedLiteral()
}

func (Conjunction) sealedLiteral()       {}
func (Disjunction) sealedLiteral()       {}
func (Predicate) sealedLiteral()         {}
func (NegativePredicate) sealedLiteral() {}
func (Constraint) sealedLiteral()        {}

type LiteralElement interface {
	String() string
	sealedLiteralElement()
}

// func (Record) sealedLiteralElement()     {} // needs LiteralRecord
func (Symbol) sealedLiteralElement()     {}
func (Unsigned) sealedLiteralElement()   {}
func (Any) sealedLiteralElement()        {}
func (Identifier) sealedLiteralElement() {}
func (Predicate) sealedLiteralElement()  {}

func (any Any) String() string {
	return "_"
}

func (identifier Identifier) String() string {
	return string(identifier)
}

func (conjunction Conjunction) String() string {
	literals := make([]string, len(conjunction))
	for i, literal := range conjunction {
		literals[i] = literal.String()
	}

	return strings.Join(literals, ", ")
}

func (disjunction Disjunction) String() string {
	literals := make([]string, len(disjunction))
	for i, literal := range disjunction {
		literals[i] = literal.String()
	}

	return fmt.Sprintf("(%s)", strings.Join(literals, "; "))
}

func (predicate Predicate) String() string {
	builder := strings.Builder{}

	builder.WriteString(predicate.Name)
	builder.WriteString("(")

	for i, element := range predicate.Elements {
		if i != 0 {
			builder.WriteString(", ")
		}

		builder.WriteString(ElementString(element))
	}

	builder.WriteString(")")

	return builder.String()
}

func (predicate NegativePredicate) String() string {
	return "!" + Predicate(predicate).String()
}

func (constraint Constraint) String() string {
	return constraint.Left.String() + constraint.Operator + constraint.Right.String()
}
