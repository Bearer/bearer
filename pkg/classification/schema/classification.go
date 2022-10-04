package schema

import (
	"github.com/bearer/curio/pkg/parser/datatype"
)

type ClassifiedDatatype struct {
	*datatype.DataType
	Properties     map[string]ClassifiedDatatype
	Classification Classification `json:"classification"`
}

type Classification struct {
	Name string
}
