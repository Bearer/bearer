package filelist_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/bearer/bearer/pkg/commands/process/filelist"
	"github.com/bearer/bearer/pkg/commands/process/filelist/files"
	"github.com/bearer/bearer/pkg/commands/process/settings"
	flagtypes "github.com/bearer/bearer/pkg/flag/types"
	"github.com/bearer/bearer/pkg/report/coverage"
	fileutil "github.com/bearer/bearer/pkg/util/file"
	"github.com/hhatto/gocloc"
	"github.com/stretchr/testify/assert"
)

func TestFileList(t *testing.T) {
	type input struct {
		projectPath string
		config      settings.Config
	}

	type testCase struct {
		Name      string
		Input     input
		Want      *files.List
		WantError bool
	}

	tests := []testCase{
		{
			Name: "Find files - standard - happy path",
			Input: input{
				projectPath: filepath.Join("testdata", "happy_path", "standard"),
				config: settings.Config{
					Worker: settings.WorkerOptions{
						FileSizeMaximum:           100000,
						TimeoutFileBytesPerSecond: 1,
					},
				},
			},
			Want: &files.List{
				Files: []files.File{
					{
						FilePath: "user.go",
						Timeout:  0,
					},
				},
				Coverage: coverage.Discovery{Files: []string{"user.go"}, Excluded: []coverage.Exclusion{}, Errors: []string{}, Complete: true},
			},
		},
		{
			Name: "Find files - standard single file - happy path",
			Input: input{
				projectPath: filepath.Join("testdata", "happy_path", "standard", "user.go"),
				config: settings.Config{
					Worker: settings.WorkerOptions{
						FileSizeMaximum:           100000,
						TimeoutFileBytesPerSecond: 1,
					},
				},
			},
			Want: &files.List{
				Files: []files.File{
					{
						FilePath: ".",
						Timeout:  0,
					},
				},
				Coverage: coverage.Discovery{Files: []string{"."}, Excluded: []coverage.Exclusion{}, Errors: []string{}, Complete: true},
			},
		},
		{
			Name: "Find files - skip - files - happy path",
			Input: input{
				projectPath: filepath.Join("testdata", "happy_path", "skip"),
				config: settings.Config{
					Scan: flagtypes.ScanOptions{
						SkipPath: []string{"users/admin.go"},
					},
					Worker: settings.WorkerOptions{
						FileSizeMaximum:           100000,
						TimeoutFileBytesPerSecond: 1,
					},
				},
			},
			Want: &files.List{
				Files: []files.File{
					{
						Timeout:  0,
						FilePath: "users/users.go",
					},
				},
				Coverage: coverage.Discovery{
					Files:    []string{"users/users.go"},
					Excluded: []coverage.Exclusion{{Path: "users/admin.go", Reason: "skip-path"}},
					Errors:   []string{},
					Complete: true,
				},
			},
		},
		{
			Name: "Find files - skip - dir - happy path",
			Input: input{
				projectPath: filepath.Join("testdata", "happy_path", "skip"),
				config: settings.Config{
					Scan: flagtypes.ScanOptions{
						SkipPath: []string{"users"},
					},
					Worker: settings.WorkerOptions{
						FileSizeMaximum:           100000,
						TimeoutFileBytesPerSecond: 1,
					},
				},
			},
			Want: &files.List{Coverage: coverage.Discovery{
				Files:    []string{},
				Excluded: []coverage.Exclusion{{Path: "users", Reason: "skip-path"}},
				Errors:   []string{},
				Complete: true,
			}},
		},
		{
			Name: "Find files - skip - dir - happy path",
			Input: input{
				projectPath: filepath.Join("testdata", "happy_path", "skip"),
				config: settings.Config{
					Scan: flagtypes.ScanOptions{
						SkipPath: []string{"users"},
					},
					Worker: settings.WorkerOptions{
						FileSizeMaximum:           100000,
						TimeoutFileBytesPerSecond: 1,
					},
				},
			},
			Want: &files.List{Coverage: coverage.Discovery{
				Files:    []string{},
				Excluded: []coverage.Exclusion{{Path: "users", Reason: "skip-path"}},
				Errors:   []string{},
				Complete: true,
			}},
		},
	}

	for _, testCase := range tests {
		dummyGoclocLanguage := gocloc.Language{}
		dummyGoclocResult := gocloc.Result{
			Total: &dummyGoclocLanguage,
			Files: map[string]*gocloc.ClocFile{
				"minifiedJs.min.js": {Code: 2},
				"users.go":          {Code: 6},
				"user.go":           {Code: 0},
				"admin.go":          {Code: 6},
			},
			Languages:     map[string]*gocloc.Language{},
			MaxPathLength: 0,
		}
		t.Run(testCase.Name, func(t *testing.T) {
			output, err := filelist.Discover(nil, testCase.Input.projectPath, &dummyGoclocResult, testCase.Input.config)

			if testCase.WantError {
				if err == nil {
					t.Errorf("expected test case to error")
				}
				return
			}

			if err != nil {
				t.Errorf("classifier returned error %s", err)
			}

			assert.Equal(t, testCase.Want, output)
		})
	}

}

