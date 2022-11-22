package interfaces

import (
	"github.com/bearer/curio/pkg/report/values"
)

type Type string

const (
	TypeURL  Type = "url"
	TypePath Type = "path"
)

type Interface struct {
	Type         Type          `json:"type" yaml:"type"`
	Value        *values.Value `json:"value" yaml:"value"`
	VariableName string        `json:"variable_name,omitempty" yaml:"variable_name,omitempty"`
}
