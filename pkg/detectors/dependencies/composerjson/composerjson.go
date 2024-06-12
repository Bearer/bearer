package composerjson

import (
	"github.com/rs/zerolog/log"
	"github.com/smacker/go-tree-sitter/javascript"

	"github.com/bearer/bearer/pkg/detectors/dependencies/depsbase"
	"github.com/bearer/bearer/pkg/parser"
	"github.com/bearer/bearer/pkg/util/file"
	"github.com/bearer/bearer/pkg/util/stringutil"
)

var language = javascript.GetLanguage()
var queryDependencies = parser.QueryMustCompile(language, `
(pair
	key: (string) @helper_require
    (#match? @helper_require "^\"require\"$")
    value: (object) @param_value
)
`)

func Discover(f *file.FileInfo) (report *depsbase.DiscoveredDependency) {
	report = &depsbase.DiscoveredDependency{}
	report.Provider = "composerjson"
	report.Language = "php"
	report.PackageManager = "packagist"
	tree, err := parser.ParseFile(f, f.Path, javascript.GetLanguage())
	if err != nil {
		log.Error().Msgf("%s: there was an error while parsing the composer file: %s", report.Provider, err.Error())
		return nil
	}
	defer tree.Close()

	captures := tree.QueryMustPass(queryDependencies)

	for _, capture := range captures {
		if stringutil.StripQuotes(capture["helper_require"].Content()) != "require" {
			continue
		}
		dependecyPairs := capture["param_value"]
		for i := 0; i < dependecyPairs.ChildCount(); i++ {
			dependecyPair := dependecyPairs.Child(i)
			pkgName := dependecyPair.ChildByFieldName("key")
			if pkgName == nil {
				return nil
			}
			version := dependecyPair.ChildByFieldName("value")
			if version == nil {
				return nil
			}

			report.Dependencies = append(report.Dependencies, depsbase.Dependency{
				Name:    stringutil.StripQuotes(pkgName.Content()),
				Version: stringutil.StripQuotes(version.Content()),
				Line:    int64(pkgName.StartLineNumber()),
				Column:  int64(pkgName.Column()),
			})
		}
	}

	return report
}
