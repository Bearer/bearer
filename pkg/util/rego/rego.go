package rego

import (
	"context"
	"fmt"

	"github.com/open-policy-agent/opa/rego"
)

type Module struct {
	Name    string
	Content string
}

func RunQuery(query string, input interface{}, modules []Module) (rego.Vars, error) {
	ctx := context.TODO()

	options := []func(r *rego.Rego){rego.Query(query)}
	for _, module := range modules {
		options = append(options, rego.Module(module.Name, module.Content))
	}

	r := rego.New(options...)
	regoQuery, err := r.PrepareForEval(ctx)
	if err != nil {
		return nil, err
	}

	// Create a prepared query that can be evaluated.
	rs, err := regoQuery.Eval(ctx, rego.EvalInput(input))
	if err != nil {
		return nil, err
	}

	if len(rs) != 1 {
		return nil, fmt.Errorf("expected single result from query got %d results %#v:\n%s", len(rs), rs, query)
	}

	return rs[0].Bindings, nil
}
