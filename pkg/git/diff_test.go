package git_test

import (
	"os"
	"path"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/bearer/bearer/pkg/git"
)

var _ = Describe("Diff", func() {
	var tempDir, baseSHA string
	filename := "foo.txt"

	BeforeEach(func() {
		var err error

		tempDir, err = os.MkdirTemp("", "diff-test")
		Expect(err).To(BeNil())

		runGit(tempDir, "init", ".")
		writeFile(tempDir, filename, "1\n2\n3")
		addAndCommit(tempDir)

		baseSHA, err = git.GetCurrentCommit(tempDir)
		Expect(err).To(BeNil())
		Expect(baseSHA).NotTo(BeEmpty())
	})

	AfterEach(func() {
		if tempDir != "" {
			Expect(os.RemoveAll(tempDir)).To(Succeed())
		}
	})

	When("a file was added", func() {
		BeforeEach(func() {
			writeFile(tempDir, "new.txt", "abc")
			addAndCommit(tempDir)
		})

		It("returns the expected result", func() {
			Expect(git.Diff(tempDir, baseSHA)).To(ConsistOf([]git.FilePatch{{
				ToPath: "new.txt",
				Chunks: []git.Chunk{{
					From: git.ChunkRange{LineNumber: 0, LineCount: 0},
					To:   git.ChunkRange{LineNumber: 1, LineCount: 1},
				}},
			}}))
		})
	})

	When("a file was removed", func() {
		BeforeEach(func() {
			Expect(os.Remove(path.Join(tempDir, filename))).To(Succeed())
			addAndCommit(tempDir)
		})

		It("returns the expected result", func() {
			Expect(git.Diff(tempDir, baseSHA)).To(ConsistOf([]git.FilePatch{{
				FromPath: filename,
				Chunks: []git.Chunk{{
					From: git.ChunkRange{LineNumber: 1, LineCount: 3},
					To:   git.ChunkRange{LineNumber: 0, LineCount: 0},
				}},
			}}))
		})
	})

	When("a file was renamed", func() {
		toPath := "to.txt"

		BeforeEach(func() {
			Expect(os.Rename(path.Join(tempDir, filename), path.Join(tempDir, toPath))).To(Succeed())
			addAndCommit(tempDir)
		})

		It("returns the expected result", func() {
			Expect(git.Diff(tempDir, baseSHA)).To(ConsistOf([]git.FilePatch{
				{FromPath: filename, ToPath: toPath},
			}))
		})
	})

	When("paths contain characters requiring quoting", func() {
		fromPath := "from\t.txt"
		toPath := "to\t.txt"

		BeforeEach(func() {
			writeFile(tempDir, fromPath, "1\n2")
			addAndCommit(tempDir)

			var err error
			baseSHA, err = git.GetCurrentCommit(tempDir)
			Expect(err).To(BeNil())

			Expect(os.Rename(path.Join(tempDir, fromPath), path.Join(tempDir, toPath))).To(Succeed())
			addAndCommit(tempDir)
		})

		It("decodes the paths correctly", func() {
			Expect(git.Diff(tempDir, baseSHA)).To(ConsistOf([]git.FilePatch{
				{FromPath: fromPath, ToPath: toPath},
			}))
		})

		fromPath = "from bar.txt"
		toPath = "to foo.txt"

		It("decodes the paths correctly with whitespace in both from and to path", func() {
			Expect(git.Diff(tempDir, baseSHA)).To(ConsistOf([]git.FilePatch{
				{FromPath: fromPath, ToPath: toPath},
			}))
		})

		fromPath = "from bar.txt"
		toPath = "to.txt"

		It("decodes the paths correctly with whitespace in from path", func() {
			Expect(git.Diff(tempDir, baseSHA)).To(ConsistOf([]git.FilePatch{
				{FromPath: fromPath, ToPath: toPath},
			}))
		})
	})

	When("a file contains changes", func() {
		BeforeEach(func() {
			writeFile(tempDir, filename, "x\ny\n2\nd")
			addAndCommit(tempDir)
		})

		It("returns the expected result", func() {
			Expect(git.Diff(tempDir, baseSHA)).To(ConsistOf([]git.FilePatch{{
				FromPath: filename,
				ToPath:   filename,
				Chunks: []git.Chunk{
					{
						From: git.ChunkRange{LineNumber: 1, LineCount: 1},
						To:   git.ChunkRange{LineNumber: 1, LineCount: 2},
					},
					{
						From: git.ChunkRange{LineNumber: 3, LineCount: 1},
						To:   git.ChunkRange{LineNumber: 4, LineCount: 1},
					},
				},
			}}))
		})
	})

	When("a file contains a single line change (line count omitted from diff output)", func() {
		BeforeEach(func() {
			writeFile(tempDir, filename, "x\n2\n3")
			addAndCommit(tempDir)
			runGit(tempDir, "diff", "--unified=0", baseSHA)
		})

		It("returns the correct line counts", func() {
			Expect(git.Diff(tempDir, baseSHA)).To(ConsistOf([]git.FilePatch{{
				FromPath: filename,
				ToPath:   filename,
				Chunks: []git.Chunk{{
					From: git.ChunkRange{LineNumber: 1, LineCount: 1},
					To:   git.ChunkRange{LineNumber: 1, LineCount: 1},
				}},
			}}))
		})
	})
})

