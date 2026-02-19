package bytereplacer_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bearer/bearer/pkg/scanner/detectors/customrule/patternquery/builder/bytereplacer"
)

func TestReplacer_Replace(t *testing.T) {
	t.Run("sequential order does not fail", func(t *testing.T) {
		replacer := bytereplacer.New([]byte("hello world"))
		require.NoError(t, replacer.Replace(0, 1, nil))
		assert.NoError(t, replacer.Replace(1, 2, []byte("foo")))
	})

	t.Run("out of order returns error", func(t *testing.T) {
		replacer := bytereplacer.New([]byte("hello world"))
		require.NoError(t, replacer.Replace(0, 2, nil))
		assert.ErrorContains(t, replacer.Replace(1, 2, []byte("foo")), "replacements must be made in sequential order")
	})
}

func TestResult_NoReplacements(t *testing.T) {
	replacer := bytereplacer.New([]byte("hello world"))
	result := replacer.Done()

	t.Run("Changed returns false", func(t *testing.T) {
		assert.False(t, result.Changed())
	})

	t.Run("Value returns original", func(t *testing.T) {
		assert.Equal(t, []byte("hello world"), result.Value())
	})

	t.Run("Translate returns original offset", func(t *testing.T) {
		assert.Equal(t, 0, result.Translate(0))
		assert.Equal(t, 5, result.Translate(5))
		assert.Equal(t, 10, result.Translate(10))
	})
}

func TestResult_NoopReplacements(t *testing.T) {
	replacer := bytereplacer.New([]byte("hello world"))
	replacer.Replace(0, 5, []byte("hello")) // nolint:errcheck
	replacer.Replace(6, 6, nil)             // nolint:errcheck
	result := replacer.Done()

	t.Run("Changed returns false", func(t *testing.T) {
		assert.False(t, result.Changed())
	})

	t.Run("Value returns original", func(t *testing.T) {
		assert.Equal(t, []byte("hello world"), result.Value())
	})

	t.Run("Translate returns original offset", func(t *testing.T) {
		assert.Equal(t, 0, result.Translate(0))
		assert.Equal(t, 5, result.Translate(5))
		assert.Equal(t, 10, result.Translate(10))
	})
}

func TestResult_WithReplacements(t *testing.T) {
	replacer := bytereplacer.New([]byte("hello world"))
	replacer.Replace(0, 5, []byte("hi"))          // nolint:errcheck
	replacer.Replace(5, 5, []byte("!"))           // nolint:errcheck
	replacer.Replace(6, 11, []byte("testing123")) // nolint:errcheck
	result := replacer.Done()

	t.Run("Changed returns true", func(t *testing.T) {
		assert.True(t, result.Changed())
	})

	t.Run("Value returns updated value", func(t *testing.T) {
		assert.Equal(t, []byte("hi! testing123"), result.Value())
	})

	t.Run("Translate returns new offset", func(t *testing.T) {
		assert.Equal(t, 0, result.Translate(0))
		assert.Equal(t, 3, result.Translate(5))
		assert.Equal(t, 14, result.Translate(11))
	})
}
