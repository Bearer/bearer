package file

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func inputFixture(t *testing.T) (string, *InputScope, *BoundInput) {
	t.Helper()
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "input.go"), []byte("package synthetic\n"), 0600); err != nil {
		t.Fatal(err)
	}
	s, err := OpenInputScope(dir)
	if err != nil {
		t.Fatal(err)
	}
	b, err := s.Acquire("input.go", nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { s.Close() })
	return dir, s, b
}
func replaceInput(t *testing.T, path, other string) {
	t.Helper()
	if err := os.Rename(path, path+".held"); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(other, path); err != nil {
		t.Fatal(err)
	}
}
func TestCoverageF4AcquireNoFollow(t *testing.T) {
	for _, kind := range []string{"leaf", "component", "root", "fifo", "escape"} {
		t.Run(kind, func(t *testing.T) {
			dir, s, _ := inputFixture(t)
			outside := t.TempDir()
			os.WriteFile(filepath.Join(outside, "input.go"), []byte("package synthetic\n"), 0600)
			relative := "input.go"
			switch kind {
			case "leaf":
				replaceInput(t, filepath.Join(dir, relative), filepath.Join(outside, relative))
			case "component":
				os.Mkdir(filepath.Join(dir, "sub"), 0700)
				replaceInput(t, filepath.Join(dir, "sub"), outside)
				relative = "sub/input.go"
			case "root":
				replaceInput(t, dir, outside)
			case "fifo":
				relative = "dir"
				os.Mkdir(filepath.Join(dir, relative), 0700)
			case "escape":
				relative = "../input.go"
			}
			if b, err := s.Acquire(relative, nil); err == nil {
				b.Close()
				t.Fatal("unsafe acquisition accepted")
			}
		})
	}
}
func TestCoverageF4ReplacementDuringAcquire(t *testing.T) {
	dir, s, b := inputFixture(t)
	expect := b.Expectation()
	outside := t.TempDir()
	os.WriteFile(filepath.Join(outside, "input.go"), []byte("package synthetic\n"), 0600)
	replaceInput(t, filepath.Join(dir, "input.go"), filepath.Join(outside, "input.go"))
	if _, err := s.Acquire("input.go", expect); err == nil {
		t.Fatal("replacement satisfied held expectation")
	}
}
func TestCoverageF4BoundReadAfterReplacement(t *testing.T) {
	for _, same := range []bool{false, true} {
		t.Run(map[bool]string{false: "different", true: "identical"}[same], func(t *testing.T) {
			dir, _, b := inputFixture(t)
			reader, err := b.Reader()
			if err != nil {
				t.Fatal(err)
			}
			outside := filepath.Join(t.TempDir(), "outside.go")
			text := []byte("outside sentinel\n")
			if same {
				text = []byte("package synthetic\n")
			}
			os.WriteFile(outside, text, 0600)
			replaceInput(t, filepath.Join(dir, "input.go"), outside)
			data, err := io.ReadAll(reader)
			if err == nil || len(data) != 0 {
				t.Fatalf("replacement read %q %v", data, err)
			}
			if _, err = b.ReadBytes(); err == nil {
				t.Fatal("second consumer accepted replacement")
			}
			if b.Close() == nil {
				t.Fatal("terminal replacement accepted")
			}
		})
	}
}
func TestCoverageF4FrozenBytesAndDigest(t *testing.T) {
	_, _, b := inputFixture(t)
	data, err := b.ReadBytes()
	if err != nil {
		t.Fatal(err)
	}
	data[0] = 'X'
	again, _ := b.ReadBytes()
	if again[0] == 'X' {
		t.Fatal("shared mutable bytes")
	}
	a, _ := b.Reader()
	c, _ := b.Reader()
	one := make([]byte, 1)
	a.Read(one)
	all, _ := io.ReadAll(c)
	if !bytes.Equal(all, again) {
		t.Fatal("shared reader cursor")
	}
	if len(b.Digest()) != 64 {
		t.Fatal("missing digest")
	}
	expectation := b.Expectation()
	expectation.SHA256 = string(bytes.Repeat([]byte("0"), 64))
	if _, err = b.scope.Acquire("input.go", expectation); err == nil {
		t.Fatal("digest mismatch accepted")
	}
}
func TestCoverageF4AcquisitionAndCloseFailures(t *testing.T) {
	t.Run("missing", func(t *testing.T) {
		_, s, _ := inputFixture(t)
		if _, err := s.Acquire("missing", nil); !errors.Is(err, os.ErrNotExist) {
			t.Fatalf("missing: %v", err)
		}
		if s.Err() == nil {
			t.Fatal("failure not latched")
		}
	})
	t.Run("read-and-close", func(t *testing.T) {
		_, s, b := inputFixture(t)
		b.handles[len(b.handles)-1].file.Close()
		if _, err := b.ReadBytes(); err == nil {
			t.Fatal("closed descriptor read accepted")
		}
		if b.Close() == nil || s.Close() == nil {
			t.Fatal("close failure lost")
		}
	})
	t.Run("unreadable", func(t *testing.T) {
		dir, s, _ := inputFixture(t)
		os.WriteFile(filepath.Join(dir, "denied"), []byte("no"), 0000)
		if _, err := s.Acquire("denied", nil); err == nil {
			t.Fatal("unreadable fixture precondition failed")
		}
	})
}
func TestCoverageF4SingleFileIdentity(t *testing.T) {
	dir := t.TempDir()
	name := filepath.Join(dir, "input.go")
	os.WriteFile(name, []byte("package synthetic\n"), 0600)
	s, err := OpenInputScope(name)
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	b, err := s.Acquire(".", nil)
	if err != nil {
		t.Fatal(err)
	}
	if b.Expectation().Relative != "." {
		t.Fatal("logical identity changed")
	}
	if _, err = s.Acquire("input.go", nil); err == nil {
		t.Fatal("single-file selector broadened")
	}
}
func TestCoverageF4ContextAndLifetime(t *testing.T) {
	t.Run("logical-parent-alias", func(t *testing.T) {
		outer := t.TempDir()
		real := filepath.Join(outer, "real")
		os.Mkdir(real, 0700)
		os.Mkdir(filepath.Join(real, "project"), 0700)
		os.Symlink(real, filepath.Join(outer, "alias"))
		logical := filepath.Join(outer, "alias/project")
		os.WriteFile(filepath.Join(logical, "input.go"), []byte("package synthetic\n"), 0600)
		os.WriteFile(filepath.Join(logical, "manage.py"), []byte("context = True\n"), 0600)
		err := IterateFilesListWithExclusions(logical, []string{"input.go"}, false, true, func(p *Path) (bool, error) { _, err := p.Join("manage.py").ReadBytes(); return true, err }, func(f *FileInfo) error {
			if f.AbsolutePath != filepath.Join(logical, "input.go") {
				t.Fatalf("logical path rewritten: %s", f.AbsolutePath)
			}
			_, err := f.ReadBytes()
			return err
		}, nil)
		if err != nil {
			t.Fatal(err)
		}
	})

	dir, s, _ := inputFixture(t)
	os.WriteFile(filepath.Join(dir, "manage.py"), []byte("print('synthetic')\n"), 0600)
	p := &Path{AbsolutePath: dir, scope: s}
	data, err := p.Join("manage.py").ReadBytes()
	if err != nil || len(data) == 0 {
		t.Fatalf("context %q %v", data, err)
	}
	if err = s.Close(); err != nil {
		t.Fatal(err)
	}
	if _, err = p.Join("manage.py").ReadBytes(); err == nil {
		t.Fatal("closed scope reopened pathname")
	}
}

