package source

import (
	"github.com/bearer/curio/pkg/util/file"
)

// Source represents a part of a source file that is referenced in the scan report.
type Source struct {
	Filename     string  `json:"filename"`
	Language     string  `json:"language"`
	LanguageType string  `json:"language_type"`
	LineNumber   *int    `json:"line_number"`
	ColumnNumber *int    `json:"column_number"`
	Text         *string `json:"text"`
}

func New(fileInfo *file.FileInfo, file *file.Path, lineNumber, columnNumber int, text string) Source {
	var sourceLineNumber *int
	if lineNumber != 0 {
		sourceLineNumber = &lineNumber
	}

	var sourceColumnNumber *int
	if columnNumber != 0 {
		sourceColumnNumber = &columnNumber
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
		Filename:     file.RelativePath,
		Language:     language,
		LanguageType: languageType,
		LineNumber:   sourceLineNumber,
		ColumnNumber: sourceColumnNumber,
		Text:         sourceText,
	}
}
