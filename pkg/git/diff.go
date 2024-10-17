package git

import (
	"context"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/bearer/bearer/pkg/util/linescanner"
)

type ChunkRange struct {
	LineNumber,
	LineCount int
}

type Chunks []Chunk

type Chunk struct {
	From ChunkRange
	To   ChunkRange
}

type FilePatch struct {
	FromPath,
	ToPath string
	Chunks Chunks
}

const SRC_PREFIX = "caf67b1d-c9b7-44ba-9b85-bd94727547af"
const DST_PREFIX = "7dcfe536-96a9-49c5-8b86-7caa3d9e5e9c"

func Diff(rootDir, baseRef string) ([]FilePatch, error) {
	var result []FilePatch

	err := captureCommand(
		context.TODO(),
		rootDir,
		[]string{
			"diff",
			"--unified=0",
			"--first-parent",
			"--find-renames",
			"--break-rewrites",
			"--src-prefix=" + SRC_PREFIX,
			"--dst-prefix=" + DST_PREFIX,
			"--no-color",
			baseRef,
			"--",
		},
		func(stdout io.Reader) error {
			var err error
			result, err = parseDiff(linescanner.New(stdout))
			if err != nil {
				return err
			}

			return nil
		},
	)

	return result, err
}

func parseDiff(scanner *linescanner.Scanner) ([]FilePatch, error) {
	var result []FilePatch
	var fromPath, toPath string
	var chunks []Chunk

	flush := func() {
		if fromPath == "" && toPath == "" {
			return
		}

		result = append(result, FilePatch{
			FromPath: fromPath,
			ToPath:   toPath,
			Chunks:   chunks,
		})

		fromPath = ""
		toPath = ""
		chunks = nil
	}

	for scanner.Scan() {
		line := scanner.Text()

		var err error
		switch {
		case strings.HasPrefix(line, "diff --git"):
			flush()

			fromPath, toPath, err = parseDiffHeader(line)
			if err != nil {
				return nil, err
			}
		case strings.HasPrefix(line, "new file"):
			fromPath = ""
		case strings.HasPrefix(line, "deleted file"):
			toPath = ""
		case strings.HasPrefix(line, "@@"):
			chunk, err := parseChunkHeader(line)
			if err != nil {
				return nil, err
			}

			chunks = append(chunks, chunk)
		}

	}

	flush()

	return result, scanner.Err()
}

func parseDiffHeader(value string) (string, string, error) {
	rawFromPath, rawToPath := splitPaths(value)
	fromPath, err := unquoteFilename(rawFromPath)
	if err != nil {
		return "", "", fmt.Errorf("error parsing header 'from' path: %w", err)
	}

	toPath, err := unquoteFilename(rawToPath)
	if err != nil {
		return "", "", fmt.Errorf("error parsing header 'to' path: %w", err)
	}

	return fromPath, toPath, nil
}

func splitPaths(value string) (fromPath string, toPath string) {
	split1 := strings.Split(value, SRC_PREFIX)[1]
	split2 := strings.Split(split1, DST_PREFIX)

	fromPath = strings.TrimSpace(split2[0])
	toPath = strings.TrimSpace(split2[1])

	// handle trailing whitespaces and missing quotation marks
	// e.g. `foo\\t.txt \"`
	if strings.Contains(fromPath, "\"") {
		fromPath = "\"" + strings.TrimSpace(strings.ReplaceAll(fromPath, "\"", "")) + "\""
	}

	if strings.Contains(toPath, "\"") {
		toPath = "\"" + strings.ReplaceAll(toPath, "\"", "") + "\""
	}

	return fromPath, toPath
}

func parseChunkHeader(value string) (Chunk, error) {
	parts := strings.Split(value, " ")

	fromRange, err := parseRange(parts[1])
	if err != nil {
		return Chunk{}, fmt.Errorf("failed to parse chunk 'from' range: %w", err)
	}

	toRange, err := parseRange(parts[2])
	if err != nil {
		return Chunk{}, fmt.Errorf("failed to parse chunk 'to' range: %w", err)
	}

	return Chunk{From: fromRange, To: toRange}, nil
}

func parseRange(value string) (ChunkRange, error) {
	parts := strings.Split(value[1:], ",")

	lineNumber, err := strconv.Atoi(parts[0])
	if err != nil {
		return ChunkRange{}, fmt.Errorf("error decoding line number: %w", err)
	}

	count := 1
	if len(parts) > 1 {
		var err error
		count, err = strconv.Atoi(parts[1])
		if err != nil {
			return ChunkRange{}, fmt.Errorf("error decoding line count: %w", err)
		}
	}

	return ChunkRange{LineNumber: lineNumber, LineCount: count}, nil
}

func (chunks Chunks) TranslateRange(baseRange ChunkRange) ChunkRange {
	baseStartLine := baseRange.LineNumber
	startLine := baseStartLine
	if startChunk := chunks.getClosestChunk(baseStartLine); startChunk != nil {
		if baseStartLine > startChunk.From.EndLineNumber() {
			startLine = baseStartLine + startChunk.EndDelta()
		} else {
			startLine = startChunk.To.LineNumber
		}
	}

	baseEndLine := baseRange.EndLineNumber()
	endLine := baseEndLine
	if endChunk := chunks.getClosestChunk(baseEndLine); endChunk != nil {
		if baseEndLine > endChunk.From.EndLineNumber() {
			endLine = baseEndLine + endChunk.EndDelta()
		} else {
			endLine = endChunk.To.EndLineNumber()
		}
	}

	lineCount := endLine - startLine + 1
	if endLine == 0 {
		lineCount = 0
	}

	return ChunkRange{LineNumber: startLine, LineCount: lineCount}
}

func (chunks Chunks) getClosestChunk(baseLineNumber int) *Chunk {
	var result *Chunk

	for i, chunk := range chunks {
		if chunk.From.StartLineNumber() > baseLineNumber {
			break
		}

		result = &chunks[i]
	}

	return result
}

func (chunk Chunk) EndDelta() int {
	return chunk.To.EndLineNumber() - chunk.From.EndLineNumber()
}

func (chunkRange ChunkRange) StartLineNumber() int {
	if chunkRange.LineCount == 0 {
		return chunkRange.LineNumber + 1
	}

	return chunkRange.LineNumber
}

func (chunkRange ChunkRange) EndLineNumber() int {
	return chunkRange.StartLineNumber() + chunkRange.LineCount - 1
}

func (chunkRange ChunkRange) Overlap(other ChunkRange) bool {
	return chunkRange.LineNumber <= other.EndLineNumber() && chunkRange.EndLineNumber() >= other.LineNumber
}
