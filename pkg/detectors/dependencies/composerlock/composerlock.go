package composerlock

import (
	"github.com/rs/zerolog/log"
	"github.com/smacker/go-tree-sitter/javascript"

	"github.com/bearer/curio/pkg/detectors/dependencies/depsbase"
	"github.com/bearer/curio/pkg/parser"
	"github.com/bearer/curio/pkg/util/file"
	"github.com/bearer/curio/pkg/util/stringutil"
)

var language = javascript.GetLanguage()
var queryDependencies = parser.QueryMustCompile(language, `
(pair 
	key: (string) @helper_require
    (#match? @helper_require "^\"packages\"$")
    value: (array [(
    	object (
        	(pair
            	key: (string) @helper_name
                (#match? @helper_name "^\"name\"$")
                value: (string) @param_name
            )
        	(pair
            	key: (string) @helper_version
                (#match? @helper_version "^\"version\"$")
                value: (string) @param_version
            )
        )
    )]) 
)`)

func Discover(f *file.FileInfo) (report *depsbase.DiscoveredDependency) {
	report = &depsbase.DiscoveredDependency{}
	report.Provider = "composerlock"
	report.Language = "PHP"
	tree, err := parser.ParseFile(f, f.Path, javascript.GetLanguage())
	if err != nil {
		log.Error().Msgf("%s: there was an error while parsing the composer lock file: %s", report.Provider, err.Error())
		return nil
	}
	defer tree.Close()

	captures := tree.QueryMustPass(queryDependencies)

	for _, capture := range captures {
		if stringutil.StripQuotes(capture["helper_require"].Content()) != "packages" ||
			stringutil.StripQuotes(capture["helper_name"].Content()) != "name" ||
			stringutil.StripQuotes(capture["helper_version"].Content()) != "version" {
			continue
		}

		report.Dependencies = append(report.Dependencies, depsbase.Dependency{
			Name:    stringutil.StripQuotes(capture["param_name"].Content()),
			Version: stringutil.StripQuotes(capture["param_version"].Content()),
			Line:    int64(capture["param_name"].LineNumber()),
			Column:  int64(capture["param_name"].Column()),
		})
	}

	return report
}
