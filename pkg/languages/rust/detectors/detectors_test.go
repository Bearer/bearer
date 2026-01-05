package detectors_test

import (
	"testing"

	"github.com/bearer/bearer/pkg/languages/rust"
	"github.com/bearer/bearer/pkg/scanner/detectors/testhelper"
)

func TestRustObjects(t *testing.T) {
	runTest(t, "object_struct", "object", "testdata/struct.rs")
	runTest(t, "object_enum", "object", "testdata/enum.rs")
}

func TestRustString(t *testing.T) {
	runTest(t, "string_literal", "string", "testdata/string.rs")
	runTest(t, "string_raw", "string", "testdata/raw_string.rs")
}

func runTest(t *testing.T, name, detectorType, fileName string) {
	testhelper.RunTest(t, name, rust.Get(), detectorType, fileName)
}

