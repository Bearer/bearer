package query

import (
	"github.com/bearer/bearer/new/language/implementation"
	builderinput "github.com/bearer/bearer/new/language/patternquery/builder/input"
	"github.com/bearer/bearer/pkg/ast/languages/ruby/patterns"
)

type Query struct {
}

func Compile(langImplementation implementation.Implementation, ruleName, input string) (*Query, error) {
	processedInput, inputParams, err := builderinput.Process(langImplementation, input)
	if err != nil {
		return nil, err
	}

	if err := patterns.CompileRule(walker, inputParams, ruleName, []byte(processedInput), rootNode, writer); err != nil {
		return nil, err
	}

	return &Query{}, nil
}
