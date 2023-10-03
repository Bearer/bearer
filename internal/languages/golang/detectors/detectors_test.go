package detectors_test

import (
	"testing"

	"github.com/bearer/bearer/internal/languages/golang"
	"github.com/bearer/bearer/internal/scanner/detectors/testhelper"
)

func TestGoObjects(t *testing.T) {
	runTest(t, "object_class", "object", "testdata/class.go")
	runTest(t, "object_no_class", "object", "testdata/no_class.go")
}

func TestGoString(t *testing.T) {
	runTest(t, "string", "string", "testdata/string.go")
}

func runTest(t *testing.T, name, detectorType, fileName string) {
	testhelper.RunTest(t, name, golang.Get(), detectorType, fileName)
}
