package values

import (
	"fmt"
	"strings"

	"github.com/bearer/curio/pkg/report/variables"
	"github.com/rs/zerolog/log"
)

type Value struct {
	Parts []Part `json:"parts"`
}

type Part interface {
	GetVariableReferences() []*VariableReference
	GetParts() []Part
	// Pattern is a string representation for debugging
	Pattern() string
}

type PartType string

const (
	PartTypeString            PartType = "string"
	PartTypeVariableReference PartType = "variable_reference"
	PartTypeUnknown           PartType = "unknown"
)

type String struct {
	Type  PartType `json:"type"`
	Value string   `json:"value"`
}

type VariableReference struct {
	Type       PartType             `json:"type"`
	Identifier variables.Identifier `json:"identifier"`
}

type Unknown struct {
	Type  PartType `json:"type"`
	Parts []Part   `json:"parts"`
}

func New() *Value {
	return &Value{}
}

func (value *Value) AppendPart(part Part) {
	value.Parts = append(value.Parts, part)
}

func (value *Value) AppendString(text string) {
	if len(value.Parts) != 0 {
		lastPart := value.Parts[len(value.Parts)-1]
		if stringPart, ok := lastPart.(*String); ok {
			stringPart.Value += text
			return
		}
	}

	value.AppendPart(NewStringPart(text))
}

func (value *Value) AppendVariableReference(variableType variables.Type, name string) {
	value.AppendPart(NewVariableReferencePart(variableType, name))
}

func (value *Value) AppendUnknown(parts []Part) {
	value.AppendPart(NewUnknownPart(parts))
}

// IsUnknown returns whether all parts of the value are unknown
func (value *Value) IsUnknown() bool {
	for _, part := range value.Parts {
		if _, isUnknown := part.(*Unknown); !isUnknown {
			return false
		}
	}

	return true
}

func (value *Value) GetParts() []Part {
	var result []Part

	for _, part := range value.Parts {
		result = append(result, part.GetParts()...)
	}

	return result
}

func (value *Value) GetVariableReferences() []*VariableReference {
	var result []*VariableReference

	for _, part := range value.Parts {
		result = append(result, part.GetVariableReferences()...)
	}

	return result
}

func (value *Value) Pattern() string {
	patterns := make([]string, len(value.Parts))
	for i, part := range value.Parts {
		patterns[i] = part.Pattern()
	}

	return strings.Join(patterns, "")
}

func (value *Value) Append(other *Value) {
	for _, part := range other.Parts {
		if stringPart, ok := part.(*String); ok {
			value.AppendString(stringPart.Value)
		} else if variableReferencePart, ok := part.(*VariableReference); ok {
			identifier := variableReferencePart.Identifier
			value.AppendVariableReference(identifier.Type, identifier.Name)
		} else if unknownPart, ok := part.(*Unknown); ok {
			value.AppendUnknown(unknownPart.Parts)
		} else {
			log.Error().Msg("unexpected value part type")
		}
	}
}

func NewStringPart(text string) *String {
	return &String{Type: PartTypeString, Value: text}
}

func (part *String) GetParts() []Part {
	return []Part{part}
}

func (part *String) GetVariableReferences() []*VariableReference {
	return nil
}

func (part *String) Pattern() string {
	return part.Value
}

func NewVariableReferencePart(variableType variables.Type, name string) *VariableReference {
	return &VariableReference{
		Type:       PartTypeVariableReference,
		Identifier: variables.Identifier{Type: variableType, Name: name},
	}
}

func (part *VariableReference) GetParts() []Part {
	return []Part{part}
}

func (part *VariableReference) GetVariableReferences() []*VariableReference {
	return []*VariableReference{part}
}

func (part *VariableReference) Pattern() string {
	return fmt.Sprintf("${variable:%s:%s}", part.Identifier.Type, part.Identifier.Name)
}

func NewUnknownPart(parts []Part) *Unknown {
	return &Unknown{
		Type:  PartTypeUnknown,
		Parts: parts,
	}
}

func (part *Unknown) GetParts() []Part {
	return part.Parts
}

func (part *Unknown) GetVariableReferences() []*VariableReference {
	var result []*VariableReference

	for _, unknownPart := range part.Parts {
		if variableReferencePart, ok := unknownPart.(*VariableReference); ok {
			result = append(result, variableReferencePart)
		}
	}

	return result
}

func (part *Unknown) Pattern() string {
	var patterns []string
	for _, unknownPart := range part.Parts {
		patterns = append(patterns, unknownPart.Pattern())
	}
	return fmt.Sprintf("${unknown:%s}", strings.Join(patterns, "|"))
}
