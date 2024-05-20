package pomxml

import (
	"github.com/rs/zerolog/log"

	"github.com/bearer/bearer/pkg/detectors/dependencies/depsbase"
	"github.com/bearer/bearer/pkg/parser"
	"github.com/bearer/bearer/pkg/parser/sitter/xml"
	"github.com/bearer/bearer/pkg/util/file"
	"github.com/bearer/bearer/pkg/util/stringutil"
)

var language = xml.GetLanguage()

var queryDependencies = parser.QueryMustCompile(language, `
(element
	(STag
		(Name) @helper_dependency
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

			if child.Type() != "content" {
				continue
			}

			for j := 0; j < child.ChildCount(); j++ {
				childElement := child.Child(j)
				if childElement.Type() != "element" {
					continue
				}

				tag := ""
				tagContent := ""

				for k := 0; k < childElement.ChildCount(); k++ {
					subElement := childElement.Child(k)
					switch subElement.Type() {
					case "STag":
						tag = subElement.Child(0).Content()
					case "content":
						tagContent += subElement.Content()
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
