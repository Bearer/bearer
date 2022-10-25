package dependencies_test

import (
	"testing"

	"github.com/bearer/curio/pkg/classification/dependencies"
	reportdependencies "github.com/bearer/curio/pkg/report/dependencies"
	"github.com/bearer/curio/pkg/report/detections"

	"github.com/stretchr/testify/assert"
)

type testCase struct {
	Name  string
	Input detections.Detection
	Want  *dependencies.Classification
}

func TestDependencies(t *testing.T) {
	tests := []testCase{
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
				RecipeMatch: true,
				RecipeName:  "Stripe",
			},
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
				RecipeMatch: true,
				RecipeName:  "PostgreSQL",
			},
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
			Want: nil,
		},
	}

	classifier := dependencies.NewDefault()

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
