package schemahelper

import (
	"github.com/bearer/bearer/pkg/report/schema"
	"github.com/bearer/bearer/pkg/report/source"
)

type Schema struct {
	Source source.Source
	Value  schema.Schema
}
