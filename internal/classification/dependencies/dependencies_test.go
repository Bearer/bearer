package dependencies_test

import (
	"testing"

	"github.com/bearer/bearer/internal/classification/dependencies"
	reportdependencies "github.com/bearer/bearer/internal/report/dependencies"
	"github.com/bearer/bearer/internal/report/detections"
	"github.com/bearer/bearer/internal/report/source"
	"github.com/bearer/bearer/internal/util/classify"

	"github.com/stretchr/testify/assert"
)

func TestDependencies(t *testing.T) {
	tests := []struct {
		Name          string
		Input         detections.Detection
		Want          *dependencies.Classification
		ShouldSucceed bool
	}{
		{
			Name: "Dependency match",
			Input: detections.Detection{
				Value: reportdependencies.Dependency{
					Group:          "",
					Name:           "stripe",
					Version:        "v1.1.1",
					PackageManager: "rubygems",
				},
				Type: detections.TypeDependency,
			},
			Want: &dependencies.Classification{
				RecipeMatch:   true,
				RecipeName:    "Stripe",
				RecipeType:    "external_service",
				RecipeSubType: "third_party",
				RecipeUUID:    "c24b836a-d035-49dc-808f-1912f16f690d",
				Decision: classify.ClassificationDecision{
					State:  classify.Valid,
					Reason: "recipe_match",
				},
			},
			ShouldSucceed: true,
		},
		{
			Name: "Dependency match with group (Java case)",
			Input: detections.Detection{
				Value: reportdependencies.Dependency{
					Group:          "org.postgresql",
					Name:           "postgresql",
					Version:        "v2.1.1",
					PackageManager: "maven",
				},
				Type: detections.TypeDependency,
			},
			Want: &dependencies.Classification{
				RecipeMatch:   true,
				RecipeName:    "PostgreSQL",
				RecipeType:    "data_store",
				RecipeSubType: "database",
				RecipeUUID:    "428ff7dd-22ea-4e80-8755-84c70cf460db",
				Decision: classify.ClassificationDecision{
					State:  classify.Valid,
					Reason: "recipe_match",
				},
			},
			ShouldSucceed: true,
		},
		{
			Name: "No dependency match",
			Input: detections.Detection{
				Value: reportdependencies.Dependency{
					Group:          "",
					Name:           "my-non-matching-dependency",
					Version:        "v2.1.1",
					PackageManager: "rubygems",
				},
				Type: detections.TypeDependency,
			},
			Want:          nil,
			ShouldSucceed: true,
		},
		{
			Name: "Vendored detection",
			Input: detections.Detection{
				Source: source.Source{
					Filename: "vendor/vendor.js",
				},
				Value: reportdependencies.Dependency{
					Group:          "",
					Name:           "stripe",
					Version:        "v1.1.1",
					PackageManager: "npm",
				},
				Type: detections.TypeDependency,
			},
			Want:          nil,
			ShouldSucceed: true,
		},
		{
			Name: "Non-dependency detection",
			Input: detections.Detection{
				Value: 12345,
				Type:  detections.TypeDependency,
			},
			Want:          nil,
			ShouldSucceed: false,
		},
	}

	classifier := dependencies.NewDefault()

	for _, testCase := range tests {
		t.Run(testCase.Name, func(t *testing.T) {
			output, err := classifier.Classify(testCase.Input)
			if err != nil && testCase.ShouldSucceed {
				t.Errorf("classifier returned error %s", err)
			}

			assert.Equal(t, testCase.Want, output.Classification)
		})
	}
}
