package dataflow

import (
	"fmt"

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
	case flag.FormatEmpty:
		output, err = outputhandler.ReportJSON(f.ReportData.Dataflow)
	case flag.FormatJSON:
		output, err = outputhandler.ReportJSON(f.ReportData.Dataflow)
	case flag.FormatYAML:
		output, err = outputhandler.ReportYAML(f.ReportData.Dataflow)
	default:
		err = fmt.Errorf(`--report flag "%s" is not supported`, f.Config.Report.Report)
	}

	return output, err
}
