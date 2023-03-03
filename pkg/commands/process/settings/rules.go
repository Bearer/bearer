package settings

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/bearer/bearer/pkg/flag"
	"github.com/rs/zerolog/log"
	"golang.org/x/exp/maps"
	"gopkg.in/yaml.v3"
)

var (
	defaultRuleType          = "risk"
	defaultAuxiliaryRuleType = "verifier"
)

func loadRules(externalRuleDirs []string, options flag.RuleOptions) (map[string]*Rule, map[string]*Rule, error) {
	definitions := make(map[string]RuleDefinition)
	builtInDefinitions := make(map[string]RuleDefinition)

	if err := loadRuleDefinitions(definitions, rulesFs); err != nil {
		return nil, nil, fmt.Errorf("error loading default rules: %s", err)
	}
	// add default documentation urls for default rules
	for id, definition := range definitions {
		if definition.Metadata.DocumentationUrl == "" {
			definitions[id].Metadata.DocumentationUrl = "https://docs.bearer.com/reference/rules/" + id
		}
	}

	if err := loadRuleDefinitions(builtInDefinitions, builtInRulesFs); err != nil {
		return nil, nil, fmt.Errorf("error loading default built-in rules: %s", err)
	}

	for _, dir := range externalRuleDirs {
		if strings.HasPrefix(dir, "~/") {
			dirname, _ := os.UserHomeDir()
			dir = filepath.Join(dirname, dir[2:])
		}
		log.Debug().Msgf("loading external rules from: %s", dir)
		if err := loadRuleDefinitions(definitions, os.DirFS(dir)); err != nil {
			return nil, nil, fmt.Errorf("error loading external rules from %s: %w", dir, err)
		}
	}

	if err := validateRuleOptionIDs(options, definitions, builtInDefinitions); err != nil {
		return nil, nil, err
	}

	enabledRules := getEnabledRules(options, definitions, nil)
	builtInRules := getEnabledRules(options, builtInDefinitions, enabledRules)

	return buildRules(builtInDefinitions, builtInRules), buildRules(definitions, enabledRules), nil
}

func loadRuleDefinitions(definitions map[string]RuleDefinition, dir fs.FS) error {
	return fs.WalkDir(dir, ".", func(path string, dirEntry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if dirEntry.IsDir() {
			return nil
		}

		filename := dirEntry.Name()
		ext := filepath.Ext(filename)

		if ext != ".yaml" && ext != ".yml" {
			return nil
		}

		entry, err := fs.ReadFile(dir, path)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %w", path, err)
		}

		var ruleDefinition RuleDefinition
		err = yaml.Unmarshal(entry, &ruleDefinition)
		if err != nil {
			return fmt.Errorf("failed to unmarshal rule %s: %w", path, err)
		}

		id := ruleDefinition.Metadata.ID

		if _, exists := definitions[id]; exists {
			return fmt.Errorf("duplicate rule ID %s", id)
		}

		definitions[id] = ruleDefinition

		return nil
	})
}

func validateRuleOptionIDs(
	options flag.RuleOptions,
	definitions map[string]RuleDefinition,
	builtInDefinitions map[string]RuleDefinition,
) error {
	var invalidRuleIDs []string

	for id := range options.OnlyRule {
		_, existsInDefinition := definitions[id]
		_, existsInBuiltInDefinition := builtInDefinitions[id]

		if !existsInBuiltInDefinition && !existsInDefinition {
			invalidRuleIDs = append(invalidRuleIDs, id)
		}
	}

	for id := range options.SkipRule {
		_, existsInDefinition := definitions[id]
		_, existsInBuiltInDefinition := builtInDefinitions[id]

		if !existsInBuiltInDefinition && !existsInDefinition {
			invalidRuleIDs = append(invalidRuleIDs, id)
		}
	}

	if len(invalidRuleIDs) > 0 {
		return fmt.Errorf("invalid rule IDs in only/skip option: %s", strings.Join(invalidRuleIDs, ","))
	}

	return nil
}

func getEnabledRules(options flag.RuleOptions, definitions map[string]RuleDefinition, rules map[string]struct{}) map[string]struct{} {
	enabledRules := make(map[string]struct{})

	for _, definition := range definitions {
		id := definition.Metadata.ID

		if definition.Disabled {
			continue
		}

		for ruleId := range rules {
			enabledRules[ruleId] = struct{}{}
		}

		if len(options.OnlyRule) > 0 && !options.OnlyRule[id] {
			continue
		}

		if options.SkipRule[id] {
			continue
		}

		enabledRules[id] = struct{}{}

		for _, dependencyID := range definition.Detectors {
			enabledRules[dependencyID] = struct{}{}
		}

	}

	return enabledRules
}

func buildRules(definitions map[string]RuleDefinition, enabledRules map[string]struct{}) map[string]*Rule {
	rules := make(map[string]*Rule)

	for _, definition := range definitions {
		id := definition.Metadata.ID

		if _, enabled := enabledRules[id]; !enabled {
			continue
		}

		ruleType := definition.Type
		if len(ruleType) == 0 {
			ruleType = defaultRuleType
		}

		rules[id] = &Rule{
			Id:                      id,
			Type:                    ruleType,
			AssociatedRecipe:        definition.Metadata.AssociatedRecipe,
			Trigger:                 definition.Trigger,
			SkipDataTypes:           definition.SkipDataTypes,
			OnlyDataTypes:           definition.OnlyDataTypes,
			Severity:                mapSeverityKeysToCategories(definition.Severity),
			Description:             definition.Metadata.Description,
			RemediationMessage:      definition.Metadata.RemediationMessage,
			Stored:                  definition.Stored,
			Detectors:               definition.Detectors,
			Processors:              definition.Processors,
			AutoEncrytPrefix:        definition.AutoEncrytPrefix,
			CWEIDs:                  definition.Metadata.CWEIDs,
			Languages:               definition.Languages,
			ParamParenting:          definition.ParamParenting,
			Patterns:                definition.Patterns,
			DocumentationUrl:        definition.Metadata.DocumentationUrl,
			OmitParentContent:       definition.OmitParentContent,
			TriggerRuleOnPresenceOf: definition.TriggerRuleOnPresenceOf,
		}

		for _, auxiliaryDefinition := range definition.Auxiliary {
			rules[auxiliaryDefinition.Id] = &Rule{
				Type:           defaultAuxiliaryRuleType,
				Languages:      definition.Languages,
				ParamParenting: auxiliaryDefinition.ParamParenting,
				Patterns:       auxiliaryDefinition.Patterns,
				Stored:         auxiliaryDefinition.Stored,
			}
		}
	}

	return rules
}

func mapSeverityKeysToCategories(ruleSeverity map[string]string) map[string]string {
	// translate data category attributes to data category names
	for _, key := range maps.Keys(ruleSeverity) {
		switch key {
		case "PD":
			ruleSeverity["Personal Data"] = ruleSeverity[key]
		case "PDS":
			ruleSeverity["Personal Data (Sensitive)"] = ruleSeverity[key]
		default:
		}
	}

	return ruleSeverity
}
