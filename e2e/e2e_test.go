package e2e_test

import (
	"testing"

	"github.com/bearer/bearer/e2e/internal/testhelper"
	"github.com/bradleyjkemp/cupaloy"
	"github.com/stretchr/testify/assert"
)

func TestCache(t *testing.T) {
	arguments := []string{
		"scan",
		"e2e/testdata/logger",
		"--format", "yaml",
		"--disable-version-check",
		"--disable-default-rules",
		"--external-rule-dir", "e2e/testdata/rules",
		"--exit-code=0",
	}

	noCacheStdOut, noCacheStdErr := testhelper.ExecuteTest(
		testhelper.NewTestCase(
			"no_cache",
			arguments,
			testhelper.TestCaseOptions{
				DisplayStdErr: true,
			},
		),
		t,
	)

	withCacheStdOut, withCacheStdErr := testhelper.ExecuteTest(
		testhelper.NewTestCase(
			"with_cache",
			arguments,
			testhelper.TestCaseOptions{
				DisplayStdErr: true,
				IgnoreForce:   true,
			},
		),
		t,
	)

	assert.NotContains(t, noCacheStdErr, "Cached data used")
	assert.Contains(t, withCacheStdErr, "Cached data used")
	assert.Equal(t, noCacheStdOut, withCacheStdOut)
	cupaloy.SnapshotT(t, withCacheStdOut)
}
