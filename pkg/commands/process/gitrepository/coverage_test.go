package gitrepository

import (
	"context"
	"fmt"
	"github.com/bearer/bearer/pkg/commands/process/filelist/ignore"
	"github.com/bearer/bearer/pkg/commands/process/settings"
	flagtypes "github.com/bearer/bearer/pkg/flag/types"
	"github.com/bearer/bearer/pkg/report/coverage"
	"github.com/hhatto/gocloc"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestCoverageGitDiscoveryIsExplicitlyIncomplete(t *testing.T) {
	discovery := unsupportedCoverage()
	if discovery.Complete {
		t.Fatal("git discovery coverage must not claim complete")
	}
	if len(discovery.Errors) != 1 {
		t.Fatalf("errors = %d, want one explicit unsupported-coverage error", len(discovery.Errors))
	}
}

func TestCoverageCurrentFileAccounting(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "input.go"), []byte("synthetic"), 0600); err != nil {
		t.Fatal(err)
	}
	r := &Repository{targetPath: dir, gitTargetPath: ".", config: settings.Config{Worker: settings.WorkerOptions{TimeoutFileBytesPerSecond: 100000, TimeoutFileMinimum: time.Second, TimeoutFileMaximum: 5 * time.Second}}}
	checker := ignore.New(dir, r.config)
	loc := &gocloc.Result{Files: map[string]*gocloc.ClocFile{}}
	f, reason, err := r.fileForWithCoverage(checker, loc, "input.go")
	if err != nil || reason != "" || f == nil || f.FilePath != "input.go" {
		t.Fatalf("file=%v reason=%q err=%v", f, reason, err)
	}
	if _, _, err = r.fileForWithCoverage(checker, loc, "absent.go"); err == nil {
		t.Fatal("missing input silently excluded")
	}
	if err = os.WriteFile(filepath.Join(dir, "skip.min.js"), []byte("synthetic"), 0600); err != nil {
		t.Fatal(err)
	}
	f, reason, err = r.fileForWithCoverage(checker, loc, "skip.min.js")
	if err != nil || f != nil || reason != "minified-js" {
		t.Fatalf("exclusion=%v %q %v", f, reason, err)
	}
}

func TestCoverageRepositoryListFilesDiffs(t *testing.T) {
	repository, _ := newCoverageRepository(t)

	list, err := repository.ListFiles(nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	var paths []string
	for _, file := range list.Files {
		paths = append(paths, file.FilePath)
	}
	if want := []string{"add.go", "modify.go", "rename_new.go"}; !reflect.DeepEqual(paths, want) {
		t.Fatalf("files = %v, want %v", paths, want)
	}
	if want := []string{"modify.go", "rename_old.go"}; len(list.BaseFiles) != len(want) || list.BaseFiles[0].FilePath != want[0] || list.BaseFiles[1].FilePath != want[1] {
		t.Fatalf("base files = %#v, want paths %v", list.BaseFiles, want)
	}
	if want := map[string]string{"rename_old.go": "rename_new.go"}; !reflect.DeepEqual(list.Renames, want) {
		t.Fatalf("renames = %v, want %v", list.Renames, want)
	}
	if list.Coverage.Complete {
		t.Fatal("diff coverage unexpectedly claims complete")
	}
}

func TestCoverageWithBaseBranchRestoresHead(t *testing.T) {
	tests := []struct {
		name    string
		bodyErr error
	}{
		{name: "success"},
		{name: "callback error", bodyErr: fmt.Errorf("callback failed")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository, head := newCoverageRepository(t)
			called := false
			err := repository.WithBaseBranch(func() error {
				called = true
				if got := gitOutput(t, repository.context.RootDir, "rev-parse", "HEAD"); got != repository.context.BaseCommitHash {
					t.Fatalf("callback HEAD = %s, want base %s", got, repository.context.BaseCommitHash)
				}
				return tt.bodyErr
			})
			if !called {
				t.Fatal("callback was not called")
			}
			if !sameError(err, tt.bodyErr) {
				t.Fatalf("error = %v, want %v", err, tt.bodyErr)
			}
			if got := gitOutput(t, repository.context.RootDir, "rev-parse", "HEAD"); got != head {
				t.Fatalf("restored HEAD = %s, want %s", got, head)
			}
			if got := gitOutput(t, repository.context.RootDir, "branch", "--show-current"); got != "main" {
				t.Fatalf("restored branch = %s, want main", got)
			}
		})
	}
}

func TestCoverageRepositoryListFilesInvalidLocalObject(t *testing.T) {
	dir := t.TempDir()
	gitRun(t, dir, nil, "init", "--quiet")
	repository := &Repository{
		ctx:           context.Background(),
		targetPath:    dir,
		gitTargetPath: ".",
		context:       &Context{RootDir: dir, BaseCommitHash: "deadbeefdeadbeefdeadbeefdeadbeefdeadbeef"},
	}

	if _, err := repository.ListFiles(nil, nil); err == nil {
		t.Fatal("invalid local base object unexpectedly succeeded")
	} else if !strings.Contains(err.Error(), "fatal") && !strings.Contains(err.Error(), "invalid") && !strings.Contains(err.Error(), "bad object") {
		t.Fatalf("invalid-object error = %v", err)
	}
}

