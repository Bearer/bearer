package git_test

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bearer/bearer/pkg/git"
)

func setupDiffTest(t *testing.T) (string, string) {
	t.Helper()

	tempDir, err := os.MkdirTemp("", "diff-test")
	require.NoError(t, err)

	t.Cleanup(func() {
		_ = os.RemoveAll(tempDir)
	})

	runGit(t, tempDir, "init", ".")
	writeFile(t, tempDir, "foo.txt", "1\n2\n3")
	addAndCommit(t, tempDir)

	baseSHA, err := git.GetCurrentCommit(tempDir)
	require.NoError(t, err)
	require.NotEmpty(t, baseSHA)

	return tempDir, baseSHA
}

func TestDiff_FileAdded(t *testing.T) {
	tempDir, baseSHA := setupDiffTest(t)

	writeFile(t, tempDir, "new.txt", "abc")
	addAndCommit(t, tempDir)

	result, err := git.Diff(tempDir, baseSHA)
	require.NoError(t, err)
	assert.Equal(t, []git.FilePatch{{
		ToPath: "new.txt",
		Chunks: []git.Chunk{{
			From: git.ChunkRange{LineNumber: 0, LineCount: 0},
			To:   git.ChunkRange{LineNumber: 1, LineCount: 1},
		}},
	}}, result)
}

func TestDiff_FileRemoved(t *testing.T) {
	tempDir, baseSHA := setupDiffTest(t)

	require.NoError(t, os.Remove(path.Join(tempDir, "foo.txt")))
	addAndCommit(t, tempDir)

	result, err := git.Diff(tempDir, baseSHA)
	require.NoError(t, err)
	assert.Equal(t, []git.FilePatch{{
		FromPath: "foo.txt",
		Chunks: []git.Chunk{{
			From: git.ChunkRange{LineNumber: 1, LineCount: 3},
			To:   git.ChunkRange{LineNumber: 0, LineCount: 0},
		}},
	}}, result)
}

func TestDiff_FileRenamed(t *testing.T) {
	tempDir, baseSHA := setupDiffTest(t)

	require.NoError(t, os.Rename(path.Join(tempDir, "foo.txt"), path.Join(tempDir, "to.txt")))
	addAndCommit(t, tempDir)

	result, err := git.Diff(tempDir, baseSHA)
	require.NoError(t, err)
	assert.Equal(t, []git.FilePatch{
		{FromPath: "foo.txt", ToPath: "to.txt"},
	}, result)
}

func TestDiff_QuotedPathsWithTabs(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "diff-test")
	require.NoError(t, err)
	t.Cleanup(func() { _ = os.RemoveAll(tempDir) })

	runGit(t, tempDir, "init", ".")

	fromPath := "from\t.txt"
	toPath := "to\t.txt"

	writeFile(t, tempDir, fromPath, "1\n2")
	addAndCommit(t, tempDir)

	baseSHA, err := git.GetCurrentCommit(tempDir)
	require.NoError(t, err)

	require.NoError(t, os.Rename(path.Join(tempDir, fromPath), path.Join(tempDir, toPath)))
	addAndCommit(t, tempDir)

	result, err := git.Diff(tempDir, baseSHA)
	require.NoError(t, err)
	assert.Equal(t, []git.FilePatch{
		{FromPath: fromPath, ToPath: toPath},
	}, result)
}

func TestDiff_QuotedPathsWithSpacesInBoth(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "diff-test")
	require.NoError(t, err)
	t.Cleanup(func() { _ = os.RemoveAll(tempDir) })

	runGit(t, tempDir, "init", ".")

	fromPath := "from bar.txt"
	toPath := "to foo.txt"

	writeFile(t, tempDir, fromPath, "1\n2")
	addAndCommit(t, tempDir)

	baseSHA, err := git.GetCurrentCommit(tempDir)
	require.NoError(t, err)

	require.NoError(t, os.Rename(path.Join(tempDir, fromPath), path.Join(tempDir, toPath)))
	addAndCommit(t, tempDir)

	result, err := git.Diff(tempDir, baseSHA)
	require.NoError(t, err)
	assert.Equal(t, []git.FilePatch{
		{FromPath: fromPath, ToPath: toPath},
	}, result)
}

