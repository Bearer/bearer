package yamlconfig

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/smacker/go-tree-sitter/yaml"

	sitter "github.com/smacker/go-tree-sitter"

	"github.com/bearer/bearer/pkg/detectors/types"
	"github.com/bearer/bearer/pkg/parser"
	"github.com/bearer/bearer/pkg/parser/interfaces"
	"github.com/bearer/bearer/pkg/report"
	reportinterface "github.com/bearer/bearer/pkg/report/interfaces"
	"github.com/bearer/bearer/pkg/util/file"

	"github.com/bearer/bearer/pkg/parser/sitter/config_variables"
	"github.com/bearer/bearer/pkg/report/detectors"
	"github.com/bearer/bearer/pkg/report/values"
	"github.com/bearer/bearer/pkg/report/variables"
)

var (
	interfacesQuery = parser.QueryMustCompile(yaml.GetLanguage(), `
		(block_mapping (block_mapping_pair
			key: (_) @key
			value: (flow_node [(plain_scalar) (single_quote_scalar) (double_quote_scalar)]) @value) @definition)

		(block_sequence (block_sequence_item (flow_node [(plain_scalar) (single_quote_scalar) (double_quote_scalar)]) @keyValue @definition))
	`)

	filenamePattern      = regexp.MustCompile(`\.ya?ml(\.|$)`)
	i18nExclusionPattern = regexp.MustCompile(`(locales?|translations?)`)

	environmentVariablePattern = regexp.MustCompile(`^[A-Z0-9_]+$`)
)

type detector struct{}

func New() types.Detector {
	return &detector{}
}

func (detector *detector) AcceptDir(dir *file.Path) (bool, error) {
	return true, nil
}

func (detector *detector) ProcessFile(file *file.FileInfo, dir *file.Path, report report.Report) (bool, error) {
	if file.Language != "YAML" && !filenamePattern.MatchString(file.Base) {
		return false, nil
	}

	if i18nExclusionPattern.MatchString(file.RelativePath) {
		return false, nil
	}

	bytes, err := os.ReadFile(file.AbsolutePath)

	if err != nil {
		return false, err
	}

	var j2Commands = regexp.MustCompile(`{(%|#).*?(%|#)}`)
	bytes = []byte(j2Commands.ReplaceAllString(string(bytes), ""))
	bytes = []byte(strings.ReplaceAll(string(bytes), "\t", "  "))
	tree, err := parser.ParseBytes(file, file.Path, bytes, yaml.GetLanguage(), 0)

	if err != nil {
		return false, err
	}

	defer tree.Close()

	return true, extractInterfaces(report, tree)
}

func extractInterfaces(report report.Report, tree *parser.Tree) error {
	return tree.Query(interfacesQuery, func(captures parser.Captures) error {
		var value string

		definitionNode := captures["definition"]
		key := ""

		if valueNode := captures["value"]; valueNode != nil {
			value = captures["value"].Content()
			key = captures["key"].Content()
		}

		if keyValueNode := captures["keyValue"]; keyValueNode != nil {
			keyValue := stripQuotes(keyValueNode.Content())
			splitKeyValue := strings.SplitN(keyValue, "=", 2)
			if len(splitKeyValue) == 2 {
				key = splitKeyValue[0]
				value = splitKeyValue[1]
			} else {
				value = keyValue
			}
		}

		parsedValue, err := parseValue(stripQuotes(value))
		if err != nil {
			return err
		}

		if interfaceType, isInterface := interfaces.GetType(parsedValue, false); isInterface {
			report.AddInterface(detectors.DetectorYamlConfig, reportinterface.Interface{
				Type:         interfaceType,
				Value:        parsedValue,
				VariableName: key,
			}, definitionNode.Source(true))
		}

		return nil
	})
}

func parseValue(text string) (*values.Value, error) {
	bytes := []byte(text)
	rootNode, err := sitter.ParseCtx(context.Background(), bytes, config_variables.GetLanguage())
	if err != nil {
		return nil, err
	}

	n := int(rootNode.ChildCount())

	value := values.New()
	for i := 0; i < n; i++ {
		child := rootNode.Child(i)
		if !child.IsNamed() {
			continue
		}

		partText := child.Content(bytes)
		switch child.Type() {
		case "literal":
			value.AppendString(partText)
		case "variable":
			variableType := variables.VariableTemplate
			if environmentVariablePattern.MatchString(partText) {
				variableType = variables.VariableEnvironment
			}

			value.AppendVariableReference(variableType, partText)
		case "unknown", "ERROR":
			value.AppendUnknown(nil)
		default:
			return nil, fmt.Errorf("unexpected node type: %s for %s", child.Type(), text)
		}
	}

	return value, nil
}

func stripQuotes(value string) string {
	return strings.Trim(value, `"'`)
}
