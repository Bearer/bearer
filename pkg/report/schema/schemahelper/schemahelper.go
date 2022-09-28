package schemahelper

import (
	"github.com/bearer/curio/pkg/report/schema"
	"github.com/bearer/curio/pkg/report/source"
)

type Schema struct {
	Source source.Source
	Value  schema.Schema
}
