package buildgradle

import (
	"regexp"

	grdlparser "github.com/bearer/curio/pkg/detectors/dependencies/buildgradle/parser"
	"github.com/bearer/curio/pkg/detectors/dependencies/depsbase"
	"github.com/bearer/curio/pkg/util/file"
)

var (
	pattern           *regexp.Regexp
	stateDependencies = "stateDependencies"
)

func init() {
	pattern = regexp.MustCompile(`implementation\(?\s+(?:'|")(?P<dependency>.*):(?P<name>.*):(?P<version>.*)(?:'|")`)
}

// Discover parses build.gradle file and add discovered dependencies to report
func Discover(file *file.FileInfo) (report *depsbase.DiscoveredDependency) {
	return grdlparser.Discover(file)
}
