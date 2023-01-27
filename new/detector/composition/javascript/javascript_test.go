package javascript_test

import (
	_ "embed"
	"testing"

	"github.com/bearer/curio/new/detector/composition/javascript"
	"github.com/bearer/curio/new/detector/composition/testhelper"
)

//go:embed testdata/rule.yml
var consoleLogRule []byte

func TestJavascript(t *testing.T) {

	testhelper.RunTest(t, map[string][]byte{
		"logger": consoleLogRule,
	}, "./testdata/testcases", javascript.New)
}
