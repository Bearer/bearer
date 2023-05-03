package generic_test

import (
	"testing"

	"github.com/bearer/bearer/new/detector/composition/ruby"
	"github.com/bearer/bearer/new/detector/implementation/testhelper"
)

func TestDatatypeDetector(t *testing.T) {
	runTest(t, "datatype", "datatype", "testdata/datatype.rb")
}

func TestInsecureUrlDetector(t *testing.T) {
	runTest(t, "insecure_url", "insecure_url", "testdata/insecureurl.rb")
}

func runTest(t *testing.T, name string, detectorType, fileName string) {
	testhelper.RunTest(t, name, ruby.New, detectorType, fileName)
}