func newCoverageRepository(t *testing.T) (*Repository, string) {
	t.Helper()
	dir := t.TempDir()
	gitRun(t, dir, nil, "init", "--quiet")
	gitRun(t, dir, nil, "config", "user.name", "coverage-test")
	gitRun(t, dir, nil, "config", "user.email", "coverage-test@example.invalid")

	base := map[string]string{
		"delete.go":     "package old\n",
		"modify.go":     "package old\nvar value = 1\n",
		"rename_old.go": "package old\nvar renamed = true\n",
	}
	head := map[string]string{
		"add.go":        "package add\n",
		"modify.go":     "package new\nvar value = 2\n",
		"rename_new.go": "package old\nvar renamed = true\n",
	}
	baseTree := writeCoverageTree(t, dir, base)
	baseCommit := gitRun(t, dir, nil, "commit-tree", baseTree, "-m", "base")
	headTree := writeCoverageTree(t, dir, head)
	headCommit := gitRun(t, dir, nil, "commit-tree", headTree, "-p", baseCommit, "-m", "head")
	gitRun(t, dir, nil, "update-ref", "refs/heads/main", headCommit)
	gitRun(t, dir, nil, "symbolic-ref", "HEAD", "refs/heads/main")
	gitRun(t, dir, nil, "read-tree", "--reset", "-u", headCommit)

	return &Repository{
		ctx:           context.Background(),
		config:        settings.Config{Scan: flagtypes.ScanOptions{Diff: true}, Worker: settings.WorkerOptions{TimeoutFileBytesPerSecond: 100000, TimeoutFileMinimum: time.Second, TimeoutFileMaximum: 5 * time.Second}},
		targetPath:    dir,
		gitTargetPath: ".",
		context: &Context{
			RootDir:           dir,
			CurrentBranch:     "main",
			CurrentCommitHash: headCommit,
			BaseCommitHash:    baseCommit,
		},
	}, headCommit
}

func writeCoverageTree(t *testing.T, dir string, files map[string]string) string {
	t.Helper()
	index := filepath.Join(t.TempDir(), "index")
	for path, contents := range files {
		blob := gitRunInput(t, dir, nil, contents, "hash-object", "-w", "--stdin")
		gitRun(t, dir, map[string]string{"GIT_INDEX_FILE": index}, "update-index", "--add", "--cacheinfo", "100644,"+blob+","+path)
	}
	return gitRun(t, dir, map[string]string{"GIT_INDEX_FILE": index}, "write-tree")
}

func gitRun(t *testing.T, dir string, extraEnv map[string]string, args ...string) string {
	t.Helper()
	return gitRunInput(t, dir, extraEnv, "", args...)
}

func gitRunInput(t *testing.T, dir string, extraEnv map[string]string, input string, args ...string) string {
	t.Helper()
	command := exec.Command("git", args...)
	command.Dir = dir
	command.Stdin = strings.NewReader(input)
	command.Env = os.Environ()
	for key, value := range extraEnv {
		command.Env = append(command.Env, key+"="+value)
	}
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("git %s: %v\n%s", strings.Join(args, " "), err, output)
	}
	return strings.TrimSpace(string(output))
}

func gitOutput(t *testing.T, dir string, args ...string) string {
	t.Helper()
	return gitRun(t, dir, nil, args...)
}

func sameError(got, want error) bool {
	if got == nil || want == nil {
		return got == want
	}
	return got.Error() == want.Error()
}

func TestCoverageCurrentTreeAcquisitionFailure(t *testing.T) {
	repository, _ := newCoverageRepository(t)
	repository.context.BaseCommitHash = ""
	repository.context.CurrentCommitHash = strings.Repeat("e", 40)
	if list, err := repository.ListFiles(nil, nil); err == nil {
		t.Fatalf("invalid current tree acquired: %+v", list)
	}
}
func TestCoverageSubdirectoryNamespaces(t *testing.T) {
	repository, _ := newCoverageRepository(t)
	dir := repository.context.RootDir
	tree := writeCoverageTree(t, dir, map[string]string{"a.go": "package root\n", "src/a.go": "package nested\n", "src/foo..go": "package nested\n"})
	head := gitRun(t, dir, nil, "commit-tree", tree, "-m", "synthetic target namespace")
	gitRun(t, dir, nil, "read-tree", "--reset", "-u", head)
	repository.context.CurrentCommitHash = head
	repository.context.BaseCommitHash = ""
	repository.targetPath = filepath.Join(dir, "src")
	repository.gitTargetPath = "src"
	list, err := repository.ListFiles(nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(list.Coverage.Files, []string{"a.go", "foo..go"}) {
		t.Fatalf("paths=%v", list.Coverage.Files)
	}
	if len(list.Coverage.Excluded) != 1 || list.Coverage.Excluded[0].Scope != "repository" || list.Coverage.Excluded[0].Path != "a.go" {
		t.Fatalf("exclusions=%+v", list.Coverage.Excluded)
	}
	record := coverage.Record{Schema: coverage.Schema, RunID: coverage.Digest(nil), ConfigurationSHA256: coverage.Digest(nil), Discovery: list.Coverage, Complete: true, Outcomes: []coverage.Outcome{}}
	for _, name := range list.Coverage.Files {
		record.Outcomes = append(record.Outcomes, coverage.Outcome{Path: name, SHA256: coverage.Digest(nil), Status: "success"})
	}
	if err := record.Validate(list.Coverage.Files); err != nil {
		t.Fatal(err)
	}
	if err := record.Validate([]string{"../a.go"}); err == nil {
		t.Fatal("outside-target promise accepted")
	}
}
