package rules

import (
	"bytes"
	"fmt"
	"net/http"
	"path/filepath"
	"sort"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/xeipuuv/gojsonschema"
	"sigs.k8s.io/yaml"

	"github.com/bearer/bearer/pkg/commands/process/settings"
	flagtypes "github.com/bearer/bearer/pkg/flag/types"
	"github.com/bearer/bearer/pkg/report/customdetectors"
	"github.com/bearer/bearer/pkg/util/output"
	"github.com/bearer/bearer/pkg/util/set"
)

const SCHEMA_URL = "https://raw.githubusercontent.com/Bearer/bearer-rules/main/scripts/rule_schema.json"

var builtinRuleIDs = []string{
	"datatype",
	"insecure_url",
	"string_literal",
}

func validateCustomRuleSchema(entry []byte, filename string) string {
	validationStr := &strings.Builder{}
	fmt.Fprintf(validationStr, "Failed to load %s\nValidating against %s\n\n", filename, SCHEMA_URL)
	schema, err := loadSchema(SCHEMA_URL)
	if err != nil {
		validationStr.WriteString("Could not load schema to validate")
		return validationStr.String()
	}

	jsonData, err := yaml.YAMLToJSON(entry)
	if err != nil {
		validationStr.WriteString("File format is invalid")
		return validationStr.String()
	}

	result, err := validateData(jsonData, schema)
	if err != nil {
		validationStr.WriteString("Could not apply validation")
		return validationStr.String()
	}

	if result.Valid() {
		validationStr.WriteString("Format of appears valid but could not be loaded")
	} else {
		fmt.Fprintf(validationStr, "%s validation issues found:\n", filename)
		for _, desc := range result.Errors() {
			fmt.Fprintf(validationStr, "- %s\n", desc)
		}
		fmt.Print("\n")
	}
	return validationStr.String()
}

func loadSchema(url string) (*gojsonschema.Schema, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var buf bytes.Buffer
	_, err = buf.ReadFrom(response.Body)
	if err != nil {
		return nil, err
	}

	loader := gojsonschema.NewStringLoader(buf.String())
	return gojsonschema.NewSchema(loader)
}

func validateData(data []byte, schema *gojsonschema.Schema) (*gojsonschema.Result, error) {
	loader := gojsonschema.NewStringLoader(string(data))
	return schema.Validate(loader)
}

func validateCustomRuleDefinition(allDefinitions map[string]settings.RuleDefinition, definition *settings.RuleDefinition) bool {
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

func validateRuleOptionIDs(
	options flagtypes.RuleOptions,
	definitions map[string]settings.RuleDefinition,
	builtInDefinitions map[string]settings.RuleDefinition,
) error {
	var invalidRuleIDs []string

	for id := range options.OnlyRule {
		_, existsInDefinition := definitions[id]
		_, existsInBuiltInDefinition := builtInDefinitions[id]

		if !existsInBuiltInDefinition && !existsInDefinition {
			invalidRuleIDs = append(invalidRuleIDs, id)
		}
	}
	var invalidSkipRuleIDs []string
	for id := range options.SkipRule {
		_, existsInDefinition := definitions[id]
		_, existsInBuiltInDefinition := builtInDefinitions[id]

		if !existsInBuiltInDefinition && !existsInDefinition {
			invalidSkipRuleIDs = append(invalidSkipRuleIDs, id)
		}
	}

	if len(invalidSkipRuleIDs) > 0 {
		sort.Strings(invalidSkipRuleIDs)
		output.StdErrLog(fmt.Sprintf("Warning: rule IDs %s given to be skipped but were not found", strings.Join(invalidSkipRuleIDs, ",")))
	}
	if len(invalidRuleIDs) > 0 {
		return fmt.Errorf("invalid rule IDs in only option: %s", strings.Join(invalidRuleIDs, ","))
	}

	return nil
}

func isRuleFile(path string) bool {
	if strings.Contains(path, ".snapshots") {
		return false
	}

	ext := filepath.Ext(path)
	if ext != ".yaml" && ext != ".yml" {
		return false
	}

	return strings.Contains(path, "/")
}
