package dataflow

import (
	"github.com/bearer/bearer/internal/commands/process/settings"
	"github.com/bearer/bearer/internal/flag"
	outputtypes "github.com/bearer/bearer/internal/report/output/types"
	outputhandler "github.com/bearer/bearer/internal/util/output"
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

func (f Formatter) Format(format string) (output string, err error) {
	switch format {
	case flag.FormatEmpty, flag.FormatJSON:
		return outputhandler.ReportJSON(f.ReportData.Dataflow)
	case flag.FormatYAML:
		return outputhandler.ReportYAML(f.ReportData.Dataflow)
	}

	return output, err
}
