package schema

import (
	"github.com/bearer/curio/pkg/report/detectors"
	"github.com/bearer/curio/pkg/report/source"
)

const (
	SimpleTypeFunction = "function"
	SimpleTypeObject   = "object"
	SimpleTypeNumber   = "number"
	SimpleTypeDate     = "date"
	SimpleTypeString   = "string"
	SimpleTypeBool     = "boolean"
	SimpleTypeBinary   = "binary"
	SimpleTypeUknown   = "unknown"
)

type Schema struct {
	ObjectName      string `json:"object_name"`
	ObjectUUID      string `json:"object_uuid"`
	FieldName       string `json:"field_name"`
	FieldUUID       string `json:"field_uuid"`
	FieldType       string `json:"field_type"`
	SimpleFieldType string `json:"field_type_simple"`
}

type ReportSchema interface {
	AddSchema(detectorType detectors.Type, schema Schema, source source.Source)
}
