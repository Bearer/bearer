package javascript_test

import (
	_ "embed"
	"testing"

	"github.com/bearer/curio/new/detector/composition/javascript"
	"github.com/bearer/curio/new/detector/composition/testhelper"
)

//go:embed testdata/insecureURL.yml
var insecureURLRule []byte

//go:embed testdata/datatype.yml
var datatypeRule []byte

//go:embed testdata/class.yml
var classNameRule []byte

func TestJavascriptClass(t *testing.T) {
	testhelper.RunTest(t, map[string][]byte{
		"class_name": classNameRule,
	}, "testdata/testcases/class", javascript.New)
}

func TestJavascriptFlow(t *testing.T) {
	testhelper.RunTest(t, map[string][]byte{
		"logger": datatypeRule,
	}, "testdata/testcases/flow", javascript.New)
}

func TestJavascriptObject(t *testing.T) {
	testhelper.RunTest(t, map[string][]byte{
		"logger": datatypeRule,
	}, "testdata/testcases/object", javascript.New)
}

func TestJavascriptObjectMemberExpression(t *testing.T) {
	testhelper.RunTest(t, map[string][]byte{
		"logger": datatypeRule,
	}, "testdata/testcases/object-member-expression", javascript.New)
	// TODO: test case nested with ()
}

func TestJavascriptObjectSubscriptExpression(t *testing.T) {
	testhelper.RunTest(t, map[string][]byte{
		"logger": datatypeRule,
	}, "testdata/testcases/object-subscript-expression", javascript.New)
	// TODO test case nested with ()
}

func TestJavascriptString(t *testing.T) {
	testhelper.RunTest(t, map[string][]byte{
		"logger": insecureURLRule,
	}, "testdata/testcases/string", javascript.New)
}
