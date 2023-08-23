package html

import (
	"github.com/bearer/bearer/internal/detectors/javascript"
	"github.com/bearer/bearer/internal/detectors/types"
	"github.com/bearer/bearer/internal/parser"
	"github.com/bearer/bearer/internal/parser/interfaces"
	"github.com/bearer/bearer/internal/parser/nodeid"
	html "github.com/bearer/bearer/internal/parser/sitter/html2"
	"github.com/bearer/bearer/internal/report"
	"github.com/bearer/bearer/internal/report/detectors"
	interfacetype "github.com/bearer/bearer/internal/report/interfaces"
	"github.com/bearer/bearer/internal/report/values"
	"github.com/bearer/bearer/internal/util/file"
)

var (
	language    = html.GetLanguage()
	scriptQuery = parser.QueryMustCompile(language, `
	(
		script_element (
				(start_tag) @tag
				(raw_text) @raw
			)
	) @script
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
	if file.Language != "HTML" &&
		file.Language != "Vue" &&
		file.Language != "HTML+ERB" &&
		file.Language != "HTML+Django" &&
		file.Language != "Mustache" &&
		file.Language != "Jinja" &&
		file.Language != "Twig" &&
		file.Language != "Handlebars" &&
		file.Language != "HTML+PHP" &&
		file.Language != "HTML+Razor" &&
		file.Language != "EJS" &&
		file.Language != "Liquid" {
		return false, nil
	}

	tree, err := parser.ParseFile(file, file.Path, language)

	if err != nil {
		return false, err
	}
	defer tree.Close()

	return true, extractScripts(report, tree, file, detector.idGenerator)
}

func ProcessRaw(file *file.FileInfo, report report.Report, input []byte, offset int, idGenerator nodeid.Generator) (bool, error) {
	tree, err := parser.ParseBytes(file, file.Path, input, language, offset)
	if err != nil {
		return false, err
	}

	defer tree.Close()

	return true, extractScripts(report, tree, file, idGenerator)
}

func extractScripts(report report.Report, tree *parser.Tree, file *file.FileInfo, idGenerator nodeid.Generator) error {
	return tree.Query(scriptQuery, func(captures parser.Captures) error {
		tagNode := captures["tag"]
		var srcNode *parser.Node
		var typeValue string

		for i := 0; i < tagNode.ChildCount(); i++ {
			child := tagNode.Child(i)
			if child.Type() == "tag_name" {
				continue
			}

			if child.Type() == "attribute" {
				if child.Child(0).Content() == "src" {
					srcNode = child
				}
				if child.Child(0).Content() == "type" {
					target := child.Child(1).Child(0)
					if target == nil {
						continue
					}

					typeValue = target.Content()
				}
			}
		}

		if srcNode != nil {
			parsedValue, err := parseValue(srcNode)

			if err != nil {
				return err
			}

			if interfaceType, isInterface := interfaces.GetType(parsedValue, true); isInterface {
				report.AddInterface(detectors.DetectorHTML, interfacetype.Interface{
					Value: parsedValue,
					Type:  interfaceType,
				}, tagNode.Source(true))
			}

			return nil
		}

		if typeValue == "" || typeValue == "text/javascript" {
			offset := captures["raw"].StartLineNumber() - 1
			_, err := javascript.ProcessRaw(
				captures["raw"].Content(),
				report,
				file,
				offset,
				idGenerator,
			)

			if err != nil {
				return err
			}
		}

		return nil
	})
}

func parseValue(node *parser.Node) (*values.Value, error) {
	n := int(node.ChildCount())

	value := values.New()
	for i := 0; i < n; i++ {
		child := node.Child(i)

		if child.Type() == "quoted_attribute_value" {
			if child.FirstChild() == nil {
				continue
			}

			value.AppendString(child.FirstChild().Content())
		}
	}

	return value, nil
}
