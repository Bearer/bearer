package risks_test

import (
	"os"
	"testing"

	"github.com/bearer/curio/pkg/commands/process/settings"
	"github.com/bearer/curio/pkg/report/customdetectors"
	"github.com/bearer/curio/pkg/report/output/dataflow"
	"github.com/bearer/curio/pkg/report/output/dataflow/types"
	"github.com/bearer/curio/pkg/report/output/detectors"
	globaltypes "github.com/bearer/curio/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestDataflowRisks(t *testing.T) {
	config := settings.Config{
		CustomDetector: map[string]settings.Rule{
			"detect_ruby_logger": {
				Stored: false,
				Type:   customdetectors.TypeRisk,
			},
			"ruby_leak": {
				Stored: true,
				Type:   customdetectors.TypeRisk,
			},
		},
	}
	testCases := []struct {
		Name        string
		Config      settings.Config
		FileContent string
		Want        []interface{}
	}{
		{
			Name:        "single detection",
			Config:      config,
			FileContent: `{"id": "1", "type": "custom_classified", "detector_type":"detect_ruby_logger", "source": {"filename": "./users.rb", "line_number": 25}, "value": {"field_name": "User_name", "classification": {"data_type": {"name": "Username", "uuid": "123", "category_uuid": "456"} ,"decision":{"state": "valid"}}}}`,
			Want: []interface{}{
				types.RiskDetector{
					DetectorID: "detect_ruby_logger",
					DataTypes: []types.RiskDatatype{
						{
							Name:   "Username",
							Stored: false,
							Locations: []types.RiskLocation{
								{Filename: "./users.rb", LineNumber: 25},
							},
						},
					},
				},
			},
		},
		{
			Name:        "single detection - no classification",
			Config:      config,
			FileContent: `{"id": "1", "type": "custom_classified", "detector_type":"detect_ruby_logger", "source": {"filename": "./users.rb", "line_number": 25}, "value": {"field_name": "User_name"}}`,
			Want:        []interface{}{},
		},
		{
			Name:   "single detection - duplicates",
			Config: config,
			FileContent: `{"id": "1", "type": "custom_classified", "detector_type":"detect_ruby_logger", "source": {"filename": "./users.rb", "line_number": 25}, "value": {"field_name": "User_name", "classification": {"data_type": {"name": "Username", "uuid": "123", "category_uuid": "456"} ,"decision":{"state": "valid"}}}}
		{"id": "2", "type": "custom_classified", "detector_type":"detect_ruby_logger", "source": {"filename": "./users.rb", "line_number": 25}, "value": {"field_name": "User_name", "classification": {"data_type": {"name": "Username", "uuid": "123", "category_uuid": "456"} ,"decision":{"state": "valid"}}}}`,
			Want: []interface{}{
				types.RiskDetector{
					DetectorID: "detect_ruby_logger",
					DataTypes: []types.RiskDatatype{
						{
							Name:   "Username",
							Stored: false,
							Locations: []types.RiskLocation{
								{Filename: "./users.rb", LineNumber: 25},
							},
						},
					},
				},
			},
		},
		{
			Name:        "single detection - stored",
			Config:      config,
			FileContent: `{"id": "1", "type": "custom_classified", "detector_type":"ruby_leak", "source": {"filename": "./users.rb", "line_number": 25}, "value": {"field_name": "User_name", "classification": {"data_type": {"name": "Username", "uuid": "123", "category_uuid": "456"} ,"decision":{"state": "valid"}}}}`,
			Want: []interface{}{
				types.RiskDetector{
					DetectorID: "ruby_leak",
					DataTypes: []types.RiskDatatype{
						{
							Name:   "Username",
							Stored: true,
							Locations: []types.RiskLocation{
								{Filename: "./users.rb", LineNumber: 25},
							},
						},
					},
				},
			},
		},
		{
			Name:   "single detection - multiple occurences - deterministic output",
			Config: config,
			FileContent: `{"id": "1", "type": "custom_classified", "detector_type":"detect_ruby_logger", "source": {"filename": "./users.rb", "line_number": 25}, "value": {"field_name": "User_name", "classification": {"data_type": {"name": "Username", "uuid": "123", "category_uuid": "456"} ,"decision":{"state": "valid"}}}}
					{"id": "2", "type": "custom_classified", "detector_type":"detect_ruby_logger", "source": {"filename": "./users.rb", "line_number": 2}, "value": {"field_name": "User_name", "classification": {"data_type": {"name": "Username", "uuid": "123", "category_uuid": "456"} ,"decision":{"state": "valid"}}}}`,
			Want: []interface{}{
				types.RiskDetector{
					DetectorID: "detect_ruby_logger",
					DataTypes: []types.RiskDatatype{
						{
							Name:   "Username",
							Stored: false,
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
			Name:   "multiple detections - same detector - deterministic output",
			Config: config,
			FileContent: `{"id": "1", "type": "custom_classified", "detector_type":"detect_ruby_logger", "source": {"filename": "./users.rb", "line_number": 25}, "value": {"field_name": "User_name", "classification": {"data_type": {"name": "Username", "uuid": "123", "category_uuid": "456"} ,"decision":{"state": "valid"}}}}
		{"id": "2", "type": "custom_classified", "detector_type":"detect_ruby_logger", "source": {"filename": "./address.rb", "line_number": 2}, "value": {"field_name": "User_name", "classification": {"data_type": {"name": "Physical Address", "uuid": "123", "category_uuid": "456"} ,"decision":{"state": "valid"}}}}`,
			Want: []interface{}{
				types.RiskDetector{
					DetectorID: "detect_ruby_logger",
					DataTypes: []types.RiskDatatype{
						{
							Name:   "Physical Address",
							Stored: false,
							Locations: []types.RiskLocation{
								{Filename: "./address.rb", LineNumber: 2},
							},
						},
						{
							Name:   "Username",
							Stored: false,
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

			detections, err := detectors.GetOutput(globaltypes.Report{
				Path: file.Name(),
			}, test.Config)
			if err != nil {
				t.Fatalf("failed to get detectors output %s", err)
				return
			}

			dataflow, err := dataflow.GetOutput(detections, test.Config, false)
			if err != nil {
				t.Fatalf("failed to get detectors output %s", err)
				return
			}

			assert.Equal(t, test.Want, dataflow.Risks)
		})
	}
}
