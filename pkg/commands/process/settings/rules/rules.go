package rules

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"

	"github.com/bearer/bearer/pkg/commands/process/settings"
	"github.com/bearer/bearer/pkg/engine"
	flagtypes "github.com/bearer/bearer/pkg/flag/types"
	"github.com/bearer/bearer/pkg/report/customdetectors"
	"github.com/bearer/bearer/pkg/util/set"
	"github.com/bearer/bearer/pkg/version_check"
)

const (
	defaultRuleType          = customdetectors.TypeRisk
	defaultAuxiliaryRuleType = customdetectors.TypeVerifier
)

type LoadRulesResult struct {
	BuiltInRules       map[string]*settings.Rule
	Rules              map[string]*settings.Rule
	LoadedRuleCount    int
	CacheUsed          bool
	BearerRulesVersion string
}

//go:embed built_in/*
var builtInRulesFS embed.FS

func Load(
	externalRuleDirs []string,
	options flagtypes.RuleOptions,
	versionMeta *version_check.VersionMeta,
	engine engine.Engine,
	force bool,
	foundLanguageIDs []string,
) (
	result LoadRulesResult,
	err error,
) {
	definitions := make(map[string]settings.RuleDefinition)
	builtInDefinitions := make(map[string]settings.RuleDefinition)

	if versionMeta.Rules.Version != nil {
		result.BearerRulesVersion = *versionMeta.Rules.Version
	}

	log.Debug().Msg("Loading rules")

	count := 0

	remoteCount, err := loadDefinitionsFromRemote(definitions, options, versionMeta)
	if err != nil {
		return result, fmt.Errorf("error loading remote rules: %w", err)
	}

	count += remoteCount

	if _, err := loadCustomDefinitions(builtInDefinitions, true, builtInRulesFS, nil); err != nil {
		return result, fmt.Errorf("error loading built-in rules: %w", err)
	}

	for _, dir := range externalRuleDirs {
		if strings.HasPrefix(dir, "~/") {
			dirname, _ := os.UserHomeDir()
			dir = filepath.Join(dirname, dir[2:])
		}

		log.Debug().Msgf("loading external rules from: %s", dir)
		externalCount, err := loadCustomDefinitions(definitions, false, os.DirFS(dir), foundLanguageIDs)
		if err != nil {
			return result, fmt.Errorf("external rules %w", err)
		}

		count += externalCount
	}

	if err := validateRuleOptionIDs(options, definitions, builtInDefinitions); err != nil {
		return result, err
	}

	enabledRules := getEnabledRules(options, definitions, nil)
	builtInRules := getEnabledRules(options, builtInDefinitions, enabledRules)

	result.Rules = BuildRules(definitions, enabledRules)
	result.BuiltInRules = BuildRules(builtInDefinitions, builtInRules)
	result.LoadedRuleCount = count

	for _, definition := range definitions {
		id := definition.Metadata.ID
		if _, enabled := enabledRules[id]; enabled {
			definitionYAML, err := yaml.Marshal(&definition)
			if err != nil {
				return result, err
			}

			if err := engine.LoadRule(string(definitionYAML)); err != nil {
				return result, fmt.Errorf("engine failed to load rule %s: %w", id, err)
			}
		}
	}

	return result, nil
}

func getFilterRuleReferences(definition *settings.RuleDefinition) set.Set[string] {
	result := set.New[string]()

	var addFilter func(filter settings.PatternFilter)

	addPatterns := func(patterns []settings.RulePattern) {
		for _, pattern := range patterns {
			for _, filter := range pattern.Filters {
				addFilter(filter)
			}
		}
	}

	addFilter = func(filter settings.PatternFilter) {
		if filter.Detection != "" {
			result.Add(filter.Detection)
		}

		if filter.Not != nil {
			addFilter(*filter.Not)
		}

		for _, subFilter := range filter.Either {
			addFilter(subFilter)
		}
	}

	addPatterns(definition.Patterns)
	for _, auxiliaryDefinition := range definition.Auxiliary {
		addPatterns(auxiliaryDefinition.Patterns)
	}

	return result
}

func getSanitizers(definition *settings.RuleDefinition) set.Set[string] {
	result := set.New[string]()

	if definition.SanitizerRuleID != "" {
		result.Add(definition.SanitizerRuleID)
	}

	for _, auxiliaryDefinition := range definition.Auxiliary {
		if auxiliaryDefinition.SanitizerRuleID != "" {
			result.Add(auxiliaryDefinition.SanitizerRuleID)
		}
	}

	return result
}

