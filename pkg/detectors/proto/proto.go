package proto

import (
	"strings"

	"github.com/bearer/bearer/pkg/detectors/types"
	"github.com/bearer/bearer/pkg/parser"
	"github.com/bearer/bearer/pkg/parser/sitter/proto"
	"github.com/bearer/bearer/pkg/report/detectors"
	"github.com/bearer/bearer/pkg/report/schema"
	"github.com/bearer/bearer/pkg/util/file"
	"github.com/gertd/go-pluralize"

	"github.com/bearer/bearer/pkg/parser/nodeid"
	parserschema "github.com/bearer/bearer/pkg/parser/schema"
	reporttypes "github.com/bearer/bearer/pkg/report"
)

var (
	language         = proto.GetLanguage()
	protoSchemaQuery = parser.QueryMustCompile(language, `
	(
		message
		(message_name (identifier) @object_name)
			(message_body
				(
					field (
						(type) @field_type
						(identifier) @field_name
					)
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
	if file.Language != "Protocol Buffer" {
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

	err = tree.Query(protoSchemaQuery, func(captures parser.Captures) error {
		objectNode := captures["object_name"]
		objectName := stripQuotes(objectNode.Content())

		columnType := stripQuotes(captures["field_type"].Content())

		fieldNode := captures["field_name"]
		fieldName := stripQuotes(fieldNode.Content())

		objectUUID := uuidHolder.Assign(objectNode.ID(), detector.idGenerator)
		fieldUUID := uuidHolder.Assign(fieldNode.ID(), detector.idGenerator)

		normalizedObjectName := ""
		normalizedFieldName := ""
		pluralizer := pluralize.NewClient()
		normalizedObjectName = pluralizer.Singular(strings.ToLower(objectName))
		normalizedFieldName = pluralizer.Singular(strings.ToLower(fieldName))

		currentSchema := schema.Schema{
			ObjectName:           objectName,
			ObjectUUID:           objectUUID,
			FieldName:            fieldName,
			FieldUUID:            fieldUUID,
			FieldType:            columnType,
			SimpleFieldType:      convertToSimpleType(columnType),
			NormalizedObjectName: normalizedObjectName,
			NormalizedFieldName:  normalizedFieldName,
		}

		if report.SchemaGroupShouldClose(objectName) {
			report.SchemaGroupEnd(detector.idGenerator)
		}

		if !report.SchemaGroupIsOpen() {
			source := objectNode.Source(true)
			report.SchemaGroupBegin(
				detectors.DetectorProto,
				objectNode,
				currentSchema,
				&source,
				nil,
			)
		}
		source := fieldNode.Source(true)
		report.SchemaGroupAddItem(
			fieldNode,
			currentSchema,
			&source,
		)

		return nil
	})

	report.SchemaGroupEnd(detector.idGenerator)

	return err
}

func stripQuotes(value string) string {
	return strings.Trim(value, `"'`)
}

func convertToSimpleType(value string) string {
	numberMap := []string{"double", "float", "int32", "in64", "uint32", "uint64", "sint32", "sint64", "fixed32", "fixed64", "sfixed32", "sfixed64"}
	for _, typeValue := range numberMap {
		if typeValue == value {
			return schema.SimpleTypeNumber
		}
	}

	if value == "string" {
		return schema.SimpleTypeString
	}

	if value == "bytes" {
		return schema.SimpleTypeBinary
	}

	if value == "bool" {
		return schema.SimpleTypeBool
	}

	return schema.SimpleTypeObject
}
