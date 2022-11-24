package policies

import (
	"encoding/json"

	"github.com/bearer/curio/pkg/classification/db"
	"github.com/bearer/curio/pkg/commands/process/settings"
	"github.com/bearer/curio/pkg/util/output"
	"github.com/bearer/curio/pkg/util/rego"
	"github.com/rs/zerolog/log"

	"github.com/bearer/curio/pkg/report/output/dataflow"
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
	// policy results grouped by severity (critical, high, ...)
	result := make(map[string][]PolicyResult)

	for _, policy := range config.Policies {
		input := PolicyInput{
			PolicyId:       policy.Id,
			Dataflow:       dataflow,
			DataCategories: db.Default().DataCategories,
		}

		log.Debug().Msgf("running policy")
		bytes, _ := json.Marshal(input)
		log.Debug().Msgf("")
		output.StdOutLogger().Msgf("%s", string(bytes))
		log.Debug().Msgf("")
		log.Debug().Msgf("input is %#v", input)
		log.Debug().Msgf("dataflow is: %#v", input.Dataflow)
		// Create a prepared query that can be evaluated.
		rs, err := rego.RunQuery(policy.Query,
			PolicyInput{
				PolicyId:       policy.Id,
				Dataflow:       dataflow,
				DataCategories: db.Default().DataCategories,
			},
			policy.Modules.ToRegoModules())
		if err != nil {
			return nil, err
		}

		if len(rs) > 0 {
			jsonRes, err := json.Marshal(rs)
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