var _ = Describe("ChunkRange", func() {
	Describe("StartLineNumber", func() {
		It("returns the line number", func() {
			Expect(git.ChunkRange{LineNumber: 2, LineCount: 1}.StartLineNumber()).To(Equal(2))
		})

		When("there are no lines in the range", func() {
			It("returns the next line after the line number", func() {
				Expect(git.ChunkRange{LineNumber: 2, LineCount: 0}.StartLineNumber()).To(Equal(3))
			})
		})
	})

	Describe("EndLineNumber", func() {
		It("returns the (inclusive) end line number", func() {
			Expect(git.ChunkRange{LineNumber: 2, LineCount: 1}.EndLineNumber()).To(Equal(2))
			Expect(git.ChunkRange{LineNumber: 3, LineCount: 2}.EndLineNumber()).To(Equal(4))
		})

		When("there are no lines in the range", func() {
			It("returns the start line number", func() {
				Expect(git.ChunkRange{LineNumber: 2, LineCount: 0}.EndLineNumber()).To(Equal(2))
			})
		})
	})

	Describe("Overlap", func() {
		When("the ranges are equal", func() {
			a := git.ChunkRange{LineNumber: 1, LineCount: 2}
			b := git.ChunkRange{LineNumber: 1, LineCount: 2}

			It("is true", func() {
				Expect(a.Overlap(b)).To(BeTrue())
			})
		})

		When("B overlaps A's start", func() {
			a := git.ChunkRange{LineNumber: 2, LineCount: 2}
			b := git.ChunkRange{LineNumber: 1, LineCount: 2}

			It("is true", func() {
				Expect(a.Overlap(b)).To(BeTrue())
			})
		})

		When("B overlaps A's end", func() {
			a := git.ChunkRange{LineNumber: 1, LineCount: 2}
			b := git.ChunkRange{LineNumber: 2, LineCount: 2}

			It("is true", func() {
				Expect(a.Overlap(b)).To(BeTrue())
			})
		})

		When("B is before A", func() {
			a := git.ChunkRange{LineNumber: 2, LineCount: 2}
			b := git.ChunkRange{LineNumber: 1, LineCount: 1}

			It("is false", func() {
				Expect(a.Overlap(b)).To(BeFalse())
			})
		})

		When("B is after A", func() {
			a := git.ChunkRange{LineNumber: 1, LineCount: 2}
			b := git.ChunkRange{LineNumber: 3, LineCount: 2}

			It("is false", func() {
				Expect(a.Overlap(b)).To(BeFalse())
			})
		})
	})
})

