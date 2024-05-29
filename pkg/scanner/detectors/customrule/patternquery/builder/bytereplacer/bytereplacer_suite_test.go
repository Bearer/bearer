package bytereplacer_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestFilters(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ByteReplacer Suite")
}
