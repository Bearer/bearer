package detectors

import (
	"fmt"
	"path/filepath"
	"runtime/debug"

	"github.com/bearer/curio/pkg/commands/process/settings"
	"github.com/bearer/curio/pkg/detectors/beego"
	"github.com/bearer/curio/pkg/detectors/csharp"
	"github.com/bearer/curio/pkg/detectors/custom"
	"github.com/bearer/curio/pkg/detectors/dependencies"
	"github.com/bearer/curio/pkg/detectors/django"
	"github.com/bearer/curio/pkg/detectors/dotnet"
	"github.com/bearer/curio/pkg/detectors/envfile"
	"github.com/bearer/curio/pkg/detectors/gitleaks"
	"github.com/bearer/curio/pkg/detectors/golang"
	"github.com/bearer/curio/pkg/detectors/graphql"
	"github.com/bearer/curio/pkg/detectors/html"
	"github.com/bearer/curio/pkg/detectors/ipynb"
	"github.com/bearer/curio/pkg/detectors/java"
	"github.com/bearer/curio/pkg/detectors/javascript"
	"github.com/bearer/curio/pkg/detectors/openapi"
	"github.com/bearer/curio/pkg/detectors/php"
	"github.com/bearer/curio/pkg/detectors/proto"
	"github.com/bearer/curio/pkg/detectors/python"
	"github.com/bearer/curio/pkg/detectors/rails"
	"github.com/bearer/curio/pkg/detectors/ruby"
	"github.com/bearer/curio/pkg/detectors/simple"
	"github.com/bearer/curio/pkg/detectors/tsx"
	"github.com/bearer/curio/pkg/detectors/typescript"
	"github.com/bearer/curio/pkg/detectors/yamlconfig"
	"github.com/bearer/curio/pkg/util/file"

	"github.com/bearer/curio/pkg/detectors/spring"
	"github.com/bearer/curio/pkg/detectors/sql"
	"github.com/bearer/curio/pkg/detectors/symfony"
	"github.com/bearer/curio/pkg/detectors/types"
	"github.com/bearer/curio/pkg/parser/nodeid"
	"github.com/rs/zerolog/log"

	reporttypes "github.com/bearer/curio/pkg/report"
	reportdetectors "github.com/bearer/curio/pkg/report/detectors"
)

type InitializedDetector struct {
	Type reportdetectors.Type
	types.Detector
}

type activeDetector struct {
	*file.Path
	reporttypes.Report
}

var customDetector = InitializedDetector{reportdetectors.DetectorCustom, custom.New(&nodeid.UUIDGenerator{})}

func SetupCustomDetector(config map[string]settings.Rule) error {
	detector := customDetector.Detector.(*custom.Detector)
	return detector.CompileRules(config)
}

