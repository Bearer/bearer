package integration_test

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bearer/curio/integration/internal/testhelper"
	"github.com/bearer/curio/pkg/commands"
	"github.com/bearer/curio/pkg/commands/process/balancer/filelist"
	"github.com/bearer/curio/pkg/commands/process/settings"
	"github.com/bearer/curio/pkg/commands/process/settings/rules"
	"github.com/bearer/curio/pkg/commands/process/worker"
	"github.com/bearer/curio/pkg/commands/process/worker/work"
	"github.com/bearer/curio/pkg/flag"
	reportoutput "github.com/bearer/curio/pkg/report/output"
	"github.com/bearer/curio/pkg/types"
	"github.com/bearer/curio/pkg/util/output"
	"github.com/bradleyjkemp/cupaloy"

	"github.com/rs/zerolog/log"
)

var rulesFs = &rules.Rules

func TestTest(t *testing.T) {
	getRunner(t).runTest(t, "../../pkg/commands/process/settings/rules/"+"javascript/lang/logger")
}

var runner *Runner

type Runner struct {
	config settings.Config
	worker worker.Worker
}

func getRunner(t *testing.T) *Runner {
	if runner != nil {
		return runner
	}

	commands.ScanFlags.BindForConfigInit(commands.NewScanCommand())
	configFlags, err := commands.ScanFlags.ToOptions([]string{})
	if err != nil {
		t.Fatalf("failed to generate default flags: %s", err)
	}
	configFlags.Format = flag.FormatYAML
	configFlags.Report = flag.ReportSummary
	configFlags.Quiet = true

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
		t.Fatalf("failed to discover files: %e", err)
	}

	if len(files) == 0 {
		t.Fatal("no scannable files found")
	}

	for _, file := range files {
		ext := filepath.Ext(file.FilePath)
		testName := strings.TrimSuffix(file.FilePath, ext) + ".yml"

		t.Run(testName, func(t *testing.T) {
			runner.ScanSingleFile(t, testDataPath, file, snapshotsPath)
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

	runner.worker.Scan(work.ProcessRequest{
		Files:      []work.File{fileRelativePath},
		ReportPath: detectorsReportPath,
		Repository: work.Repository{
			Dir: testDataPath,
		},
	})

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

func buildRulesTestCase(name, reportType, ruleID, filename string) testhelper.TestCase {
	arguments := []string{
		"scan",
		filepath.Join("pkg", "commands", "process", "settings", "rules", filename),
		"--report=" + reportType,
		"--format=yaml",
		"--only-rule=" + ruleID,
	}
	options := testhelper.TestCaseOptions{}

	return testhelper.NewTestCase(name, arguments, options)
}

func runRulesTest(folderPath, format, ruleID string, t *testing.T) {
	snapshotDirectory := "../../pkg/commands/process/settings/rules/" + folderPath + "/.snapshots"

	testDataDir := fmt.Sprintf("%s/testdata", folderPath)

	log.Debug().Msgf("%s", testDataDir)

	testdataDirEntries, err := rulesFs.ReadDir(testDataDir)
	if err != nil {
		t.Fatalf("failed to read rules/%s dir %e", folderPath, err)
	}

	dataflowTests := []testhelper.TestCase{}
	for _, testdataFile := range testdataDirEntries {
		name := testdataFile.Name()

		testName := strings.Replace(fmt.Sprintf("%s_%s_%s", format, folderPath, name), "/", "_", -1)
		dataflowTests = append(dataflowTests,
			buildRulesTestCase(
				testName,
				format,
				ruleID,
				fmt.Sprintf("%s/testdata/%s", folderPath, name),
			),
		)
	}

	testhelper.RunTestsWithSnapshotSubdirectory(t, dataflowTests, snapshotDirectory)
}
