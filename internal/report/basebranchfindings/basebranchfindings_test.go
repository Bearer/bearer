package basebranchfindings_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/bearer/bearer/internal/report/basebranchfindings"
	"github.com/bearer/bearer/internal/report/basebranchfindings/types"
)

func TestLineRangeOverlap(t *testing.T) {
	tests := []struct {
		Name string
		A,
		B types.LineRange
		Expected bool
	}{
		{
			"equal",
			types.LineRange{Start: 1, End: 2},
			types.LineRange{Start: 1, End: 2},
			true,
		},
		{
			"B overlap A start",
			types.LineRange{Start: 2, End: 3},
			types.LineRange{Start: 1, End: 2},
			true,
		},
		{
			"B overlap A end",
			types.LineRange{Start: 1, End: 2},
			types.LineRange{Start: 2, End: 3},
			true,
		},
		{
			"B before A",
			types.LineRange{Start: 2, End: 3},
			types.LineRange{Start: 1, End: 1},
			false,
		},
		{
			"B after A",
			types.LineRange{Start: 1, End: 2},
			types.LineRange{Start: 3, End: 4},
			false,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.Name, func(t *testing.T) {
			assert.Equal(t, testCase.Expected, testCase.A.Overlap(testCase.B))
		})
	}
}

func TestChunksTranslateRange(t *testing.T) {
	chunks := basebranchfindings.NewChunks()
	chunks.Add(types.ChunkAdd, 1)
	chunks.Add(types.ChunkEqual, 2)
	chunks.Add(types.ChunkRemove, 2)
	chunks.Add(types.ChunkEqual, 2)
	chunks.Add(types.ChunkRemove, 1)
	chunks.Add(types.ChunkAdd, 2)
	chunks.Add(types.ChunkRemove, 1)

	tests := []struct {
		Name string
		BaseLineRange,
		Expected types.LineRange
	}{
		{
			"equal after add",
			types.LineRange{Start: 1, End: 2},
			types.LineRange{Start: 2, End: 3},
		},
		{
			"remove inbetween equals",
			types.LineRange{Start: 2, End: 5},
			types.LineRange{Start: 3, End: 4},
		},
		{
			"equal overlapping remove",
			types.LineRange{Start: 2, End: 3},
			types.LineRange{Start: 3, End: 3},
		},
		{
			"in remove",
			types.LineRange{Start: 3, End: 3},
			types.LineRange{Start: 4, End: 3},
		},
		{
			"in remove followed by add",
			types.LineRange{Start: 7, End: 7},
			types.LineRange{Start: 6, End: 7},
		},
		{
			"add between removes",
			types.LineRange{Start: 7, End: 8},
			types.LineRange{Start: 6, End: 7},
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.Name, func(t *testing.T) {
			actual := chunks.TranslateRange(testCase.BaseLineRange.Start, testCase.BaseLineRange.End)
			assert.Equal(t, testCase.Expected, actual)
		})
	}
}
