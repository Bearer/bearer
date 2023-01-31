package javascript_test

import (
	_ "embed"
	"testing"

	"github.com/bearer/curio/new/detector/composition/javascript"
	"github.com/bearer/curio/new/detector/composition/testhelper"
)

//go:embed testdata/console_log.yml
var consoleLogRule []byte

//go:embed testdata/class.yml
var classNameRule []byte

func TestJavascriptObject(t *testing.T) {
	testhelper.RunTest(t, map[string][]byte{
		"logger": consoleLogRule,
	}, "testdata/testcases/object", javascript.New)
}

func TestJavascriptFlow(t *testing.T) {
	testhelper.RunTest(t, map[string][]byte{
		"logger": consoleLogRule,
	}, "testdata/testcases/flow", javascript.New)
}

func TestJavascriptClass(t *testing.T) {
	testhelper.RunTest(t, map[string][]byte{
		"class_name": classNameRule,
	}, "testdata/testcases/class", javascript.New)
}
