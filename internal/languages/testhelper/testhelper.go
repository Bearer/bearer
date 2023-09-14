package testhelper

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bradleyjkemp/cupaloy"
	"github.com/hhatto/gocloc"
	"github.com/rs/zerolog"
	"gopkg.in/yaml.v3"

	"github.com/bearer/bearer/internal/commands"
	"github.com/bearer/bearer/internal/commands/process/filelist"
	"github.com/bearer/bearer/internal/commands/process/filelist/files"
	"github.com/bearer/bearer/internal/commands/process/orchestrator/work"
	"github.com/bearer/bearer/internal/commands/process/orchestrator/worker"
	"github.com/bearer/bearer/internal/commands/process/settings"
	"github.com/bearer/bearer/internal/flag"
	"github.com/bearer/bearer/internal/report/output"
	"github.com/bearer/bearer/internal/types"
	util "github.com/bearer/bearer/internal/util/output"
	"github.com/bearer/bearer/internal/version_check"
)

type Runner struct {
	config settings.Config
	worker worker.Worker
}

func GetRunner(t *testing.T, ruleBytes []byte, lang string) *Runner {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	err := commands.ScanFlags.BindForConfigInit(commands.NewScanCommand())
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

	meta := &version_check.VersionMeta{
		Rules: version_check.RuleVersionMeta{
			Packages: make(map[string]string),
		},
		Binary: version_check.BinaryVersionMeta{
			Latest:  true,
			Message: "",
		},
	}
	config, err := settings.FromOptions(configFlags, meta)
	if err != nil {
		t.Fatalf("failed to generate default scan settings: %s", err)
	}

	config.Rules = getRulesFromYaml(t, ruleBytes)

	worker := worker.Worker{}
	err = worker.Setup(config)
	if err != nil {
		t.Fatalf("failed to setup scan worker: %s", err)
	}

	runner := &Runner{
		worker: worker,
		config: config,
	}

	return runner
}

func getRulesFromYaml(t *testing.T, ruleBytes []byte) map[string]*settings.Rule {
	var ruleDefinition settings.RuleDefinition
	err := yaml.Unmarshal(ruleBytes, &ruleDefinition)
	if err != nil {
		t.Fatalf("failed to unmarshal rule %s", err)
	}

	rules := map[string]settings.RuleDefinition{
		ruleDefinition.Metadata.ID: ruleDefinition,
	}
	enabledRules := map[string]struct{}{
		ruleDefinition.Metadata.ID: {},
	}

	return settings.BuildRules(rules, enabledRules)
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
		myfile := file
		ext := filepath.Ext(file.FilePath)
		testName := "/" + strings.TrimSuffix(file.FilePath, ext) + ".yml"
		t.Run(testName, func(t *testing.T) {
			runner.scanSingleFile(t, testdataPath, myfile, snapshotPath)
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

	_, err = runner.worker.Scan(context.Background(), work.ProcessRequest{
		File:       fileRelativePath,
		ReportPath: detectorsReportPath,
		Repository: work.Repository{
			Dir: testDataPath,
		},
	})

	if err != nil {
		t.Fatalf("failed to do scan %s", err)
	}

	runner.config.Scan.Target = testDataPath
	reportData, err := output.GetData(
		types.Report{
			Path: detectorsReportPath,
		},
		runner.config,
		nil,
	)
	if err != nil {
		t.Fatalf("failed to get output: %s", err)
	}

	report, err := util.ReportYAML(reportData.FindingsBySeverity)
	if err != nil {
		t.Fatalf("failed to encoded to yaml: %s", err)
	}

	cupaloyCopy := cupaloy.NewDefaultConfig().WithOptions(cupaloy.SnapshotSubdirectory(snapshotsPath))
	cupaloyCopy.SnapshotT(t, report)
}