func TestCoverageF4ResidualCompletedReconciliation(t *testing.T) {
	for _, mode := range []string{"valid", "leaf-identical", "component-identical", "content", "missing", "single"} {
		t.Run(mode, func(t *testing.T) {
			root := t.TempDir()
			outside := t.TempDir()
			relative := "sub/a.go"
			body := []byte("package synthetic\n")
			for _, dir := range []string{root, outside} {
				if err := os.MkdirAll(filepath.Join(dir, "sub"), 0700); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(filepath.Join(dir, relative), body, 0600); err != nil {
					t.Fatal(err)
				}
			}
			target := root
			selector := relative
			if mode == "single" {
				target = filepath.Join(root, relative)
				selector = "."
			}
			s, err := OpenInputScope(target)
			if err != nil {
				t.Fatal(err)
			}
			defer s.Close()
			b, err := s.Acquire(selector, nil)
			if err != nil {
				t.Fatal(err)
			}
			if err = b.Close(); err != nil {
				t.Fatal(err)
			}
			for _, h := range b.handles {
				if _, e := h.file.Stat(); !errors.Is(e, os.ErrClosed) {
					t.Fatalf("descriptor not released: %v", e)
				}
			}
			switch mode {
			case "leaf-identical", "single":
				replaceInput(t, filepath.Join(root, relative), filepath.Join(outside, relative))
			case "component-identical":
				replaceInput(t, filepath.Join(root, "sub"), filepath.Join(outside, "sub"))
			case "content":
				if err = os.WriteFile(filepath.Join(root, relative), []byte("package changed\n"), 0600); err != nil {
					t.Fatal(err)
				}
			case "missing":
				if err = os.Remove(filepath.Join(root, relative)); err != nil {
					t.Fatal(err)
				}
			}
			err = s.Close()
			t.Logf("mode=%s final=%v checks=%d bytes=%d", mode, err, b.validationReads, b.validationBytes)
			if (err == nil) != (mode == "valid") {
				t.Fatalf("final reconciliation=%v", err)
			}
		})
	}
}

func TestCoverageF4ResidualReaderLifetime(t *testing.T) {
	for _, method := range []string{"Read", "ReadAt", "Seek", "WriteTo", "ReadByte", "ReadRune", "UnreadByte", "UnreadRune"} {
		t.Run(method, func(t *testing.T) {
			_, s, b := inputFixture(t)
			r, err := b.Reader()
			if err != nil {
				t.Fatal(err)
			}
			if err = r.Close(); err != nil {
				t.Fatal(err)
			}
			switch method {
			case "Read":
				_, err = r.Read(make([]byte, 1))
			case "ReadAt":
				_, err = r.ReadAt(make([]byte, 1), 0)
			case "Seek":
				_, err = r.Seek(0, io.SeekStart)
			case "WriteTo":
				var dst bytes.Buffer
				_, err = io.Copy(&dst, r)
				if dst.Len() != 0 {
					t.Fatal("closed reader emitted bytes")
				}
			case "ReadByte":
				_, err = r.ReadByte()
			case "ReadRune":
				_, _, err = r.ReadRune()
			case "UnreadByte":
				err = r.UnreadByte()
			case "UnreadRune":
				err = r.UnreadRune()
			}
			if err == nil || s.Close() == nil {
				t.Fatalf("closed %s did not fail terminal", method)
			}
		})
	}
	t.Run("writer-fault", func(t *testing.T) {
		_, s, b := inputFixture(t)
		r, err := b.Reader()
		if err != nil {
			t.Fatal(err)
		}
		if _, err = r.WriteTo(residualBadWriter{}); err == nil || s.Close() == nil {
			t.Fatal("writer fault lost")
		}
	})
}

