package pomxml

import (
	"github.com/bearer/bearer/pkg/detectors/dependencies/depsbase"
	"github.com/bearer/bearer/pkg/parser"
	xml "github.com/bearer/bearer/pkg/parser/sitter/xml2"
	"github.com/bearer/bearer/pkg/util/file"
	"github.com/bearer/bearer/pkg/util/stringutil"
	"github.com/rs/zerolog/log"
)

var language = xml.GetLanguage()

var queryDependencies = parser.QueryMustCompile(language, `
(element
	(start_tag
		(tag_name) @helper_dependency
		(#match? @helper_dependency "^dependency$")
	)
 ) @param_dependency
`)

func Discover(f *file.FileInfo) (report *depsbase.DiscoveredDependency) {
	report = &depsbase.DiscoveredDependency{}
	report.Provider = "pom-xml"
	report.Language = "Java"
	report.PackageManager = "maven"
	tree, err := parser.ParseFile(f, f.Path, language)
	if err != nil {
		log.Error().Msgf("%s: there was an error while parsing the file: %s", report.Provider, err.Error())
		return nil
	}
	defer tree.Close()

	captures := tree.QueryConventional(queryDependencies)
	for _, capture := range captures {
		var groupId, artifactId, version string

		dependencyNode := capture["param_dependency"]
		for i := 0; i < dependencyNode.ChildCount(); i++ {
			child := dependencyNode.Child(i)

			if child.Type() != "element" {
				continue
			}

			tag := ""
			tagContent := ""

			for j := 0; j < child.ChildCount(); j++ {
				elementChild := child.Child(j)

				if elementChild.Type() == "start_tag" {
					tag = elementChild.Child(0).Content()
				}

				if elementChild.Type() == "text" {
					tagContent = elementChild.Content()
				}
			}

			switch tag {
			case "groupId":
				groupId = tagContent
			case "artifactId":
				artifactId = tagContent
			case "version":
				version = tagContent
			}
		}

		report.Dependencies = append(report.Dependencies, depsbase.Dependency{
			Name:    artifactId,
			Group:   groupId,
			Version: stringutil.StripQuotes(version),
			Line:    int64(capture["param_dependency"].StartLineNumber()),
			Column:  int64(capture["param_dependency"].Column()),
		})
	}

	return report
}
