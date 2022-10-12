package interfaces_test

import (
	"testing"

	"github.com/bearer/curio/pkg/classification/db"
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

func TestInterface(t *testing.T) {
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
				RecipeName:  "Stripe",
				RecipeMatch: true,
				Decision: interfaces.ClassificationDecision{
					State:  interfaces.Valid,
					Reason: "asdf",
				},
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

	classifier := interfaces.NewDefault()

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

type recipeMatchTestCase struct {
	Name         string
	DetectionURL string
	RecipeURLs   []string
	Want         *interfaces.RecipeURLMatch
}

func TestFindMatchingRecipeUrl(t *testing.T) {
	tests := []recipeMatchTestCase{
		{
			Name:         "when multiple recipes match",
			DetectionURL: "https://api.eu-west.example.com",
			RecipeURLs: []string{
				"https://api.*.example.com",
				"https://api.eu-west.example.com",
			},
			Want: &interfaces.RecipeURLMatch{
				RecipeURL:        "https://api.eu-west.example.com",
				DetectionURLPart: "https://api.eu-west.example.com",
			},
		},
		{
			Name:         "when no recipes match",
			DetectionURL: "http://no-match.example.com",
			RecipeURLs: []string{
				"https://api.*.example.com",
				"https://api.eu-west.example.com",
			},
			Want: nil,
		},
		{
			Name:         "when multiple recipes with the same url length match and one has a wildcard",
			DetectionURL: "https://api.1.example.com",
			RecipeURLs: []string{
				"https://api.1.example.com",
				"https://api.*.example.com",
			},
			Want: &interfaces.RecipeURLMatch{
				RecipeURL:        "https://api.1.example.com",
				DetectionURLPart: "https://api.1.example.com",
			},
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.Name, func(t *testing.T) {
			classifier := interfaces.New(interfaces.Config{
				Recipes: []db.Recipe{
					{
						URLS: testCase.RecipeURLs,
					},
				},
			})

			output, err := classifier.FindMatchingRecipeUrl(
				testCase.DetectionURL,
			)
			if err != nil {
				t.Errorf("FindMatchingRecipeUrl returned error %s", err)
			}

			assert.Equal(t, testCase.Want, output)
		})
	}
}
