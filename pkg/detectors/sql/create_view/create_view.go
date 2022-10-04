package createview

import (
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/bearer/curio/pkg/detectors/sql/util"
	"github.com/bearer/curio/pkg/parser"
	"github.com/bearer/curio/pkg/parser/nodeid"
	"github.com/bearer/curio/pkg/util/file"

	"github.com/bearer/curio/pkg/parser/sitter/sql"
	reporttypes "github.com/bearer/curio/pkg/report"
	createview "github.com/bearer/curio/pkg/report/create_view"
	"github.com/bearer/curio/pkg/report/detectors"
)

var createTableRegexp = regexp.MustCompile(`(?i)(create( or replace)?( or update)? view)`)
var fromTableRegexp = regexp.MustCompile(`(?i)\sfrom\s`)

var language = sql.GetLanguage()
var viewNameQuery = parser.QueryMustCompile(language, `
	(create_view_statement 
		(dotted_name
			(identifier) @param_schema_name
			(identifier) @param_view_name
		) @param_identifier
	) @param_view
`)
var viewFieldsQuery = parser.QueryMustCompile(language, `
	(dotted_name
		(identifier) @param_table_name
		(identifier) @param_field_name
	) @param_identifier
`)

var viewFromSchemaQuery = parser.QueryMustCompile(language, `
	(from_clause
		(dotted_name
			(identifier) @param_table_name
			(identifier) @param_field_name
		) @param_identifier
	)
`)

var viewFromJoinQuery = parser.QueryMustCompile(language, `
	(join_clause
		(dotted_name
			(identifier) @param_table_name
			(identifier) @param_field_name
		) @param_identifier
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
			chunkBytes = []byte("CREATE VIEW" + chunk)
		}

		tree, err := parser.ParseBytes(file, file.Path, chunkBytes, language, lineOffset)
		if err != nil {
			return err
		}
		defer tree.Close()

		var view = createview.View{}
		// add view name
		captures := tree.QueryConventional(viewNameQuery)
		for _, capture := range captures {
			view.Source = capture["param_view"].Source(false)
			view.SchemaName = util.StripQuotes(capture["param_schema_name"].Content())
			view.ViewName = util.StripQuotes(capture["param_view_name"].Content())
		}

		// add view fields
		captures = tree.QueryConventional(viewFieldsQuery)
		for _, capture := range captures {
			_, err := capture["param_identifier"].FindParent("create_view_statement")
			if err != nil {
				continue
			}
			selectNode, err := capture["param_identifier"].FindParent("select_clause_body")
			if err != nil {
				continue
			}

			// verify select is inside view body
			selectClause := selectNode.Parent()
			if selectClause == nil || selectClause.Type() != "select_clause" {
				continue
			}
			selectStatement := selectClause.Parent()
			if selectStatement == nil || selectStatement.Type() != "select_statement" {
				continue
			}
			viewBody := selectStatement.Parent()
			if viewBody == nil || viewBody.Type() != "view_body" {
				continue
			}

			view.Fields = append(view.Fields, &createview.Field{
				TableName: util.StripQuotes(capture["param_table_name"].Content()),
				FieldName: util.StripQuotes(capture["param_field_name"].Content()),
				Source:    capture["param_identifier"].Source(false),
			})
		}

		fromFields, err := extractFromFields(file, string(chunkBytes), lineOffset)
		if err != nil {
			continue
		}

		view.From = fromFields

		if len(view.From) == 0 && len(view.Fields) == 0 {
			continue
		}

		report.AddCreateView(detectors.DetectorSQL, view)

		lineOffset = lineOffset + strings.Count(chunk, "\n")
	}
	return nil
}

func extractFromFields(file *file.FileInfo, Command string, lineOffset int) (fromFields []*createview.Table, err error) {
	chunks := fromTableRegexp.Split(Command, -1)
	if len(chunks) < 2 {
		return
	}

	var tree *parser.Tree
	fromOffset := lineOffset + strings.Count(chunks[0], "\n")

	fromClause := "CREATE VIEW test AS SELECT a \nFROM " + chunks[1]
	tree, err = parser.ParseBytes(file, file.Path, []byte(fromClause), language, fromOffset)
	if err != nil {
		return nil, err
	}

	// add from identifiers
	captures := tree.QueryConventional(viewFromSchemaQuery)
	for _, capture := range captures {
		_, err := capture["param_identifier"].FindParent("create_view_statement")
		if err != nil {
			continue
		}

		_, err = capture["param_identifier"].FindParent("from_clause")
		if err != nil {
			continue
		}

		fromField := &createview.Table{
			TableName: util.StripQuotes(capture["param_table_name"].Content()),
			FieldName: util.StripQuotes(capture["param_field_name"].Content()),
			Source:    capture["param_identifier"].Source(false),
		}

		alias := capture["param_identifier"].Parent().Child(1)

		if alias != nil && alias.Type() == "alias" {
			fromField.Alias = alias.Child(0).Content()
		}

		fromFields = append(fromFields, fromField)
	}

	// add join identifiers
	captures = tree.QueryConventional(viewFromJoinQuery)
	for _, capture := range captures {
		_, err := capture["param_identifier"].FindParent("create_view_statement")
		if err != nil {
			continue
		}

		_, err = capture["param_identifier"].FindParent("join_clause")
		if err != nil {
			continue
		}

		joinField := &createview.Table{
			TableName: util.StripQuotes(capture["param_table_name"].Content()),
			FieldName: util.StripQuotes(capture["param_field_name"].Content()),
			Source:    capture["param_identifier"].Source(false),
		}

		alias := capture["param_identifier"].Parent().Child(1)
		if alias != nil && alias.Type() == "alias" {
			joinField.Alias = alias.Child(0).Content()
		}

		fromFields = append(fromFields, joinField)
	}

	return
}
