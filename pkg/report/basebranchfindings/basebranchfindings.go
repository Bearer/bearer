package basebranchfindings

import (
	"github.com/bearer/bearer/pkg/commands/process/filelist/files"
	"github.com/bearer/bearer/pkg/report/basebranchfindings/types"
	"golang.org/x/exp/slices"
)

type key struct {
	RuleID   string
	Filename string
}

type Findings struct {
	fileList *files.List
	chunks   map[string][]chunk
	items    map[key][]types.LineRange
}

type chunk struct {
	ChangeType types.ChangeType
	LineCount  int
	HeadLine,
	BaseLine int
}

type Chunks struct {
	items []chunk
}

func New(fileList *files.List) *Findings {
	return &Findings{
		fileList: fileList,
		chunks:   make(map[string][]chunk),
		items:    make(map[key][]types.LineRange),
	}
}

func (findings Findings) Add(ruleID string, baseFilename string, baseStartLine, baseEndLine int) {
	filename := findings.fileList.Renames[baseFilename]
	if filename == "" {
		filename = baseFilename
	}

	fileChunks := findings.fileList.Chunks[filename]
	key := key{
		RuleID:   ruleID,
		Filename: filename,
	}

	findings.items[key] = append(findings.items[key], fileChunks.TranslateRange(baseStartLine, baseEndLine))
}

func (findings Findings) Consume(ruleID string, filename string, startLine, endLine int) bool {
	key := key{
		RuleID:   ruleID,
		Filename: filename,
	}

	lineRange := types.LineRange{Start: startLine, End: endLine}

	for i, findingLineRange := range findings.items[key] {
		if findingLineRange.Overlap(lineRange) {
			slices.Delete(findings.items[key], i, i+1)
			return true
		}
	}

	return false
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
