package bytereplacer_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/bearer/bearer/internal/scanner/detectors/customrule/patternquery/builder/bytereplacer"
)

var _ = Describe("Replacer", func() {
	var replacer *bytereplacer.Replacer

	BeforeEach(func(ctx SpecContext) {
		replacer = bytereplacer.New([]byte("hello world"))
	})

	Describe("Replace", func() {
		When("replacements are made in sequential order", func() {
			BeforeEach(func(ctx SpecContext) {
				Expect(replacer.Replace(0, 1, nil)).To(Succeed())
			})

			It("does not fail", func(ctx SpecContext) {
				Expect(replacer.Replace(1, 2, []byte("foo"))).To(Succeed())
			})
		})

		When("replacements are made out of order", func() {
			BeforeEach(func(ctx SpecContext) {
				Expect(replacer.Replace(0, 2, nil)).To(Succeed())
			})

			It("returns an error", func(ctx SpecContext) {
				Expect(replacer.Replace(1, 2, []byte("foo"))).To(
					MatchError(ContainSubstring("replacements must be made in sequential order")),
				)
			})
		})
	})
})

var _ = Describe("Result", func() {
	var replacer *bytereplacer.Replacer
	original := []byte("hello world")

	BeforeEach(func(ctx SpecContext) {
		replacer = bytereplacer.New(original)
	})

	When("no replacements are made", func() {
		var result *bytereplacer.Result

		BeforeEach(func(ctx SpecContext) {
			result = replacer.Done()
		})

		Describe("Changed", func() {
			It("returns false", func(ctx SpecContext) {
				Expect(result.Changed()).To(BeFalse())
			})
		})

		Describe("Value", func() {
			It("returns the original value", func(ctx SpecContext) {
				Expect(result.Value()).To(Equal(original))
			})
		})

		Describe("Translate", func() {
			It("returns the original offset", func(ctx SpecContext) {
				Expect(result.Translate(0)).To(Equal(0))
				Expect(result.Translate(5)).To(Equal(5))
				Expect(result.Translate(10)).To(Equal(10))
			})
		})
	})

	When("noop replacements are made", func() {
		var result *bytereplacer.Result

		BeforeEach(func(ctx SpecContext) {
			replacer.Replace(0, 5, []byte("hello")) // nolint:errcheck
			replacer.Replace(6, 6, nil)             // nolint:errcheck
			result = replacer.Done()
		})

		Describe("Changed", func() {
			It("returns false", func(ctx SpecContext) {
				Expect(result.Changed()).To(BeFalse())
			})
		})

		Describe("Value", func() {
			It("returns the original value", func(ctx SpecContext) {
				Expect(result.Value()).To(Equal(original))
			})
		})

		Describe("Translate", func() {
			It("returns the original offset", func(ctx SpecContext) {
				Expect(result.Translate(0)).To(Equal(0))
				Expect(result.Translate(5)).To(Equal(5))
				Expect(result.Translate(10)).To(Equal(10))
			})
		})
	})

	When("replacements are made", func() {
		var result *bytereplacer.Result

		BeforeEach(func(ctx SpecContext) {
			replacer.Replace(0, 5, []byte("hi"))          // nolint:errcheck
			replacer.Replace(5, 5, []byte("!"))           // nolint:errcheck
			replacer.Replace(6, 11, []byte("testing123")) // nolint:errcheck
			result = replacer.Done()
		})

		Describe("Changed", func() {
			It("returns true", func(ctx SpecContext) {
				Expect(result.Changed()).To(BeTrue())
			})
		})

		Describe("Value", func() {
			It("returns the updated value", func(ctx SpecContext) {
				Expect(result.Value()).To(Equal([]byte("hi! testing123")))
			})
		})

		Describe("Translate", func() {
			It("returns the new offset", func(ctx SpecContext) {
				Expect(result.Translate(0)).To(Equal(0))
				Expect(result.Translate(5)).To(Equal(3))
				Expect(result.Translate(11)).To(Equal(14))
			})
		})
	})
})
