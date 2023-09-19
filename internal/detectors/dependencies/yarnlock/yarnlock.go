package yarnlock

import (
	"bytes"
	"os"
	"regexp"
	"strings"

	"github.com/bearer/bearer/internal/detectors/dependencies/depsbase"
	"github.com/bearer/bearer/internal/util/file"
	"github.com/bearer/bearer/internal/util/linescanner"
	"github.com/rs/zerolog/log"
)

const supportedFileFormatVersion = "1"

var documentVersionRegexp *regexp.Regexp
var dependencySectionRegexp *regexp.Regexp
var dependencyNameRegexpSlit *regexp.Regexp
var dependencyLockedVersionRegexp *regexp.Regexp

func Discover(f *file.FileInfo) (report *depsbase.DiscoveredDependency) {
	report = &depsbase.DiscoveredDependency{}
	report.Provider = "yarn.lock"
	report.Language = "JavaScript"
	report.PackageManager = "npm"

	fileBytes, err := os.ReadFile(f.AbsolutePath)
	if err != nil {
		log.Error().Msgf("%s: there was an error while opening the file: %s", report.Provider, err.Error())
		return nil
	}

	depfile := linescanner.New(bytes.NewBuffer(fileBytes))

	var version string
	var dep depsbase.Dependency
	for depfile.Scan() {
		line := depfile.Text()

		switch {
		// Empty line

		case len(line) == 0:
			continue

			// Comment line

		case line[:1] == "#":
			// Look for file format version
			matches := documentVersionRegexp.FindStringSubmatch(line)
			if len(matches) == 0 {
				continue
			}
			version = matches[1]

			// Dependency Section (first line)

		case dependencySectionRegexp.MatchString(line):
			// TODO: do this only once at beginning of parsing
			if version != supportedFileFormatVersion {
				return nil
			}
			matches := dependencySectionRegexp.FindStringSubmatch(line)
			dep.Name = strings.Trim(dependencyNameRegexpSlit.Split(matches[1], -1)[0], `"`)
			dep.Line = int64(depfile.LineNumber())
			dep.Column = 0
			// Dependency Section (line showing locked version)

		case dependencyLockedVersionRegexp.MatchString(line):
			matches := dependencyLockedVersionRegexp.FindStringSubmatch(line)
			dep.Version = strings.Trim(matches[1], `"`)
			report.Dependencies = append(report.Dependencies, dep)
			dep = depsbase.Dependency{}
		}
	}

	return report
}

func init() {
	// compile regexps
	documentVersionRegexp = regexp.MustCompile(`\# yarn lockfile v(\d)`)
	dependencySectionRegexp = regexp.MustCompile(`^(\S.*):$`)
	dependencyNameRegexpSlit = regexp.MustCompile(`\b@`)
	dependencyLockedVersionRegexp = regexp.MustCompile(`^\s+version\s+"(\S+)"`)
}