type residualBadWriter struct{}

func (residualBadWriter) Write([]byte) (int, error) { return 0, errors.New("synthetic output failure") }

func TestCoverageF4ResidualBoundedValidation(t *testing.T) {
	root := t.TempDir()
	body := bytes.Repeat([]byte("a"), 2*1000*1000)
	target := filepath.Join(root, "large.txt")
	if err := os.WriteFile(target, body, 0600); err != nil {
		t.Fatal(err)
	}
	s, err := OpenInputScope(root)
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	b, err := s.Acquire("large.txt", nil)
	if err != nil {
		t.Fatal(err)
	}
	r, err := b.Reader()
	if err != nil {
		t.Fatal(err)
	}
	before := b.validationReads
	read := 0
	buf := make([]byte, 4096)
	for {
		n, e := r.Read(buf)
		read += n
		if e == io.EOF {
			break
		}
		if e != nil {
			t.Fatal(e)
		}
	}
	if read != len(body) || b.validationReads != before {
		t.Fatalf("stream read=%d whole-file checks=%d before=%d", read, b.validationReads, before)
	}
	if err = r.Close(); err != nil {
		t.Fatal(err)
	}
	if err = b.Close(); err != nil {
		t.Fatal(err)
	}
	if err = s.Close(); err != nil {
		t.Fatal(err)
	}
	t.Logf("bytes=%d full_checks=%d full_validation_bytes=%d", read, b.validationReads, b.validationBytes)
	if b.validationReads != 4 || b.validationBytes != uint64(4*len(body)) {
		t.Fatal("unbounded or missing full reconciliation work")
	}
}

func TestCoverageReportLifetimeDescriptorCloseFailure(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "input.go"), []byte("source"), 0600); err != nil {
		t.Fatal(err)
	}
	scope, err := OpenInputScope(root)
	if err != nil {
		t.Fatal(err)
	}
	b, err := scope.Acquire("input.go", nil)
	if err != nil {
		t.Fatal(err)
	}
	if err = b.Close(); err != nil {
		t.Fatal(err)
	}
	if err = scope.anchor.file.Close(); err != nil {
		t.Fatal(err)
	}
	if err = scope.Close(); err == nil {
		t.Fatal("descriptor close/reconciliation failure lost")
	}
}

func TestCoverageStatisticsReaderErrors(t *testing.T) {
	for _, mode := range []string{"ReadAt", "Seek", "UnreadRune", "UnreadByte", "valid"} {
		t.Run(mode, func(t *testing.T) {
			_, scope, b := inputFixture(t)
			r, err := b.Reader()
			if err != nil {
				t.Fatal(err)
			}
			switch mode {
			case "ReadAt":
				_, err = r.ReadAt(make([]byte, 1), -1)
			case "Seek":
				_, err = r.Seek(0, 999)
			case "UnreadRune":
				err = r.UnreadRune()
			case "UnreadByte":
				err = r.UnreadByte()
			case "valid":
				if _, err = r.ReadByte(); err != nil {
					t.Fatal(err)
				}
				if err = r.UnreadByte(); err != nil {
					t.Fatal(err)
				}
				if _, _, err = r.ReadRune(); err != nil {
					t.Fatal(err)
				}
				if err = r.UnreadRune(); err != nil {
					t.Fatal(err)
				}
				if _, err = r.Seek(0, io.SeekStart); err != nil {
					t.Fatal(err)
				}
				if _, err = r.ReadAt(make([]byte, 1), 0); err != nil {
					t.Fatal(err)
				}
				if _, err = io.ReadAll(r); err != nil {
					t.Fatal(err)
				}
			}
			terminal := scope.Close()
			if mode == "valid" {
				if terminal != nil {
					t.Fatal(terminal)
				}
			} else if err == nil || terminal == nil {
				t.Fatalf("operation=%v terminal=%v", err, terminal)
			}
		})
	}
}
func TestCoverageStatisticsSampleContext(t *testing.T) {
	for _, mode := range []string{"valid", "component", "identical-leaf", "missing", "directory", "post-close-change", "forced-close", "single"} {
		t.Run(mode, func(t *testing.T) {
			root := t.TempDir()
			name := "part/script"
			full := filepath.Join(root, name)
			if err := os.MkdirAll(filepath.Dir(full), 0700); err != nil {
				t.Fatal(err)
			}
			body := bytes.Repeat([]byte("sample text\n"), 10000)
			if err := os.WriteFile(full, body, 0600); err != nil {
				t.Fatal(err)
			}
			target := root
			if mode == "single" {
				target = full
				name = "."
			}
			scope, err := OpenInputScope(target)
			if err != nil {
				t.Fatal(err)
			}
			defer scope.Close()
			if mode == "component" || mode == "identical-leaf" {
				outside := t.TempDir()
				other := filepath.Join(outside, "script")
				if err = os.WriteFile(other, body, 0600); err != nil {
					t.Fatal(err)
				}
				if mode == "component" {
					replaceInput(t, filepath.Dir(full), outside)
				} else {
					replaceInput(t, full, other)
				}
			}
			if mode == "missing" {
				if err = os.Remove(full); err != nil {
					t.Fatal(err)
				}
			}
			if mode == "directory" {
				if err = os.Remove(full); err != nil {
					t.Fatal(err)
				}
				if err = os.Mkdir(full, 0700); err != nil {
					t.Fatal(err)
				}
			}
			sample, info, err := scope.ReportSample(name, 32)
			if mode == "component" || mode == "identical-leaf" || mode == "missing" || mode == "directory" {
				if err == nil || scope.Close() == nil {
					t.Fatal("unsafe acquisition accepted")
				}
				return
			}
			if err != nil || len(sample) != 32 || info.Size() != int64(len(body)) {
				t.Fatalf("sample=%d info=%v err=%v", len(sample), info, err)
			}
			if len(scope.inputs) != 1 || scope.inputs[0].sampleLimit != 32 || len(scope.inputs[0].data) != 0 {
				t.Fatal("sample ownership or context bound incorrect")
			}
			b := scope.inputs[0]
			if b.validationBytes > 128 {
				t.Fatalf("unbounded sample work: %d", b.validationBytes)
			}
			if mode == "post-close-change" {
				if err = os.WriteFile(full, bytes.Repeat([]byte("x"), len(body)), 0600); err != nil {
					t.Fatal(err)
				}
			}
			if mode == "forced-close" {
				if err = scope.anchor.file.Close(); err != nil {
					t.Fatal(err)
				}
			}
			err = scope.Close()
			if mode == "post-close-change" || mode == "forced-close" {
				if err == nil {
					t.Fatal("final fault accepted")
				}
			} else if err != nil {
				t.Fatal(err)
			}
		})
	}
}

