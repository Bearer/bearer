package interfaces_test

import (
	"testing"

	"github.com/bearer/curio/pkg/classification/interfaces"
	"github.com/bearer/curio/pkg/report"
	reportinterfaces "github.com/bearer/curio/pkg/report/interfaces"
	"github.com/bearer/curio/pkg/report/values"
	"github.com/stretchr/testify/assert"
)

type testCase struct {
	Name  string
	Input report.Detection
	Want  *interfaces.Classification
}

func TestSchema(t *testing.T) {
	tests := []testCase{
		{
			Name: "simple path",
			Input: report.Detection{
				Value: reportinterfaces.Interface{
					Type: reportinterfaces.TypeURL,
					Value: &values.Value{
						Parts: []values.Part{
							&values.String{
								Type:  values.PartTypeString,
								Value: "http://",
							},
							&values.String{
								Type:  values.PartTypeString,
								Value: "api.stripe.com",
							},
						},
					},
				},
			},
			Want: &interfaces.Classification{
				RecipeName: "stripe",
			},
		},
		{
			Name: "simple path - no match",
			Input: report.Detection{
				Value: reportinterfaces.Interface{
					Type: reportinterfaces.TypeURL,
					Value: &values.Value{
						Parts: []values.Part{
							&values.String{
								Type:  values.PartTypeString,
								Value: "http://",
							},
							&values.String{
								Type:  values.PartTypeString,
								Value: "api.example.com",
							},
						},
					},
				},
			},
			Want: nil,
		},
	}

	classifier := interfaces.New(interfaces.Config{})

	for _, testCase := range tests {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Skip("interfaces not implemented")

			output, err := classifier.Classify(testCase.Input)
			if err != nil {
				t.Errorf("classifier returned error %s", err)
			}

			assert.Equal(t, testCase.Want, output.Classification)
		})
	}
}
