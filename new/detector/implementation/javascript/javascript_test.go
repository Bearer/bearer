package ruby_test

import (
	"testing"

	"github.com/bearer/bearer/new/detector/composition/javascript"
	"github.com/bearer/bearer/new/detector/implementation/testhelper"
)

func TestJavascriptObjectDetector(t *testing.T) {
	runTest(t, "object_class", "object", "testdata/object_class.js")
	runTest(t, "object_object", "object", "testdata/object_object.js")
	runTest(t, "object_spread", "object", "testdata/object_spread.js")
	runTest(t, "object_projection", "object", "testdata/object_projection.js")
}

func TestJavascriptStringDetector(t *testing.T) {
	runTest(t, "string_assign_eq", "string", "testdata/string_assign_eq.js")
	runTest(t, "string_literal", "string", "testdata/string_literal.js")
	runTest(t, "string_non_literal", "string", "testdata/string_non_literal.js")
}

func TestJavascriptTypes(t *testing.T) {
	runTest(t, "typed_object_const", "object", "testdata/typed_object_const.ts")
	runTest(t, "typed_object_let", "object", "testdata/typed_object_let.ts")
	runTest(t, "typed_object_var", "object", "testdata/typed_object_var.ts")
}

func TestJavascriptFileTypes(t *testing.T) {
	runTest(t, "file_type_tsx", "object", "testdata/file_type.tsx")
	runTest(t, "file_type_jsx", "object", "testdata/file_type.jsx")
}

func runTest(t *testing.T, name, detectorType, fileName string) {
	testhelper.RunTest(t, name, javascript.New, detectorType, fileName)
}
