package policies

import (
	"github.com/bearer/curio/pkg/classification/db"
	"github.com/bearer/curio/pkg/commands/process/settings"
	"github.com/bearer/curio/pkg/util/rego"
	regolib "github.com/open-policy-agent/opa/rego"

	"github.com/bearer/curio/pkg/report/output/dataflow"
)

type PolicyInput struct {
	PolicyName        string             `json:"policy_name" yaml:"policy_name"`
	PolicyId          string             `json:"policy_id" yaml:"policy_id"`
	PolicyDescription string             `json:"policy_description" yaml:"policy_description"`
	Dataflow          *dataflow.DataFlow `json:"dataflow" yaml:"dataflow"`
	DataCategories    []db.DataCategory  `json:"data_categories" yaml:"data_categories"`
}

func GetOutput(dataflow *dataflow.DataFlow, config settings.Config) (output []regolib.Vars, err error) {
	for _, policy := range config.Policies {
		result, err := rego.RunQuery(policy.Query,
			PolicyInput{
				PolicyName:        policy.Name,
				PolicyDescription: policy.Description,
				PolicyId:          policy.Id,
				Dataflow:          dataflow,
				DataCategories:    db.Default().DataCategories,
			},
			policy.Modules.ToRegoModules())
		if err != nil {
			return nil, err
		}

		output = append(output, result)
	}

	return output, nil
}
