package testhelper

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/bradleyjkemp/cupaloy"
	"github.com/hhatto/gocloc"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"

	"github.com/bearer/bearer/pkg/classification"
	"github.com/bearer/bearer/pkg/commands"
	"github.com/bearer/bearer/pkg/commands/process/filelist"
	"github.com/bearer/bearer/pkg/commands/process/filelist/files"
	"github.com/bearer/bearer/pkg/commands/process/settings"
	settingsloader "github.com/bearer/bearer/pkg/commands/process/settings/loader"
	"github.com/bearer/bearer/pkg/commands/process/settings/rules"
	"github.com/bearer/bearer/pkg/detectors"
	engine "github.com/bearer/bearer/pkg/engine"
	engineimpl "github.com/bearer/bearer/pkg/engine/implementation"
	"github.com/bearer/bearer/pkg/flag"
	"github.com/bearer/bearer/pkg/report/output"
	"github.com/bearer/bearer/pkg/report/writer"
	"github.com/bearer/bearer/pkg/scanner"
	"github.com/bearer/bearer/pkg/scanner/language"
	"github.com/bearer/bearer/pkg/types"
	util "github.com/bearer/bearer/pkg/util/output"
	"github.com/bearer/bearer/pkg/util/set"
	"github.com/bearer/bearer/pkg/version_check"
)

type Runner struct {
	engine     engine.Engine
	config     settings.Config
	classifier *classification.Classifier
	scanner    *scanner.Scanner
}

func GetRunner(t *testing.T, ruleBytes []byte, lang language.Language) *Runner {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out: os.Stderr,
		FormatTimestamp: func(i interface{}) string {
			timestamp, _ := time.Parse(time.RFC3339, i.(string))
			return timestamp.Format("2006-01-02 15:04:05")
		},
	})

	engine := engineimpl.New([]language.Language{lang})
	config := buildConfig(t, engine, ruleBytes)

	if err := engine.Initialize("trace"); err != nil {
		t.Fatalf("failed to initialize engine: %s", err)
	}

	classifier, err := classification.NewClassifier(&classification.Config{Config: config})
	if err != nil {
		t.Fatalf("failed to create classifier: %s", err)
	}

	scanner, err := scanner.New(engine, classifier.Schema, config.Rules)
	if err != nil {
		t.Fatalf("failed to create scanner: %s", err)
	}

	runner := &Runner{
		engine:     engine,
		config:     config,
		classifier: classifier,
		scanner:    scanner,
	}
	runtime.SetFinalizer(runner, func(runner *Runner) {
		runner.scanner.Close()
		runner.engine.Close()
	})

	return runner
}

func (runner *Runner) RunTest(t *testing.T, testdataPath string, snapshotPath string) {
	dummyGoclocLanguage := gocloc.Language{}
	dummyGoclocResult := gocloc.Result{
		Total: &dummyGoclocLanguage,
		Files: map[string]*gocloc.ClocFile{
			testdataPath: {Code: 10},
		},
		Languages:     map[string]*gocloc.Language{},
		MaxPathLength: 0,
	}

	absTestdataPath, err := filepath.Abs(testdataPath)
	if err != nil {
		t.Fatalf("failed to get absolute target: %s", err)
	}

	fileList, err := filelist.Discover(nil, absTestdataPath, &dummyGoclocResult, runner.config)
	if err != nil {
		t.Fatalf("failed to discover files: %s", err)
	}

	if len(fileList.Files) == 0 {
		t.Fatal("no scannable files found")
	}

	for _, file := range fileList.Files {
		testName := strings.TrimSuffix(file.FilePath, filepath.Ext(file.FilePath))
		t.Run(testName, func(tt *testing.T) {
			runner.scanSingleFile(tt, testdataPath, file, snapshotPath)
		})
	}
}

func (runner *Runner) scanSingleFile(t *testing.T, testDataPath string, fileRelativePath files.File, snapshotsPath string) {
	detectorsReportFile, err := os.CreateTemp("", "report.jsonl")
	if err != nil {
		t.Fatalf("failed to create tmp report file: %s", err)
	}
	defer detectorsReportFile.Close()

	detectorsReportPath := detectorsReportFile.Name()
	if err != nil {
		t.Fatalf("failed to get absolute path of report file: %s", err)
	}

	if err = detectors.Extract(
		context.Background(),
		testDataPath,
		fileRelativePath.FilePath,
		&writer.Detectors{
			Classifier: runner.classifier,
			File:       detectorsReportFile,
		},
		nil,
		[]string{"sast"},
		runner.scanner,
		false,
	); err != nil {
		t.Fatalf("failed to do scan %s", err)
	}

	runner.config.Scan.Target = testDataPath
	reportData, err := output.GetData(
		types.Report{
			Path:     detectorsReportPath,
			HasFiles: true,
		},
		runner.config,
		nil,
		nil,
	)
	if err != nil {
		t.Fatalf("failed to get output: %s", err)
	}

	report, err := util.ReportYAML(reportData.FindingsBySeverity)
	if err != nil {
		t.Fatalf("failed to encoded to yaml: %s", err)
	}

	cupaloy.NewDefaultConfig().WithOptions(
		cupaloy.SnapshotSubdirectory(snapshotsPath),
		cupaloy.SnapshotFileExtension(".yml"),
	).SnapshotT(t, report)
}

func buildConfig(t *testing.T, engine engine.Engine, ruleBytes []byte) settings.Config {
	err := commands.ScanFlags.BindForConfigInit(commands.NewScanCommand(nil))
	if err != nil {
		t.Fatalf("failed to bind flags: %s", err)
	}

	configFlags, err := commands.ScanFlags.ToOptions([]string{})
	if err != nil {
		t.Fatalf("failed to generate default flags: %s", err)
	}
	configFlags.Format = flag.FormatYAML
	configFlags.Report = flag.ReportSecurity
	configFlags.Quiet = true
	configFlags.DisableDefaultRules = true
	configFlags.ExternalRuleDir = []string{}
	configFlags.DisableVersionCheck = true
	configFlags.IgnoreGit = true

	meta := &version_check.VersionMeta{
		Rules: version_check.RuleVersionMeta{
			Packages: make(map[string]string),
		},
		Binary: version_check.BinaryVersionMeta{
			Latest:  true,
			Message: "",
		},
	}

	rules := getRulesFromYaml(t, ruleBytes)
	languageIDs := set.New[string]()
	for _, rule := range rules {
		languageIDs.AddAll(rule.Languages)
	}

	config, err := settingsloader.FromOptions(configFlags, meta, engine, languageIDs.Items())
	if err != nil {
		t.Fatalf("failed to generate default scan settings: %s", err)
	}

	config.Rules = rules

	return config
}

func getRulesFromYaml(t *testing.T, ruleBytes []byte) map[string]*settings.Rule {
	var ruleDefinition settings.RuleDefinition
	err := yaml.Unmarshal(ruleBytes, &ruleDefinition)
	if err != nil {
		t.Fatalf("failed to unmarshal rule %s", err)
	}

	definitions := map[string]settings.RuleDefinition{
		ruleDefinition.Metadata.ID: ruleDefinition,
	}
	enabledRules := map[string]struct{}{
		ruleDefinition.Metadata.ID: {},
	}

	return rules.BuildRules(definitions, enabledRules)
}
