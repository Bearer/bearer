package dependencies

import (
	"github.com/bearer/bearer/pkg/detectors/dependencies/buildgradle"
	"github.com/bearer/bearer/pkg/detectors/dependencies/composerjson"
	"github.com/bearer/bearer/pkg/detectors/dependencies/composerlock"
	"github.com/bearer/bearer/pkg/detectors/dependencies/depsbase"
	"github.com/bearer/bearer/pkg/detectors/dependencies/gemfile"
	"github.com/bearer/bearer/pkg/detectors/dependencies/gosum"
	"github.com/bearer/bearer/pkg/detectors/dependencies/ivy"
	"github.com/bearer/bearer/pkg/detectors/dependencies/mvnplugin"
	"github.com/bearer/bearer/pkg/detectors/dependencies/npm"
	"github.com/bearer/bearer/pkg/detectors/dependencies/nuget"
	packageconfig "github.com/bearer/bearer/pkg/detectors/dependencies/package-config"
	packagejson "github.com/bearer/bearer/pkg/detectors/dependencies/package-json"
	paketdependencies "github.com/bearer/bearer/pkg/detectors/dependencies/paket-dependencies"
	"github.com/bearer/bearer/pkg/detectors/dependencies/pipdeptree"
	"github.com/bearer/bearer/pkg/detectors/dependencies/piplock"
	"github.com/bearer/bearer/pkg/detectors/dependencies/poetry"
	pomxml "github.com/bearer/bearer/pkg/detectors/dependencies/pom-xml"
	projectjson "github.com/bearer/bearer/pkg/detectors/dependencies/project-json"
	"github.com/bearer/bearer/pkg/detectors/dependencies/pyproject"
	"github.com/bearer/bearer/pkg/detectors/dependencies/requirements"
	"github.com/bearer/bearer/pkg/detectors/dependencies/yarnlock"
	"github.com/bearer/bearer/pkg/detectors/types"
	"github.com/bearer/bearer/pkg/report"
	"github.com/bearer/bearer/pkg/report/dependencies"
	"github.com/bearer/bearer/pkg/report/detectors"
	"github.com/bearer/bearer/pkg/report/source"
	"github.com/bearer/bearer/pkg/util/file"
)

type detector struct{}

func New() types.Detector {
	return &detector{}
}

func (detector *detector) AcceptDir(dir *file.Path) (bool, error) {
	return true, nil
}

func (detector *detector) ProcessFile(file *file.FileInfo, dir *file.Path, report report.Report) (bool, error) {
	switch file.Base {
	case "Gemfile.lock":
		return discoverDependency(report, file, gemfile.Discover)
	case "package.json":
		return discoverDependency(report, file, packagejson.Discover)
	case "yarn.lock":
		return discoverDependency(report, file, yarnlock.Discover)
	case "maven-dependencies.json", "gemnasium-maven-plugin.json", "gradle-dependencies.json":
		return discoverDependency(report, file, mvnplugin.Discover)
	case "Pipfile.lock":
		return discoverDependency(report, file, piplock.Discover)
	case "package-lock.json", "npm-shrinkwrap.json":
		return discoverDependency(report, file, npm.Discover)
	case "packages.lock.json":
		return discoverDependency(report, file, nuget.Discover)
	case "go.sum":
		return discoverDependency(report, file, gosum.Discover)
	case "project.json":
		return discoverDependency(report, file, projectjson.Discover)
	case "packages.config":
		return discoverDependency(report, file, packageconfig.Discover)
	case "paket.dependencies":
		return discoverDependency(report, file, paketdependencies.Discover)
	case "ivy-report.xml":
		return discoverDependency(report, file, ivy.Discover)
	case "composer.lock":
		return discoverDependency(report, file, composerlock.Discover)
	case "composer.json":
		return discoverDependency(report, file, composerjson.Discover)
	case "pipdeptree.json":
		return discoverDependency(report, file, pipdeptree.Discover)
	case "poetry.lock":
		return discoverDependency(report, file, poetry.Discover)
	case "pyproject.toml":
		return discoverDependency(report, file, pyproject.Discover)
	case "pom.xml":
		return discoverDependency(report, file, pomxml.Discover)
	case "requirements.txt":
		return discoverDependency(report, file, requirements.Discover)
	case "build.gradle":
		return discoverDependency(report, file, buildgradle.Discover)
	}

	return false, nil
}

func discoverDependency(report report.Report, file *file.FileInfo, discover func(file *file.FileInfo) (report *depsbase.DiscoveredDependency)) (bool, error) {
	result := discover(file)

	if result == nil {
		return true, nil
	}

	for _, dep := range result.Dependencies {
		columnNumber := int(dep.Column)
		lineNumber := int(dep.Line)
		report.AddDependency(detectors.Type(result.Provider), dependencies.Dependency{
			Group:          dep.Group,
			Name:           dep.Name,
			Version:        dep.Version,
			PackageManager: result.PackageManager,
		}, source.Source{
			Language:     file.Language,
			LanguageType: file.LanguageTypeString(),
			Filename:     file.RelativePath,
			ColumnNumber: &columnNumber,
			LineNumber:   &lineNumber,
		})
	}

	return true, nil
}
