package flag

type Severity int

var (
	FormatJSON = "json"
	FormatYAML = "yaml"

	ReportDetectors = "detectors"
	ReportDataFlow  = "dataflow"
	ReportPolicies  = "policies"
	ReportStats     = "stats"
)

var (
	FormatFlag = Flag{
		Name:       "format",
		ConfigName: "report.format",
		Shorthand:  "f",
		Value:      "",
		Usage:      "Specify report format (json, yaml)",
	}
	ReportFlag = Flag{
		Name:       "report",
		ConfigName: "report.report",
		Value:      ReportDetectors,
		Usage:      "Specify the kind of report (detectors, dataflow, policies, stats)",
	}
	OutputFlag = Flag{
		Name:       "output",
		ConfigName: "report.output",
		Value:      "",
		Usage:      "Specify output path for report",
	}
)

type ReportFlagGroup struct {
	Format *Flag
	Report *Flag
	Output *Flag
}

type ReportOptions struct {
	Format string `mapstructure:"format" json:"format" yaml:"format"`
	Report string `mapstructure:"report" json:"report" yaml:"report"`
	Output string `mapstructure:"output" json:"output" yaml:"output"`
}

func NewReportFlagGroup() *ReportFlagGroup {
	return &ReportFlagGroup{
		Format: &FormatFlag,
		Report: &ReportFlag,
		Output: &OutputFlag,
	}
}

func (f *ReportFlagGroup) Name() string {
	return "Report"
}

func (f *ReportFlagGroup) Flags() []*Flag {
	return []*Flag{
		f.Format,
		f.Report,
		f.Output,
	}
}

func (f *ReportFlagGroup) ToOptions() ReportOptions {
	return ReportOptions{
		Format: getString(f.Format),
		Report: getString(f.Report),
		Output: getString(f.Output),
	}
}
