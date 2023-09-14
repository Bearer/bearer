package ipynb

import (
	"strconv"
	"strings"

	"github.com/bearer/bearer/internal/detectors/python"
	"github.com/bearer/bearer/internal/detectors/types"
	"github.com/bearer/bearer/internal/parser"
	"github.com/bearer/bearer/internal/parser/nodeid"
	"github.com/bearer/bearer/internal/report"
	"github.com/bearer/bearer/internal/util/file"
	"github.com/bearer/bearer/internal/util/stringutil"
	jslang "github.com/smacker/go-tree-sitter/javascript"
)

var (
	language   = jslang.GetLanguage()
	inputQuery = parser.QueryMustCompile(language, `
	(_
		(
		  pair
			key:
				(string) @helperCellTypeKey
				(#match? @helperCellTypeKey "^\"cell_type\"$")
			 value:
				(string) @helperCellTypeValue
				(#match? @helperCellTypeValue "^\"code\"$")
		)
		(
		  pair
			key:
				(string) @helperInputKey
				(#match? @helperInputKey "^\"input\"$")
			 value:
				(array) @param_input_value
		)
		(
		  pair
			key:
				(string) @helperLanguageKey
				(#match? @helperLanguageKey "^\"language\"$")
			 value:
				(string) @helperLanguageValue
				(#match? @helperLanguageValue "^\"python\"$")
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

func (detector *detector) ProcessFile(file *file.FileInfo, dir *file.Path, report report.Report) (bool, error) {
	if file.Language != "Jupyter Notebook" {
		return false, nil
	}

	tree, err := parser.ParseFile(file, file.Path, language)
	if err != nil {
		return false, err
	}
	defer tree.Close()

	return true, detector.extractScripts(report, tree, file)
}

func (detector *detector) extractScripts(report report.Report, tree *parser.Tree, file *file.FileInfo) error {
	captures := tree.QueryMustPass(inputQuery)

	for _, capture := range captures {

		if stringutil.StripQuotes(capture["helperCellTypeKey"].Content()) != `cell_type` ||
			stringutil.StripQuotes(capture["helperCellTypeValue"].Content()) != `code` ||
			stringutil.StripQuotes(capture["helperInputKey"].Content()) != `input` ||
			stringutil.StripQuotes(capture["helperLanguageKey"].Content()) != `language` ||
			stringutil.StripQuotes(capture["helperLanguageValue"].Content()) != `python` {
			continue
		}

		input := capture["param_input_value"]

		pythonScript := ""

		for i := 0; i < input.ChildCount(); i++ {
			stringLine := input.Child(i).Content()
			stringLine = strings.TrimRight(stringLine, ",")
			stringLineParsed, _ := strconv.Unquote(stringLine)
			pythonScript += stringLineParsed
		}

		_, err := python.ProcessRaw(file, report, []byte(pythonScript), input.StartLineNumber(), detector.idGenerator)

		if err != nil {
			return err
		}
	}

	return nil
}
