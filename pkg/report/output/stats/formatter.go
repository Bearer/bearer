package stats

import (
	"github.com/bearer/bearer/pkg/commands/process/settings"
	"github.com/bearer/bearer/pkg/flag"
	outputtypes "github.com/bearer/bearer/pkg/report/output/types"
	outputhandler "github.com/bearer/bearer/pkg/util/output"
)

type Formatter struct {
	ReportData *outputtypes.ReportData
	Config     settings.Config
}

func NewFormatter(reportData *outputtypes.ReportData, config settings.Config) *Formatter {
	return &Formatter{
		ReportData: reportData,
		Config:     config,
	}
}

func (f Formatter) Format(format string) (output *string, err error) {
	switch format {
	case flag.FormatEmpty, flag.FormatJSON:
		return outputhandler.ReportJSON(f.ReportData.Stats)
	case flag.FormatYAML:
		return outputhandler.ReportYAML(f.ReportData.Stats)
	}

	return output, err
}
