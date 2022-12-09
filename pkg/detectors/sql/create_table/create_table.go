package createtable

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/bearer/curio/pkg/detectors/sql/util"
	"github.com/bearer/curio/pkg/parser"
	"github.com/bearer/curio/pkg/parser/nodeid"
	"github.com/bearer/curio/pkg/util/file"

	parserschema "github.com/bearer/curio/pkg/parser/schema"

	"github.com/bearer/curio/pkg/parser/sitter/sql"
	reporttypes "github.com/bearer/curio/pkg/report"
	"github.com/bearer/curio/pkg/report/detectors"
	"github.com/bearer/curio/pkg/report/schema"
)

var createTableRegexp = regexp.MustCompile(`(?i)(create table)`)
var language = sql.GetLanguage()
var sqlDatabaseSchemaQuery = parser.QueryMustCompile(language, `
(
	(create_table_statement
		(dotted_name
			(identifier) @schema_name
			(identifier) @table_name
		)
		(table_parameters
			(table_column
				name: (identifier) @column_name
				type: [
					((_)? (type) @column_type)
					(array_type (type) @column_type)
					((type) @column_type)
				]
			)
		)
	)
)
`)

func Detect(file *file.FileInfo, report reporttypes.Report, idGenerator nodeid.Generator) error {
	f, err := os.Open(file.Path.AbsolutePath)
	if err != nil {
		return err
	}
	fileBytes, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	defer f.Close()
	// our sql tree sitter parser tends to error sometimes mid file causing us to partially parse file
	// with this hack we increase our parsing percentage
	chunks := createTableRegexp.Split(string(fileBytes), -1)

	lineOffset := 0
	for i, chunk := range chunks {
		chunkBytes := []byte(chunk)
		if i != 0 {
			chunkBytes = []byte("CREATE TABLE" + chunk)
		}

		tree, err := parser.ParseBytes(file, file.Path, chunkBytes, language, lineOffset)
		if err != nil {
			return err
		}
		defer tree.Close()

		uuidHolder := parserschema.NewUUIDHolder()

		err = tree.Query(sqlDatabaseSchemaQuery, func(captures parser.Captures) error {
			tableNode := captures["table_name"]
			tableName := util.StripQuotes(tableNode.Content())
			columnNode := captures["column_name"]

			columnType := util.StripQuotes(captures["column_type"].Content())

			columnName := util.StripQuotes(columnNode.Content())

			if columnName == "CONSTRAINT" {
				return nil
			}

			objectUUID := uuidHolder.Assign(tableNode.ID(), idGenerator)
			fieldUUID := uuidHolder.Assign(columnNode.ID(), idGenerator)

			currentSchema := schema.Schema{
				ObjectName:      tableName,
				ObjectUUID:      objectUUID,
				FieldName:       columnName,
				FieldUUID:       fieldUUID,
				FieldType:       columnType,
				SimpleFieldType: util.ConvertToSimpleType(columnType),
			}

			if report.SchemaGroupShouldClose(tableName) {
				report.SchemaGroupEnd(idGenerator)
			}

			if !report.SchemaGroupIsOpen() {
				source := tableNode.Source(true)
				report.SchemaGroupBegin(
					detectors.DetectorSQL,
					tableNode,
					currentSchema,
					&source,
					nil,
				)
			}
			source := columnNode.Source(true)
			report.SchemaGroupAddItem(
				columnNode,
				currentSchema,
				&source,
			)

			return nil
		})

		report.SchemaGroupEnd(idGenerator)

		if err != nil {
			report.AddError(file.RelativePath, fmt.Errorf("create table error: %s", err))
		}

		lineOffset = lineOffset + strings.Count(chunk, "\n")
	}
	return nil
}
