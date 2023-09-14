package schemahelper

import (
	"github.com/bearer/bearer/internal/report/schema"
	"github.com/bearer/bearer/internal/report/source"
)

type Schema struct {
	Source source.Source
	Value  schema.Schema
}
