package slices_test

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"

	sliceutil "github.com/bearer/bearer/pkg/util/slices"
)

func TestExcept(t *testing.T) {
	slice := []string{"a", "b", "b"}

	t.Run("when the slice contains the value", func(t *testing.T) {
		t.Run("returns a slice without any occurrences of the value", func(t *testing.T) {
			assert.Equal(t, []string{"a"}, sliceutil.Except(slice, "b"))
		})

		t.Run("leaves the original slice unchanged", func(t *testing.T) {
			sliceutil.Except(slice, "b")
			assert.Equal(t, []string{"a", "b", "b"}, slice)
		})
	})

	t.Run("when the slice does NOT contain the value", func(t *testing.T) {
		t.Run("returns a copy of the original slice", func(t *testing.T) {
			new := sliceutil.Except(slice, "not-there")
			assert.Equal(t, slice, new)

			new = slices.Delete(new, 0, 1)
			assert.NotEqual(t, slice, new)
		})
	})
}