func TestDiscoverLanguageFilter(t *testing.T) {
	dir := t.TempDir()
	for name, content := range map[string]string{
		"a.rb": "puts 1\n",
		"b.py": "print(1)\n",
		"c.go": "package main\n",
	} {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o600); err != nil {
			t.Fatalf("failed to write %s: %s", name, err)
		}
	}

	// gocloc result covering all languages (no --language restriction)
	fullGocloc := &gocloc.Result{
		Files: map[string]*gocloc.ClocFile{
			filepath.Join(dir, "a.rb"): {Lang: "Ruby", Code: 1},
			filepath.Join(dir, "b.py"): {Lang: "Python", Code: 1},
			filepath.Join(dir, "c.go"): {Lang: "Go", Code: 1},
		},
	}

	// gocloc result as it would be after being restricted to ruby upstream: only
	// the ruby file remains
	rubyOnlyGocloc := &gocloc.Result{
		Files: map[string]*gocloc.ClocFile{
			filepath.Join(dir, "a.rb"): {Lang: "Ruby", Code: 1},
		},
	}

	baseConfig := settings.Config{
		Worker: settings.WorkerOptions{
			FileSizeMaximum:           100000,
			TimeoutFileBytesPerSecond: 1,
		},
	}

	restrictedConfig := baseConfig
	restrictedConfig.Scan = flagtypes.ScanOptions{Language: []string{"ruby"}}

	filePaths := func(list *files.List) []string {
		paths := make([]string, 0, len(list.Files))
		for _, f := range list.Files {
			paths = append(paths, f.FilePath)
		}
		return paths
	}

	t.Run("no restriction scans every language", func(t *testing.T) {
		output, err := filelist.Discover(nil, dir, fullGocloc, baseConfig)
		assert.NoError(t, err)
		assert.ElementsMatch(t, []string{"a.rb", "b.py", "c.go"}, filePaths(output))
	})

	t.Run("restriction keeps only files in the pre-filtered gocloc result", func(t *testing.T) {
		output, err := filelist.Discover(nil, dir, rubyOnlyGocloc, restrictedConfig)
		assert.NoError(t, err)
		assert.ElementsMatch(t, []string{"a.rb"}, filePaths(output))
	})

	t.Run("restriction skips files absent from the gocloc result", func(t *testing.T) {
		if err := os.WriteFile(filepath.Join(dir, "notes.txt"), []byte("hello\n"), 0o600); err != nil {
			t.Fatalf("failed to write notes.txt: %s", err)
		}

		output, err := filelist.Discover(nil, dir, rubyOnlyGocloc, restrictedConfig)
		assert.NoError(t, err)
		assert.ElementsMatch(t, []string{"a.rb"}, filePaths(output))
	})
}

