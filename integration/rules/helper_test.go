package integration_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/bearer/bearer/pkg/commands"
	"github.com/bearer/bearer/pkg/commands/process/balancer/filelist"
	"github.com/bearer/bearer/pkg/commands/process/settings"
	"github.com/bearer/bearer/pkg/commands/process/worker"
	"github.com/bearer/bearer/pkg/commands/process/worker/work"
	"github.com/bearer/bearer/pkg/flag"
	reportoutput "github.com/bearer/bearer/pkg/report/output"
	"github.com/bearer/bearer/pkg/types"
	"github.com/bearer/bearer/pkg/util/output"
	"github.com/bradleyjkemp/cupaloy"
)

var runner *Runner
var mu sync.Mutex

type Runner struct {
	config settings.Config
	worker worker.Worker
}

func getRunner(t *testing.T) *Runner {
	mu.Lock()
	defer mu.Unlock()
	if runner != nil {
		return runner
	}

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
	configFlags.OnlyRule = map[string]bool{
		"javascript_lang_file_generation": true,
	}

	config, err := settings.FromOptions(configFlags)
	if err != nil {
		t.Fatalf("failed to generate default scan settings: %s", err)
	}

	worker := worker.Worker{}
	err = worker.Setup(config)
	if err != nil {
		t.Fatalf("failed to setup scan worker: %s", err)
	}

	runner = &Runner{
		worker: worker,
		config: config,
	}

	return runner
}

func (runner *Runner) runTest(t *testing.T, projectPath string) {
	testDataPath := projectPath + "/testdata/"
	snapshotsPath := projectPath + "/.snapshots/"

	files, err := filelist.Discover(testDataPath, runner.config)
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
			runner.ScanSingleFile(t, testDataPath, myfile, snapshotsPath)
		})
	}
}

func (runner *Runner) ScanSingleFile(t *testing.T, testDataPath string, fileRelativePath work.File, snapshotsPath string) {
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

	outputBuffer := bytes.NewBuffer(nil)
	logger := output.PlainLogger(outputBuffer)

	err = reportoutput.ReportYAML(types.Report{
		Path: detectorsReportPath,
	}, logger, runner.config)
	if err != nil {
		t.Fatalf("failed to generate report yaml: %s", err)
	}

	cupaloyCopy := cupaloy.NewDefaultConfig().WithOptions(cupaloy.SnapshotSubdirectory(snapshotsPath))
	cupaloyCopy.SnapshotT(t, outputBuffer.String())
}
