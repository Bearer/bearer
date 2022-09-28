package knex

import "github.com/bearer/curio/pkg/report/frameworks"

const TypeFunction frameworks.Type = "knex_function"
const TypeSchema frameworks.Type = "knex_schema"

type Function struct {
	DataType string `json:"data_type"`
}

type Schema struct {
	DataType     string `json:"data_type"`
	PropertyName string `json:"property_name"`
}
