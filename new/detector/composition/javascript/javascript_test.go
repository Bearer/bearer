package javascript_test

import (
	_ "embed"
	"testing"

	"github.com/bearer/bearer/new/detector/composition/testhelper"
)

//go:embed testdata/insecureURL.yml
var insecureURLRule []byte

//go:embed testdata/datatype.yml
var datatypeRule []byte

//go:embed testdata/deconstructing.yml
var deconstructingRule []byte

func TestFlow(t *testing.T) {
	t.Parallel()
	testhelper.GetRunner(t, datatypeRule, "Javascript").RunTest(t, "./testdata/testcases/flow", ".snapshots/flow/")
}

func TestObjectDeconstructing(t *testing.T) {
	t.Parallel()
	testhelper.GetRunner(t, deconstructingRule, "Javascript").RunTest(t, "./testdata/testcases/object-deconstructing", ".snapshots/object-deconstructing/")
}

func TestString(t *testing.T) {
	t.Parallel()
	testhelper.GetRunner(t, insecureURLRule, "Javascript").RunTest(t, "./testdata/testcases/string", ".snapshots/string/")
}
