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

func TestDataflowDataType(t *testing.T) {
	testCases := []struct {
		Name        string
		FileContent string
		Want        *dataflow.DataFlow
	}{
		{
			Name:        "single detection",
			FileContent: `{"type": "schema", "detector_type":"ruby", "source": {"filename": "./users.rb", "line_number": 25}, "value": {"field_name": "User"}}`,
			Want: &dataflow.DataFlow{
				Datatypes: []types.Datatype{
					{
						Name: "User",
						Detectors: []types.DatatypeDetector{
							{
								Name:   "ruby",
								Stored: true,
								Locations: []types.DatatypeLocation{
									{Filename: "./users.rb", LineNumber: 25},
								},
							},
						},
					},
				},
			},
		},
		{
			Name: "single detection - duplicates",
			FileContent: `{"type": "schema", "detector_type":"ruby", "source": {"filename": "./users.rb", "line_number": 25}, "value": {"field_name": "User"}}
{"type": "schema", "detector_type":"ruby", "source": {"filename": "./users.rb", "line_number": 25}, "value": {"field_name": "User"}}`,
			Want: &dataflow.DataFlow{
				Datatypes: []types.Datatype{
					{
						Name: "User",
						Detectors: []types.DatatypeDetector{
							{
								Name:   "ruby",
								Stored: true,
								Locations: []types.DatatypeLocation{
									{Filename: "./users.rb", LineNumber: 25},
								},
							},
						},
					},
				},
			},
		},
		{
			Name: "single detection - with wierd data in report",
			FileContent: `{"type": "schema", "detector_type":"ruby", "source": {"filename": "./users.rb", "line_number": 25}, "value": {"field_name": "User"}}
{"user": true }`,
			Want: &dataflow.DataFlow{
				Datatypes: []types.Datatype{
					{
						Name: "User",
						Detectors: []types.DatatypeDetector{
							{
								Name:   "ruby",
								Stored: true,
								Locations: []types.DatatypeLocation{
									{Filename: "./users.rb", LineNumber: 25},
								},
							},
						},
					},
				},
			},
		},
		{
			Name: "multiple detections - with same object name - deterministic output",
			FileContent: `{"type": "schema", "detector_type":"ruby", "source": {"filename": "./users.rb", "line_number": 25}, "value": {"field_name": "User"}}
			{"type": "schema", "detector_type":"csharp", "source": {"filename": "./users.cs", "line_number": 12}, "value": {"field_name": "User"}}`,
			Want: &dataflow.DataFlow{
				Datatypes: []types.Datatype{
					{
						Name: "User",
						Detectors: []types.DatatypeDetector{
							{
								Name:   "csharp",
								Stored: true,
								Locations: []types.DatatypeLocation{
									{Filename: "./users.cs", LineNumber: 12},
								},
							},
							{
								Name:   "ruby",
								Stored: true,
								Locations: []types.DatatypeLocation{
									{Filename: "./users.rb", LineNumber: 25},
								},
							},
						},
					},
				},
			},
		},
		{
			Name: "multiple detections - with different names",
			FileContent: `{"type": "schema", "detector_type":"ruby", "source": {"filename": "./users.rb", "line_number": 25}, "value": {"field_name": "User"}}
			{"type": "schema", "detector_type":"csharp", "source": {"filename": "./users.cs", "line_number": 12}, "value": {"field_name": "Address"}}`,
			Want: &dataflow.DataFlow{
				Datatypes: []types.Datatype{
					{
						Name: "Address",
						Detectors: []types.DatatypeDetector{
							{
								Name:   "csharp",
								Stored: true,
								Locations: []types.DatatypeLocation{
									{Filename: "./users.cs", LineNumber: 12},
								},
							},
						},
					},
					{
						Name: "User",
						Detectors: []types.DatatypeDetector{
							{
								Name:   "ruby",
								Stored: true,
								Locations: []types.DatatypeLocation{
									{Filename: "./users.rb", LineNumber: 25},
								},
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

			assert.Equal(t, test.Want, dataflow)
		})
	}
}
