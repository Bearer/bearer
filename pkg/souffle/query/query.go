package query

import (
	"github.com/bearer/bearer/new/language/implementation"
)

type Query struct {
	ruleName string
}

func Compile(langImplementation implementation.Implementation, ruleName, input string) (*Query, error) {
	return &Query{ruleName: ruleName}, nil
}
