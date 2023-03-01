package ruby_test

import (
	"testing"

	"github.com/bearer/bearer/new/detector/composition/javascript"
	"github.com/bearer/bearer/new/detector/implementation/testhelper"
)

func TestJavascriptStringDetector(t *testing.T) {
	runTest(t, "string_literal", "string", "testdata/string_literal.js")
	runTest(t, "string_non_literal", "string", "testdata/string_non_literal.js")
}

func runTest(t *testing.T, name, detectorType, fileName string) {
	testhelper.RunTest(t, name, javascript.New, detectorType, fileName)
}
