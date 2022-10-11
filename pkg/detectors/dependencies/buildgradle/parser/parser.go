package grdlparser

import (
	"bytes"
	"os"
	"regexp"

	"github.com/bearer/curio/pkg/detectors/dependencies/depsbase"
	"github.com/bearer/curio/pkg/util/file"
	"github.com/bearer/curio/pkg/util/linescanner"
	"github.com/rs/zerolog/log"
)

var (
	pattern             *regexp.Regexp
	patternInlineParams *regexp.Regexp
	groupRegexp         *regexp.Regexp
	nameRegexp          *regexp.Regexp
	versionRegexp       *regexp.Regexp
)

func init() {
	pattern = regexp.MustCompile(`(?:implementation|api|compile)\(?\s+(?:'|")(?P<dependency>.*):(?P<name>.*):(?P<version>.*)(?:'|")`)
	patternInlineParams = regexp.MustCompile(`(?:implementation|api|compile)\(?\s+(.*group:.*)`)
	groupRegexp = regexp.MustCompile(`group:[^'"]*(?:'|")([^'"]+)(?:'|")`)
	nameRegexp = regexp.MustCompile(`name:[^'"]*(?:'|")([^'"]+)(?:'|")`)
	versionRegexp = regexp.MustCompile(`version:[^'"]*(?:'|")([^'"]+)(?:'|")`)
}

// Discover parses build.gradle file and add discovered dependencies to report
func Discover(file *file.FileInfo) (report *depsbase.DiscoveredDependency) {
	report = &depsbase.DiscoveredDependency{}
	report.Provider = "gradleparser"
	report.Language = "Java"
	report.PackageManager = "maven"

	fileBytes, err := os.ReadFile(file.AbsolutePath)
	if err != nil {
		log.Error().Msgf("%s: there was an error while opening the file: %s", report.Provider, err.Error())
		return nil
	}

	scanner := linescanner.New(bytes.NewBuffer(fileBytes))

	for scanner.Scan() {
		line := scanner.Text()

		switch {

		case pattern.MatchString(line):
			for _, e := range pattern.FindAllSubmatch([]byte(line), -1) {
				group := string(e[1])
				name := string(e[2])

				report.Dependencies = append(
					report.Dependencies,
					depsbase.Dependency{
						Name:    name,
						Group:   group,
						Version: string(e[3]),
						Line:    int64(scanner.LineNumber()),
						Column:  0,
					})
			}
		case patternInlineParams.MatchString(line):
			name := ""
			group := ""
			version := ""

			if nameRegexp.MatchString(line) {
				name = string(nameRegexp.FindSubmatch([]byte(line))[1])
			}
			if versionRegexp.MatchString(line) {
				version = string(versionRegexp.FindSubmatch([]byte(line))[1])
			}
			if groupRegexp.MatchString(line) {
				group = string(groupRegexp.FindSubmatch([]byte(line))[1])
			}

			report.Dependencies = append(report.Dependencies, depsbase.Dependency{
				Name:    name,
				Group:   group,
				Version: version,
				Line:    int64(scanner.LineNumber()),
				Column:  0,
			})
		}
	}

	return report
}
