package testhelper

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bearer/bearer/pkg/commands"
	"github.com/bearer/bearer/pkg/commands/process/balancer/filelist"
	"github.com/bearer/bearer/pkg/commands/process/settings"
	"github.com/bearer/bearer/pkg/commands/process/settings/rules"
	"github.com/bearer/bearer/pkg/commands/process/worker"
	"github.com/bearer/bearer/pkg/commands/process/worker/work"
	"github.com/bearer/bearer/pkg/flag"
	"github.com/bearer/bearer/pkg/report/output"
	reportoutput "github.com/bearer/bearer/pkg/report/output"
	"github.com/bearer/bearer/pkg/types"
	"github.com/bradleyjkemp/cupaloy"
	"gopkg.in/yaml.v3"
)

type Runner struct {
	config settings.Config
	worker worker.Worker
}

func GetRunner(t *testing.T, ruleBytes []byte, lang string) *Runner {
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

	config, err := settings.FromOptions(configFlags, []string{lang})
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

func getRulesFromYaml(t *testing.T, ruleBytes []byte) map[string]*rules.Rule {
	var ruleDefinition rules.RuleDefinition
	err := yaml.Unmarshal(ruleBytes, &ruleDefinition)
	if err != nil {
		t.Fatalf("failed to unmarshal rule %s", err)
	}

	definitions := map[string]rules.RuleDefinition{
		ruleDefinition.Metadata.ID: ruleDefinition,
	}
	enabledRules := map[string]struct{}{
		ruleDefinition.Metadata.ID: {},
	}

	rules, err := settings.BuildRules(definitions, enabledRules)
	if err != nil {
		t.Fatalf("failed to build rules %s", err)
	}

	return rules
}

func (runner *Runner) RunTest(t *testing.T, testdataPath string, snapshotPath string) {
	files, err := filelist.Discover(testdataPath, runner.config)
	if err != nil {
		t.Fatalf("failed to discover files: %s", err)
	}

	if len(files) == 0 {
		t.Fatal("no scannable files found")
	}

	for _, file := range files {
		myfile := file
		ext := filepath.Ext(file.FilePath)
		testName := strings.TrimSuffix(file.FilePath, ext) + ".yml"
		t.Run(testName, func(t *testing.T) {
			runner.scanSingleFile(t, testdataPath, myfile, snapshotPath)
		})
	}
}

func (runner *Runner) scanSingleFile(t *testing.T, testDataPath string, fileRelativePath work.File, snapshotsPath string) {
	detectorsReportFile, err := os.CreateTemp("", "report.jsonl")
	if err != nil {
		t.Fatalf("failed to create tmp report file: %s", err)
	}
	defer detectorsReportFile.Close()

	detectorsReportPath := detectorsReportFile.Name()
	if err != nil {
		t.Fatalf("failed to get absolute path of report file: %s", err)
	}

	err = runner.worker.Scan(work.ProcessRequest{
		Files:      []work.File{fileRelativePath},
		ReportPath: detectorsReportPath,
		Repository: work.Repository{
			Dir: testDataPath,
		},
	})

	if err != nil {
		t.Fatalf("failed to do scan %s", err)
	}

	runner.config.Scan.Target = testDataPath
	detections, _, _ := output.GetOutput(
		types.Report{
			Path: detectorsReportPath,
		},
		runner.config,
	)
	report, _ := reportoutput.ReportYAML(detections)

	cupaloyCopy := cupaloy.NewDefaultConfig().WithOptions(cupaloy.SnapshotSubdirectory(snapshotsPath))
	cupaloyCopy.SnapshotT(t, *report)
}