func TestDiff_QuotedPathsWithSpaceInFrom(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "diff-test")
	require.NoError(t, err)
	t.Cleanup(func() { _ = os.RemoveAll(tempDir) })

	runGit(t, tempDir, "init", ".")

	fromPath := "from bar.txt"
	toPath := "to.txt"

	writeFile(t, tempDir, fromPath, "1\n2")
	addAndCommit(t, tempDir)

	baseSHA, err := git.GetCurrentCommit(tempDir)
	require.NoError(t, err)

	require.NoError(t, os.Rename(path.Join(tempDir, fromPath), path.Join(tempDir, toPath)))
	addAndCommit(t, tempDir)

	result, err := git.Diff(tempDir, baseSHA)
	require.NoError(t, err)
	assert.Equal(t, []git.FilePatch{
		{FromPath: fromPath, ToPath: toPath},
	}, result)
}

func TestDiff_FileContainsChanges(t *testing.T) {
	tempDir, baseSHA := setupDiffTest(t)

	writeFile(t, tempDir, "foo.txt", "x\ny\n2\nd")
	addAndCommit(t, tempDir)

	result, err := git.Diff(tempDir, baseSHA)
	require.NoError(t, err)
	assert.Equal(t, []git.FilePatch{{
		FromPath: "foo.txt",
		ToPath:   "foo.txt",
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
	}}, result)
}

func TestDiff_SingleLineChange(t *testing.T) {
	tempDir, baseSHA := setupDiffTest(t)

	writeFile(t, tempDir, "foo.txt", "x\n2\n3")
	addAndCommit(t, tempDir)
	runGit(t, tempDir, "diff", "--unified=0", baseSHA)

	result, err := git.Diff(tempDir, baseSHA)
	require.NoError(t, err)
	assert.Equal(t, []git.FilePatch{{
		FromPath: "foo.txt",
		ToPath:   "foo.txt",
		Chunks: []git.Chunk{{
			From: git.ChunkRange{LineNumber: 1, LineCount: 1},
			To:   git.ChunkRange{LineNumber: 1, LineCount: 1},
		}},
	}}, result)
}

func TestChunkRange_StartLineNumber(t *testing.T) {
	t.Run("returns the line number", func(t *testing.T) {
		assert.Equal(t, 2, git.ChunkRange{LineNumber: 2, LineCount: 1}.StartLineNumber())
	})

	t.Run("when there are no lines in the range returns the next line", func(t *testing.T) {
		assert.Equal(t, 3, git.ChunkRange{LineNumber: 2, LineCount: 0}.StartLineNumber())
	})
}

func TestChunkRange_EndLineNumber(t *testing.T) {
	t.Run("returns the inclusive end line number", func(t *testing.T) {
		assert.Equal(t, 2, git.ChunkRange{LineNumber: 2, LineCount: 1}.EndLineNumber())
		assert.Equal(t, 4, git.ChunkRange{LineNumber: 3, LineCount: 2}.EndLineNumber())
	})

	t.Run("when there are no lines in the range returns the start line number", func(t *testing.T) {
		assert.Equal(t, 2, git.ChunkRange{LineNumber: 2, LineCount: 0}.EndLineNumber())
	})
}

func TestChunkRange_Overlap(t *testing.T) {
	t.Run("equal ranges overlap", func(t *testing.T) {
		a := git.ChunkRange{LineNumber: 1, LineCount: 2}
		b := git.ChunkRange{LineNumber: 1, LineCount: 2}
		assert.True(t, a.Overlap(b))
	})

	t.Run("B overlaps A start", func(t *testing.T) {
		a := git.ChunkRange{LineNumber: 2, LineCount: 2}
		b := git.ChunkRange{LineNumber: 1, LineCount: 2}
		assert.True(t, a.Overlap(b))
	})

	t.Run("B overlaps A end", func(t *testing.T) {
		a := git.ChunkRange{LineNumber: 1, LineCount: 2}
		b := git.ChunkRange{LineNumber: 2, LineCount: 2}
		assert.True(t, a.Overlap(b))
	})

	t.Run("B is before A", func(t *testing.T) {
		a := git.ChunkRange{LineNumber: 2, LineCount: 2}
		b := git.ChunkRange{LineNumber: 1, LineCount: 1}
		assert.False(t, a.Overlap(b))
	})

	t.Run("B is after A", func(t *testing.T) {
		a := git.ChunkRange{LineNumber: 1, LineCount: 2}
		b := git.ChunkRange{LineNumber: 3, LineCount: 2}
		assert.False(t, a.Overlap(b))
	})
}

