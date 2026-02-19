package policies

import (
	"embed"
	"fmt"

	"gopkg.in/yaml.v3"

	"github.com/bearer/bearer/pkg/commands/process/settings"
)

//go:embed *.rego
var policiesFS embed.FS

func Load() (map[string]*settings.Policy, error) {
	policies, err := loadDefaultPolicies()
	if err != nil {
		return nil, err
	}

	for _, policy := range policies {
		for _, module := range policy.Modules {
			if module.Path != "" {
				content, err := policiesFS.ReadFile(module.Path)
				if err != nil {
					return nil, err
				}

				module.Content = string(content)
			}
		}
	}

	return policies, nil
}

//go:embed policies.yml
var defaultPolicies []byte

func loadDefaultPolicies() (map[string]*settings.Policy, error) {
	policiesByType := make(map[string]*settings.Policy)

	var policies []*settings.Policy
	err := yaml.Unmarshal(defaultPolicies, &policies)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal policy file %s", err)
	}

	for _, policy := range policies {
		policiesByType[policy.Type] = policy
	}

	return policiesByType, nil
}
