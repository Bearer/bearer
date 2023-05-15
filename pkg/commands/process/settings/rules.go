package settings

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/bearer/bearer/new/language"
	languagetypes "github.com/bearer/bearer/new/language/types"
	"github.com/bearer/bearer/pkg/commands/process/settings/rules"
	"github.com/bearer/bearer/pkg/flag"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

var (
	defaultRuleType          = "risk"
	defaultAuxiliaryRuleType = "verifier"
)

func GetSupportedRuleLanguages() map[string]bool {
	return map[string]bool{
		"ruby":       true,
		"javascript": true,
		"typescript": true,
	}
}

func RefreshRules(config Config, externalRuleDirs []string, options flag.RuleOptions, foundLanguages []string) (err error) {
	result, err := loadRules(externalRuleDirs, options, foundLanguages, true)
	config.BuiltInRules = result.BuiltInRules
	config.Rules = result.Rules
	config.BearerRulesVersion = result.BearerRulesVersion

	return
}

func loadRules(
	externalRuleDirs []string,
	options flag.RuleOptions,
	foundLanguages []string,
	force bool) (
	result LoadRulesResult,
	err error,
) {
	definitions := make(map[string]rules.RuleDefinition)
	builtInDefinitions := make(map[string]rules.RuleDefinition)
	ruleLanguages := make(map[string]bool)

	if !options.DisableDefaultRules {
		bearerRulesDir := bearerRulesDir()
		if !force && cachedRulesExist(bearerRulesDir) {
			result.CacheUsed = true
			err := filepath.WalkDir(bearerRulesDir, func(filePath string, d fs.DirEntry, err error) error {
				if !d.IsDir() {
					file, err := os.Open(filepath.Join(bearerRulesDir, d.Name()))
					if err != nil {
						return err
					}
					if ruleLanguages, err = ReadRuleDefinitions(definitions, file); err != nil {
						return err
					}
				}
				return nil
			})

			if err != nil {
				return result, fmt.Errorf("error loading rules from cache: %s", err)
			}

			supportedLanguages := GetSupportedRuleLanguages()
			for _, foundLang := range foundLanguages {
				if !supportedLanguages[foundLang] {
					// no rule support for this language e.g. CSS, plain text
					continue
				}

				if !ruleLanguages[foundLang] {
					definitions = make(map[string]rules.RuleDefinition)
					result.CacheUsed = false // re-cache rules
				}
			}
		}

		if !result.CacheUsed {
			if err := cleanupRuleDirFiles(bearerRulesDir); err != nil {
				return result, fmt.Errorf("error cleaning rules cache: %s", err)
			}

			tagVersion, err := LoadRuleDefinitionsFromGitHub(definitions, foundLanguages)
			if err != nil {
				return result, fmt.Errorf("error loading rules: %s", err)
			}

			result.BearerRulesVersion = tagVersion
		}

		// add default documentation urls for default rules
		for id, definition := range definitions {
			if definition.Metadata.DocumentationUrl == "" {
				definitions[id].Metadata.DocumentationUrl = "https://docs.bearer.com/reference/rules/" + id
			}
		}
	}

	if err := loadRuleDefinitionsFromDir(builtInDefinitions, buildInRulesFs); err != nil {
		return result, fmt.Errorf("error loading built-in rules: %w", err)
	}

	for _, dir := range externalRuleDirs {
		if strings.HasPrefix(dir, "~/") {
			dirname, _ := os.UserHomeDir()
			dir = filepath.Join(dirname, dir[2:])
		}
		log.Debug().Msgf("loading external rules from: %s", dir)
		if err := loadRuleDefinitionsFromDir(definitions, os.DirFS(dir)); err != nil {
			return result, fmt.Errorf("error loading external rules from %s: %w", dir, err)
		}
	}

	if err := validateRuleOptionIDs(options, definitions, builtInDefinitions); err != nil {
		return result, err
	}

	enabledRules := getEnabledRules(options, definitions, nil)
	builtInRules := getEnabledRules(options, builtInDefinitions, enabledRules)

	result.Rules, err = BuildRules(definitions, enabledRules)
	if err != nil {
		return result, err
	}

	result.BuiltInRules, err = BuildRules(builtInDefinitions, builtInRules)
	if err != nil {
		return result, err
	}

	return result, nil
}