func TestCoverageDiscoveryRecordsExclusions(t *testing.T) {
	dir := t.TempDir()
	if err := os.Mkdir(filepath.Join(dir, "skip"), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "keep.go"), []byte("package keep\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "skip", "gone.go"), []byte("package gone\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	result := &gocloc.Result{Files: map[string]*gocloc.ClocFile{
		filepath.Join(dir, "keep.go"):         {Lang: "Go", Code: 1},
		filepath.Join(dir, "skip", "gone.go"): {Lang: "Go", Code: 1},
	}}
	list, err := filelist.Discover(nil, dir, result, settings.Config{
		Scan:   flagtypes.ScanOptions{SkipPath: []string{"skip"}},
		Worker: settings.WorkerOptions{FileSizeMaximum: 100000, TimeoutFileBytesPerSecond: 1},
	})
	assert.NoError(t, err)
	assert.True(t, list.Coverage.Complete)
	assert.Empty(t, list.Coverage.Errors)
	assert.Equal(t, []string{"keep.go"}, list.Coverage.Files)
	assert.Equal(t, []coverage.Exclusion{{Path: "skip", Reason: "skip-path"}}, list.Coverage.Excluded)
}

func TestCoverageDiscoverConfiguredEligibilityMatchesProcessingRules(t *testing.T) {
	dir := t.TempDir()
	files := map[string][]byte{
		"keep.go":          []byte("package keep\n"),
		"tests/fixture.go": []byte("package fixture\n"),
		"fixtures/data.go": []byte("package fixture\n"),
		"binary.dat":       []byte{'b', 0, 'i', 'n', 'a', 'r', 'y'},
		"image.png":        []byte("not really an image"),
		"gitignored.go":    []byte("package ignored\n"),
		".gitignore":       []byte("gitignored.go\n"),
	}
	for name, contents := range files {
		path := filepath.Join(dir, name)
		if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, contents, 0o600); err != nil {
			t.Fatal(err)
		}
	}

	config := settings.Config{
		Scan: flagtypes.ScanOptions{SkipTest: true},
		Worker: settings.WorkerOptions{
			FileSizeMaximum:           100000,
			TimeoutFileBytesPerSecond: 1,
		},
	}
	list, err := filelist.Discover(nil, dir, &gocloc.Result{Files: map[string]*gocloc.ClocFile{}}, config)
	assert.NoError(t, err)
	assert.ElementsMatch(t, []string{".gitignore", "keep.go"}, coverageFiles(list))
	assert.ElementsMatch(t, []coverage.Exclusion{
		{Path: "binary.dat", Reason: "file-type-or-gitignore"},
		{Path: "fixtures", Reason: "built-in-path"},
		{Path: "gitignored.go", Reason: "file-type-or-gitignore"},
		{Path: "image.png", Reason: "file-type-or-gitignore"},
		{Path: "tests", Reason: "test-name"},
	}, list.Coverage.Excluded)

	config.Scan.SkipGitIgnore = true
	list, err = filelist.Discover(nil, dir, &gocloc.Result{Files: map[string]*gocloc.ClocFile{}}, config)
	assert.NoError(t, err)
	assert.Contains(t, coverageFiles(list), "gitignored.go")
}

func coverageFiles(list *files.List) []string {
	return list.Coverage.Files
}

func TestCoverageIterateFilesListConfiguredEligibilityReadErrors(t *testing.T) {
	err := fileutil.IterateFilesListWithExclusions(
		t.TempDir(),
		[]string{"missing.go"},
		false,
		false,
		func(*fileutil.Path) (bool, error) { return true, nil },
		func(*fileutil.FileInfo) error { return nil },
		nil,
	)
	assert.Error(t, err)
}