func getEnabledRules(
	options flagtypes.RuleOptions,
	definitions map[string]settings.RuleDefinition,
	rules map[string]struct{},
) map[string]struct{} {
	enabledRules := make(map[string]struct{})

	for ruleId := range rules {
		enabledRules[ruleId] = struct{}{}
	}

	var enableRule func(definition settings.RuleDefinition)
	enableRule = func(definition settings.RuleDefinition) {
		if definition.Disabled {
			return
		}

		id := definition.Metadata.ID

		if _, alreadyEnabled := enabledRules[id]; alreadyEnabled {
			return
		}

		enabledRules[id] = struct{}{}

		for _, dependencyID := range definition.Detectors {
			enabledRules[dependencyID] = struct{}{}
		}

		for _, importedRuleID := range definition.Imports {
			if importedDefinition, exists := definitions[importedRuleID]; exists {
				enableRule(importedDefinition)
			}
		}
	}

	for _, definition := range definitions {
		id := definition.Metadata.ID

		if len(options.OnlyRule) > 0 && !options.OnlyRule[id] {
			continue
		}

		if options.SkipRule[id] {
			continue
		}

		enableRule(definition)
	}

	return enabledRules
}

func BuildRules(
	definitions map[string]settings.RuleDefinition,
	enabledRules map[string]struct{},
) map[string]*settings.Rule {
	rules := make(map[string]*settings.Rule)

	for _, definition := range definitions {
		id := definition.Metadata.ID

		if _, enabled := enabledRules[id]; !enabled {
			continue
		}

		ruleType := definition.Type
		if len(ruleType) == 0 {
			ruleType = defaultRuleType
		}

		// build rule trigger
		ruleTrigger := settings.RuleTrigger{
			MatchOn:           settings.PRESENCE,
			DataTypesRequired: false,
		}

		if definition.Trigger != nil {
			if definition.Trigger.MatchOn != nil {
				ruleTrigger.MatchOn = *definition.Trigger.MatchOn
			}
			if definition.Trigger.DataTypesRequired != nil {
				ruleTrigger.DataTypesRequired = *definition.Trigger.DataTypesRequired
			}

			// concat any required detections
			ruleTrigger.RequiredDetections = definition.Trigger.RequiredDetections
			if definition.Trigger.RequiredDetection != nil {
				ruleTrigger.RequiredDetections = append(ruleTrigger.RequiredDetections, *definition.Trigger.RequiredDetection)
			}
		}

		isLocal := false
		for _, rulePattern := range definition.Patterns {
			if strings.Contains(rulePattern.Pattern, "$<DATA_TYPE>") {
				isLocal = true
				break
			}
		}

		rules[id] = &settings.Rule{
			Id:                 id,
			Type:               ruleType,
			AssociatedRecipe:   definition.Metadata.AssociatedRecipe,
			Trigger:            ruleTrigger,
			IsLocal:            isLocal,
			SkipDataTypes:      definition.SkipDataTypes,
			OnlyDataTypes:      definition.OnlyDataTypes,
			Severity:           definition.Severity,
			Description:        definition.Metadata.Description,
			RemediationMessage: definition.Metadata.RemediationMessage,
			Stored:             definition.Stored,
			Detectors:          definition.Detectors,
			Processors:         definition.Processors,
			AutoEncrytPrefix:   definition.AutoEncrytPrefix,
			CWEIDs:             definition.Metadata.CWEIDs,
			Languages:          definition.Languages,
			ParamParenting:     definition.ParamParenting,
			Patterns:           definition.Patterns,
			SanitizerRuleID:    definition.SanitizerRuleID,
			DocumentationUrl:   definition.Metadata.DocumentationUrl,
			HasDetailedContext: definition.HasDetailedContext,
			DependencyCheck:    definition.DependencyCheck,
			Dependency:         definition.Dependency,
		}

		for _, auxiliaryDefinition := range definition.Auxiliary {
			rules[auxiliaryDefinition.Id] = &settings.Rule{
				Id:              auxiliaryDefinition.Id,
				Type:            defaultAuxiliaryRuleType,
				Languages:       definition.Languages,
				ParamParenting:  auxiliaryDefinition.ParamParenting,
				Patterns:        auxiliaryDefinition.Patterns,
				SanitizerRuleID: auxiliaryDefinition.SanitizerRuleID,
				Stored:          auxiliaryDefinition.Stored,
				IsAuxilary:      true,
			}
		}
	}

	return rules
}
