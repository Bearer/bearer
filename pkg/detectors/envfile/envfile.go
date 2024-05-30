package envfile

import (
	"regexp"
	"strings"

	"github.com/smacker/go-tree-sitter/bash"

	"github.com/bearer/bearer/pkg/detectors/types"
	"github.com/bearer/bearer/pkg/parser"
	"github.com/bearer/bearer/pkg/parser/interfaces"
	"github.com/bearer/bearer/pkg/report"
	"github.com/bearer/bearer/pkg/report/detectors"
	reportinterface "github.com/bearer/bearer/pkg/report/interfaces"
	"github.com/bearer/bearer/pkg/report/values"
	"github.com/bearer/bearer/pkg/util/file"
)

var (
	language = bash.GetLanguage()

	variablesQuery = parser.QueryMustCompile(language, `
		(variable_assignment
			name: (variable_name) @name
			value: (_) @value) @detection
	`)

	filenamePattern = regexp.MustCompile(`\.env`)
)

type detector struct{}

func New() types.Detector {
	return &detector{}
}

func (detector *detector) AcceptDir(dir *file.Path) (bool, error) {
	return true, nil
}

func (detector *detector) ProcessFile(file *file.FileInfo, dir *file.Path, report report.Report) (bool, error) {
	if !filenamePattern.MatchString(file.Base) {
		return false, nil
	}

	tree, err := parser.ParseFile(file, file.Path, language)
	if err != nil {
		return false, err
	}
	defer tree.Close()

	return true, tree.Query(variablesQuery, func(captures parser.Captures) error {
		detectionNode := captures["detection"]
		name := captures["name"].Content()
		value := stripQuotes(captures["value"].Content())

		parsedValue := values.New()
		parsedValue.AppendString(value)

		interfaceType, isInterface := interfaces.GetTypeWithKey(name, parsedValue)
		if isInterface {
			report.AddInterface(detectors.DetectorEnvFile, reportinterface.Interface{
				Value:        parsedValue,
				Type:         interfaceType,
				VariableName: name,
			}, detectionNode.Source(true))
		}

		return nil
	})
}

func stripQuotes(value string) string {
	return strings.Trim(value, `"'`)
}
