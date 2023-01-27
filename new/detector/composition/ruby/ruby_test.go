package ruby_test

import (
	_ "embed"
	"testing"

	"github.com/bearer/curio/new/detector/composition/ruby"
	"github.com/bearer/curio/new/detector/composition/testhelper"
)

//go:embed testdata/rule.yml
var loggerRule []byte

func TestRuby(t *testing.T) {
	testhelper.RunTest(t, map[string][]byte{
		"logger": loggerRule,
	}, "./testdata/testcases", ruby.New)
}
