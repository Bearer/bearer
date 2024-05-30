package django

import (
	"errors"
	"strings"

	"github.com/bearer/bearer/pkg/detectors/types"
	"github.com/bearer/bearer/pkg/parser"
	"github.com/bearer/bearer/pkg/report"
	"github.com/bearer/bearer/pkg/report/detectors"
	"github.com/bearer/bearer/pkg/report/frameworks/django"
	"github.com/bearer/bearer/pkg/util/file"

	"github.com/smacker/go-tree-sitter/python"
)

var (
	foundErr = errors.New("found")

	language = python.GetLanguage()

	projectQuery = parser.QueryMustCompile(language, `
		(import_from_statement
			module_name: (dotted_name) @module
			(#match? @module "^django(\\..*|$)"))

		(import_statement
			name: (dotted_name) @module
			(#match? @module "^django(\\..*|$)"))
	`)

	databasesQuery = parser.QueryMustCompile(language, `
		(assignment
			left: (identifier) @variable
			right: (dictionary
							(pair
								key: (string) @name
								value: (dictionary) @attributes))
			(#eq? @variable "DATABASES"))
	`)

	databaseAttributesQuery = parser.QueryMustCompile(language, `
		(dictionary
			(pair
				key: (string) @key
				value: (string) @value))
	`)
)

type detector struct{}

func New() types.Detector {
	return &detector{}
}

func (detector *detector) AcceptDir(dir *file.Path) (bool, error) {
	manageFile := dir.Join("manage.py")
	if exists := manageFile.Exists(); !exists {
		return false, nil
	}

	tree, err := parser.ParseFile(nil, manageFile, language)
	if err != nil {
		return false, err
	}

	err = tree.Query(projectQuery, func(_ parser.Captures) error {
		return foundErr
	})

	if err == foundErr {
		return true, nil
	}

	return false, err
}

func (detector *detector) ProcessFile(file *file.FileInfo, dir *file.Path, report report.Report) (bool, error) {
	if file.Base != "settings.py" {
		return false, nil
	}

	tree, err := parser.ParseFile(file, file.Path, language)
	if err != nil {
		return false, err
	}
	defer tree.Close()

	err = tree.Query(databasesQuery, func(captures parser.Captures) error {
		nameNode := captures["name"]
		name := stripQuotes(nameNode.Content())

		engine := ""

		err := captures["attributes"].Query(
			databaseAttributesQuery,
			func(attrCaptures parser.Captures) error {
				key := stripQuotes(attrCaptures["key"].Content())
				if key == "ENGINE" {
					engine = stripQuotes(attrCaptures["value"].Content())
				}

				return nil
			},
		)
		report.AddFramework(detectors.DetectorDjango, django.TypeDatabase, django.Database{
			Name:   name,
			Engine: engine,
		}, nameNode.Source(false))

		return err
	})

	// Allow "python" detector to process file
	return false, err
}

func stripQuotes(value string) string {
	return strings.Trim(value, `'"`)
}
