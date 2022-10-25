package interfaces_test

import (
	"testing"

	"github.com/bearer/curio/pkg/classification/db"
	"github.com/bearer/curio/pkg/classification/interfaces"

	"github.com/bearer/curio/pkg/report/detections"

	reportinterfaces "github.com/bearer/curio/pkg/report/interfaces"
	"github.com/bearer/curio/pkg/report/values"
	"github.com/bearer/curio/pkg/util/url"
	"github.com/stretchr/testify/assert"
)

type testCase struct {
	Name  string
	Input detections.Detection
	Want  *interfaces.Classification
}

func TestInterface(t *testing.T) {
	tests := []testCase{
		{
			Name: "when there is a matching recipe",
			Input: detections.Detection{
				Value: reportinterfaces.Interface{
					Type: reportinterfaces.TypeURL,
					Value: &values.Value{
						Parts: []values.Part{
							&values.String{
								Type:  values.PartTypeString,
								Value: "https://",
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
				URL:         "https://api.stripe.com",
				RecipeName:  "Stripe",
				RecipeMatch: true,
				Decision: interfaces.ClassificationDecision{
					State:  url.Valid,
					Reason: "recipe_match",
				},
			},
		},
		{
			Name: "when it matches an internal domain",
			Input: detections.Detection{
				Value: reportinterfaces.Interface{
					Type: reportinterfaces.TypeURL,
					Value: &values.Value{
						Parts: []values.Part{
							&values.String{
								Type:  values.PartTypeString,
								Value: "https://",
							},
							&values.String{
								Type:  values.PartTypeString,
								Value: "my.internal.domain.com",
							},
						},
					},
				},
			},
			Want: &interfaces.Classification{
				URL: "https://my.internal.domain.com",
				Decision: interfaces.ClassificationDecision{
					State:  url.Valid,
					Reason: "internal_domain_and_subdomain",
				},
			},
		},
		{
			Name: "when there is a matching recipe with a wildcard",
			Input: detections.Detection{
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
								Value: "*.stripe.com",
							},
						},
					},
				},
			},
			Want: &interfaces.Classification{
				URL:         "http://*.stripe.com",
				RecipeName:  "Stripe",
				RecipeMatch: true,
				Decision: interfaces.ClassificationDecision{
					State:  url.Potential,
					Reason: "recipe_match_with_wildcard",
				},
			},
		},
		{
			Name: "when there is a recipe with a path",
			Input: detections.Detection{
				Value: reportinterfaces.Interface{
					Type: reportinterfaces.TypeURL,
					Value: &values.Value{
						Parts: []values.Part{
							&values.String{
								Type:  values.PartTypeString,
								Value: "googleapis.com",
							},
							&values.String{
								Type:  values.PartTypeString,
								Value: "/auth/spreadsheets/",
							},
						},
					},
				},
			},
			Want: &interfaces.Classification{
				URL:         "https://googleapis.com/auth/spreadsheets",
				RecipeName:  "Google Spreadsheets",
				RecipeMatch: true,
				Decision: interfaces.ClassificationDecision{
					State:  url.Valid,
					Reason: "recipe_match",
				},
			},
		},
		{
			Name: "TLD is not allowed",
			Input: detections.Detection{
				Value: reportinterfaces.Interface{
					Type: reportinterfaces.TypeURL,
					Value: &values.Value{
						Parts: []values.Part{
							&values.String{
								Type:  values.PartTypeString,
								Value: "https://",
							},
							&values.String{
								Type:  values.PartTypeString,
								Value: "example.id",
							},
						},
					},
				},
			},
			Want: &interfaces.Classification{
				URL: "https://example.id",
				Decision: interfaces.ClassificationDecision{
					State:  url.Invalid,
					Reason: "tld_error",
				},
			},
		},
		{
			Name: "excluded domain",
			Input: detections.Detection{
				Value: reportinterfaces.Interface{
					Type: reportinterfaces.TypeURL,
					Value: &values.Value{
						Parts: []values.Part{
							&values.String{
								Type:  values.PartTypeString,
								Value: "https://",
							},
							&values.String{
								Type:  values.PartTypeString,
								Value: "github.com",
							},
						},
					},
				},
			},
			Want: &interfaces.Classification{
				URL: "https://github.com",
				Decision: interfaces.ClassificationDecision{
					State:  url.Invalid,
					Reason: "excluded_domains_error",
				},
			},
		},
	}

	classifier, err := interfaces.New(
		interfaces.Config{
			Recipes:         db.Default().Recipes,
			InternalDomains: []string{"https://my.internal.domain.com"},
		},
	)
	if err != nil {
		t.Errorf("Error initializing interface %s", err)
	}

	for _, testCase := range tests {
		t.Run(testCase.Name, func(t *testing.T) {
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
	RecipeName   string
	DetectionURL string
	RecipeURLs   []string
	Want         *interfaces.RecipeURLMatch
}

func TestFindMatchingRecipeUrl(t *testing.T) {
	tests := []recipeMatchTestCase{
		{
			Name:         "when multiple recipes match",
			DetectionURL: "https://api.eu-west.example.com",
			RecipeName:   "Example API",
			RecipeURLs: []string{
				"https://api.*.example.com",
				"https://api.eu-west.example.com",
			},
			Want: &interfaces.RecipeURLMatch{
				RecipeName:       "Example API",
				RecipeURL:        "https://api.eu-west.example.com",
				DetectionURLPart: "https://api.eu-west.example.com",
			},
		},
		{
			Name:         "when no recipes match",
			DetectionURL: "http://no-match.example.com",
			RecipeName:   "Example API",
			RecipeURLs: []string{
				"https://api.*.example.com",
				"https://api.eu-west.example.com",
			},
			Want: nil,
		},
		{
			Name:         "when multiple recipes with the same url length match and one has a wildcard",
			DetectionURL: "https://api.1.example.com",
			RecipeName:   "Example API",
			RecipeURLs: []string{
				"https://api.1.example.com",
				"https://api.*.example.com",
			},
			Want: &interfaces.RecipeURLMatch{
				RecipeName:       "Example API",
				RecipeURL:        "https://api.1.example.com",
				DetectionURLPart: "https://api.1.example.com",
			},
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.Name, func(t *testing.T) {
			classifier, err := interfaces.New(interfaces.Config{
				Recipes: []db.Recipe{
					{
						Name: testCase.RecipeName,
						URLS: testCase.RecipeURLs,
					},
				},
			})

			if err != nil {
				t.Errorf("Error initializing interface %s", err)
			}

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
