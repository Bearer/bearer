package pipdeptree

import (
	"github.com/bearer/curio/pkg/detectors/dependencies/depsbase"
	"github.com/bearer/curio/pkg/parser"
	"github.com/bearer/curio/pkg/util/file"
	"github.com/bearer/curio/pkg/util/stringutil"
	"github.com/rs/zerolog/log"
	"github.com/smacker/go-tree-sitter/javascript"
)

var language = javascript.GetLanguage()

var queryDependencies = parser.QueryMustCompile(language, `
(object
	(pair
		key: (string) @helper_version
		(#match? @helper_version "^\"installed_version\"$")
		value: (string) @param_version
	)
	(pair
		key: (string) @helper_dependency
		(#match? @helper_dependency "^\"key\"$")
		value: (string) @param_dependency
	)
)
`)

func Discover(f *file.FileInfo) (report *depsbase.DiscoveredDependency) {
	report = &depsbase.DiscoveredDependency{}
	report.Provider = "pipdeptree"
	report.Language = "Python"
	report.PackageManager = "pypi"
	tree, err := parser.ParseFile(f, f.Path, language)
	if err != nil {
		log.Error().Msgf("%s: there was an error while parsing the file: %s", report.Provider, err.Error())
		return nil
	}
	defer tree.Close()

	captures := tree.QueryMustPass(queryDependencies)
	for _, capture := range captures {
		if stringutil.StripQuotes(capture["helper_version"].Content()) != "installed_version" ||
			stringutil.StripQuotes(capture["helper_dependency"].Content()) != "key" {
			continue
		}

		report.Dependencies = append(report.Dependencies, depsbase.Dependency{
			Name:    stringutil.StripQuotes(capture["param_dependency"].Content()),
			Version: stringutil.StripQuotes(capture["param_version"].Content()),
			Line:    int64(capture["param_dependency"].LineNumber()),
			Column:  int64(capture["param_dependency"].Column()),
		})
	}

	return report
}
