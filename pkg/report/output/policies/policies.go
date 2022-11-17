package policies

import (
	"github.com/bearer/curio/pkg/commands/process/settings"
	"github.com/bearer/curio/pkg/util/rego"
	regolib "github.com/open-policy-agent/opa/rego"

	"github.com/bearer/curio/pkg/report/output/dataflow"
)

func GetOutput(dataflow *dataflow.DataFlow, config settings.Config) (output []regolib.Vars, err error) {
	for _, policy := range config.Policies {
		result, err := rego.RunQuery(policy.Query, dataflow, policy.Modules.ToRegoModules())
		if err != nil {
			return nil, err
		}

		output = append(output, result)
	}

	return output, nil
}
