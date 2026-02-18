package paths_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/bearer/bearer/pkg/commands/process/settings"
	"github.com/bearer/bearer/pkg/report/output/dataflow"
	"github.com/bearer/bearer/pkg/report/output/dataflow/types"
	"github.com/bearer/bearer/pkg/report/output/detectors"
	outputtypes "github.com/bearer/bearer/pkg/report/output/types"
	globaltypes "github.com/bearer/bearer/pkg/types"
)

func TestDataflowPaths(t *testing.T) {
	config := settings.Config{}
	var lineNumber = new(int)
	*lineNumber = 558

	testCases := []struct {
		Name        string
		Config      settings.Config
		FileContent string
		Want        []types.Path
	}{
		{
			Name:        "OpenAPI paths",
			Config:      config,
			FileContent: `{ "detector_type": "openapi", "source": { "end_column_number": 8, "end_line_number": 558, "filename": "testdata/v3yaml/petstore-openapi.yaml", "full_filename": "", "language": "YAML", "language_type": "data", "start_column_number": 5, "start_line_number": 558, "text": "get" }, "type": "operation", "value": { "path": "/user/*", "type": "GET", "url": [ { "url": "{protocol}://api.example.com", "variables": [ { "Name": "protocol", "Values": [ "http", "https" ] } ] }, { "url": "https://{environment}.example.com/v2", "variables": [ { "Name": "environment", "Values": [ "api", "api.dev", "api.staging" ] } ] }, { "url": "{server}/v1", "variables": [ { "Name": "server", "Values": [ "https://api.example.com" ] } ] } ] } }`,
			Want: []types.Path{
				{
					DetectorName: "openapi",
					Detections: []*types.Detection{
						{
							FullFilename: "testdata/v3yaml/petstore-openapi.yaml",
							FullName:     "testdata/v3yaml/petstore-openapi.yaml",
							LineNumber:   lineNumber,
							Path:         "/user/*",
							HttpMethod:   "GET",
							Urls: []string{
								"{protocol}://api.example.com",
								"https://{environment}.example.com/v2",
								"{server}/v1",
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

			output := &outputtypes.ReportData{}
			if err = detectors.AddReportData(output, globaltypes.Report{
				Path: file.Name(),
			}, test.Config); err != nil {
				t.Fatalf("failed to get detectors output %s", err)
				return
			}

			if err = dataflow.AddReportData(output, test.Config, false, true); err != nil {
				t.Fatalf("failed to get dataflow output %s", err)
				return
			}

			assert.Equal(t, test.Want, output.Dataflow.Paths)
		})
	}
}
