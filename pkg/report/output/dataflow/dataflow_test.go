package dataflow_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bearer/bearer/pkg/commands/process/settings"
	"github.com/bearer/bearer/pkg/report/output/dataflow"
	outputtypes "github.com/bearer/bearer/pkg/report/output/types"
)

func TestAddReportDataIncludesLanguageStats(t *testing.T) {
	testCases := []struct {
		name     string
		hasFiles bool
	}{
		{
			name:     "no files in report",
			hasFiles: false,
		},
		{
			name:     "report with files",
			hasFiles: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			reportData := &outputtypes.ReportData{
				FoundLanguages:     map[string]int32{"Ruby": 123},
				LanguageFiles:      map[string]int32{"Ruby": 7},
				TotalLanguageFiles: 7,
			}

			err := dataflow.AddReportData(reportData, settings.Config{}, false, testCase.hasFiles)
			require.NoError(t, err)

			require.NotNil(t, reportData.Dataflow)
			require.Len(t, reportData.Dataflow.Languages, 1)
			require.Equal(t, "Ruby", reportData.Dataflow.Languages[0].Language)
			require.Equal(t, int32(123), reportData.Dataflow.Languages[0].Lines)
			require.Equal(t, int32(7), reportData.Dataflow.Languages[0].Files)
			require.Equal(t, int32(7), reportData.Dataflow.TotalFiles)
		})
	}
}
