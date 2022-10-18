package flag

import (
	"io"
)

type Severity int

var (
	FormatJSON      = "json"
	FormatJSONLines = "jsonlines"
)

var (
	FormatFlag = Flag{
		Name:       "format",
		ConfigName: "format",
		Shorthand:  "f",
		Value:      FormatJSON,
		Usage:      "format (json)",
	}
)

type ReportFlagGroup struct {
	Format *Flag
}

type ReportOptions struct {
	Format string
}

func NewReportFlagGroup() *ReportFlagGroup {
	return &ReportFlagGroup{
		Format: &FormatFlag,
	}
}

func (f *ReportFlagGroup) Name() string {
	return "Report"
}

func (f *ReportFlagGroup) Flags() []*Flag {
	return []*Flag{f.Format}
}

func (f *ReportFlagGroup) ToOptions(out io.Writer) ReportOptions {
	return ReportOptions{
		Format: getString(f.Format),
	}
}