func Registrations() []InitializedDetector {
	// The order of these is important, the first one to claim a file will win

	detectors := []InitializedDetector{
		{reportdetectors.DetectorCustom, customDetector},
		// gitleaks detector
		{reportdetectors.DetectorGitleaks, gitleaks.New(&nodeid.UUIDGenerator{})},
		//dependencies detectors
		{reportdetectors.DetectorDependencies, dependencies.New()},

		{reportdetectors.DetectorBeego, beego.New()},
		{reportdetectors.DetectorGo, golang.New(&nodeid.UUIDGenerator{})},

		{reportdetectors.DetectorDjango, django.New()},
		{reportdetectors.DetectorPython, python.New(&nodeid.UUIDGenerator{})},

		{reportdetectors.DetectorDotnet, dotnet.New()},
		{reportdetectors.DetectorCSharp, csharp.New(&nodeid.UUIDGenerator{})},

		{reportdetectors.DetectorOpenAPI, openapi.New(&nodeid.UUIDGenerator{})},

		{reportdetectors.DetectorEnvFile, envfile.New()},

		{reportdetectors.DetectorJavascript, javascript.New(&nodeid.UUIDGenerator{})},
		{reportdetectors.DetectorTsx, tsx.New(&nodeid.UUIDGenerator{})},
		{reportdetectors.DetectorTypescript, typescript.New(&nodeid.UUIDGenerator{})},

		{reportdetectors.DetectorRails, rails.New(&nodeid.UUIDGenerator{})},
		{reportdetectors.DetectorRuby, ruby.New(&nodeid.UUIDGenerator{})},

		{reportdetectors.DetectorSpring, spring.New()},
		{reportdetectors.DetectorJava, java.New(&nodeid.UUIDGenerator{})},

		{reportdetectors.DetectorSymfony, symfony.New()},
		{reportdetectors.DetectorPHP, php.New(&nodeid.UUIDGenerator{})},

		{reportdetectors.DetectorYamlConfig, yamlconfig.New()},

		{reportdetectors.DetectorSQL, sql.New(&nodeid.UUIDGenerator{})},
		{reportdetectors.DetectorProto, proto.New(&nodeid.UUIDGenerator{})},
		{reportdetectors.DetectorGraphQL, graphql.New(&nodeid.UUIDGenerator{})},

		{reportdetectors.DetectorHTML, html.New(&nodeid.UUIDGenerator{})},

		{reportdetectors.DetectorIPYNB, ipynb.New(&nodeid.UUIDGenerator{})},

		{reportdetectors.DetectorSimple, simple.New()},
	}

	return detectors
}

func Extract(
	path string,
	files []string,
	report reporttypes.Report,
) error {
	return ExtractWithDetectors(path, files, report, Registrations())
}

func ExtractWithDetectors(
	rootDir string,
	files []string,
	report reporttypes.Report,
	allDetectors []InitializedDetector,
) error {

	activeDetectors := make(map[InitializedDetector]activeDetector)

	if err := file.IterateFilesList(
		rootDir,
		files,
		func(dir *file.Path) (bool, error) {
			for _, detector := range allDetectors {
				active, isActive := activeDetectors[detector]

				if isActive && !isParentedBy(active.AbsolutePath, dir.AbsolutePath) {
					delete(activeDetectors, detector)
					isActive = false
				}

				if !isActive {
					activate, err := detector.Detector.AcceptDir(dir)
					if err != nil {
						report.AddError(dir.RelativePath, fmt.Errorf("accept dir failed for detector %s: %s", detector.Type, err))
						continue
					}

					if activate {
						activeDetectors[detector] = activeDetector{Path: dir, Report: report}
					}
				}
			}

			return true, nil
		},
		func(file *file.FileInfo) error {
			recovery := func() {
				if r := recover(); r != nil {
					log.Printf("error recovered %s %s", r, debug.Stack())
					report.AddError(file.Path.RelativePath, fmt.Errorf("skipping file: due to panic %s, %s", r, string(debug.Stack())))
				}
			}
			defer recovery()

			log.Debug().Msgf("processing file %s", file.AbsolutePath)
			for _, detector := range allDetectors {
				active, isActive := activeDetectors[detector]
				if !isActive {
					continue
				}

				if file.Size() == 0 {
					continue
				}

				wasConsumed, err := detector.Detector.ProcessFile(file, active.Path, active.Report)
				if err != nil {
					log.Debug().Msgf("failed to process file %s for detector %s: %s", file.RelativePath, detector.Type, err)
					report.AddError(file.RelativePath, fmt.Errorf("failed to process file for detector %s: %s", detector.Type, err))
					continue
				}

				if wasConsumed {
					break
				}
			}

			return nil
		},
	); err != nil {
		return err
	}

	return nil
}

func isParentedBy(rootPath, path string) bool {
	relativePath, err := filepath.Rel(rootPath, path)
	if err != nil {
		return false
	}

	// directory can have one letter names like "R"
	if len(relativePath) == 1 {
		return true
	} else {
		return relativePath[:2] != ".."
	}
}
