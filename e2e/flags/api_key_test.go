package flags_test

import (
	"path/filepath"
	"testing"

	"github.com/bearer/bearer/e2e/internal/testhelper"
)

func TestApiKeyFlags(t *testing.T) {
	t.Parallel()
	arguments := []string{
		"scan",
		filepath.Join("e2e", "flags", "testdata", "simple"),
		"--disable-version-check",
		"--disable-default-rules",
		"--api-key",
		"123",
		"--format",
		"json",
	}
	tests := []testhelper.TestCase{
		testhelper.NewTestCase("bad-api-key-with-stderr", arguments, testhelper.TestCaseOptions{DisplayStdErr: true, IgnoreForce: false}),
		testhelper.NewTestCase("bad-api-key", arguments, testhelper.TestCaseOptions{DisplayStdErr: false, IgnoreForce: false}),
	}

	testhelper.RunTests(t, tests)
}
