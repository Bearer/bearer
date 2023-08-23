package settings

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/xeipuuv/gojsonschema"
)

const SCHEMA_URL = "https://raw.githubusercontent.com/Bearer/bearer-rules/main/scripts/rule_schema.json"

func ValidateRule(entry []byte, filename string) string {
	validationStr := &strings.Builder{}
	validationStr.WriteString(fmt.Sprintf("Failed to load %s\nValidating against %s\n\n", filename, SCHEMA_URL))
	schema, err := loadSchema(SCHEMA_URL)
	if err != nil {
		validationStr.WriteString("Could not load schema to validate")
		return validationStr.String()
	}

	jsonData, err := yaml.YAMLToJSON(entry)
	if err != nil {
		validationStr.WriteString("File fdormat is invalid")
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
		validationStr.WriteString(fmt.Sprintf("%s validation issues found:\n", filename))
		for _, desc := range result.Errors() {
			validationStr.WriteString(fmt.Sprintf("- %s\n", desc))
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
