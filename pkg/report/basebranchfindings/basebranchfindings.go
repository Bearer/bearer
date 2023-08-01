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
