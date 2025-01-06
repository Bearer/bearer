package schema

import (
	"github.com/bearer/bearer/pkg/parser"
	"github.com/bearer/bearer/pkg/parser/nodeid"
	"github.com/bearer/bearer/pkg/report/detectors"
	"github.com/bearer/bearer/pkg/report/source"
)

const (
	SimpleTypeFunction = "function"
	SimpleTypeObject   = "object"
	SimpleTypeNumber   = "number"
	SimpleTypeDate     = "date"
	SimpleTypeString   = "string"
	SimpleTypeBool     = "boolean"
	SimpleTypeBinary   = "binary"
	SimpleTypeUnknown  = "unknown"
)

type Schema struct {
	ObjectName           string      `json:"object_name" yaml:"object_name"`
	ObjectUUID           string      `json:"-" yaml:"-"`
	FieldName            string      `json:"field_name" yaml:"field_name"`
	FieldUUID            string      `json:"-" yaml:"-"`
	FieldType            string      `json:"field_type" yaml:"field_type"`
	SimpleFieldType      string      `json:"field_type_simple" yaml:"field_type_simple"`
	Classification       interface{} `json:"classification,omitempty" yaml:"classification,omitempty"`
	Source               *Source     `json:"source,omitempty" yaml:"source,omitempty"`
	NormalizedObjectName string      `json:"normalized_object_name,omitempty" yaml:"normalized_object_name,omitempty"`
	NormalizedFieldName  string      `json:"normalized_field_name,omitempty" yaml:"normalized_field_name,omitempty"`
}

type Source struct {
	// This is the starting line number, the very beginning of what's used by the custom detection
	StartLineNumber   int    `json:"start_line_number,omitempty" yaml:"start_line_number,omitempty"`
	StartColumnNumber int    `json:"start_column_number,omitempty" yaml:"start_column_number,omitempty"`
	EndLineNumber     int    `json:"end_line_number,omitempty" yaml:"end_line_number,omitempty"`
	EndColumnNumber   int    `json:"end_column_number,omitempty" yaml:"end_column_number,omitempty"`
	Content           string `json:"content,omitempty" yaml:"content,omitempty"`
}

type ReportSchema interface {
	SchemaGroupBegin(detectorType detectors.Type, node *parser.Node, schema Schema, source *source.Source, parent *parser.Node)
	SchemaGroupIsOpen() bool
	SchemaGroupShouldClose(tableName string) bool
	SchemaGroupAddItem(node *parser.Node, schema Schema, source *source.Source)
	SchemaGroupEnd(idGenerator nodeid.Generator)
}
