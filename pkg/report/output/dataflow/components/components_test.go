package components_test

import (
	"os"
	"testing"

	"github.com/bearer/bearer/pkg/commands/process/settings"
	"github.com/bearer/bearer/pkg/report/output/dataflow"
	"github.com/bearer/bearer/pkg/report/output/dataflow/types"
	"github.com/bearer/bearer/pkg/report/output/detectors"
	globaltypes "github.com/bearer/bearer/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestDataflowComponents(t *testing.T) {
	testCases := []struct {
		Name        string
		FileContent string
		Want        []types.Component
	}{
		{
			Name:        "single detection - dependency",
			FileContent: `{	"detector_type": "gemfile-lock", "type": "dependency_classified", "source": {"filename": "Gemfile.lock", "line_number": 258, "start_line_number": 258}, "classification": { "Decision": { "state": "valid" }, "recipe_name": "Stripe", "recipe_match": true, "recipe_type": "external_service", "recipe_sub_type": "third_party"}}`,
			Want: []types.Component{
				{
					Name:    "Stripe",
					Type:    "external_service",
					SubType: "third_party",
					Locations: []types.ComponentLocation{
						{
							Detector:     "gemfile-lock",
							FullFilename: "Gemfile.lock",
							Filename:     "Gemfile.lock",
							LineNumber:   258,
						},
					},
				},
			},
		},
		{
			Name:        "single detection - dependency - no classification",
			FileContent: `{	"detector_type": "gemfile-lock", "type": "dependency_classified", "source": {"filename": "Gemfile.lock", "line_number": 258, "start_line_number": 258}}`,
			Want:        []types.Component{},
		},
		{
			Name:        "single detection - interface",
			FileContent: `{	"detector_type": "ruby", "type": "interface_classified", "source": {"filename": "billing.rb", "line_number": 2, "start_line_number": 2}, "classification": { "Decision": { "state": "valid" }, "recipe_name": "Stripe", "recipe_match": true, "recipe_type": "external_service", "recipe_sub_type": "third_party"}}`,
			Want: []types.Component{
				{
					Name:    "Stripe",
					Type:    "external_service",
					SubType: "third_party",
					Locations: []types.ComponentLocation{
						{
							Detector:     "ruby",
							FullFilename: "billing.rb",
							Filename:     "billing.rb",
							LineNumber:   2,
						},
					},
				},
			},
		},
		{
			Name:        "single detection - interface - no classification",
			FileContent: `{	"detector_type": "ruby", "type": "interface_classified", "source": {"filename": "billing.rb", "line_number": 2, "start_line_number": 2}}`,
			Want:        []types.Component{},
		},
		{
			Name: "single detection - duplicates",
			FileContent: `{	"detector_type": "ruby", "type": "interface_classified", "source": {"filename": "billing.rb", "line_number": 2, "start_line_number": 2}, "classification": { "Decision": { "state": "valid" }, "recipe_name": "Stripe", "recipe_match": true, "recipe_type": "external_service", "recipe_sub_type": "third_party"}}
{ "detector_type": "ruby", "type": "interface_classified", "source": {"filename": "billing.rb", "line_number": 2, "start_line_number": 2}, "classification": { "Decision": { "state": "valid" }, "recipe_name": "Stripe", "recipe_match": true, "recipe_type": "external_service", "recipe_sub_type": "third_party"}}`,
			Want: []types.Component{
				{
					Name:    "Stripe",
					Type:    "external_service",
					SubType: "third_party",
					Locations: []types.ComponentLocation{
						{
							Detector:     "ruby",
							Filename:     "billing.rb",
							FullFilename: "billing.rb",
							LineNumber:   2,
						},
					},
				},
			},
		},
		{
			Name: "multiple detections - deterministic output",
			FileContent: `{	"detector_type": "ruby", "type": "interface_classified", "source": {"filename": "billing.rb", "line_number": 2, "start_line_number": 2}, "classification": { "Decision": { "state": "valid" }, "recipe_name": "Stripe", "recipe_type": "external_service", "recipe_sub_type": "third_party", "recipe_uuid": "123-abc", "recipe_match": true}}
{"detector_type": "gemfile-lock", "type": "dependency_classified", "source": {"filename": "Gemfile.lock", "line_number": 258, "start_line_number": 258}, "classification": { "Decision": { "state": "valid" }, "recipe_name": "Stripe", "recipe_type": "external_service", "recipe_sub_type": "third_party", "recipe_uuid": "123-abc", "recipe_match": true}}`,
			Want: []types.Component{
				{
					Name:    "Stripe",
					Type:    "external_service",
					SubType: "third_party",
					Locations: []types.ComponentLocation{
						{
							Detector:     "gemfile-lock",
							FullFilename: "Gemfile.lock",
							Filename:     "Gemfile.lock",
							LineNumber:   258,
						},
						{
							Detector:     "ruby",
							FullFilename: "billing.rb",
							Filename:     "billing.rb",
							LineNumber:   2,
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

			detections, _, err := detectors.GetOutput(globaltypes.Report{
				Path: file.Name(),
			}, settings.Config{})
			if err != nil {
				t.Fatalf("failed to get detectors output %s", err)
				return
			}

			dataflow, _, err := dataflow.GetOutput(detections, settings.Config{}, false)
			if err != nil {
				t.Fatalf("failed to get detectors output %s", err)
				return
			}

			assert.Equal(t, test.Want, dataflow.Components)
		})
	}
}
