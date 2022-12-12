package reportadder

import (
	"regexp"
	"sort"
	"strings"

	"github.com/bearer/curio/pkg/parser"
	"github.com/bearer/curio/pkg/parser/nodeid"
	reporttypes "github.com/bearer/curio/pkg/report"
	"github.com/bearer/curio/pkg/report/detectors"
	"github.com/bearer/curio/pkg/report/operations"
	"github.com/bearer/curio/pkg/report/schema"
	"github.com/bearer/curio/pkg/report/schema/schemahelper"
	"github.com/bearer/curio/pkg/util/file"
	"github.com/bearer/curio/pkg/util/stringutil"
)

var regexpPathVariable = regexp.MustCompile(`\{.+\}`)

type SortableSchema struct {
	Key parser.Node
	Value *schemahelper.Schema
}
type SortableSchemaList []SortableSchema
func (s SortableSchemaList) Len() int {
	return len(s)
}
func (s SortableSchemaList) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s SortableSchemaList) Less(i, j int) bool {
	lineNumberA := s[i].Value.Source.LineNumber
	lineNumberB := s[j].Value.Source.LineNumber
	return *lineNumberA < *lineNumberB
}

func AddSchema(file *file.FileInfo, report reporttypes.Report, foundValues map[parser.Node]*schemahelper.Schema, idGenerator nodeid.Generator) {
	// we need sorted schemas so our reports are consistent and repeatable
	sortedSchemas := make(SortableSchemaList, len(foundValues))
	i := 0
	for k, v := range(foundValues) {
		sortedSchemas[i] = SortableSchema{Key: k, Value: v}
		i++
	}
	sort.Sort(sortedSchemas)

	for _, sortableSchema := range sortedSchemas {
		node := sortableSchema.Key
		schema := sortableSchema.Value

		objectName := stringutil.StripQuotes(schema.Value.ObjectName)

		schema.Source.Language = file.Language
		schema.Source.LanguageType = file.LanguageTypeString()
		schema.Value.FieldName = stringutil.StripQuotes(schema.Value.FieldName)
		schema.Value.FieldType = stringutil.StripQuotes(schema.Value.FieldType)
		schema.Value.ObjectName = objectName
		schema.Value.SimpleFieldType = convertSchema(schema.Value.FieldType)

		if report.SchemaGroupShouldClose(objectName) {
			report.SchemaGroupEnd(idGenerator)
		}

		if !report.SchemaGroupIsOpen() {
			report.SchemaGroupBegin(
				detectors.DetectorOpenAPI,
				nil,
				schema.Value,
				&schema.Source,
				nil,
			)
		}

		fieldNode := node
		report.SchemaGroupAddItem(
			&fieldNode,
			schema.Value,
			&schema.Source,
		)
	}

	report.SchemaGroupEnd(idGenerator)
}

func convertSchema(value string) string {
	switch value {
	case "string":
		return schema.SimpleTypeString
	case "number":
		return schema.SimpleTypeNumber
	case "integer":
		return schema.SimpleTypeNumber
	case "boolean":
		return schema.SimpleTypeBool
	default:
		return schema.SimpleTypeObject
	}
}

func standardizeOperationType(input string) (output string) {
	input = strings.ToUpper(input)
	supportedvalues := []string{operations.TypeGet, operations.TypeDelete, operations.TypePost, operations.TypePut}

	for _, v := range supportedvalues {
		if input == v {
			return v
		}
	}

	return operations.TypeOther
}

func standardizeOperationPath(input string) (output string) {
	return regexpPathVariable.ReplaceAllString(input, "*")
}
