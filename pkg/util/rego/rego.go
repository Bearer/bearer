package rego

import (
	"context"
	"fmt"

	"github.com/open-policy-agent/opa/rego"
	"github.com/rs/zerolog/log"
)

type Module struct {
	Name    string
	Content string
}

func RunQuery(query string, input interface{}, modules []Module) (rego.Vars, error) {
	ctx := context.TODO()
	log.Error().Msg("RG: setting options")

	options := []func(r *rego.Rego){rego.Query(query)}
	for _, module := range modules {
		options = append(options, rego.Module(module.Name, module.Content))
	}

	r := rego.New(options...)
	log.Error().Msg("RG: compiling query")
	regoQuery, err := r.PrepareForEval(ctx)
	if err != nil {
		return nil, err
	}
	log.Error().Msg("RG: evaling")

	// Create a prepared query that can be evaluated.
	rs, err := regoQuery.Eval(ctx, rego.EvalInput(input))
	if err != nil {
		return nil, err
	}
	log.Error().Msg("RG: done eval")

	if len(rs) != 1 {
		return nil, fmt.Errorf("expected single result from query got %d results %#v", len(rs), rs)
	}

	return rs[0].Bindings, nil
}
