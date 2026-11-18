package ignore

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/bearer/bearer/pkg/commands/process/settings"
	flagtypes "github.com/bearer/bearer/pkg/flag/types"
	"github.com/hhatto/gocloc"
)

func TestCoverageDecisionReasons(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "skip.go")
	if err := os.WriteFile(path, []byte("package skip\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	checker := New(dir, settings.Config{Scan: flagtypes.ScanOptions{SkipPath: []string{"skip.go"}}})
	reason, err := checker.Decision(dir, path, &gocloc.Result{Files: map[string]*gocloc.ClocFile{path: {Lang: "Go"}}}, info)
	if err != nil {
		t.Fatal(err)
	}
	if reason != "skip-path" {
		t.Fatalf("reason = %q, want skip-path", reason)
	}
}

func TestCoverageDecisionConfiguredEligibilityAndInputErrors(t *testing.T) {
	dir := t.TempDir()
	if err := os.Mkdir(filepath.Join(dir, "src"), 0o700); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, "src", "plain.go")
	if err := os.WriteFile(path, []byte("package plain\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	checker := New(dir, settings.Config{Scan: flagtypes.ScanOptions{SkipTest: true}})
	reason, err := checker.Decision(dir, path, &gocloc.Result{Files: map[string]*gocloc.ClocFile{}}, info)
	if err != nil || reason != "" {
		t.Fatalf("plain decision = %q, %v", reason, err)
	}

	if err := os.WriteFile(filepath.Join(dir, "src", "thing_test.go"), []byte("package test\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	testPath := filepath.Join(dir, "src", "thing_test.go")
	testInfo, err := os.Stat(testPath)
	if err != nil {
		t.Fatal(err)
	}
	reason, err = checker.Decision(dir, testPath, &gocloc.Result{Files: map[string]*gocloc.ClocFile{}}, testInfo)
	if err != nil || reason != "test-name" {
		t.Fatalf("test decision = %q, %v", reason, err)
	}

	if _, err := checker.Decision(dir, filepath.Join(dir, "missing.go"), &gocloc.Result{Files: map[string]*gocloc.ClocFile{}}, info); err == nil {
		t.Fatal("missing input unexpectedly became an exclusion")
	}

	// A single-file target uses the target itself as projectPath; preserve that
	// existing entrypoint shape while applying the shared eligibility check.
	single := New(path, settings.Config{})
	reason, err = single.Decision(path, path, &gocloc.Result{Files: map[string]*gocloc.ClocFile{}}, info)
	if err != nil || reason != "" {
		t.Fatalf("single-file decision = %q, %v", reason, err)
	}
}
