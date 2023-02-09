package ruby_test

import (
	"testing"

	"github.com/bearer/curio/new/detector/composition/ruby"
	"github.com/bearer/curio/new/detector/implementation/testhelper"
)

func TestRubyObjectDetector(t *testing.T) {
	runTest(t, "object_assignment", "object", "testdata/object_assignment.rb")
	runTest(t, "object_chain", "object", "testdata/object_chain.rb")
	runTest(t, "object_class", "object", "testdata/object_class.rb")
	runTest(t, "object_hash", "object", "testdata/object_hash.rb")
	runTest(t, "object_parent_pair", "object", "testdata/object_parent_pair.rb")
}

func TestRubyPropertyDetector(t *testing.T) {
	runTest(t, "property_accessor", "property", "testdata/property_accessor.rb")
	runTest(t, "property_method", "property", "testdata/property_method.rb")
	runTest(t, "property_pair", "property", "testdata/property_pair.rb")
}

func runTest(t *testing.T, name string, detectorType, fileName string) {
	testhelper.RunTest(t, name, ruby.New, detectorType, fileName)
}
