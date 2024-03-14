package detectors

import (
	"context"
	"fmt"
	"path/filepath"
	"runtime/debug"
	"slices"

	"github.com/rs/zerolog/log"

	"github.com/bearer/bearer/internal/commands/process/settings"
	"github.com/bearer/bearer/internal/detectors/beego"
	"github.com/bearer/bearer/internal/detectors/csharp"
	"github.com/bearer/bearer/internal/detectors/custom"
	"github.com/bearer/bearer/internal/detectors/dependencies"
	"github.com/bearer/bearer/internal/detectors/django"
	"github.com/bearer/bearer/internal/detectors/dotnet"
	"github.com/bearer/bearer/internal/detectors/envfile"
	"github.com/bearer/bearer/internal/detectors/gitleaks"
	"github.com/bearer/bearer/internal/detectors/golang"
	"github.com/bearer/bearer/internal/detectors/graphql"
	"github.com/bearer/bearer/internal/detectors/html"
	"github.com/bearer/bearer/internal/detectors/ipynb"
	"github.com/bearer/bearer/internal/detectors/java"
	"github.com/bearer/bearer/internal/detectors/javascript"
	"github.com/bearer/bearer/internal/detectors/openapi"
	"github.com/bearer/bearer/internal/detectors/php"
	"github.com/bearer/bearer/internal/detectors/proto"
	"github.com/bearer/bearer/internal/detectors/python"
	"github.com/bearer/bearer/internal/detectors/rails"
	"github.com/bearer/bearer/internal/detectors/ruby"
	"github.com/bearer/bearer/internal/detectors/simple"
	"github.com/bearer/bearer/internal/detectors/spring"
	"github.com/bearer/bearer/internal/detectors/sql"
	"github.com/bearer/bearer/internal/detectors/symfony"
	"github.com/bearer/bearer/internal/detectors/tsx"
	"github.com/bearer/bearer/internal/detectors/typescript"
	"github.com/bearer/bearer/internal/detectors/yamlconfig"
	"github.com/bearer/bearer/internal/scanner"
	"github.com/bearer/bearer/internal/scanner/stats"
	"github.com/bearer/bearer/internal/util/file"

	"github.com/bearer/bearer/internal/detectors/types"
	"github.com/bearer/bearer/internal/parser/nodeid"

	reporttypes "github.com/bearer/bearer/internal/report"
	reportdetectors "github.com/bearer/bearer/internal/report/detectors"
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

func SetupLegacyDetector(config map[string]*settings.Rule) error {
	detector := customDetector.Detector.(*custom.Detector)

	return detector.CompileRules(config)
}

func Registrations(scanners []string) []InitializedDetector {
	// The order of these is important, the first one to claim a file will win
	detectors := []InitializedDetector{}

	if slices.Contains(scanners, "secrets") {
		// Enable GitLeaks
		detectors = append(
			detectors,
			InitializedDetector{
				reportdetectors.DetectorGitleaks, gitleaks.New(&nodeid.UUIDGenerator{}),
			},
		)
	}

	if slices.Contains(scanners, "sast") {
		detectors = append(
			detectors,
			[]InitializedDetector{
				{reportdetectors.DetectorCustom, customDetector},

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
			}...,
		)
	}

	return detectors
}

func Extract(
	ctx context.Context,
	path string,
	filename string,
	report reporttypes.Report,
	fileStats *stats.FileStats,
	enabledScanners []string,
	sastScanner *scanner.Scanner,
	skipTest bool,
) error {
	return ExtractWithDetectors(
		ctx,
		path,
		filename,
		report,
		fileStats,
		Registrations(enabledScanners),
		sastScanner,
		skipTest,
	)
}

func ExtractWithDetectors(
	ctx context.Context,
	rootDir string,
	filename string,
	report reporttypes.Report,
	fileStats *stats.FileStats,
	allDetectors []InitializedDetector,
	sastScanner *scanner.Scanner,
	skipTest bool,
) error {

	activeDetectors := make(map[InitializedDetector]activeDetector)

	if err := file.IterateFilesList(
		rootDir,
		[]string{filename},
		skipTest,
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
					log.Printf("file %s -> error recovered %s", file.AbsolutePath, r)
					log.Print(string(debug.Stack()))
					report.AddError(file.Path.RelativePath, fmt.Errorf("skipping file: due to panic %s", r))
				}
			}
			defer recovery()

			if err := sastScanner.Scan(ctx, report, fileStats, file); err != nil {
				log.Debug().Msgf("failed to scan file %s: %s", file.RelativePath, err)
				report.AddError(file.RelativePath, fmt.Errorf("failed to scan file: %s", err))
			}

			for _, detector := range allDetectors {
				if ctx.Err() != nil {
					return ctx.Err()
				}

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
