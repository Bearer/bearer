package policies

import (
	"context"

	"github.com/bearer/curio/pkg/classification/db"
	"github.com/bearer/curio/pkg/commands/process/settings"
	"github.com/bearer/curio/pkg/report/output/dataflow"
	"github.com/open-policy-agent/opa/rego"
	"github.com/rs/zerolog/log"
)

type PolicyInput struct {
	Dataflow       *dataflow.DataFlow `json:"dataflow"`
	DataCategories []db.DataCategory  `json:"data_categories"`
}

func GetOutput(dataflow *dataflow.DataFlow, config settings.Config) ([]rego.Vars, error) {
	ctx := context.TODO()

	var result []rego.Vars

	for _, policy := range config.Policies {
		options := []func(r *rego.Rego){rego.Query(policy.Query)}
		for _, module := range policy.Modules {
			options = append(options, rego.Module(module.Name, module.Content))
		}

		r := rego.New(options...)
		query, err := r.PrepareForEval(ctx)
		if err != nil {
			return nil, err
		}

		// Create a prepared query that can be evaluated.
		rs, err := query.Eval(
			ctx,
			rego.EvalInput(
				PolicyInput{
					Dataflow:       dataflow,
					DataCategories: db.Default().DataCategories,
				},
			),
		)
		if err != nil {
			return nil, err
		}

		log.Debug().Msgf("result %#v", rs)

		result = append(result, rs[0].Bindings)
	}

	// Create a prepared query that can be evaluated.
	return result, nil
}