func loadRuleDefinitionsFromDir(definitions map[string]rules.RuleDefinition, dir fs.FS) error {
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

		var ruleDefinition rules.RuleDefinition
		err = yaml.Unmarshal(entry, &ruleDefinition)
		if err != nil {
			return fmt.Errorf("failed to unmarshal rule %s: %w", path, err)
		}

		if ruleDefinition.Metadata == nil {
			log.Debug().Msgf("Rule file has invalid metadata %s", path)
			return nil
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
	definitions map[string]rules.RuleDefinition,
	builtInDefinitions map[string]rules.RuleDefinition,
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

func getEnabledRules(options flag.RuleOptions, definitions map[string]rules.RuleDefinition, rules map[string]struct{}) map[string]struct{} {
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

func BuildRules(definitions map[string]rules.RuleDefinition, enabledRules map[string]struct{}) (map[string]*rules.Rule, error) {
	result := make(map[string]*rules.Rule)

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
		ruleTrigger := rules.RuleTrigger{
			MatchOn:           rules.PRESENCE,
			DataTypesRequired: false,
		}
		if definition.Trigger != nil {
			if definition.Trigger.MatchOn != nil {
				ruleTrigger.MatchOn = *definition.Trigger.MatchOn
			}
			if definition.Trigger.DataTypesRequired != nil {
				ruleTrigger.DataTypesRequired = *definition.Trigger.DataTypesRequired
			}
			if definition.Trigger.RequiredDetection != nil {
				ruleTrigger.RequiredDetection = definition.Trigger.RequiredDetection
			}
		}

		isLocal := false
		for _, rulePattern := range definition.Patterns {
			if strings.Contains(rulePattern.Pattern, "$<DATA_TYPE>") {
				isLocal = true
				break
			}
		}

		lang, err := language.Get(definition.Languages[0])
		if err != nil {
			return nil, fmt.Errorf("error loading rule %s: %w", id, err)
		}

		patterns, err := compilePatterns(lang, definition.Patterns)
		if err != nil {
			return nil, fmt.Errorf("error loading rule %s: %w", id, err)
		}

		result[id] = &rules.Rule{
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
			Patterns:           patterns,
			DocumentationUrl:   definition.Metadata.DocumentationUrl,
			HasDetailedContext: definition.HasDetailedContext,
		}

		for _, auxiliaryDefinition := range definition.Auxiliary {
			auxPatterns, err := compilePatterns(lang, auxiliaryDefinition.Patterns)
			if err != nil {
				return nil, fmt.Errorf("error loading rule %s: %w", id, err)
			}

			result[auxiliaryDefinition.Id] = &rules.Rule{
				Id:             auxiliaryDefinition.Id,
				Type:           defaultAuxiliaryRuleType,
				Languages:      definition.Languages,
				ParamParenting: auxiliaryDefinition.ParamParenting,
				Patterns:       auxPatterns,
				Stored:         auxiliaryDefinition.Stored,
				IsAuxilary:     true,
			}
		}
	}

	return result, nil
}

func cachedRulesExist(bearerRulesDir string) bool {
	_, err := os.Stat(bearerRulesDir)
	return err == nil
}

func cleanupRuleDirFiles(bearerRulesDir string) error {
	return os.RemoveAll(bearerRulesDir)
}

func bearerRulesDir() string {
	return filepath.Join(os.TempDir(), "bearer-rules")
}

func compilePatterns(lang languagetypes.Language, sourcePatterns []rules.RuleDefinitionPattern) ([]rules.RulePattern, error) {
	patterns := make([]rules.RulePattern, len(sourcePatterns))

	for i, sourcePattern := range sourcePatterns {
		// query, err := lang.CompilePatternQuery(sourcePattern.Pattern)
		// if err != nil {
		// 	return nil, fmt.Errorf("error compiling pattern %s: %s", sourcePattern.Pattern, err)
		// }

		// log.Error().Msgf("Q: %#v", query)

		patterns[i] = rules.RulePattern{
			Pattern: sourcePattern.Pattern,
			Filters: sourcePattern.Filters,
			// Query:   query,
		}
	}

	return patterns, nil
}
