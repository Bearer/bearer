package nuget

import (
	"github.com/bearer/bearer/pkg/detectors/dependencies/depsbase"
	"github.com/bearer/bearer/pkg/parser"
	"github.com/bearer/bearer/pkg/util/file"
	"github.com/bearer/bearer/pkg/util/stringutil"
	"github.com/rs/zerolog/log"
	"github.com/smacker/go-tree-sitter/javascript"
)

var language = javascript.GetLanguage()

var queryDependencies = parser.QueryMustCompile(language, `
(pair
	key: (string) @helper_dependencies
    (#match? @helper_dependencies "^\"dependencies\"$")
    value:
    	(object
    		(pair
        		value:
                	(object
                        (pair
                          key: (string) @param_dependency
                          value: (object
                              (pair
                                  key: (string) @helper_resolved
                                  (#match? @helper_resolved "^\"resolved\"$")
                                  value: (string) @param_version
                              )
                          )
                        )
        			)
        	)
    	)
)
`)

//	dependencies: {
//		name: version
//	}
var queryNestedDependencies = parser.QueryMustCompile(language, `
(pair
	key: (string) @helper_dependencies
    (#match? @helper_dependencies "^\"dependencies\"$")
    value: (object
    	(pair
        	key: (string) @param_dependency
            value: (string) @param_version
    	)
    )
)
`)

func Discover(f *file.FileInfo) (report *depsbase.DiscoveredDependency) {
	report = &depsbase.DiscoveredDependency{}
	report.Provider = "nuget"
	report.Language = "C#"
	report.PackageManager = "nuget"
	tree, err := parser.ParseFile(f, f.Path, language)
	if err != nil {
		log.Error().Msgf("%s: there was an error while parsing the file: %s", report.Provider, err.Error())
		return nil
	}
	defer tree.Close()

	captures := tree.QueryMustPass(queryDependencies)
	for _, capture := range captures {
		if stringutil.StripQuotes(capture["helper_dependencies"].Content()) != "dependencies" ||
			stringutil.StripQuotes(capture["helper_resolved"].Content()) != "resolved" {
			continue
		}

		report.Dependencies = append(report.Dependencies, depsbase.Dependency{
			Name:    stringutil.StripQuotes(capture["param_dependency"].Content()),
			Version: stringutil.StripQuotes(capture["param_version"].Content()),
			Line:    int64(capture["param_dependency"].LineNumber()),
			Column:  int64(capture["param_dependency"].Column()),
		})
	}

	captures = tree.QueryMustPass(queryNestedDependencies)
	for _, capture := range captures {
		if stringutil.StripQuotes(capture["helper_dependencies"].Content()) != "dependencies" {
			continue
		}

		report.Dependencies = append(report.Dependencies, depsbase.Dependency{
			Name:    stringutil.StripQuotes(capture["param_dependency"].Content()),
			Version: stringutil.StripQuotes(capture["param_version"].Content()),
			Line:    int64(capture["param_dependency"].LineNumber()),
			Column:  int64(capture["param_dependency"].Column()),
		})
	}

	report.Dependencies = filterDuplicates(report.Dependencies)

	return report
}

func filterDuplicates(initial []depsbase.Dependency) (filtered []depsbase.Dependency) {
	keys := make(map[string]bool)

	filtered = make([]depsbase.Dependency, 0)

	for _, pkg := range initial {
		key := pkg.Name + "@" + pkg.Version
		if _, found := keys[key]; !found {
			keys[key] = true
			filtered = append(filtered, pkg)
		}
	}

	return
}
