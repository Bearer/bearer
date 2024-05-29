package ivy

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/rs/zerolog/log"

	"github.com/bearer/bearer/pkg/detectors/dependencies/depsbase"
	"github.com/bearer/bearer/pkg/parser"
	"github.com/bearer/bearer/pkg/parser/sitter/xml"
	"github.com/bearer/bearer/pkg/util/file"
	"github.com/bearer/bearer/pkg/util/stringutil"
)

var language = xml.GetLanguage()

var query = `
(element
	(_
    	(tag_name) @helper_module
        (#match? @helper_module "^module$")
        (attribute
        	(attribute_name) @helper_organisation
            (#match? @helper_organisation "^organisation$")
			(attribute_value) @param_organisation_name
        )
        (attribute
        	(attribute_name) @helper_organisation_name
            (#match? @helper_organisation_name "^name$")
			(attribute_value) @param_module_name
        )
    )
    (element
      	(_
            (tag_name) @helper_revision
            (#match? @helper_revision "^revision$")
            (attribute
                (attribute_name) @helper_revision_name
                (#match? @helper_revision_name "^name$")
                (attribute_value) @param_version
            )
      	)
    )
)
`

var queryDependencies = parser.QueryMustCompile(language, fmt.Sprintf(`
(document
	(element
    	(element
      		(start_tag
              (tag_name) @helper_dependencies
              (#match? @helper_dependencies "^dependencies$")
          	)

            %s
      )
    )
)
`, query))

var regexRemoveStyleSheet = regexp.MustCompile(`<\?xml-stylesheet.+`)

// formattedName is responsible for formatting module names which are postfixed with the scala version in the ivy report
// e.g. when sbt-dependency-graph plugin outputs:
// <module organisation="org.parboiled" name="parboiled-scala_2.13">...</module>
func formattedName(inputname string) string {
	parts := strings.Split(inputname, "_")
	var name string
	if len(parts) == 1 {
		name = parts[0]
	} else {
		name = strings.Join(parts[0:len(parts)-1], "")
	}
	return name
}

func Discover(f *file.FileInfo) (report *depsbase.DiscoveredDependency) {
	report = &depsbase.DiscoveredDependency{}
	report.Provider = "ivy"
	report.Language = "Java"
	report.PackageManager = "maven"

	bytes, err := os.ReadFile(f.AbsolutePath)
	if err != nil {
		log.Error().Msgf("%s: there was an error while reading the file: %s", report.Provider, err.Error())
	}

	formattedBytes := regexRemoveStyleSheet.ReplaceAll(bytes, []byte(""))

	tree, err := parser.ParseBytes(f, f.Path, formattedBytes, language, 1)

	if err != nil {
		log.Error().Msgf("%s: there was an error while parsing the file: %s", report.Provider, err.Error())
		return nil
	}
	defer tree.Close()

	captures := tree.QueryMustPass(queryDependencies)
	for _, capture := range captures {
		if stringutil.StripQuotes(capture["helper_module"].Content()) != "module" ||
			stringutil.StripQuotes(capture["helper_organisation"].Content()) != "organisation" ||
			stringutil.StripQuotes(capture["helper_organisation_name"].Content()) != "name" ||
			stringutil.StripQuotes(capture["helper_revision"].Content()) != "revision" ||
			stringutil.StripQuotes(capture["helper_revision_name"].Content()) != "name" {
			continue
		}

		moduleName := stringutil.StripQuotes(capture["param_module_name"].Content())
		organisationName := stringutil.StripQuotes(capture["param_organisation_name"].Content())
		version := stringutil.StripQuotes(capture["param_version"].Content())

		report.Dependencies = append(report.Dependencies, depsbase.Dependency{
			Name:    formattedName(moduleName),
			Group:   organisationName,
			Version: version,
			Line:    int64(capture["param_module_name"].StartLineNumber()),
			Column:  int64(capture["param_module_name"].Column()),
		})
	}

	return
}
