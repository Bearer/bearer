package source

import (
	"github.com/bearer/bearer/pkg/util/file"
)

// Source represents a part of a source file that is referenced in the scan report.
type Source struct {
	Filename          string  `json:"filename" yaml:"filename"`
	FullFilename      string  `json:"full_filename" yaml:"full_filename"`
	Language          string  `json:"language" yaml:"language"`
	LanguageType      string  `json:"language_type" yaml:"language_type"`
	StartLineNumber   *int    `json:"start_line_number" yaml:"start_line_number"`
	StartColumnNumber *int    `json:"start_column_number" yaml:"start_column_number"`
	EndLineNumber     *int    `json:"end_line_number" yaml:"end_line_number"`
	EndColumnNumber   *int    `json:"end_column_number" yaml:"end_column_number"`
	Text              *string `json:"text" yaml:"text"`
}

func New(
	fileInfo *file.FileInfo,
	file *file.Path,
	startLineNumber,
	startColumnNumber int,
	endLineNumber,
	endColumnNumber int,
	text string,
) Source {
	var sourceStartLineNumber *int
	if startLineNumber != 0 {
		sourceStartLineNumber = &startLineNumber
	}

	var sourceStartColumnNumber *int
	if startColumnNumber != 0 {
		sourceStartColumnNumber = &startColumnNumber
	}

	var sourceEndLineNumber *int
	if startLineNumber != 0 {
		sourceEndLineNumber = &endLineNumber
	}

	var sourceEndColumnNumber *int
	if startColumnNumber != 0 {
		sourceEndColumnNumber = &endColumnNumber
	}

	var sourceText *string
	if text != "" {
		sourceText = &text
	}

	language := ""
	languageType := ""
	if fileInfo != nil {
		language = fileInfo.Language
		languageType = fileInfo.LanguageTypeString()
	}

	return Source{
		Filename:          file.RelativePath,
		Language:          language,
		LanguageType:      languageType,
		StartLineNumber:   sourceStartLineNumber,
		StartColumnNumber: sourceStartColumnNumber,
		EndLineNumber:     sourceEndLineNumber,
		EndColumnNumber:   sourceEndColumnNumber,
		Text:              sourceText,
	}
}
