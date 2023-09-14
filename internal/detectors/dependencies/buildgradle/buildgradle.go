package buildgradle

import (
	grdlparser "github.com/bearer/bearer/internal/detectors/dependencies/buildgradle/parser"
	"github.com/bearer/bearer/internal/detectors/dependencies/depsbase"
	"github.com/bearer/bearer/internal/util/file"
)

// Discover parses build.gradle file and add discovered dependencies to report
func Discover(file *file.FileInfo) (report *depsbase.DiscoveredDependency) {
	return grdlparser.Discover(file)
}
