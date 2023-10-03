package detectors_test

import (
	"testing"

	"github.com/bearer/bearer/internal/scanner/detectors/testhelper"
)

func Test<Language>Objects(t *testing.T) {
	runTest(t, "object_class", "object", "testdata/class.<language>")
	runTest(t, "object_no_class", "object", "testdata/no_class.<language>")
}

func Test<Language>String(t *testing.T) {
	runTest(t, "string", "string", "testdata/string.<language>")
}

func runTest(t *testing.T, name, detectorType, fileName string) {
	testhelper.RunTest(t, name, <language>.Get(), detectorType, fileName)
}
