package settings

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"

	"github.com/bearer/bearer/internal/flag"
	"github.com/bearer/bearer/internal/report/customdetectors"
	"github.com/bearer/bearer/internal/util/output"
	"github.com/bearer/bearer/internal/util/set"
	"github.com/bearer/bearer/internal/version_check"
)

const (
	defaultRuleType          = customdetectors.TypeRisk
	defaultAuxiliaryRuleType = customdetectors.TypeVerifier
)

var (
	builtinRuleIDs = []string{
		"datatype",
		"insecure_url",
		"string_literal",
	}
)

func GetSupportedRuleLanguages() map[string]bool {
	return map[string]bool{
		"java":       true,
		"ruby":       true,
		"javascript": true,
		"typescript": true,
	}
}

func loadRules(
	externalRuleDirs []string,
	options flag.RuleOptions,
	versionMeta *version_check.VersionMeta,
	force bool) (
	result LoadRulesResult,
	err error,
) {
	definitions := make(map[string]RuleDefinition)
	builtInDefinitions := make(map[string]RuleDefinition)

	log.Debug().Msg("Loading rules")

	if versionMeta.Rules.Version != nil {
		result.BearerRulesVersion = *versionMeta.Rules.Version

		urls := make([]string, 0, len(versionMeta.Rules.Packages))
		for _, value := range versionMeta.Rules.Packages {
			log.Debug().Msgf("Added rule package URL %s", value)
			urls = append(urls, value)
		}

		err = LoadRuleDefinitionsFromUrls(definitions, urls)
		if err != nil {
			output.Fatal(fmt.Sprintf("Error loading rules: %s", err))
			// sysexit
		}
	} else {
		log.Debug().Msg("No rule packages found")
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
			return result, fmt.Errorf("external rules %w", err)
		}
	}

	if err := validateRuleOptionIDs(options, definitions, builtInDefinitions); err != nil {
		return result, err
	}

	enabledRules := getEnabledRules(options, definitions, nil)
	builtInRules := getEnabledRules(options, builtInDefinitions, enabledRules)

	result.Rules = BuildRules(definitions, enabledRules)
	result.BuiltInRules = BuildRules(builtInDefinitions, builtInRules)

	return result, nil
}

func loadRuleDefinitionsFromDir(definitions map[string]RuleDefinition, dir fs.FS) error {
	loadedDefinitions := make(map[string]RuleDefinition)
	if err := fs.WalkDir(dir, ".", func(path string, dirEntry fs.DirEntry, err error) error {
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
			output.StdErrLog(ValidateRule(entry, filename))
			return fmt.Errorf("rule file was invalid")
		}

		if ruleDefinition.Metadata == nil {
			log.Debug().Msgf("rule file has invalid metadata %s", path)
			return nil
		}

		id := ruleDefinition.Metadata.ID
		if id == "" {
			log.Debug().Msgf("rule file missing metadata.id %s", path)
			return nil
		}

		if _, exists := loadedDefinitions[id]; exists {
			return fmt.Errorf("duplicate rule ID %s", id)
		}

		loadedDefinitions[id] = ruleDefinition

		return nil
	}); err != nil {
		return err
	}

	for id, definition := range loadedDefinitions {
		if validateRuleDefinition(loadedDefinitions, &definition) {
			definitions[id] = definition
		}
	}

	return nil
}

func validateRuleDefinition(allDefinitions map[string]RuleDefinition, definition *RuleDefinition) bool {
	metadata := definition.Metadata

	valid := true
	fail := func(message string) {
		valid = false
		log.Debug().Msgf("%s: %s", metadata.ID, message)
	}

	visibleRuleIDs := set.New[string]()
	visibleRuleIDs.Add(metadata.ID)
	visibleRuleIDs.AddAll(builtinRuleIDs)

	for _, importedID := range definition.Imports {
		visibleRuleIDs.Add(importedID)

		importedDefinition, exists := allDefinitions[importedID]

		if !exists {
			fail(fmt.Sprintf("import of unknown rule '%s'", importedID))
			continue
		}

		if importedDefinition.Type != customdetectors.TypeShared {
			fail(fmt.Sprintf("imported rule '%s' is not of type 'shared'", importedID))
		}
	}

	for _, auxiliaryDefinition := range definition.Auxiliary {
		visibleRuleIDs.Add(auxiliaryDefinition.Id)
	}

	for _, filterRuleID := range getFilterRuleReferences(definition).Items() {
		if !visibleRuleIDs.Has(filterRuleID) {
			fail(fmt.Sprintf("filter references invalid or non-imported rule '%s'", filterRuleID))
		}
	}

	for _, sanitizerRuleID := range getSanitizers(definition).Items() {
		if !visibleRuleIDs.Has(sanitizerRuleID) {
			fail(fmt.Sprintf("sanitizer references invalid or non-imported rule '%s'", sanitizerRuleID))
		}
	}

	if metadata.ID == "" {
		fail("metadata.id must be specified")
	}

	if definition.Type == customdetectors.TypeShared {
		metadata := definition.Metadata
		if metadata != nil {
			if metadata.CWEIDs != nil {
				fail("cwe ids cannot be specified for a shared rule")
			}

			if metadata.RemediationMessage != "" {
				fail("remediation message cannot be specified for a shared rule")
			}
		}

		if definition.Severity != "" {
			fail("severity cannot be specified for a shared rule")
		}
	}

	if !valid {
		log.Debug().Msgf("%s ignored due to validation errors", metadata.ID)
	}

	return valid
}

func getFilterRuleReferences(definition *RuleDefinition) set.Set[string] {
	result := set.New[string]()

	var addFilter func(filter PatternFilter)

	addPatterns := func(patterns []RulePattern) {
		for _, pattern := range patterns {
			for _, filter := range pattern.Filters {
				addFilter(filter)
			}
		}
	}

	addFilter = func(filter PatternFilter) {
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

func getSanitizers(definition *RuleDefinition) set.Set[string] {
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

	for ruleId := range rules {
		enabledRules[ruleId] = struct{}{}
	}

	var enableRule func(definition RuleDefinition)
	enableRule = func(definition RuleDefinition) {
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

func BuildRules(definitions map[string]RuleDefinition, enabledRules map[string]struct{}) map[string]*Rule {
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

		// build rule trigger
		ruleTrigger := RuleTrigger{
			MatchOn:           PRESENCE,
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

		rules[id] = &Rule{
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
			rules[auxiliaryDefinition.Id] = &Rule{
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

func bearerRulesDir() string {
	return filepath.Join(os.TempDir(), "bearer-rules")
}
