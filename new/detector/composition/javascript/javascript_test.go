package javascript_test

import (
	_ "embed"
	"testing"

	"github.com/bearer/bearer/new/detector/composition/javascript"
	"github.com/bearer/bearer/new/detector/composition/testhelper"
)

//go:embed testdata/insecureURL.yml
var insecureURLRule []byte

//go:embed testdata/datatype.yml
var datatypeRule []byte

//go:embed testdata/deconstructing.yml
var deconstructingRule []byte

//go:embed testdata/class.yml
var classNameRule []byte

func TestClass(t *testing.T) {
	testhelper.RunTest(t, map[string][]byte{
		"class_name": classNameRule,
	}, "testdata/testcases/class", javascript.New)
}

func TestFlow(t *testing.T) {
	testhelper.RunTest(t, map[string][]byte{
		"logger": datatypeRule,
	}, "testdata/testcases/flow", javascript.New)
}

func TestObject(t *testing.T) {
	testhelper.RunTest(t, map[string][]byte{
		"logger": datatypeRule,
	}, "testdata/testcases/object", javascript.New)
}

func TestObjectMemberExpression(t *testing.T) {
	testhelper.RunTest(t, map[string][]byte{
		"logger": datatypeRule,
	}, "testdata/testcases/object-member-expression", javascript.New)
}

func TestObjectSubscriptExpression(t *testing.T) {
	testhelper.RunTest(t, map[string][]byte{
		"logger": datatypeRule,
	}, "testdata/testcases/object-subscript-expression", javascript.New)
}

func TestObjectDeconstructing(t *testing.T) {
	testhelper.RunTest(t, map[string][]byte{
		"logger": deconstructingRule,
	}, "testdata/testcases/object-deconstructing", javascript.New)
}

func TestString(t *testing.T) {
	testhelper.RunTest(t, map[string][]byte{
		"logger": insecureURLRule,
	}, "testdata/testcases/string", javascript.New)
}