// A rejected expectation must win over even a failing context acquisition.
// The valid expectation control proves that the same context really is reached.
func TestCoverageExpectedContextOrdering(t *testing.T) {
	for _, mode := range []string{"wrong-root", "wrong-identity", "wrong-relative", "wrong-digest", "valid-context-failure"} {
		t.Run(mode, func(t *testing.T) {
			root, _, parent := inputFixture(t)
			expected := parent.Expectation()
			outside := filepath.Join(t.TempDir(), "outside-ignore")
			if err := os.WriteFile(outside, []byte("input.go\n"), 0600); err != nil {
				t.Fatal(err)
			}
			if err := os.Symlink(outside, filepath.Join(root, ".gitignore")); err != nil {
				t.Fatal(err)
			}
			want := "invalid input expectation"
			switch mode {
			case "wrong-root":
				expected.Root += "/different"
			case "wrong-identity":
				expected.RootIdentity.Inode++
			case "wrong-relative":
				expected.Relative = "other.go"
			case "wrong-digest":
				expected.SHA256 = strings.Repeat("0", 64)
				want = "input digest mismatch"
			}
			visits, directories := 0, 0
			err := IterateFilesListWithInput(root, []string{"input.go"}, false, false,
				func(*Path) (bool, error) { directories++; return true, nil },
				func(*FileInfo) error { visits++; return nil }, nil, expected)
			t.Logf("mode=%s visits=%d directories=%d error=%v", mode, visits, directories, err)
			if err == nil || visits != 0 || directories != 0 {
				t.Fatalf("invalid ordering/completion: %v", err)
			}
			if mode == "valid-context-failure" {
				if strings.Contains(err.Error(), "expectation") || strings.Contains(err.Error(), "digest mismatch") {
					t.Fatalf("valid expectation did not reach failing context: %v", err)
				}
			} else if !strings.Contains(err.Error(), want) {
				t.Fatalf("context consumed before expected-input rejection: %v; want %s", err, want)
			}
		})
	}
}
func TestCoverageExpectedContextReuseAndCompatibility(t *testing.T) {
	for _, mode := range []string{"valid", "absent-context", "excluded", "single", "no-expectation", "missing-input"} {
		t.Run(mode, func(t *testing.T) {
			root, _, parent := inputFixture(t)
			expected := parent.Expectation()
			target, logical := root, "input.go"
			if mode == "single" {
				target = filepath.Join(root, logical)
				logical = "."
				scope, err := OpenInputScope(target)
				if err != nil {
					t.Fatal(err)
				}
				defer scope.Close()
				selected, err := scope.Acquire(logical, nil)
				if err != nil {
					t.Fatal(err)
				}
				expected = selected.Expectation()
			}
			if mode != "absent-context" {
				content := "# synthetic context\n"
				if mode == "excluded" {
					content = "input.go\n"
				}
				if err := os.WriteFile(filepath.Join(root, ".gitignore"), []byte(content), 0600); err != nil {
					t.Fatal(err)
				}
			}
			if mode == "no-expectation" {
				expected = nil
			}
			if mode == "missing-input" {
				if err := os.Remove(filepath.Join(root, "input.go")); err != nil {
					t.Fatal(err)
				}
			}
			visits, exclusions := 0, 0
			var owned *InputScope
			err := IterateFilesListWithInput(target, []string{logical}, false, false, func(*Path) (bool, error) { return true, nil }, func(info *FileInfo) error {
				visits++
				owned = info.scope
				if info.RelativePath != logical {
					t.Errorf("logical path=%q", info.RelativePath)
				}
				count := 0
				for _, b := range owned.inputs {
					if b.relative == logical && b.sampleLimit == 0 {
						count++
					}
				}
				if count != 1 {
					t.Errorf("selected acquisition not reused: %d", count)
				}
				body, readErr := info.ReadBytes()
				if readErr != nil {
					return readErr
				}
				if string(body) != "package synthetic\n" {
					t.Errorf("changed selected bytes")
				}
				return nil
			}, func(string) { exclusions++ }, expected)
			t.Logf("mode=%s visits=%d exclusions=%d error=%v", mode, visits, exclusions, err)
			if mode == "missing-input" {
				if err == nil || visits != 0 {
					t.Fatal("missing acquisition accepted")
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if mode == "excluded" {
				if visits != 0 || exclusions != 1 {
					t.Fatal("exclusion became processing")
				}
			} else if visits != 1 || exclusions != 0 {
				t.Fatal("valid input not visited")
			}
			if owned != nil && !owned.closed {
				t.Fatal("scope leaked")
			}
		})
	}
}

func TestCoverageContextExistsConfinement(t *testing.T) {
	for _, mode := range []string{"file", "directory", "absent-leaf", "absent-component", "leaf-link", "component-link", "escape", "nul", "not-directory", "closed-anchor", "closed-scope", "root-replaced", "single-context"} {
		t.Run(mode, func(t *testing.T) {
			root, scope, selected := inputFixture(t)
			marker := filepath.Join(root, "context", "marker")
			if err := os.Mkdir(filepath.Dir(marker), 0700); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(marker, []byte("metadata only\n"), 0600); err != nil {
				t.Fatal(err)
			}
			wantPresent, wantFault := true, false
			switch mode {
			case "directory":
				marker = filepath.Dir(marker)
			case "absent-leaf":
				marker = filepath.Join(root, "context", "absent")
				wantPresent = false
			case "absent-component":
				marker = filepath.Join(root, "absent", "marker")
				wantPresent = false
			case "leaf-link", "component-link":
				outside := t.TempDir()
				os.WriteFile(filepath.Join(outside, "marker"), []byte("metadata only\n"), 0600)
				if mode == "leaf-link" {
					replaceInput(t, marker, filepath.Join(outside, "marker"))
				} else {
					replaceInput(t, filepath.Dir(marker), outside)
				}
				wantPresent = false
				wantFault = true
			case "escape":
				marker = filepath.Join(root, "..", "outside")
				wantPresent = false
				wantFault = true
			case "nul":
				marker = filepath.Join(root, "bad\x00marker")
				wantPresent = false
				wantFault = true
			case "not-directory":
				marker = filepath.Join(marker, "child")
				wantPresent = false
				wantFault = true
			case "closed-anchor":
				if err := scope.anchor.file.Close(); err != nil {
					t.Fatal(err)
				}
				wantPresent = false
				wantFault = true
			case "closed-scope":
				if err := scope.Close(); err != nil {
					t.Fatal(err)
				}
				wantPresent = false
				wantFault = true
			case "root-replaced":
				replaceInput(t, root, t.TempDir())
				wantPresent = false
				wantFault = true
			case "single-context":
				if err := scope.Close(); err != nil {
					t.Fatal(err)
				}
				var err error
				scope, err = OpenInputScope(filepath.Join(root, "input.go"))
				if err != nil {
					t.Fatal(err)
				}
				t.Cleanup(func() { scope.Close() })
				selected, err = scope.Acquire(".", nil)
				if err != nil {
					t.Fatal(err)
				}
			}
			p := &Path{AbsolutePath: marker, scope: scope}
			got := p.Exists()
			if got != wantPresent {
				t.Fatalf("exists=%v want=%v", got, wantPresent)
			}
			if (scope.Err() != nil) != wantFault {
				t.Fatalf("fault=%v wanted=%v", scope.Err(), wantFault)
			}
			if !wantFault {
				for i := 0; i < 3; i++ {
					if p.Exists() != got {
						t.Fatal("stable repeated observation changed")
					}
				}
				if len(scope.existence) != 1 || len(scope.inputs) != 1 {
					t.Fatal("metadata duplicated or became processing input")
				}
				if err := selected.Close(); err != nil {
					t.Fatal(err)
				}
			}
			terminal := scope.Close()
			t.Logf("mode=%s observed=%v terminal=%v", mode, got, terminal)
			if (terminal != nil) != wantFault {
				t.Fatalf("terminal=%v", terminal)
			}
		})
	}
}

func TestCoverageContextExistsReconciliation(t *testing.T) {
	for _, mode := range []string{"stable-file", "stable-absence", "content-only", "removed", "created", "leaf-replaced", "component-replaced", "absent-prefix-created", "kind-changed", "close-failure"} {
		t.Run(mode, func(t *testing.T) {
			root, scope, selected := inputFixture(t)
			dir := filepath.Join(root, "context")
			if err := os.Mkdir(dir, 0700); err != nil {
				t.Fatal(err)
			}
			marker := filepath.Join(dir, "marker")
			if err := os.WriteFile(marker, []byte("before\n"), 0600); err != nil {
				t.Fatal(err)
			}
			absent := mode == "stable-absence" || mode == "created" || mode == "absent-prefix-created"
			if absent {
				if err := os.Remove(marker); err != nil {
					t.Fatal(err)
				}
			}
			if mode == "absent-prefix-created" {
				if err := os.Remove(dir); err != nil {
					t.Fatal(err)
				}
			}
			path := &Path{AbsolutePath: marker, scope: scope}
			if path.Exists() == absent || scope.Err() != nil {
				t.Fatalf("observation fixture: %v", scope.Err())
			}
			if err := selected.Close(); err != nil {
				t.Fatal(err)
			}
			if err := scope.Validate(); err != nil {
				t.Fatal(err)
			}
			switch mode {
			case "content-only":
				os.WriteFile(marker, []byte("after content change\n"), 0600)
			case "removed":
				os.Remove(marker)
			case "created":
				os.WriteFile(marker, []byte("before\n"), 0600)
			case "leaf-replaced":
				outside := t.TempDir()
				os.WriteFile(filepath.Join(outside, "marker"), []byte("before\n"), 0600)
				replaceInput(t, marker, filepath.Join(outside, "marker"))
			case "component-replaced":
				outside := t.TempDir()
				os.WriteFile(filepath.Join(outside, "marker"), []byte("before\n"), 0600)
				replaceInput(t, dir, outside)
			case "absent-prefix-created":
				os.Mkdir(dir, 0700)
			case "kind-changed":
				os.Remove(marker)
				os.Mkdir(marker, 0700)
			case "close-failure":
				if err := scope.anchor.file.Close(); err != nil {
					t.Fatal(err)
				}
			}
			wantFault := mode != "stable-file" && mode != "stable-absence" && mode != "content-only"
			validation := scope.Validate()
			terminal := scope.Close()
			t.Logf("mode=%s validation=%v terminal=%v", mode, validation, terminal)
			if (validation != nil) != wantFault || (terminal != nil) != wantFault {
				t.Fatal("retained metadata obligation lost or overclaims content")
			}
			if mode == "close-failure" && !errors.Is(terminal, os.ErrClosed) {
				t.Fatal("anchor close failure not propagated")
			}
		})
	}
}

func TestCoverageContextHandoffLifetime(t *testing.T) {
	t.Run("consumed-provenance", contextConsumedLifetime)
	for _, mode := range []string{"present", "absent", "empty", "leaf-link", "absent-created", "component-link", "after-import", "content-only", "wrong-input", "missing", "closed-parent", "delivery-close"} {
		t.Run(mode, func(t *testing.T) {
			root, child, b := inputFixture(t)
			expected := b.Expectation()
			marker := filepath.Join(root, "ctx", "marker")
			if err := os.Mkdir(filepath.Dir(marker), 0700); err != nil {
				t.Fatal(err)
			}
			absent := mode == "absent" || mode == "absent-created"
			if !absent {
				if err := os.WriteFile(marker, []byte("same"), 0600); err != nil {
					t.Fatal(err)
				}
			}
			if mode != "empty" && (&Path{AbsolutePath: marker, scope: child}).Exists() == absent {
				t.Fatal("incorrect initial observation")
			}
			receipt, err := child.ContextReceipt(expected)
			if err != nil {
				t.Fatal(err)
			}
			if mode == "empty" && len(receipt.Observations) != 0 {
				t.Fatal("invented observation")
			}
			if err = child.Close(); err != nil {
				t.Fatal(err)
			}
			parent, err := OpenInputScope(root)
			if err != nil {
				t.Fatal(err)
			}
			defer parent.Close()
			held, err := parent.Acquire("input.go", expected)
			if err != nil {
				t.Fatal(err)
			}
			if err = held.Close(); err != nil {
				t.Fatal(err)
			}
			replace := func() {
				outside := t.TempDir()
				if err := os.WriteFile(filepath.Join(outside, "marker"), []byte("same"), 0600); err != nil {
					t.Fatal(err)
				}
				if mode == "component-link" {
					if err := os.Rename(filepath.Dir(marker), filepath.Dir(marker)+".held"); err != nil {
						t.Fatal(err)
					}
					if err := os.Symlink(outside, filepath.Dir(marker)); err != nil {
						t.Fatal(err)
					}
				} else {
					if err := os.Rename(marker, marker+".held"); err != nil {
						t.Fatal(err)
					}
					if err := os.Symlink(filepath.Join(outside, "marker"), marker); err != nil {
						t.Fatal(err)
					}
				}
			}
			switch mode {
			case "leaf-link", "component-link":
				replace()
			case "absent-created":
				if err := os.WriteFile(marker, []byte("same"), 0600); err != nil {
					t.Fatal(err)
				}
			case "wrong-input":
				receipt.Input.SHA256 = strings.Repeat("0", 64)
			case "closed-parent":
				parent.Close()
			case "content-only":
				if err := os.WriteFile(marker, []byte("different content"), 0600); err != nil {
					t.Fatal(err)
				}
			}
			if mode == "missing" {
				if parent.RequireContextReceipts([]string{"input.go"}) == nil {
					t.Fatal("missing receipt accepted")
				}
				return
			}
			err = parent.ImportContextReceipt(receipt, expected)
			bad := mode == "leaf-link" || mode == "component-link" || mode == "absent-created" || mode == "wrong-input" || mode == "closed-parent"
			if bad {
				if err == nil {
					t.Fatal("invalid transfer accepted")
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if err = parent.RequireContextReceipts([]string{"input.go"}); err != nil {
				t.Fatal(err)
			}
			if mode == "after-import" {
				replace()
				if parent.Validate() == nil {
					t.Fatal("changed context accepted")
				}
				return
			}
			if mode == "delivery-close" {
				if err := parent.anchor.file.Close(); err != nil {
					t.Fatal(err)
				}
				if parent.Close() == nil {
					t.Fatal("close failure lost")
				}
				return
			}
			if err = parent.Validate(); err != nil {
				t.Fatal(err)
			}
			if err = parent.Close(); err != nil {
				t.Fatal(err)
			}
		})
	}
}
func TestCoverageContextHandoffDecode(t *testing.T) {
	t.Run("consumed-provenance", contextConsumedDecode)
	root, s, b := inputFixture(t)
	(&Path{AbsolutePath: filepath.Join(root, "absent"), scope: s}).Exists()
	r, err := s.ContextReceipt(b.Expectation())
	if err != nil {
		t.Fatal(err)
	}
	raw, err := json.Marshal(r)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = DecodeContextReceipt(raw); err != nil {
		t.Fatal(err)
	}
	cases := map[string]string{
		"missing":    strings.Replace(string(raw), `"observations":`, `"missing":`, 1),
		"alias":      strings.Replace(string(raw), `"relative":`, `"Relative":`, 1),
		"null":       strings.Replace(string(raw), `"present":false`, `"present":null`, 1),
		"duplicate":  strings.Replace(string(raw), `"present":false`, `"present":true,"present":false`, 1),
		"wrong-type": strings.Replace(string(raw), `"present":false`, `"present":"false"`, 1),
		"escape":     strings.Replace(string(raw), `"relative":"absent"`, `"relative":"../absent"`, 1),
		"partial":    string(raw[:len(raw)-1]),
		"trailing":   string(raw) + ` {}`,
	}
	for name, value := range cases {
		t.Run(name, func(t *testing.T) {
			if _, err := DecodeContextReceipt([]byte(value)); err == nil {
				t.Fatal("invalid receipt accepted")
			}
		})
	}
}

// Exercise generic acquired bytes independently of any detector filename.
func contextConsumedLifetime(t *testing.T) {
	for _, mode := range []string{"valid", "content-before-import", "leaf-before-import", "component-before-import", "after-import", "size-change", "conflicting-transfer", "closed-parent", "close-failure", "producer-close-failure"} {
		t.Run(mode, func(t *testing.T) {
			root, child, selected := inputFixture(t)
			dir := filepath.Join(root, "configuration")
			if err := os.Mkdir(dir, 0700); err != nil {
				t.Fatal(err)
			}
			for _, name := range []string{"one.conf", "two.conf"} {
				if err := os.WriteFile(filepath.Join(dir, name), []byte("original-context"), 0600); err != nil {
					t.Fatal(err)
				}
			}
			for i, name := range []string{"one.conf", "two.conf"} {
				path := &Path{AbsolutePath: filepath.Join(dir, name), scope: child}
				if i == 0 {
					if _, err := path.ReadBytes(); err != nil {
						t.Fatal(err)
					}
				} else {
					r, err := path.Reader()
					if err != nil {
						t.Fatal(err)
					}
					_, err = io.ReadAll(r)
					if err = errors.Join(err, r.Close()); err != nil {
						t.Fatal(err)
					}
				}
			}
			receipt, err := child.ContextReceipt(selected.Expectation())
			if err != nil {
				t.Fatal(err)
			}
			if len(receipt.Reads) != 2 || len(receipt.Observations) != 0 {
				t.Fatalf("wrong actual provenance reads=%d observations=%d", len(receipt.Reads), len(receipt.Observations))
			}
			raw, err := json.Marshal(receipt)
			if err != nil {
				t.Fatal(err)
			}
			receipt, err = DecodeContextReceipt(raw)
			if err != nil {
				t.Fatal(err)
			}
			if mode == "producer-close-failure" {
				if err = child.anchor.file.Close(); err != nil {
					t.Fatal(err)
				}
				if child.Close() == nil {
					t.Fatal("producer close failure disappeared")
				}
				return
			}
			if err = child.Close(); err != nil {
				t.Fatal(err)
			}
			parent, err := OpenInputScope(root)
			if err != nil {
				t.Fatal(err)
			}
			defer parent.Close()
			b, err := parent.Acquire("input.go", selected.Expectation())
			if err != nil {
				t.Fatal(err)
			}
			if err = b.Close(); err != nil {
				t.Fatal(err)
			}
			mutate := func() {
				switch mode {
				case "leaf-before-import":
					outside := filepath.Join(t.TempDir(), "identical")
					if e := os.WriteFile(outside, []byte("original-context"), 0600); e != nil {
						t.Fatal(e)
					}
					replaceInput(t, filepath.Join(dir, "one.conf"), outside)
				case "component-before-import":
					outside := t.TempDir()
					for _, n := range []string{"one.conf", "two.conf"} {
						if e := os.WriteFile(filepath.Join(outside, n), []byte("original-context"), 0600); e != nil {
							t.Fatal(e)
						}
					}
					replaceInput(t, dir, outside)
				case "size-change":
					if e := os.WriteFile(filepath.Join(dir, "one.conf"), []byte("short"), 0600); e != nil {
						t.Fatal(e)
					}
				default:
					if e := os.WriteFile(filepath.Join(dir, "one.conf"), []byte("modified-context"), 0600); e != nil {
						t.Fatal(e)
					}
				}
			}
			before := mode == "content-before-import" || mode == "leaf-before-import" || mode == "component-before-import" || mode == "size-change"
			if before {
				mutate()
			}
			if mode == "closed-parent" {
				parent.Close()
			}
			err = parent.ImportContextReceipt(receipt, selected.Expectation())
			if before || mode == "closed-parent" {
				if err == nil {
					t.Fatal("bad context transfer accepted")
				}
				if parent.Close() == nil {
					t.Fatal("transfer failure not terminal")
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			count := len(parent.inputs)
			if err = parent.ImportContextReceipt(receipt, nil); err != nil || len(parent.inputs) != count {
				t.Fatalf("duplicate import retains redundant resources: %v", err)
			}
			if mode == "conflicting-transfer" {
				receipt.Reads[0].Input.SHA256 = strings.Repeat("0", 64)
				if parent.ImportContextReceipt(receipt, nil) == nil {
					t.Fatal("conflicting handoff accepted")
				}
				if parent.Close() == nil {
					t.Fatal("conflict not terminal")
				}
				return
			}
			// A context-only capability must not be promoted to the selected processing input.
			if parent.matchesInputLocked(&receipt.Reads[0].Input) {
				t.Fatal("context satisfied selected promise")
			}
			if mode == "after-import" {
				mutate()
			}
			if mode == "close-failure" {
				if err = parent.anchor.file.Close(); err != nil {
					t.Fatal(err)
				}
			}
			validation := parent.Validate()
			terminal := parent.Close()
			t.Logf("mode=%s reads=%d retained=%d validation=%v close=%v", mode, len(receipt.Reads), count, validation, terminal)
			bad := mode == "after-import" || mode == "close-failure"
			if (validation != nil) != bad || (terminal != nil) != bad {
				t.Fatal("completed context lifetime lost")
			}
		})
	}
}
func contextConsumedDecode(t *testing.T) {
	root, s, b := inputFixture(t)
	path := filepath.Join(root, "arbitrary.conf")
	if err := os.WriteFile(path, []byte("context"), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := (&Path{AbsolutePath: path, scope: s}).ReadBytes(); err != nil {
		t.Fatal(err)
	}
	receipt, err := s.ContextReceipt(b.Expectation())
	if err != nil {
		t.Fatal(err)
	}
	raw, err := json.Marshal(receipt)
	if err != nil {
		t.Fatal(err)
	}
	if len(receipt.Reads) != 1 {
		t.Fatal("read provenance absent")
	}
	var fields map[string]json.RawMessage
	if err = json.Unmarshal(raw, &fields); err != nil {
		t.Fatal(err)
	}
	for _, mode := range []string{"missing", "null", "duplicate", "alias", "null-size", "wrong-size", "negative-size", "wrong-root", "wrong-components", "bad-digest", "unsafe-context", "conflicting-limit"} {
		t.Run(mode, func(t *testing.T) {
			var r ContextReceipt
			if e := json.Unmarshal(raw, &r); e != nil {
				t.Fatal(e)
			}
			switch mode {
			case "duplicate":
				r.Reads = append(r.Reads, r.Reads[0])
			case "negative-size":
				r.Reads[0].Size = -1
			case "wrong-root":
				r.Reads[0].Input.Root += "/other"
			case "wrong-components":
				r.Reads[0].Input.Components = nil
			case "bad-digest":
				r.Reads[0].Input.SHA256 = "bad"
			case "unsafe-context":
				r.Reads[0].Input.Relative = "../outside"
			case "conflicting-limit":
				r.Reads[0].SampleLimit = 16385
			}
			changed, e := json.Marshal(r)
			if e != nil {
				t.Fatal(e)
			}
			switch mode {
			case "missing":
				changed = bytes.Replace(changed, []byte(`"reads":`), []byte(`"omitted":`), 1)
			case "null":
				var m map[string]json.RawMessage
				json.Unmarshal(changed, &m)
				m["reads"] = json.RawMessage("null")
				changed, _ = json.Marshal(m)
			case "alias":
				changed = bytes.Replace(changed, []byte(`"sample_limit":`), []byte(`"Sample_Limit":`), 1)
			case "null-size":
				changed = bytes.Replace(changed, []byte(`"size":7`), []byte(`"size":null`), 1)
			case "wrong-size":
				changed = bytes.Replace(changed, []byte(`"size":7`), []byte(`"size":"7"`), 1)
			}
			if _, e = DecodeContextReceipt(changed); e == nil {
				t.Fatalf("malformed context reads accepted: %s", changed)
			}
		})
	}
	legacy := *receipt
	legacy.Schema = legacyContextReceiptSchema
	legacy.Reads = nil
	legacyRaw, err := json.Marshal(legacy)
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := DecodeContextReceipt(legacyRaw)
	if err != nil {
		t.Fatal("valid legacy cache receipt rejected", err)
	}
	if s.ImportContextReceipt(decoded, b.Expectation()) == nil {
		t.Fatal("legacy receipt accepted as fresh complete handoff")
	}
}
