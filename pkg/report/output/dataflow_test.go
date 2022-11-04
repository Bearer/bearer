package output_test

import (
	"os"
	"testing"

	"github.com/bearer/curio/pkg/report/output"
	"github.com/bearer/curio/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestDataflow(t *testing.T) {
	testCases := []struct {
		Name        string
		FileContent string
		Want        *output.DataFlow
	}{
		{
			Name:        "single detection",
			FileContent: `{"type": "schema", "detector_type":"ruby", "source": {"filename": "./users.rb", "line_number": 25}, "value": {"object_name": "User"}}`,
			Want: &output.DataFlow{
				Datatypes: []output.Datatype{
					{
						Name: "User",
						Detectors: []output.Detector{
							{
								Name:   "ruby",
								Stored: true,
								Locations: []output.Location{
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
			FileContent: `{"type": "schema", "detector_type":"ruby", "source": {"filename": "./users.rb", "line_number": 25}, "value": {"object_name": "User"}}
{"type": "schema", "detector_type":"ruby", "source": {"filename": "./users.rb", "line_number": 25}, "value": {"object_name": "User"}}`,
			Want: &output.DataFlow{
				Datatypes: []output.Datatype{
					{
						Name: "User",
						Detectors: []output.Detector{
							{
								Name:   "ruby",
								Stored: true,
								Locations: []output.Location{
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
			FileContent: `{"type": "schema", "detector_type":"ruby", "source": {"filename": "./users.rb", "line_number": 25}, "value": {"object_name": "User"}}
{"user": true }`,
			Want: &output.DataFlow{
				Datatypes: []output.Datatype{
					{
						Name: "User",
						Detectors: []output.Detector{
							{
								Name:   "ruby",
								Stored: true,
								Locations: []output.Location{
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
			FileContent: `{"type": "schema", "detector_type":"ruby", "source": {"filename": "./users.rb", "line_number": 25}, "value": {"object_name": "User"}}
			{"type": "schema", "detector_type":"csharp", "source": {"filename": "./users.cs", "line_number": 12}, "value": {"object_name": "User"}}`,
			Want: &output.DataFlow{
				Datatypes: []output.Datatype{
					{
						Name: "User",
						Detectors: []output.Detector{
							{
								Name:   "csharp",
								Stored: true,
								Locations: []output.Location{
									{Filename: "./users.cs", LineNumber: 12},
								},
							},
							{
								Name:   "ruby",
								Stored: true,
								Locations: []output.Location{
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
			FileContent: `{"type": "schema", "detector_type":"ruby", "source": {"filename": "./users.rb", "line_number": 25}, "value": {"object_name": "User"}}
			{"type": "schema", "detector_type":"csharp", "source": {"filename": "./users.cs", "line_number": 12}, "value": {"object_name": "Address"}}`,
			Want: &output.DataFlow{
				Datatypes: []output.Datatype{
					{
						Name: "Address",
						Detectors: []output.Detector{
							{
								Name:   "csharp",
								Stored: true,
								Locations: []output.Location{
									{Filename: "./users.cs", LineNumber: 12},
								},
							},
						},
					},
					{
						Name: "User",
						Detectors: []output.Detector{
							{
								Name:   "ruby",
								Stored: true,
								Locations: []output.Location{
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

			dataflow, err := output.GetDataFlowOutput(types.Report{
				Path: file.Name(),
			})
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, test.Want, dataflow)
		})
	}
}
