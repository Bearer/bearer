package git_test

import (
	"os"
	"path"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/bearer/bearer/internal/git"
	"github.com/bearer/bearer/internal/util/file"
)

var _ = Describe("GetRoot", func() {
	var tempDir string
	filename := "foo.txt"
	dirname := "stuff"

	BeforeEach(func() {
		var err error
		tempDir, err = os.MkdirTemp("", "diff-test")
		Expect(err).To(BeNil())
		tempDir, err = file.CanonicalPath(tempDir)
		Expect(err).To(BeNil())

		writeFile(tempDir, filename, "42")
		Expect(os.Mkdir(path.Join(tempDir, dirname), 0700)).To(Succeed())
	})

	AfterEach(func() {
		if tempDir != "" {
			Expect(os.RemoveAll(tempDir)).To(Succeed())
		}
	})

	When("the target path is in a git repository", func() {
		BeforeEach(func() {
			runGit(tempDir, "init", ".")
		})

		When("the target path is the repository root", func() {
			It("returns the root", func() {
				Expect(git.GetRoot(tempDir)).To(Equal(tempDir))
			})
		})

		When("the target path is a file", func() {
			It("returns the root", func() {
				Expect(git.GetRoot(path.Join(tempDir, filename))).To(Equal(tempDir))
			})
		})

		When("the target path is in a subfolder", func() {
			It("returns the root", func() {
				Expect(git.GetRoot(path.Join(tempDir, dirname))).To(Equal(tempDir))
			})
		})
	})

	When("the target path is NOT in a git repository", func() {
		It("returns an empty string", func() {
			Expect(git.GetRoot(tempDir)).To(BeEmpty())
		})
	})
})
