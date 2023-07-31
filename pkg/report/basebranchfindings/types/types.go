package types

type ChangeType int

const (
	ChunkAdd ChangeType = iota
	ChunkRemove
	ChunkEqual
)

type LineRange struct {
	Start,
	End int
}

type Chunks interface {
	Add(changeType ChangeType, lineCount int)
	TranslateRange(baseStartLine, baseEndLine int) LineRange
}

func (lineRange *LineRange) Overlap(other LineRange) bool {
	return lineRange.Start <= other.End && lineRange.End >= other.Start
}
