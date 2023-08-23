package types

import "github.com/bearer/bearer/pkg/util/file"

type Finding struct {
	*Rule
	LineNumber       int         `json:"line_number,omitempty" yaml:"line_number,omitempty"`
	FullFilename     string      `json:"full_filename,omitempty" yaml:"full_filename,omitempty"`
	Filename         string      `json:"filename,omitempty" yaml:"filename,omitempty"`
	DataType         *DataType   `json:"data_type,omitempty" yaml:"data_type,omitempty"`
	CategoryGroups   []string    `json:"category_groups,omitempty" yaml:"category_groups,omitempty"`
	Source           Source      `json:"source,omitempty" yaml:"source,omitempty"`
	Sink             Sink        `json:"sink,omitempty" yaml:"sink,omitempty"`
	ParentLineNumber int         `json:"parent_line_number,omitempty" yaml:"parent_line_number,omitempty"`
	ParentContent    string      `json:"snippet,omitempty" yaml:"snippet,omitempty"`
	Fingerprint      string      `json:"fingerprint,omitempty" yaml:"fingerprint,omitempty"`
	OldFingerprint   string      `json:"old_fingerprint,omitempty" yaml:"old_fingerprint,omitempty"`
	DetailedContext  string      `json:"detailed_context,omitempty" yaml:"detailed_context,omitempty"`
	CodeExtract      string      `json:"code_extract,omitempty" yaml:"code_extract,omitempty"`
	RawCodeExtract   []file.Line `json:"-" yaml:"-"`
}

type DataType struct {
	CategoryUUID string `json:"category_uuid,omitempty" yaml:"category_uuid,omitempty"`
	Name         string `json:"name,omitempty" yaml:"name,omitempty"`
}

type Rule struct {
	CWEIDs           []string `json:"cwe_ids" yaml:"cwe_ids"`
	Id               string   `json:"id" yaml:"id"`
	Title            string   `json:"title" yaml:"title"`
	Description      string   `json:"description" yaml:"description"`
	DocumentationUrl string   `json:"documentation_url" yaml:"documentation_url"`
}

type Location struct {
	Start  int    `json:"start" yaml:"start"`
	End    int    `json:"end" yaml:"end"`
	Column Column `json:"column" yaml:"column"`
}

type Source struct {
	*Location
}

type Column struct {
	Start int `json:"start" yaml:"start"`
	End   int `json:"end" yaml:"end"`
}

type Sink struct {
	*Location
	Content string `json:"content" yaml:"content"`
}
