package operationshelper

import (
	"github.com/bearer/bearer/pkg/report/operations"
	"github.com/bearer/bearer/pkg/report/source"
)

type Operation struct {
	Source source.Source
	Value  operations.Operation
}
