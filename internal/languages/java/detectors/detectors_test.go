package detectors_test

import (
	"testing"

	"github.com/bearer/bearer/internal/languages/java"
	"github.com/bearer/bearer/internal/scanner/detectors/testhelper"
)

func TestJavaObjects(t *testing.T) {
	runTest(t, "object_class", "object", "testdata/class.java")
	runTest(t, "object_no_class", "object", "testdata/no_class.java")
}

func TestJavaString(t *testing.T) {
	runTest(t, "string", "string", "testdata/string.java")
}

func runTest(t *testing.T, name, detectorType, fileName string) {
	testhelper.RunTest(t, name, java.Get(), detectorType, fileName)
}