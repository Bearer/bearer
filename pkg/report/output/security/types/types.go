package types

import (
	"fmt"
	"strings"

	"github.com/bearer/bearer/pkg/util/file"
	"github.com/fatih/color"
)

type Finding struct {
	*Rule
	LineNumber        int               `json:"line_number,omitempty" yaml:"line_number,omitempty"`
	FullFilename      string            `json:"full_filename,omitempty" yaml:"full_filename,omitempty"`
	Filename          string            `json:"filename,omitempty" yaml:"filename,omitempty"`
	DataType          *DataType         `json:"data_type,omitempty" yaml:"data_type,omitempty"`
	CategoryGroups    []string          `json:"category_groups,omitempty" yaml:"category_groups,omitempty"`
	Source            Source            `json:"source,omitempty" yaml:"source,omitempty"`
	Sink              Sink              `json:"sink,omitempty" yaml:"sink,omitempty"`
	ParentLineNumber  int               `json:"parent_line_number,omitempty" yaml:"parent_line_number,omitempty"`
	ParentContent     string            `json:"snippet,omitempty" yaml:"snippet,omitempty"`
	Fingerprint       string            `json:"fingerprint,omitempty" yaml:"fingerprint,omitempty"`
	OldFingerprint    string            `json:"old_fingerprint,omitempty" yaml:"old_fingerprint,omitempty"`
	DetailedContext   string            `json:"detailed_context,omitempty" yaml:"detailed_context,omitempty"`
	CodeExtract       string            `json:"code_extract,omitempty" yaml:"code_extract,omitempty"`
	RawCodeExtract    []file.Line       `json:"-" yaml:"-"`
	SeverityWeighting SeverityWeighting `json:"-" yaml:"-"`
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

type SeverityWeighting struct {
	RuleSeverity                   string `json:"rule_severity" yaml:"rule_severity"`
	DisplaySeverity                string `json:"display_severity" yaml:"display_severity"`
	SensitiveDataCategoryWeighting int    `json:"sensitive_data_category_weighting,omitempty" yaml:"sensitive_data_category_weighting,omitempty"`
	RuleSeverityWeighting          int    `json:"rule_severity_weighting,omitempty" yaml:"rule_severity_weighting,omitempty"`
	HasLocalDataTypes              *bool  `json:"local_data_types,omitempty" yaml:"local_data_types,omitempty"`
	FinalWeighting                 int    `json:"final_weighting,omitempty" yaml:"final_weighting,omitempty"`
}

func (f Finding) HighlightCodeExtract() string {
	result := ""
	for _, line := range f.RawCodeExtract {
		if line.Strip {
			result += color.HiBlackString(
				fmt.Sprintf(" %s %s", strings.Repeat(" ", iterativeDigitsCount(line.LineNumber)), line.Extract),
			)
		} else {
			result += color.HiMagentaString(fmt.Sprintf(" %d ", line.LineNumber))
			if line.LineNumber == f.Source.Start && line.LineNumber == f.Source.End {
				for i, char := range line.Extract {
					if i >= f.Source.Column.Start-1 && i < f.Source.Column.End-1 {
						result += color.MagentaString(fmt.Sprintf("%c", char))
					} else {
						result += color.HiMagentaString(fmt.Sprintf("%c", char))
					}
				}
			} else if line.LineNumber == f.Source.Start && line.LineNumber <= f.Source.End {
				for i, char := range line.Extract {
					if i >= f.Source.Column.Start-1 {
						result += color.MagentaString(fmt.Sprintf("%c", char))
					} else {
						result += color.HiMagentaString(fmt.Sprintf("%c", char))
					}
				}
			} else if line.LineNumber == f.Source.End && line.LineNumber >= f.Source.Start {
				for i, char := range line.Extract {
					if i <= f.Source.Column.End-1 {
						result += color.MagentaString(fmt.Sprintf("%c", char))
					} else {
						result += color.HiMagentaString(fmt.Sprintf("%c", char))
					}
				}
			} else if line.LineNumber > f.Source.Start && line.LineNumber < f.Source.End {
				result += color.MagentaString("%s", line.Extract)
			} else {
				result += color.HiMagentaString(line.Extract)
			}
		}

		if line.LineNumber != f.Sink.End {
			result += "\n"
		}
	}

	return result
}

func iterativeDigitsCount(number int) int {
	count := 0
	for number != 0 {
		number /= 10
		count += 1
	}

	return count
}
