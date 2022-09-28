package createview

import "github.com/bearer/curio/pkg/report/source"

type View struct {
	SchemaName string        `json:"schema_name"`
	ViewName   string        `json:"view_name"`
	Source     source.Source `json:"source"`
	From       []*Table      `json:"from"`
	Fields     []*Field      `json:"fields"`
}

type Field struct {
	TableName string        `json:"table_name"`
	FieldName string        `json:"field_name"`
	CommitSHA string        `json:"commit_sha"`
	Source    source.Source `json:"source"`
}

type Table struct {
	TableName string        `json:"table_name"`
	FieldName string        `json:"field_name"`
	CommitSHA string        `json:"commit_sha"`
	Source    source.Source `json:"source"`
	Alias     string        `json:"alias"`
}