func TestChunks_TranslateRange(t *testing.T) {
	t.Run("preceded by an add chunk shifts range", func(t *testing.T) {
		chunks := git.Chunks{{
			From: git.ChunkRange{LineNumber: 0, LineCount: 0},
			To:   git.ChunkRange{LineNumber: 1, LineCount: 1},
		}}
		assert.Equal(t,
			git.ChunkRange{LineNumber: 2, LineCount: 2},
			chunks.TranslateRange(git.ChunkRange{LineNumber: 1, LineCount: 2}),
		)
	})

	t.Run("at an add chunk shifts range", func(t *testing.T) {
		chunks := git.Chunks{{
			From: git.ChunkRange{LineNumber: 1, LineCount: 0},
			To:   git.ChunkRange{LineNumber: 2, LineCount: 1},
		}}
		assert.Equal(t,
			git.ChunkRange{LineNumber: 3, LineCount: 2},
			chunks.TranslateRange(git.ChunkRange{LineNumber: 2, LineCount: 2}),
		)
	})

	t.Run("surrounds a remove chunk overlaps unchanged portion", func(t *testing.T) {
		chunks := git.Chunks{{
			From: git.ChunkRange{LineNumber: 3, LineCount: 1},
			To:   git.ChunkRange{LineNumber: 2, LineCount: 0},
		}}
		assert.Equal(t,
			git.ChunkRange{LineNumber: 2, LineCount: 2},
			chunks.TranslateRange(git.ChunkRange{LineNumber: 2, LineCount: 3}),
		)
	})

	t.Run("overlaps start of remove chunk ends at removed chunk", func(t *testing.T) {
		chunks := git.Chunks{{
			From: git.ChunkRange{LineNumber: 3, LineCount: 2},
			To:   git.ChunkRange{LineNumber: 2, LineCount: 0},
		}}
		assert.Equal(t,
			git.ChunkRange{LineNumber: 2, LineCount: 1},
			chunks.TranslateRange(git.ChunkRange{LineNumber: 2, LineCount: 2}),
		)
	})

	t.Run("inside a remove chunk returns invalid range", func(t *testing.T) {
		chunks := git.Chunks{{
			From: git.ChunkRange{LineNumber: 1, LineCount: 2},
			To:   git.ChunkRange{LineNumber: 0, LineCount: 0},
		}}
		assert.Equal(t,
			git.ChunkRange{LineNumber: 0, LineCount: 0},
			chunks.TranslateRange(git.ChunkRange{LineNumber: 1, LineCount: 1}),
		)
	})

	t.Run("overlaps start of edit chunk expands to end", func(t *testing.T) {
		chunks := git.Chunks{{
			From: git.ChunkRange{LineNumber: 2, LineCount: 1},
			To:   git.ChunkRange{LineNumber: 2, LineCount: 2},
		}}
		assert.Equal(t,
			git.ChunkRange{LineNumber: 1, LineCount: 3},
			chunks.TranslateRange(git.ChunkRange{LineNumber: 1, LineCount: 2}),
		)
	})

	t.Run("overlaps end of edit chunk expands to start", func(t *testing.T) {
		chunks := git.Chunks{{
			From: git.ChunkRange{LineNumber: 1, LineCount: 2},
			To:   git.ChunkRange{LineNumber: 1, LineCount: 3},
		}}
		assert.Equal(t,
			git.ChunkRange{LineNumber: 1, LineCount: 4},
			chunks.TranslateRange(git.ChunkRange{LineNumber: 2, LineCount: 2}),
		)
	})
}
