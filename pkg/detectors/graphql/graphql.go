package graphql

import (
	"strings"

	"github.com/bearer/curio/pkg/detectors/types"
	"github.com/bearer/curio/pkg/parser"
	"github.com/bearer/curio/pkg/parser/nodeid"
	parserschema "github.com/bearer/curio/pkg/parser/schema"
	"github.com/bearer/curio/pkg/parser/sitter/graphql"
	"github.com/bearer/curio/pkg/report/detectors"
	"github.com/bearer/curio/pkg/report/schema"
	"github.com/bearer/curio/pkg/util/file"

	reporttypes "github.com/bearer/curio/pkg/report"
)

var (
	language           = graphql.GetLanguage()
	graphqlSchemaQuery = parser.QueryMustCompile(language, `
	(
		object_type_definition (name) @object_name
			(fields_definition
				(
				field_definition
					(name) @field_name
					(type) @field_type
				)
			)
	)
	`)
)

type detector struct {
	idGenerator nodeid.Generator
}

func New(idGenerator nodeid.Generator) types.Detector {
	return &detector{
		idGenerator: idGenerator,
	}
}

func (detector *detector) AcceptDir(dir *file.Path) (bool, error) {
	return true, nil
}

func (detector *detector) ProcessFile(file *file.FileInfo, dir *file.Path, report reporttypes.Report) (bool, error) {
	if file.Language != "GraphQL" {
		return false, nil
	}

	err := detector.ExtractFromSchema(file, report)

	return true, err
}

func (detector *detector) ExtractFromSchema(
	file *file.FileInfo,
	report reporttypes.Report,
) error {
	tree, err := parser.ParseFile(file, file.Path, language)
	if err != nil {
		return err
	}
	defer tree.Close()

	uuidHolder := parserschema.NewUUIDHolder()

	return tree.Query(graphqlSchemaQuery, func(captures parser.Captures) error {
		objectNode := captures["object_name"]
		objectName := stripQuotes(objectNode.Content())
		fieldType := stripQuotes(captures["field_type"].Content())

		fieldNode := captures["field_name"]
		fieldName := stripQuotes(fieldNode.Content())

		objectUUID := uuidHolder.Assign(objectNode.ID(), detector.idGenerator)
		fieldUUID := uuidHolder.Assign(fieldNode.ID(), detector.idGenerator)

		report.AddSchema(detectors.DetectorGraphQL,
			schema.Schema{
				ObjectName:      objectName,
				ObjectUUID:      objectUUID,
				FieldName:       fieldName,
				FieldType:       fieldType,
				FieldUUID:       fieldUUID,
				SimpleFieldType: convertType(fieldType),
			},
			fieldNode.Source(true),
		)

		return nil
	})
}

func stripQuotes(value string) string {
	return strings.Trim(value, `"'`)
}

func convertType(value string) string {
	simplified := strings.TrimSuffix(value, "!")
	switch simplified {
	case "Int":
		return schema.SimpleTypeNumber
	case "Float":
		return schema.SimpleTypeNumber
	case "String":
		return schema.SimpleTypeString
	case "ID":
		return schema.SimpleTypeString
	case "Boolean":
		return schema.SimpleTypeBool
	default:
		return schema.SimpleTypeObject
	}
}
