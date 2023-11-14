package basebranchfindings

import (
	"slices"

	"github.com/bearer/bearer/internal/commands/process/filelist/files"
	"github.com/bearer/bearer/internal/git"
)

type key struct {
	RuleID   string
	Filename string
}

type Findings struct {
	fileList *files.List
	chunks   map[string]git.Chunks
	items    map[key][]git.ChunkRange
}

func New(fileList *files.List) *Findings {
	return &Findings{
		fileList: fileList,
		chunks:   make(map[string]git.Chunks),
		items:    make(map[key][]git.ChunkRange),
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

	findings.items[key] = append(
		findings.items[key],
		fileChunks.TranslateRange(newRange(baseStartLine, baseEndLine)),
	)
}

func (findings Findings) Consume(ruleID string, filename string, startLine, endLine int) bool {
	key := key{
		RuleID:   ruleID,
		Filename: filename,
	}

	lineRange := newRange(startLine, endLine)

	for i, findingLineRange := range findings.items[key] {
		if findingLineRange.Overlap(lineRange) {
			slices.Delete(findings.items[key], i, i+1)
			return true
		}
	}

	return false
}

func newRange(startLine, endLine int) git.ChunkRange {
	return git.ChunkRange{LineNumber: startLine, LineCount: endLine - startLine + 1}
}
