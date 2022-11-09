package output_test

import (
	"os"
	"testing"

	"github.com/bearer/curio/pkg/commands/process/settings"
	"github.com/bearer/curio/pkg/report/output"
	"github.com/bearer/curio/pkg/report/output/dataflow"
	"github.com/bearer/curio/pkg/report/output/dataflow/types"
	globaltypes "github.com/bearer/curio/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestDataflowComponents(t *testing.T) {
	testCases := []struct {
		Name        string
		FileContent string
		Want        []types.Component
	}{
		{
			Name: "single detection - dependency",
			FileContent: `{	"detector_type": "gemfile-lock", "type": "dependency_classified", "source": {"filename": "Gemfile.lock", "line_number": 258}, "classification": { "Decision": { "state": "valid" }, "recipe_name": "stripe"}}`,
			Want: []types.Component{
				{
					Name: "stripe",
					Locations: []types.ComponentLocation{
						{
							Detector:   "gemfile-lock",
							Filename:   "Gemfile.lock",
							LineNumber: 258,
						},
					},
				},
			},
		},
		{
			Name: "single detection - interface",
			FileContent: `{	"detector_type": "ruby", "type": "interface_classified", "source": {"filename": "billing.rb", "line_number": 2}, "classification": { "Decision": { "state": "valid" }, "recipe_name": "Stripe"}}`,
			Want: []types.Component{
				{
					Name: "stripe",
					Locations: []types.ComponentLocation{
						{
							Detector:   "ruby",
							Filename:   "billing.rb",
							LineNumber: 2,
						},
					},
				},
			},
		},
		{
			Name: "single detection - duplicates",
			FileContent: `{	"detector_type": "ruby", "type": "interface_classified", "source": {"filename": "billing.rb", "line_number": 2}, "classification": { "Decision": { "state": "valid" }, "recipe_name": "Stripe"}}
{ "detector_type": "ruby", "type": "interface_classified", "source": {"filename": "billing.rb", "line_number": 2}, "classification": { "Decision": { "state": "valid" }, "recipe_name": "Stripe"}}`,
			Want: []types.Component{
				{
					Name: "stripe",
					Locations: []types.ComponentLocation{
						{
							Detector:   "ruby",
							Filename:   "billing.rb",
							LineNumber: 2,
						},
					},
				},
			},
		},
		{
			Name: "multiple detections - deterministic output",
			FileContent: `{	"detector_type": "ruby", "type": "interface_classified", "source": {"filename": "billing.rb", "line_number": 2}, "classification": { "Decision": { "state": "valid" }, "recipe_name": "Stripe"}}
{"detector_type": "gemfile-lock", "type": "dependency_classified", "source": {"filename": "Gemfile.lock", "line_number": 258}, "classification": { "Decision": { "state": "valid" }, "recipe_name": "stripe"}}`,
			Want: []types.Component{
				{
					Name: "stripe",
					Locations: []types.ComponentLocation{
						{
							Detector:   "gemfile-lock",
							Filename:   "Gemfile.lock",
							LineNumber: 258,
						},
						{
							Detector:   "ruby",
							Filename:   "billing.rb",
							LineNumber: 2,
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

			dataflow, err := dataflow.GetOuput(detections, settings.Config{})
			if err != nil {
				t.Fatalf("failed to get detectors output %s", err)
				return
			}

			assert.Equal(t, test.Want, dataflow.Components)
		})
	}
}
