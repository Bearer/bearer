package integration_test

import (
	"testing"

	"github.com/bearer/curio/integration/internal/testhelper"
)

func TestMetadataFlags(t *testing.T) {
	tests := []testhelper.TestCase{
		*testhelper.NewTestCase("help", []string{"help"}),
		*testhelper.NewTestCase("version", []string{"version"}),
	}

	testhelper.RunTests(t, tests)
}
