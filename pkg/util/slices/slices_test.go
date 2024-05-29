package slices_test

import (
	"slices"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	sliceutil "github.com/bearer/bearer/pkg/util/slices"
)

var _ = Describe("Except", func() {
	slice := []string{"a", "b", "b"}

	When("the slice contains the value", func() {
		It("returns a slice without any occurances of the value", func() {
			Expect(sliceutil.Except(slice, "b")).To(Equal([]string{"a"}))
		})

		It("leaves the original slice unchanged", func() {
			sliceutil.Except(slice, "b")

			Expect(slice).To(Equal([]string{"a", "b", "b"}))
		})
	})

	When("the slice does NOT contain the value", func() {
		It("returns a copy of the original slice", func() {
			new := sliceutil.Except(slice, "not-there")
			Expect(new).To(Equal(slice))

			new = slices.Delete(new, 0, 1)
			Expect(new).NotTo(Equal(slice))
		})
	})
})
