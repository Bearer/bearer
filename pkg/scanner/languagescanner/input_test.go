package languagescanner

import (
	"context"
	"github.com/bearer/bearer/pkg/classification/schema"
	"github.com/bearer/bearer/pkg/languages/javascript"
	"github.com/bearer/bearer/pkg/scanner/ast/query"
	types "github.com/bearer/bearer/pkg/scanner/detectors/types"
	"github.com/bearer/bearer/pkg/scanner/language"
	"github.com/bearer/bearer/pkg/util/file"
	"os"
	"path/filepath"
	"testing"
)

type coverageLanguage struct{ language.Language }

func (coverageLanguage) NewBuiltInDetectors(*schema.Classifier, *query.Set) []types.Detector {
	return nil
}
func TestCoverageF4LanguageScannerBoundInput(t *testing.T) {
	scanner, err := New(coverageLanguage{javascript.Get()}, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer scanner.Close()
	for _, mode := range []string{"valid", "leaf", "component"} {
		t.Run(mode, func(t *testing.T) {
			dir := t.TempDir()
			os.Mkdir(filepath.Join(dir, "sub"), 0700)
			name := filepath.Join(dir, "sub/input.js")
			os.WriteFile(name, []byte("const value = 1;\n"), 0600)
			err := file.IterateFilesListWithExclusions(dir, []string{"sub/input.js"}, false, true, func(*file.Path) (bool, error) { return true, nil }, func(info *file.FileInfo) error {
				if mode != "valid" {
					outside := t.TempDir()
					os.WriteFile(filepath.Join(outside, "input.js"), []byte("const value = 1;\n"), 0600)
					target, dest := name, filepath.Join(outside, "input.js")
					if mode == "component" {
						target, dest = filepath.Dir(name), outside
					}
					os.Rename(target, target+".held")
					os.Symlink(dest, target)
				}
				_, _, e := scanner.Scan(context.Background(), nil, info)
				if mode != "valid" && e == nil {
					t.Fatal("language scanner accepted replacement")
				}
				return e
			}, nil)
			if (err == nil) != (mode == "valid") {
				t.Fatalf("%s: %v", mode, err)
			}
		})
	}
}
