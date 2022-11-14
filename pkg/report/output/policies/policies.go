package policies

import (
	"context"

	"github.com/bearer/curio/pkg/commands/process/settings"
	"github.com/bearer/curio/pkg/report/output/dataflow"
	"github.com/open-policy-agent/opa/rego"
	"github.com/rs/zerolog/log"
)

func GetPolicies(dataflow *dataflow.DataFlow, config settings.Config) ([]rego.Vars, error) {
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
		rs, err := query.Eval(ctx, rego.EvalInput(dataflow))
		if err != nil {
			return nil, err
		}

		log.Debug().Msgf("result %#v", rs)

		result = append(result, rs[0].Bindings)
	}

	// Create a prepared query that can be evaluated.
	return result, nil
}

func GetDataflow(data []interface{}) (interface{}, error) {
	const dataFlowQuery = `
package example.authz


sites := [
		{
			"region": "east",
			"name": "prod",
			"servers": [
				{
					"name": "web-0",
					"hostname": "hydrogen"
				},
				{
					"name": "web-1",
					"hostname": "helium"
				},
				{
					"name": "db-0",
					"hostname": "lithium"
				}
			]
		},
		{
			"region": "west",
			"name": "smoke",
			"servers": [
				{
					"name": "web-1000",
					"hostname": "beryllium"
				},
				{
					"name": "web-1001",
					"hostname": "boron"
				},
				{
					"name": "db-1000",
					"hostname": "carbon"
				}
			]
		},
		{
			"region": "west",
			"name": "dev",
			"servers": [
				{
					"name": "web-dev",
					"hostname": "nitrogen"
				},
				{
					"name": "db-dev",
					"hostname": "oxygen"
				}
			]
		}
	]
	
hostnames[name] {
    name := sites[_].servers[_].hostname
}
`

	ctx := context.TODO()

	r := rego.New(
		rego.Query("x = data.example.authz.hostnames"), rego.Module("example.rego", dataFlowQuery))

	// Create a prepared query that can be evaluated.
	query, err := r.PrepareForEval(ctx)
	if err != nil {
		return nil, err
	}

	// Create a prepared query that can be evaluated.
	rs, err := query.Eval(ctx, rego.EvalInput(data))
	if err != nil {
		return nil, err
	}

	log.Debug().Msgf("result %#v", rs)

	return rs[0].Bindings, nil
}
