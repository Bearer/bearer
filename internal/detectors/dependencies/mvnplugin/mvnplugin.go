package mvnplugin

import (
	"github.com/bearer/bearer/internal/detectors/dependencies/depsbase"
	"github.com/bearer/bearer/internal/parser"
	"github.com/bearer/bearer/internal/util/file"
	"github.com/bearer/bearer/internal/util/stringutil"
	"github.com/rs/zerolog/log"
	"github.com/smacker/go-tree-sitter/javascript"
)

var language = javascript.GetLanguage()

var queryDependencies = parser.QueryMustCompile(language, `
(object
	(pair
		key: (string) @helper_group_id
		(#match? @helper_group_id "^\"groupId\"$")
		value: (string) @param_group_id
	)
	(pair
		key: (string) @helper_dependency
		(#match? @helper_dependency "^\"artifactId\"$")
		value: (string) @param_dependency
	)
	(pair
		key: (string) @helper_version
		(#match? @helper_version "^\"version\"$")
		value: (string) @param_version
	)
)
`)

func Discover(f *file.FileInfo) (report *depsbase.DiscoveredDependency) {
	report = &depsbase.DiscoveredDependency{}
	report.Provider = "mvnplugin"
	report.Language = "Java"
	report.PackageManager = "maven"
	tree, err := parser.ParseFile(f, f.Path, language)
	if err != nil {
		log.Error().Msgf("%s: there was an error while parsing the file: %s", report.Provider, err.Error())
		return nil
	}
	defer tree.Close()

	captures := tree.QueryMustPass(queryDependencies)
	for _, capture := range captures {
		if stringutil.StripQuotes(capture["helper_dependency"].Content()) != "artifactId" ||
			stringutil.StripQuotes(capture["helper_version"].Content()) != "version" ||
			stringutil.StripQuotes(capture["helper_group_id"].Content()) != "groupId" {
			continue
		}

		report.Dependencies = append(report.Dependencies, depsbase.Dependency{
			Name:    stringutil.StripQuotes(capture["param_dependency"].Content()),
			Group:   stringutil.StripQuotes(capture["param_group_id"].Content()),
			Version: stringutil.StripQuotes(capture["param_version"].Content()),
			Line:    int64(capture["param_dependency"].StartLineNumber()),
			Column:  int64(capture["param_dependency"].Column()),
		})
	}

	return report
}
