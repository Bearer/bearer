package output_test

import (
	"os"
	"testing"

	"github.com/bearer/curio/pkg/report/output"
	"github.com/bearer/curio/pkg/report/output/dataflow"
	"github.com/bearer/curio/pkg/report/output/dataflow/types"
	globaltypes "github.com/bearer/curio/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestDataflowRisks(t *testing.T) {
	testCases := []struct {
		Name        string
		FileContent string
		Want        []types.RiskDetector
	}{
		{
			Name:        "single detection",
			FileContent: `{"type": "custom_classified", "detector_type":"logger_leak", "source": {"filename": "./users.rb", "line_number": 25}, "value": {"field_name": "User"}}`,
			Want: []types.RiskDetector{
				{
					DetectorID: "logger_leak",
					DataTypes: []types.RiskDatatype{
						{
							Name: "User",
							Locations: []types.RiskLocation{
								{Filename: "./users.rb", LineNumber: 25},
							},
						},
					},
				},
			},
		},
		{
			Name: "single detection - duplicates",
			FileContent: `{"type": "custom_classified", "detector_type":"logger_leak", "source": {"filename": "./users.rb", "line_number": 25}, "value": {"field_name": "User"}}
{"type": "custom_classified", "detector_type":"logger_leak", "source": {"filename": "./users.rb", "line_number": 25}, "value": {"field_name": "User"}}`,
			Want: []types.RiskDetector{
				{
					DetectorID: "logger_leak",
					DataTypes: []types.RiskDatatype{
						{
							Name: "User",
							Locations: []types.RiskLocation{
								{Filename: "./users.rb", LineNumber: 25},
							},
						},
					},
				},
			},
		},
		{
			Name: "single detection - multiple occurences",
			FileContent: `{"type": "custom_classified", "detector_type":"logger_leak", "source": {"filename": "./users.rb", "line_number": 25}, "value": {"field_name": "User"}}
{"type": "custom_classified", "detector_type":"logger_leak", "source": {"filename": "./users.rb", "line_number": 2}, "value": {"field_name": "User"}}`,
			Want: []types.RiskDetector{
				{
					DetectorID: "logger_leak",
					DataTypes: []types.RiskDatatype{
						{
							Name: "User",
							Locations: []types.RiskLocation{
								{Filename: "./users.rb", LineNumber: 2},
								{Filename: "./users.rb", LineNumber: 25},
							},
						},
					},
				},
			},
		},
		{
			Name: "multiple detections - same detector",
			FileContent: `{"type": "custom_classified", "detector_type":"logger_leak", "source": {"filename": "./users.rb", "line_number": 25}, "value": {"field_name": "User"}}
{"type": "custom_classified", "detector_type":"logger_leak", "source": {"filename": "./address.rb", "line_number": 2}, "value": {"field_name": "Address"}}`,
			Want: []types.RiskDetector{
				{
					DetectorID: "logger_leak",
					DataTypes: []types.RiskDatatype{
						{
							Name: "Address",
							Locations: []types.RiskLocation{
								{Filename: "./address.rb", LineNumber: 2},
							},
						},
						{
							Name: "User",
							Locations: []types.RiskLocation{
								{Filename: "./users.rb", LineNumber: 25},
							},
						},
					},
				},
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.Name, func(t *testing.T) {
			file, err := os.CreateTemp("", "*test.jsonlines")
			if err != nil {
				t.Fatalf("failed to create tmp file for report %s", err)
				return
			}
			defer os.Remove(file.Name())
			_, err = file.Write([]byte(test.FileContent))
			if err != nil {
				t.Fatalf("failed to write to tmp file %s", err)
				return
			}
			file.Close()

			detections, err := output.GetDetectorsOutput(globaltypes.Report{
				Path: file.Name(),
			})
			if err != nil {
				t.Fatalf("failed to get detectors output %s", err)
				return
			}

			dataflow, err := dataflow.GetOuput(detections)
			if err != nil {
				t.Fatalf("failed to get detectors output %s", err)
				return
			}

			assert.Equal(t, test.Want, dataflow.Risks)
		})
	}
}
