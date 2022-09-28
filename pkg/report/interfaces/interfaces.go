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
	Type         Type          `json:"type"`
	Value        *values.Value `json:"value"`
	VariableName string        `json:"variable_name,omitempty"`
}
