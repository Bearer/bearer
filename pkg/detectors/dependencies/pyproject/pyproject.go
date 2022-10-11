package pyproject

import (
	"github.com/bearer/curio/pkg/detectors/dependencies/depsbase"
	"github.com/bearer/curio/pkg/parser"
	"github.com/bearer/curio/pkg/parser/sitter/toml"
	"github.com/bearer/curio/pkg/util/file"
	"github.com/bearer/curio/pkg/util/stringutil"
	"github.com/rs/zerolog/log"
)

var language = toml.GetLanguage()

var queryDependencies = parser.QueryMustCompile(language, `
(table
	(dotted_key) @helper_table_name
    (#match? @helper_table_name "^tool.poetry.dependencies$")
	(pair
		(bare_key) @param_dependency
		(string) @param_version
	)
)
`)

func Discover(f *file.FileInfo) (report *depsbase.DiscoveredDependency) {
	report = &depsbase.DiscoveredDependency{}
	report.Provider = "pyproject"
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
		if stringutil.StripQuotes(capture["helper_table_name"].Content()) != "tool.poetry.dependencies" {
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
