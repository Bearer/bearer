package buildgradle

import (
	grdlparser "github.com/bearer/curio/pkg/detectors/dependencies/buildgradle/parser"
	"github.com/bearer/curio/pkg/detectors/dependencies/depsbase"
	"github.com/bearer/curio/pkg/util/file"
)

// Discover parses build.gradle file and add discovered dependencies to report
func Discover(file *file.FileInfo) (report *depsbase.DiscoveredDependency) {
	return grdlparser.Discover(file)
}
