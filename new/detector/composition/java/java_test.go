package java_test

import (
	_ "embed"
	"testing"

	"github.com/bearer/bearer/new/detector/composition/java"
	"github.com/bearer/bearer/new/detector/composition/testhelper"
)

//go:embed testdata/logger.yml
var loggerRule []byte

func TestFlow(t *testing.T) {
	testhelper.RunTest(t, map[string][]byte{
		"logger": loggerRule,
	}, "testdata/testcases/flow", java.New)
}