func TestCoverageEligibilityKeepsExternalPromise(t *testing.T) {
	dir := t.TempDir()
	keep := []byte("package keep\n")
	excluded := []byte("binary\x00payload")
	for name, b := range map[string][]byte{"keep.go": keep, "binary.dat": excluded} {
		if err := os.WriteFile(filepath.Join(dir, name), b, 0600); err != nil {
			t.Fatal(err)
		}
	}
	promise := map[string]string{"keep.go": coverage.Digest(keep), "binary.dat": coverage.Digest(excluded)}
	list, err := filelist.Discover(nil, dir, &gocloc.Result{Files: map[string]*gocloc.ClocFile{}}, settings.Config{Scan: flagtypes.ScanOptions{SkipGitIgnore: true}, Worker: settings.WorkerOptions{TimeoutFileBytesPerSecond: 1}})
	if err != nil {
		t.Fatal(err)
	}
	if len(list.Files) != 1 || list.Files[0].FilePath != "keep.go" || len(list.Coverage.Excluded) != 1 {
		t.Fatalf("schedule/exclusions=%+v", list)
	}
	visited := []string{}
	if err := fileutil.IterateFilesListWithExclusions(dir, []string{list.Files[0].FilePath}, false, true, func(*fileutil.Path) (bool, error) { return true, nil }, func(f *fileutil.FileInfo) error { visited = append(visited, f.RelativePath); return nil }, func(string) { t.Fatal("scheduled input excluded during processing") }); err != nil {
		t.Fatal(err)
	}
	if len(visited) != 1 || visited[0] != "keep.go" {
		t.Fatalf("visits=%v", visited)
	}
	record := coverage.Record{Schema: coverage.Schema, RunID: coverage.Digest([]byte("run")), ConfigurationSHA256: coverage.Digest(nil), Discovery: list.Coverage, Outcomes: []coverage.Outcome{{Path: "keep.go", SHA256: coverage.Digest(keep), Status: "success"}}, Complete: true}
	if err := record.ValidateInputs(promise); err == nil {
		t.Fatal("ordinary compatible output was mistaken for coverage of external promise")
	}
	if len(promise) != 2 || promise["binary.dat"] != coverage.Digest(excluded) {
		t.Fatal("external promise changed")
	}
	record.Outcomes = nil
	if err := record.Validate(list.Coverage.Files); err == nil {
		t.Fatal("unvisited scheduled input accepted")
	}
}

func TestCoverageEligibilityInspectionFailureAndUnsafeInput(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "input.go"), []byte("package fixture\n"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(filepath.Join(dir, ".gitignore"), 0700); err != nil {
		t.Fatal(err)
	}
	if _, err := filelist.Discover(nil, dir, &gocloc.Result{Files: map[string]*gocloc.ClocFile{}}, settings.Config{}); err == nil {
		t.Fatal("unreadable ignore configuration became a successful discovery")
	}
	if _, err := fileutil.ConfiguredExclusion(dir, "../outside.go", false, true); err == nil {
		t.Fatal("outside-target eligibility accepted")
	}
	if _, err := fileutil.ConfiguredExclusion(dir, filepath.Join(dir, "input.go"), false, true); err == nil {
		t.Fatal("absolute eligibility input accepted")
	}
}

func TestCoverageF4DiscoveryCompatibility(t *testing.T) {
	dir := t.TempDir()
	name := filepath.Join(dir, "input.go")
	if err := os.WriteFile(name, []byte("package synthetic\n"), 0600); err != nil {
		t.Fatal(err)
	}
	for _, root := range []string{dir, name} {
		relative := "input.go"
		if root == name {
			relative = "."
		}
		reason, err := fileutil.ConfiguredExclusion(root, relative, false, false)
		if err != nil || reason != "" {
			t.Fatalf("valid %s: %q %v", root, reason, err)
		}
	}
	os.WriteFile(filepath.Join(dir, ".gitignore"), []byte("input.go\n"), 0600)
	if reason, err := fileutil.ConfiguredExclusion(name, ".", false, false); err != nil || reason == "" {
		t.Fatalf("single-file gitignore %q %v", reason, err)
	}
	if reason, err := fileutil.ConfiguredExclusion(name, ".", false, true); err != nil || reason != "" {
		t.Fatalf("disabled control %q %v", reason, err)
	}
	outside := filepath.Join(t.TempDir(), "outside.go")
	os.WriteFile(outside, []byte("package synthetic\n"), 0600)
	os.Remove(name)
	os.Symlink(outside, name)
	if _, err := fileutil.ConfiguredExclusion(dir, "input.go", false, true); err == nil {
		t.Fatal("discovery accepted replaced input")
	}
}
