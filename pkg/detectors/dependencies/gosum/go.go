package gosum

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/bearer/bearer/pkg/detectors/dependencies/depsbase"
	"github.com/bearer/bearer/pkg/util/file"
	"github.com/bearer/bearer/pkg/util/linescanner"
	"github.com/rs/zerolog/log"
)

func Discover(file *file.FileInfo) (report *depsbase.DiscoveredDependency) {
	report = &depsbase.DiscoveredDependency{}
	report.Provider = "gosum"
	report.Language = "go"
	report.PackageManager = "go"

	fileBytes, err := os.ReadFile(file.AbsolutePath)
	if err != nil {
		log.Error().Msgf("%s: there was an error while opening the file: %s", report.Provider, err.Error())
		return nil
	}

	set := make(map[string]bool)
	scanner := linescanner.New(bytes.NewBuffer(fileBytes))
	for scanner.Scan() {
		parts := strings.Fields(scanner.Text())
		if len(parts) != 3 {
			continue
		}
		name, version := parts[0], parts[1]
		version = strings.TrimSuffix(version, "/go.mod")
		version = strings.TrimSuffix(version, "+incompatible")
		id := fmt.Sprintf("%s:%s", name, version)
		if set[id] {
			continue
		}
		report.Dependencies = append(report.Dependencies, depsbase.Dependency{Name: name, Version: version, Line: int64(scanner.LineNumber()), Column: 0})
		set[id] = true
	}

	return report
}
