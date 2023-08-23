package requirements

import (
	"bytes"
	"os"
	"regexp"

	"github.com/bearer/bearer/internal/detectors/dependencies/depsbase"
	"github.com/bearer/bearer/internal/util/file"
	"github.com/bearer/bearer/internal/util/linescanner"
	"github.com/rs/zerolog/log"
)

var lineRegexp = regexp.MustCompile(`^\s*([^#\s~=<>]+)\s*(?:[~=<>]=?([^#,\s]+))?`)

func Discover(f *file.FileInfo) (report *depsbase.DiscoveredDependency) {
	report = &depsbase.DiscoveredDependency{}
	report.Provider = "requirements.txt"
	report.Language = "Python"
	report.PackageManager = "pypi"

	fileBytes, err := os.ReadFile(f.AbsolutePath)
	if err != nil {
		log.Error().Msgf("%s: there was an error while opening the file: %s", report.Provider, err.Error())
		return nil
	}

	scanner := linescanner.New(bytes.NewBuffer(fileBytes))
	for scanner.Scan() {
		match := lineRegexp.FindStringSubmatch(scanner.Text())
		if len(match) == 0 {
			continue
		}

		name := match[1]
		version := match[2]
		if version == "" {
			version = "unknown"
		}

		report.Dependencies = append(report.Dependencies, depsbase.Dependency{
			Name:    name,
			Version: version,
			Line:    int64(scanner.LineNumber()),
			Column:  0,
		})
	}

	return report
}
