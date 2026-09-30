package gitleaks

import (
	"archive/zip"
	"bytes"
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"github.com/bearer/bearer/pkg/report/writer"
	"github.com/bearer/bearer/pkg/util/file"
	"github.com/fatih/semgroup"
	"github.com/zricethezav/gitleaks/v8/config"
	"github.com/zricethezav/gitleaks/v8/detect"
	gre "github.com/zricethezav/gitleaks/v8/regexp"
	"github.com/zricethezav/gitleaks/v8/sources"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestCoverageF4GitleaksReaderParity(t *testing.T) {
	rule := config.Rule{RuleID: "synthetic-token", Description: "synthetic fixture", Regex: gre.MustCompile(`token_[a-z]+`), Keywords: []string{"token_"}}
	for _, mode := range []string{"plain", "allowlist", "chunk", "binary", "empty", "archive-disabled", "archive-enabled"} {
		t.Run(mode, func(t *testing.T) {
			cfg := config.Config{Rules: map[string]config.Rule{"synthetic-token": rule}, OrderedRules: []string{"synthetic-token"}, Keywords: map[string]struct{}{"token_": {}}}
			content := []byte("synthetic token_alpha\n")
			name := "input.txt"
			depth := 0
			switch mode {
			case "allowlist":
				cfg.Allowlists = []*config.Allowlist{{Paths: []*gre.Regexp{gre.MustCompile(`input.txt$`)}}}
			case "chunk":
				content = []byte(strings.Repeat("ordinary line\n", 10000) + "token_alpha\n")
			case "binary":
				content = []byte{'P', 'N', 'G', 0, 0, 0}
			case "empty":
				content = nil
			case "archive-disabled", "archive-enabled":
				name = "input.zip"
				var buf bytes.Buffer
				z := zip.NewWriter(&buf)
				f, e := z.Create("inner.txt")
				if e != nil {
					t.Fatal(e)
				}
				f.Write(content)
				if e = z.Close(); e != nil {
					t.Fatal(e)
				}
				content = buf.Bytes()
				if mode == "archive-enabled" {
					depth = 1
				}
			}
			path := filepath.Join(t.TempDir(), name)
			if err := os.WriteFile(path, content, 0600); err != nil {
				t.Fatal(err)
			}
			left := detect.NewDetector(cfg)
			a, err := left.DetectSource(context.Background(), &readerSource{content: bytes.NewReader(content), Path: path, Config: &cfg, MaxArchiveDepth: depth})
			if depth > 0 {
				// Adopted C1 compatibility boundary: expansion now fails closed.
				if !errors.Is(err, ErrArchiveExpansionUnsupported) {
					t.Fatalf("enabled expansion accepted: %v", err)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			right := detect.NewDetector(cfg)
			b, err := right.DetectSource(context.Background(), &sources.Files{Path: path, Config: &cfg, Sema: semgroup.NewGroup(context.Background(), 1), MaxArchiveDepth: depth})
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(a, b) {
				t.Fatalf("reader/path mismatch: %#v / %#v", a, b)
			}
			if (mode == "plain" || mode == "chunk" || mode == "archive-enabled") && len(a) == 0 {
				t.Fatal("positive fixture did not detect")
			}
			if (mode == "allowlist" || mode == "empty" || mode == "archive-disabled") && len(a) != 0 {
				t.Fatal("excluded fixture detected")
			}
		})
	}
}
func TestCoverageF4GitleaksReaderFailure(t *testing.T) {
	inputErr := errors.New("synthetic reader failure")
	d := detect.NewDetector(config.Config{})
	_, err := d.DetectSource(context.Background(), &readerSource{content: &failingReader{inputErr}, Path: "synthetic.txt"})
	if !errors.Is(err, inputErr) {
		t.Fatalf("lost read error: %v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := (&readerSource{content: bytes.NewReader(nil)}).Fragments(ctx, func(sources.Fragment, error) error { return nil }); !errors.Is(err, context.Canceled) {
		t.Fatalf("cancel: %v", err)
	}
}

func TestCoverageF4ResidualBoundInputProcessFile(t *testing.T) {
	for _, mode := range []string{"valid", "replaced"} {
		t.Run(mode, func(t *testing.T) {
			dir := t.TempDir()
			name := filepath.Join(dir, "fixture.txt")
			if err := os.WriteFile(name, []byte("ordinary input\n"), 0600); err != nil {
				t.Fatal(err)
			}

			detector := New(nil)
			var processErr error
			scanErr := file.IterateFilesListWithExclusions(
				dir,
				[]string{"fixture.txt"},
				false,
				true,
				func(*file.Path) (bool, error) { return true, nil },
				func(info *file.FileInfo) error {
					if mode == "replaced" {
						if err := os.Rename(name, name+".held"); err != nil {
							return err
						}
						defer func() { _ = os.Rename(name+".held", name) }()
					}
					_, processErr = detector.ProcessFile(info, nil, &writer.Detectors{File: &bytes.Buffer{}})
					return nil
				},
				nil,
			)

			if mode == "valid" {
				if processErr != nil || scanErr != nil {
					t.Fatalf("bound input ProcessFile failed: process=%v scan=%v", processErr, scanErr)
				}
				return
			}
			if processErr == nil {
				t.Fatal("ProcessFile accepted a replaced bound input")
			}
			if scanErr == nil {
				t.Fatal("bound input replacement was not latched by the input scope")
			}
		})
	}
}

func TestCoverageF4ResidualBoundInputArchiveReader(t *testing.T) {
	dir := t.TempDir()
	name := filepath.Join(dir, "fixture.zip")
	var archive bytes.Buffer
	zipWriter := zip.NewWriter(&archive)
	inner, err := zipWriter.Create("inner.txt")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := inner.Write([]byte("ordinary archive input\n")); err != nil {
		t.Fatal(err)
	}
	if err := zipWriter.Close(); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(name, archive.Bytes(), 0600); err != nil {
		t.Fatal(err)
	}

	scope, err := file.OpenInputScope(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = scope.Close() }()
	bound, err := scope.Acquire("fixture.zip", nil)
	if err != nil {
		t.Fatal(err)
	}
	reader, err := bound.Reader()
	if err != nil {
		t.Fatal(err)
	}
	defer reader.Close()
	if _, ok := interface{}(reader).(interface {
		io.ReaderAt
		io.Seeker
	}); !ok {
		t.Fatal("production bound input reader is not seekable for archive extraction")
	}

	cfg := config.Config{}
	if _, err := detect.NewDetector(cfg).DetectSource(context.Background(), &readerSource{
		content:         reader,
		Path:            "fixture.zip",
		Config:          &cfg,
		MaxArchiveDepth: 1,
	}); !errors.Is(err, ErrArchiveExpansionUnsupported) {
		t.Fatalf("enabled archive request not rejected: %v", err)
	}
}

func TestCoverageF4ResidualBoundInputArchiveApplicability(t *testing.T) {
	for _, tc := range []struct {
		name    string
		content []byte
		wantErr bool
	}{
		{name: "fixture.gz", content: []byte("not a gzip archive\n"), wantErr: true},
		{name: "fixture.zip", content: []byte("not a zip archive\n"), wantErr: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			if err := os.WriteFile(filepath.Join(dir, tc.name), tc.content, 0600); err != nil {
				t.Fatal(err)
			}

			detector := New(nil).(*detector)
			detector.gitleaksDetector.MaxArchiveDepth = 1
			var processErr error
			scanErr := file.IterateFilesListWithExclusions(
				dir,
				[]string{tc.name},
				false,
				true,
				func(*file.Path) (bool, error) { return true, nil },
				func(info *file.FileInfo) error {
					_, processErr = detector.ProcessFile(info, nil, &writer.Detectors{File: &bytes.Buffer{}})
					return nil
				},
				nil,
			)

			if scanErr != nil {
				t.Fatalf("bound input acquisition failed: %v", scanErr)
			}
			t.Logf("C1 adopted mode rejection: path=%s ProcessFile=%v scope=%v; decoder remains unqualified", tc.name, processErr, scanErr)
			if (processErr != nil) != tc.wantErr {
				t.Fatalf("ProcessFile error = %v, want error = %v", processErr, tc.wantErr)
			}
		})
	}
}

type failingReader struct{ err error }

func (r *failingReader) Read([]byte) (int, error) { return 0, r.err }

var _ io.Reader = (*failingReader)(nil)

func TestCoverageStatisticsArchiveRejection(t *testing.T) {
	var zipped bytes.Buffer
	z := zip.NewWriter(&zipped)
	w, err := z.Create("ordinary.txt")
	if err != nil {
		t.Fatal(err)
	}
	if _, err = w.Write([]byte("synthetic ordinary text")); err != nil {
		t.Fatal(err)
	}
	if err = z.Close(); err != nil {
		t.Fatal(err)
	}
	var gz bytes.Buffer
	zw := gzip.NewWriter(&gz)
	if _, err = zw.Write([]byte("synthetic ordinary text")); err != nil {
		t.Fatal(err)
	}
	if err = zw.Close(); err != nil {
		t.Fatal(err)
	}

	for _, tc := range []struct {
		name string
		body []byte
	}{{"valid.zip", zipped.Bytes()}, {"valid.gz", gz.Bytes()}, {"bad.zip", []byte("bad zip")}, {"bad.gz", []byte("bad gzip")}, {"plain.txt", []byte("ordinary")}} {
		for _, depth := range []int{0, 1, 2} {
			t.Run(fmt.Sprintf("%s/%d", tc.name, depth), func(t *testing.T) {
				root := t.TempDir()
				if err := os.WriteFile(filepath.Join(root, tc.name), tc.body, 0600); err != nil {
					t.Fatal(err)
				}
				d := New(nil).(*detector)
				d.gitleaksDetector.MaxArchiveDepth = depth
				// Direct invocation uses the existing secure public FileInfo reader.
				// The earlier iterator-only fixture never reached valid archives.
				info, err := file.FileInfoFromPath(filepath.Join(root, tc.name))
				if err != nil {
					t.Fatal(err)
				}
				entered := 0
				invoke := func() error {
					entered++
					_, e := d.ProcessFile(info, nil, &writer.Detectors{File: &bytes.Buffer{}})
					return e
				}
				returned := invoke()
				t.Logf("direct entry=%d effective_depth=%d returned=%v caller_outcome=%v", entered, d.gitleaksDetector.MaxArchiveDepth, returned, returned)
				if entered != 1 || d.gitleaksDetector.MaxArchiveDepth != depth {
					t.Fatal("adapter invocation precondition missing")
				}
				if depth > 0 {
					if !errors.Is(returned, ErrArchiveExpansionUnsupported) {
						t.Fatalf("enabled adapter request accepted: %v", returned)
					}
				} else if returned != nil {
					t.Fatalf("depth-zero input: %v", returned)
				}

				// Observe eligibility separately. An excluded input is neither
				// an adapter rejection nor successful archive processing.
				visits := 0
				var excluded []string
				terminal := file.IterateFilesListWithExclusions(root, []string{tc.name}, false, true, func(*file.Path) (bool, error) { return true, nil }, func(info *file.FileInfo) error {
					visits++
					_, e := d.ProcessFile(info, nil, &writer.Detectors{File: &bytes.Buffer{}})
					return e
				}, func(reason string) { excluded = append(excluded, reason) })
				t.Logf("iterator entry=%d exclusions=%v effective_depth=%d returned=%v", visits, excluded, depth, terminal)
				if tc.name == "valid.zip" || tc.name == "valid.gz" {
					if visits != 0 || len(excluded) != 1 || terminal != nil {
						t.Fatal("binary exclusion observation changed")
					}
				} else {
					if visits != 1 || len(excluded) != 0 {
						t.Fatal("iterator invocation missing")
					}
					if depth > 0 && !errors.Is(terminal, ErrArchiveExpansionUnsupported) {
						t.Fatalf("invoked adapter failure lost: %v", terminal)
					}
					if depth == 0 && terminal != nil {
						t.Fatal(terminal)
					}
				}
			})
		}
	}
	cfg := config.Config{Allowlists: []*config.Allowlist{{Paths: []*gre.Regexp{gre.MustCompile(`.*`)}}}}
	err = (&readerSource{Config: &cfg, Path: "allowed.zip", MaxArchiveDepth: 1}).Fragments(context.Background(), func(sources.Fragment, error) error { t.Fatal("yield on unsupported mode"); return nil })
	if !errors.Is(err, ErrArchiveExpansionUnsupported) {
		t.Fatalf("allowlist hid unsupported request: %v", err)
	}
}
