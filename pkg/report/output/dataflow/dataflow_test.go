package dataflow_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bearer/bearer/pkg/commands/process/settings"
	"github.com/bearer/bearer/pkg/report/output/dataflow"
	outputtypes "github.com/bearer/bearer/pkg/report/output/types"
)

func TestAddReportDataIncludesLanguages(t *testing.T) {
	reportData := &outputtypes.ReportData{
		FoundLanguages: map[string]int32{"Ruby": 123},
	}

	err := dataflow.AddReportData(reportData, settings.Config{}, false, false)
	require.NoError(t, err)

	require.NotNil(t, reportData.Dataflow)
	require.Equal(t, reportData.FoundLanguages, reportData.Dataflow.Languages)
}

func TestAddReportDataIncludesLanguagesWithFiles(t *testing.T) {
	reportData := &outputtypes.ReportData{
		FoundLanguages: map[string]int32{"Go": 456},
	}

	// Simulate hasFiles = true
	err := dataflow.AddReportData(reportData, settings.Config{}, true, false)
	require.NoError(t, err)

	require.NotNil(t, reportData.Dataflow)
	require.Equal(t, reportData.FoundLanguages, reportData.Dataflow.Languages)
}
