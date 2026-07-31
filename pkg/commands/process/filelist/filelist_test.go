package filelist_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/bearer/bearer/pkg/commands/process/filelist"
	"github.com/bearer/bearer/pkg/commands/process/filelist/files"
	"github.com/bearer/bearer/pkg/commands/process/settings"
	flagtypes "github.com/bearer/bearer/pkg/flag/types"
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
			Want: &files.List{},
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
			Want: &files.List{},
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
