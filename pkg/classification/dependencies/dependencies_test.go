package dependencies_test

import (
	"testing"

	"github.com/bearer/curio/pkg/classification/dependencies"
	"github.com/bearer/curio/pkg/report"
	reportdependencies "github.com/bearer/curio/pkg/report/dependencies"

	"github.com/stretchr/testify/assert"
)

type testCase struct {
	Name  string
	Input report.Detection
	Want  *dependencies.Classification
}

func TestDependencies(t *testing.T) {
	tests := []testCase{
		{
			Name: "Dependency match",
			Input: report.Detection{
				Value: reportdependencies.Dependency{
					Group:          "",
					Name:           "stripe",
					Version:        "v1.1.1",
					PackageManager: "rubygems",
				},
				Type: report.TypeDependency,
			},
			Want: &dependencies.Classification{
				RecipeMatch: true,
				RecipeName:  "stripe",
			},
		},
		{
			Name: "Dependency match with group (Java case)",
			Input: report.Detection{
				Value: reportdependencies.Dependency{
					Group:          "org.postgresql",
					Name:           "postgresql",
					Version:        "v2.1.1",
					PackageManager: "maven",
				},
				Type: report.TypeDependency,
			},
			Want: &dependencies.Classification{
				RecipeMatch: true,
				RecipeName:  "postgres",
			},
		},
		{
			Name: "No dependency match",
			Input: report.Detection{
				Value: reportdependencies.Dependency{
					Group:          "",
					Name:           "my-non-matching-dependency",
					Version:        "v2.1.1",
					PackageManager: "rubygems",
				},
				Type: report.TypeDependency,
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
