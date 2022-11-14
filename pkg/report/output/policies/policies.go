package policies

import (
	"context"

	"github.com/bearer/curio/pkg/report/output/dataflow"
	"github.com/open-policy-agent/opa/rego"
	"github.com/rs/zerolog/log"
)

const opaConfig string = `
services:
  acmecorp:
    url: https://example.com/control-plane-api/v1
    response_header_timeout_seconds: 5
    credentials:
      bearer:
        token: "bGFza2RqZmxha3NkamZsa2Fqc2Rsa2ZqYWtsc2RqZmtramRmYWxkc2tm"

labels:
  app: myapp
  region: west
  environment: production

bundles:
  authz:
    service: acmecorp
    resource: bundles/http/example/authz.tar.gz
    persist: true
    polling:
      min_delay_seconds: 60
      max_delay_seconds: 120
    signing:
      keyid: global_key
      scope: write

decision_logs:
  service: acmecorp
  reporting:
    min_delay_seconds: 300
    max_delay_seconds: 600

status:
  service: acmecorp

default_decision: /http/example/authz/allow

persistence_directory: /var/opa

keys:
  global_key:
    algorithm: RS256
    key: <PEM_encoded_public_key>
    scope: read

caching:
  inter_query_builtin_cache:
    max_size_bytes: 10000000

distributed_tracing:
  type: grpc
  address: localhost:4317
  service_name: opa
  sample_percentage: 50
  encryption: "off"
`

const dataflowQuery = `
package demo

import future.keywords

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

const dataflowQuery2 = `
package example.authz

import future.keywords.if
import future.keywords.in


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

func GetPolicies(dataflow *dataflow.DataFlow) (rego.Vars, error) {
	input := ` 
		some detector2 in input.risks
		detector2.detector_id == "detect_ruby_logger"
		locations := detector2.data_types[_].locations

		result = {"warning": {"message": "there are logger leaks detected" ,  "count": count(locations) , "locations": locations}}{
			count(locations) > 0
		}
		`
	result, err := rego.New(rego.Query(input), rego.Input(*dataflow), rego.Imports([]string{"future.keywords"})).Eval(context.Background())
	if err != nil {
		log.Debug().Msgf("got error %s", err)
		return rego.Vars{}, err
	}

	log.Debug().Msgf("result %#v", result)
	log.Debug().Msgf("result %#v")

	if len(result) > 0 {
		return result[0].Bindings, nil
	}

	return rego.Vars{}, nil
}

func GetDataflow(data []interface{}) (interface{}, error) {
	ctx := context.TODO()

	r := rego.New(
		rego.Query("x = data.example.authz.hostnames"), rego.Module("example.rego", dataflowQuery2))

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

	// result, err := rego.New(rego.Query(input), rego.Input(data), rego.ShallowInlining(false), rego.Imports([]string{"future.keywords"})).Eval(context.Background())
	// if err != nil {
	// 	log.Debug().Msgf("got error %s", err)
	// 	return rego.Vars{}, err
	// }

	return rs[0].Bindings, nil
}
