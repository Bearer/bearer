package detectors_test

import (
	"testing"

	"github.com/bearer/bearer/internal/languages/php"
	"github.com/bearer/bearer/internal/scanner/detectors/testhelper"
)

func TestJavaObjects(t *testing.T) {
	runTest(t, "object_class", "object", "testdata/class.php")
	runTest(t, "object_no_class", "object", "testdata/no_class.php")
}

func TestJavaString(t *testing.T) {
	runTest(t, "string", "string", "testdata/string.php")
}

func runTest(t *testing.T, name, detectorType, fileName string) {
	testhelper.RunTest(t, name, php.Get(), detectorType, fileName)
}
