package reportadder

import (
	"regexp"
	"sort"
	"strings"

	"github.com/bearer/curio/pkg/parser"
	reporttypes "github.com/bearer/curio/pkg/report"
	"github.com/bearer/curio/pkg/report/detectors"
	"github.com/bearer/curio/pkg/report/operations"
	"github.com/bearer/curio/pkg/report/schema"
	"github.com/bearer/curio/pkg/report/schema/schemahelper"
	"github.com/bearer/curio/pkg/util/file"
	"github.com/bearer/curio/pkg/util/stringutil"
)

var regexpPathVariable = regexp.MustCompile(`\{.+\}`)

func AddSchema(file *file.FileInfo, report reporttypes.Report, foundValues map[parser.Node]*schemahelper.Schema) {
	// we need sorted schemas so our reports are consistent and repeatable
	var sortedSchemas []*schemahelper.Schema
	for _, schema := range foundValues {
		sortedSchemas = append(sortedSchemas, schema)
	}
	sort.Slice(sortedSchemas, func(i, j int) bool {
		lineNumberA := sortedSchemas[i].Source.LineNumber
		lineNumberB := sortedSchemas[j].Source.LineNumber
		return *lineNumberA < *lineNumberB
	})
	for _, schema := range sortedSchemas {
		schema.Source.Language = file.Language
		schema.Source.LanguageType = file.LanguageTypeString()
		schema.Value.FieldName = stringutil.StripQuotes(schema.Value.FieldName)
		schema.Value.FieldType = stringutil.StripQuotes(schema.Value.FieldType)
		schema.Value.ObjectName = stringutil.StripQuotes(schema.Value.ObjectName)
		schema.Value.SimpleFieldType = convertSchema(schema.Value.FieldType)
		report.AddSchema(detectors.DetectorOpenAPI, schema.Value, schema.Source)
	}
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
