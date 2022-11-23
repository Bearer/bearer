package policies

import (
	"context"
	"encoding/json"

	"github.com/bearer/curio/pkg/classification/db"
	"github.com/bearer/curio/pkg/commands/process/settings"
	"github.com/bearer/curio/pkg/report/output/dataflow"
	"github.com/open-policy-agent/opa/rego"
	"github.com/rs/zerolog/log"
)

type PolicyInput struct {
	PolicyId       string             `json:"policy_id" yaml:"policy_id"`
	Dataflow       *dataflow.DataFlow `json:"dataflow" yaml:"dataflow"`
	DataCategories []db.DataCategory  `json:"data_categories" yaml:"data_categories"`
}

type PolicyOutput struct {
	LineNumber    string `json:"line_number,omitempty" yaml:"line_number,omitempty"`
	Filename      string `json:"filename,omitempty" yaml:"filename,omitempty"`
	CategoryGroup string `json:"category_group,omitempty" yaml:"category_group,omitempty"`
}

type PolicyResult struct {
	PolicyName        string `json:"policy_name" yaml:"policy_name"`
	PolicyDescription string `json:"policy_description" yaml:"policy_description"`
	LineNumber        string `json:"line_number,omitempty" yaml:"line_number,omitempty"`
	Filename          string `json:"filename,omitempty" yaml:"filename,omitempty"`
	CategoryGroup     string `json:"category_group,omitempty" yaml:"category_group,omitempty"`
}

func GetOutput(dataflow *dataflow.DataFlow, config settings.Config) (map[string][]PolicyResult, error) {
	ctx := context.TODO()

	// policy results grouped by severity (critical, high, ...)
	result := make(map[string][]PolicyResult)

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
					PolicyId:       policy.Id,
					Dataflow:       dataflow,
					DataCategories: db.Default().DataCategories,
				},
			),
		)
		if err != nil {
			return nil, err
		}

		log.Debug().Msgf("result %#v", rs)

		if len(rs) > 0 {
			jsonRes, err := json.Marshal(rs[0].Bindings)
			if err != nil {
				return nil, err
			}

			var policyGrouping map[string][]PolicyOutput
			err = json.Unmarshal(jsonRes, &policyGrouping)
			if err != nil {
				return nil, err
			}

			for _, severity := range []string{
				settings.LevelCritical,
				settings.LevelHigh,
				settings.LevelMedium,
				settings.LevelLow,
			} {
				for _, policyOutput := range policyGrouping[severity] {
					policyResult := PolicyResult{
						PolicyName:        policy.Name,
						PolicyDescription: policy.Description,
						Filename:          policyOutput.Filename,
						CategoryGroup:     policyOutput.CategoryGroup,
					}

					result[severity] = append(result[severity], policyResult)
				}
			}
		}
	}

	return result, nil
}
