package operationshelper

import (
	"github.com/bearer/curio/pkg/report/operations"
	"github.com/bearer/curio/pkg/report/source"
)

type Operation struct {
	Source source.Source
	Value  operations.Operation
}
