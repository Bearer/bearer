package packageconfig

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
	(_
    	(tag_name) @helper_package
        (#match? @helper_package "^package$")
        (attribute
        	(attribute_name) @helper_id
            (#match? @helper_id "^id$")
			(attribute_value) @param_dependency
        )
        (attribute
        	(attribute_name) @helper_version
            (#match? @helper_version "^version$")
			(attribute_value) @param_version
        )
    )
)
`)

func Discover(f *file.FileInfo) (report *depsbase.DiscoveredDependency) {
	report = &depsbase.DiscoveredDependency{}
	report.Provider = "package-config"
	report.Language = "csharp"
	report.PackageManager = "nuget"
	tree, err := parser.ParseFile(f, f.Path, language)
	if err != nil {
		log.Error().Msgf("%s: there was an error while parsing the file: %s", report.Provider, err.Error())
		return nil
	}
	defer tree.Close()

	captures := tree.QueryMustPass(queryDependencies)
	for _, capture := range captures {
		if stringutil.StripQuotes(capture["helper_package"].Content()) != "package" ||
			stringutil.StripQuotes(capture["helper_id"].Content()) != "id" ||
			stringutil.StripQuotes(capture["helper_version"].Content()) != "version" {
			continue
		}

		report.Dependencies = append(report.Dependencies, depsbase.Dependency{
			Name:    stringutil.StripQuotes(capture["param_dependency"].Content()),
			Version: stringutil.StripQuotes(capture["param_version"].Content()),
			Line:    int64(capture["param_dependency"].StartLineNumber()),
			Column:  int64(capture["param_dependency"].Column()),
		})
	}

	return report
}
