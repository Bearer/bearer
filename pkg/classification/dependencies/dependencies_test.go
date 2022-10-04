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
	Want  dependencies.Classification
}

func TestSchema(t *testing.T) {
	tests := []testCase{
		{
			Name: "simple path",
			Input: report.Detection{
				Value: reportdependencies.Dependency{
					Group:   "",
					Name:    "stripe",
					Version: "v1.1.1",
				},
				Type:         report.TypeDependency,
				DetectorType: reportdependencies.DetectorGemFileLock,
			},
			Want: dependencies.Classification{
				RecipeMatch: true,
				RecipeName:  "stripe",
			},
		},
	}

	classifier := dependencies.New(dependencies.Config{})

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
