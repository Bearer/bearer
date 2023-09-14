package filelist_test

import (
	"path/filepath"
	"testing"

	"github.com/bearer/bearer/internal/commands/process/filelist"
	"github.com/bearer/bearer/internal/commands/process/filelist/files"
	"github.com/bearer/bearer/internal/commands/process/settings"
	"github.com/bearer/bearer/internal/flag"
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
					Scan: flag.ScanOptions{
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
					Scan: flag.ScanOptions{
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
					Scan: flag.ScanOptions{
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
