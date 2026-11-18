package parser

import (
	"github.com/bearer/bearer/pkg/util/file"
	"github.com/smacker/go-tree-sitter/python"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCoverageF4ParserBoundInput(t *testing.T) {
	for _, mode := range []string{"valid", "leaf", "component", "context", "closed"} {
		t.Run(mode, func(t *testing.T) {
			dir := t.TempDir()
			os.Mkdir(filepath.Join(dir, "sub"), 0700)
			name := filepath.Join(dir, "sub/input.py")
			os.WriteFile(name, []byte("value = 'inside'\n"), 0600)
			os.WriteFile(filepath.Join(dir, "manage.py"), []byte("context = True\n"), 0600)
			var saved *file.FileInfo
			err := file.IterateFilesListWithExclusions(dir, []string{"sub/input.py"}, false, true, func(p *file.Path) (bool, error) {
				if mode == "context" && p.AbsolutePath == dir {
					tree, e := ParseFile(nil, p.Join("manage.py"), python.GetLanguage())
					if e != nil {
						return false, e
					}
					defer tree.Close()
					if !strings.Contains(tree.RootNode().Content(), "context") {
						t.Fatal("context bytes lost")
					}
				}
				return true, nil
			}, func(info *file.FileInfo) error {
				saved = info
				if mode == "leaf" || mode == "component" {
					outside := t.TempDir()
					os.WriteFile(filepath.Join(outside, "input.py"), []byte("value = 'inside'\n"), 0600)
					target := name
					dest := filepath.Join(outside, "input.py")
					if mode == "component" {
						target = filepath.Dir(name)
						dest = outside
					}
					os.Rename(target, target+".held")
					os.Symlink(dest, target)
				}
				tree, e := ParseFile(info, info.Path, python.GetLanguage())
				if mode == "leaf" || mode == "component" {
					if e == nil {
						tree.Close()
						t.Fatal("unsafe parser accepted replacement")
					}
					return e
				}
				if e != nil {
					return e
				}
				defer tree.Close()
				if tree.File().RelativePath != "sub/input.py" || !strings.Contains(tree.RootNode().Content(), "inside") {
					t.Fatal("logical/source compatibility")
				}
				return nil
			}, nil)
			if (err != nil) != (mode == "leaf" || mode == "component") {
				t.Fatalf("%s: %v", mode, err)
			}
			if mode == "closed" {
				if _, e := ParseFile(saved, saved.Path, python.GetLanguage()); e == nil {
					t.Fatal("parser reopened closed scope")
				}
			}
		})
	}
}
