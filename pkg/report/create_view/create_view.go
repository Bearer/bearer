package createview

import "github.com/bearer/curio/pkg/report/source"

type View struct {
	SchemaName string        `json:"schema_name" yaml:"schema_name"`
	ViewName   string        `json:"view_name" yaml:"view_name"`
	Source     source.Source `json:"source" yaml:"source"`
	From       []*Table      `json:"from" yaml:"from"`
	Fields     []*Field      `json:"fields" yaml:"fields"`
}

type Field struct {
	TableName string        `json:"table_name" yaml:"table_name"`
	FieldName string        `json:"field_name" yaml:"field_name"`
	CommitSHA string        `json:"commit_sha" yaml:"commit_sha"`
	Source    source.Source `json:"source" yaml:"source"`
}

type Table struct {
	TableName string        `json:"table_name" yaml:"table_name"`
	FieldName string        `json:"field_name" yaml:"field_name"`
	CommitSHA string        `json:"commit_sha" yaml:"commit_sha"`
	Source    source.Source `json:"source" yaml:"source"`
	Alias     string        `json:"alias" yaml:"alias"`
}
