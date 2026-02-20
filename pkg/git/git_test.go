package git_test

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bearer/bearer/pkg/git"
	"github.com/bearer/bearer/pkg/util/file"
)

func TestGetRoot(t *testing.T) {
	filename := "foo.txt"
	dirname := "stuff"

	setup := func(t *testing.T) string {
		t.Helper()

		tempDir, err := os.MkdirTemp("", "diff-test")
		require.NoError(t, err)
		tempDir, err = file.CanonicalPath(tempDir)
		require.NoError(t, err)

		writeFile(t, tempDir, filename, "42")
		require.NoError(t, os.Mkdir(path.Join(tempDir, dirname), 0700))

		t.Cleanup(func() {
			_ = os.RemoveAll(tempDir)
		})

		return tempDir
	}

	t.Run("when the target path is in a git repository", func(t *testing.T) {
		t.Run("and the target path is the repository root", func(t *testing.T) {
			tempDir := setup(t)
			runGit(t, tempDir, "init", ".")

			result, err := git.GetRoot(tempDir)
			require.NoError(t, err)
			assert.Equal(t, tempDir, result)
		})

		t.Run("and the target path is a file", func(t *testing.T) {
			tempDir := setup(t)
			runGit(t, tempDir, "init", ".")

			result, err := git.GetRoot(path.Join(tempDir, filename))
			require.NoError(t, err)
			assert.Equal(t, tempDir, result)
		})

		t.Run("and the target path is in a subfolder", func(t *testing.T) {
			tempDir := setup(t)
			runGit(t, tempDir, "init", ".")

			result, err := git.GetRoot(path.Join(tempDir, dirname))
			require.NoError(t, err)
			assert.Equal(t, tempDir, result)
		})
	})

	t.Run("when the target path is NOT in a git repository", func(t *testing.T) {
		tempDir := setup(t)

		result, err := git.GetRoot(tempDir)
		require.NoError(t, err)
		assert.Empty(t, result)
	})
}
