package basebranchfindings

import "github.com/bearer/bearer/pkg/report/basebranchfindings/types"

type chunk struct {
	ChangeType types.ChangeType
	LineCount  int
	HeadLine,
	BaseLine int
}

type Chunks struct {
	items []chunk
}

func NewChunks() types.Chunks {
	return &Chunks{}
}

func (chunks *Chunks) Add(changeType types.ChangeType, lineCount int) {
	baseLine := 1
	headLine := 1

	if len(chunks.items) != 0 {
		previousChunk := chunks.items[len(chunks.items)-1]
		baseLine, headLine = previousChunk.nextLine()
	}

	chunks.items = append(chunks.items, chunk{
		ChangeType: changeType,
		LineCount:  lineCount,
		BaseLine:   baseLine,
		HeadLine:   headLine,
	})
}

func (chunks *Chunks) TranslateRange(baseStartLine, baseEndLine int) types.LineRange {
	startLine := 0
	endLine := 0

	for i, chunk := range chunks.items {
		if chunk.ChangeType == types.ChunkAdd {
			continue
		}

		chunkBaseEndLine, chunkHeadEndLine := chunk.nextLine()

		if baseStartLine >= chunk.BaseLine && baseStartLine < chunkBaseEndLine {
			if chunk.ChangeType == types.ChunkEqual {
				startLine = baseStartLine - chunk.BaseLine + chunk.HeadLine
			} else {
				// for a removal, use the start of the next chunk
				startLine = chunkHeadEndLine
			}
		}

		if baseEndLine >= chunk.BaseLine && baseEndLine < chunkBaseEndLine {
			if chunk.ChangeType == types.ChunkEqual {
				endLine = baseEndLine - chunk.BaseLine + chunk.HeadLine
			} else {
				// for a removal, use the end of the previous chunk
				endLine = chunk.HeadLine - 1

				// but if there's an addition next then incorporate that, as it could
				// be an edit
				if i != len(chunks.items)-1 {
					nextChunk := chunks.items[i+1]
					if nextChunk.ChangeType == types.ChunkAdd {
						endLine = nextChunk.HeadLine + nextChunk.LineCount - 1
					}
				}
			}
		}

		if startLine != 0 && endLine != 0 {
			break
		}
	}

	return types.LineRange{
		Start: startLine,
		End:   endLine,
	}
}

func (chunk *chunk) nextLine() (int, int) {
	baseLine := chunk.BaseLine
	if chunk.ChangeType != types.ChunkAdd {
		baseLine += chunk.LineCount
	}

	headLine := chunk.HeadLine
	if chunk.ChangeType != types.ChunkRemove {
		headLine += chunk.LineCount
	}

	return baseLine, headLine
}
