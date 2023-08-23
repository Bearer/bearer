package operationshelper

import (
	"github.com/bearer/bearer/internal/report/operations"
	"github.com/bearer/bearer/internal/report/source"
)

type Operation struct {
	Source source.Source
	Value  operations.Operation
}
