package personal_data

import (
	"strings"

	"github.com/smacker/go-tree-sitter/ruby"

	"github.com/bearer/curio/pkg/parser"
	"github.com/bearer/curio/pkg/parser/nodeid"
	parserschema "github.com/bearer/curio/pkg/parser/schema"
	"github.com/bearer/curio/pkg/report"
	"github.com/bearer/curio/pkg/report/detectors"
	"github.com/bearer/curio/pkg/report/schema"
	"github.com/bearer/curio/pkg/util/file"
)

var (
	language = ruby.GetLanguage()

	rubyDatabaseSchemaQuery = parser.QueryMustCompile(language, `
		(call
			method: (identifier) @table_method
			arguments: (argument_list . (string) @table_name)
			block: (do_block
				parameters: (block_parameters (identifier) @block_param)
				(call
					receiver: (_) @receiver
					method: (_) @type
					arguments: (argument_list . (string) @column_name))
				(#eq @receiver @block_param))
			(#eq @table_method "create_table"))
	`)
)

func ExtractFromDatabaseSchema(
	idGenerator nodeid.Generator,
	file *file.FileInfo,
	report report.Report,
) error {
	tree, err := parser.ParseFile(file, file.Path, language)
	if err != nil {
		return err
	}
	defer tree.Close()

	uuidHolder := parserschema.NewUUIDHolder()

	return tree.Query(rubyDatabaseSchemaQuery, func(captures parser.Captures) error {
		tableNode := captures["table_name"]
		tableName := stripQuotes(tableNode.Content())
		columnName := stripQuotes(captures["column_name"].Content())
		columnTypeNode := captures["type"]
		columnType := columnTypeNode.Content()

		objectUUID := uuidHolder.Assign(tableNode.ID(), idGenerator)
		fieldUUID := uuidHolder.Assign(columnTypeNode.ID(), idGenerator)

		report.AddSchema(detectors.DetectorRails, schema.Schema{
			ObjectName:      tableName,
			ObjectUUID:      objectUUID,
			FieldName:       columnName,
			FieldUUID:       fieldUUID,
			FieldType:       columnType,
			SimpleFieldType: convertToSimpleType(columnType),
		}, columnTypeNode.Source(false))

		return nil
	})
}

func stripQuotes(value string) string {
	return strings.Trim(value, `"'`)
}

func convertToSimpleType(value string) string {
	switch value {

	case "boolean":
		return schema.SimpleTypeBool
	case "date":
		return schema.SimpleTypeDate
	case "timestamp":
		return schema.SimpleTypeDate
	case "datetime":
		return schema.SimpleTypeDate
	case "float":
		return schema.SimpleTypeNumber
	case "integer":
		return schema.SimpleTypeNumber
	case "bigint":
		return schema.SimpleTypeNumber
	case "binary":
		return schema.SimpleTypeString
	case "string":
		return schema.SimpleTypeString
	case "text":
		return schema.SimpleTypeString
	default:
		return schema.SimpleTypeObject
	}
}
