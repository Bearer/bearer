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

//go:embed testdata/scope_rule.yml
var scopeRule []byte

func TestFlow(t *testing.T) {
	testhelper.GetRunner(t, datatypeRule, "Javascript").RunTest(t, "./testdata/testcases/flow", ".snapshots/flow/")
}

func TestObjectDeconstructing(t *testing.T) {
	testhelper.GetRunner(t, deconstructingRule, "Javascript").RunTest(t, "./testdata/testcases/object-deconstructing", ".snapshots/object-deconstructing/")
}

func TestString(t *testing.T) {
	testhelper.GetRunner(t, insecureURLRule, "Javascript").RunTest(t, "./testdata/testcases/string", ".snapshots/string/")
}

func TestScope(t *testing.T) {
	testhelper.GetRunner(t, scopeRule, "Javascript").RunTest(t, "./testdata/scope", ".snapshots/")
}
