package paketdependencies

import (
	"bytes"
	"os"
	"regexp"
	"strings"

	"github.com/bearer/curio/pkg/detectors/dependencies/depsbase"
	"github.com/bearer/curio/pkg/util/file"
	"github.com/bearer/curio/pkg/util/linescanner"
	"github.com/rs/zerolog/log"
)

const (
	nugetType  = "nuget"
	gitType    = "git"
	githubType = "github"
)

var (
	githubStripBranchRegexp = regexp.MustCompile(":.*")
)

type rawDependency struct {
	Line       string
	LineNumber int64
}

type document struct {
	Dependencies []rawDependency
}

func (d rawDependency) IsValid() bool {
	if strings.HasPrefix(d.Line, "//") || strings.HasPrefix(d.Line, "source") || len(d.Line) == 0 {
		return false
	}
	return true
}

func nugetDependency(lineNumber int64, fields []string) depsbase.Dependency {
	return depsbase.Dependency{
		Name:    fields[1],
		Version: nugetVersion(fields),
		Line:    lineNumber,
		Column:  0,
	}
}

func nugetVersion(fields []string) string {
	if len(fields) == 4 {
		return fields[3]
	}

	if len(fields) == 3 {
		return fields[2]
	}

	return "N/A"
}

func gitDependency(lineNumber int64, fields []string) depsbase.Dependency {
	return depsbase.Dependency{
		Name:    fields[1],
		Version: "N/A",
		Line:    lineNumber,
		Column:  0,
	}
}

func githubDependency(lineNumber int64, fields []string) depsbase.Dependency {
	return depsbase.Dependency{
		Name:    githubStripBranchRegexp.ReplaceAllString(fields[1], ""),
		Version: "N/A",
		Line:    lineNumber,
		Column:  0,
	}
}

func Discover(file *file.FileInfo) (report *depsbase.DiscoveredDependency) {
	report = &depsbase.DiscoveredDependency{}
	report.Provider = "paket-dependencies"
	report.Language = "C#"

	result := document{}

	fileBytes, err := os.ReadFile(file.AbsolutePath)
	if err != nil {
		log.Error().Msgf("%s: there was an error while opening the file: %s", report.Provider, err.Error())
		return nil
	}

	scanner := linescanner.New(bytes.NewBuffer(fileBytes))

	for scanner.Scan() {
		line := scanner.Text()
		dep := rawDependency{Line: line, LineNumber: int64(scanner.LineNumber())}
		if dep.IsValid() {
			result.Dependencies = append(result.Dependencies, dep)
		}
	}

	for _, dependency := range result.Dependencies {
		fields := strings.Fields(dependency.Line)
		if len(fields) < 2 {
			continue
		}

		switch fields[0] {
		case nugetType:
			report.Dependencies = append(report.Dependencies, nugetDependency(dependency.LineNumber, fields))
		case gitType:
			report.Dependencies = append(report.Dependencies, gitDependency(dependency.LineNumber, fields))
		case githubType:
			report.Dependencies = append(report.Dependencies, githubDependency(dependency.LineNumber, fields))
		}
	}

	return report
}