var _ = Describe("Chunks", func() {
	Describe("TranslateRange", func() {
		When("the base range is preceded by an add chunk", func() {
			chunks := git.Chunks{{
				From: git.ChunkRange{LineNumber: 0, LineCount: 0},
				To:   git.ChunkRange{LineNumber: 1, LineCount: 1},
			}}

			It("returns a range shifted by the add", func() {
				Expect(chunks.TranslateRange(git.ChunkRange{LineNumber: 1, LineCount: 2})).To(
					Equal(git.ChunkRange{LineNumber: 2, LineCount: 2}),
				)
			})
		})

		When("the base range is at an add chunk", func() {
			chunks := git.Chunks{{
				From: git.ChunkRange{LineNumber: 1, LineCount: 0},
				To:   git.ChunkRange{LineNumber: 2, LineCount: 1},
			}}

			It("returns a range shifted by the add", func() {
				Expect(chunks.TranslateRange(git.ChunkRange{LineNumber: 2, LineCount: 2})).To(
					Equal(git.ChunkRange{LineNumber: 3, LineCount: 2}),
				)
			})
		})

		When("the base range surrounds a remove chunk", func() {
			chunks := git.Chunks{{
				From: git.ChunkRange{LineNumber: 3, LineCount: 1},
				To:   git.ChunkRange{LineNumber: 2, LineCount: 0},
			}}

			It("returns a range that still overlaps the unchanged portion by the same amount", func() {
				Expect(chunks.TranslateRange(git.ChunkRange{LineNumber: 2, LineCount: 3})).To(
					Equal(git.ChunkRange{LineNumber: 2, LineCount: 2}),
				)
			})
		})

		When("the base range overlaps the start of a remove chunk", func() {
			chunks := git.Chunks{{
				From: git.ChunkRange{LineNumber: 3, LineCount: 2},
				To:   git.ChunkRange{LineNumber: 2, LineCount: 0},
			}}

			It("returns a range that ends at the removed chunk", func() {
				Expect(chunks.TranslateRange(git.ChunkRange{LineNumber: 2, LineCount: 2})).To(
					Equal(git.ChunkRange{LineNumber: 2, LineCount: 1}),
				)
			})
		})

		When("the base range is inside a remove chunk", func() {
			chunks := git.Chunks{{
				From: git.ChunkRange{LineNumber: 1, LineCount: 2},
				To:   git.ChunkRange{LineNumber: 0, LineCount: 0},
			}}

			It("returns an invalid range (will be ignored)", func() {
				Expect(chunks.TranslateRange(git.ChunkRange{LineNumber: 1, LineCount: 1})).To(
					Equal(git.ChunkRange{LineNumber: 0, LineCount: 0}),
				)
			})
		})

		When("the base range overlaps the start of an edit chunk", func() {
			chunks := git.Chunks{{
				From: git.ChunkRange{LineNumber: 2, LineCount: 1},
				To:   git.ChunkRange{LineNumber: 2, LineCount: 2},
			}}

			It("expands the range to the end of the chunk", func() {
				Expect(chunks.TranslateRange(git.ChunkRange{LineNumber: 1, LineCount: 2})).To(
					Equal(git.ChunkRange{LineNumber: 1, LineCount: 3}),
				)
			})
		})

		When("the base range overlaps the end of an edit chunk", func() {
			chunks := git.Chunks{{
				From: git.ChunkRange{LineNumber: 1, LineCount: 2},
				To:   git.ChunkRange{LineNumber: 1, LineCount: 3},
			}}

			It("expands the range to the start of the chunk", func() {
				Expect(chunks.TranslateRange(git.ChunkRange{LineNumber: 2, LineCount: 2})).To(
					Equal(git.ChunkRange{LineNumber: 1, LineCount: 4}),
				)
			})
		})
	})
})
