package frameworks_test

import (
	"testing"

	"github.com/bearer/curio/pkg/classification/frameworks"
	"github.com/bearer/curio/pkg/report/detections"
	"github.com/bearer/curio/pkg/report/frameworks/rails"
	"github.com/bearer/curio/pkg/report/source"
	"github.com/bearer/curio/pkg/util/classify"

	"github.com/stretchr/testify/assert"
)

func TestFrameworks(t *testing.T) {
	tests := []struct{
		Name  string
		Input detections.Detection
		Want  *frameworks.Classification
		ShouldSucceed bool
	}{
		{
			Name: "Framework match for Rails cache",
			Input: detections.Detection{
				Source: source.Source{
					Filename: "config/application.rb",
					Language: "Ruby",
					LanguageType: "programming",
				},
				Value: rails.Cache{
					Type: "redis_cache_store",
				},
				Type: detections.TypeFramework,
			},
			Want: &frameworks.Classification{
				RecipeMatch: true,
				RecipeName:  "Redis",
				RecipeUUID:  "62c20409-c1bf-4be9-a859-6fe6be7b11e3",
				Decision: classify.ClassificationDecision{
					State:  classify.Valid,
					Reason: "recipe_match",
				},
			},
			ShouldSucceed: true,
		},
		{
			Name: "Framework match for Rails database",
			Input: detections.Detection{
				Source: source.Source{
					Filename: "config/database.yml",
					Language: "YAML",
					LanguageType: "config",
				},
				Value: rails.Database{
					Name: "",
					Adapter: "postgresql",
				},
				Type: detections.TypeFramework,
			},
			Want: &frameworks.Classification{
				RecipeMatch: true,
				RecipeName:  "PostgreSQL",
				RecipeUUID:  "428ff7dd-22ea-4e80-8755-84c70cf460db",
				Decision: classify.ClassificationDecision{
					State:  classify.Valid,
					Reason: "recipe_match",
				},
			},
			ShouldSucceed: true,
		},
		{
			Name: "Framework match for Rails storage",
			Input: detections.Detection{
				Source: source.Source{
					Filename: "config/storage.yml",
					Language: "YAML",
					LanguageType: "config",
				},
				Value: rails.Storage{
					Name: "amazon",
					Service: "S3",
					Encryption: "AES256",
				},
				Type: detections.TypeFramework,
			},
			Want: &frameworks.Classification{
				RecipeMatch: true,
				RecipeName:  "AWS S3",
				RecipeUUID:  "4e5a3a3a-47cd-4b0e-b0a6-fa30a0a62499",
				Decision: classify.ClassificationDecision{
					State:  classify.Valid,
					Reason: "recipe_match",
				},
			},
			ShouldSucceed: true,
		},
		{
			Name: "No framework match (unknown database adapter)",
			Input: detections.Detection{
				Value: rails.Database{
					Name: "",
					Adapter: "my-non-matching-adapter",
				},
				Type: detections.TypeFramework,
			},
			Want: nil,
			ShouldSucceed: true,
		},
		{
			Name: "Invalid detection (vendored detection source)",
			Input: detections.Detection{
				Source: source.Source{
					Filename: "vendor/vendor.rb",
				},
				Value: rails.Database{
					Name: "",
					Adapter: "postgresql",
				},
				Type: detections.TypeFramework,
			},
			Want: nil,
			ShouldSucceed: true,
		},
		{
			Name: "Non-framework detection",
			Input: detections.Detection{
				Value: 12345,
				Type: detections.TypeFramework,
			},
			Want: nil,
			ShouldSucceed: false,
		},
		{
			Name: "Rails cache store in ignore list",
			Input: detections.Detection{
				Source: source.Source{
					Filename: "config/application.rb",
					Language: "Ruby",
					LanguageType: "programming",
				},
				Value: rails.Cache{
					Type: "memory_store",
				},
				Type: detections.TypeFramework,
			},
			Want: nil,
			ShouldSucceed: true,
		},
		{
			Name: "Rails storage with name including 'test'",
			Input: detections.Detection{
				Source: source.Source{
					Filename: "config/storage.yml",
					Language: "YAML",
					LanguageType: "config",
				},
				Value: rails.Storage{
					Name: "my_test",
					Service: "S3",
					Encryption: "AES256",
				},
				Type: detections.TypeFramework,
			},
			Want: nil,
			ShouldSucceed: true,
			},
		{
			Name: "Rails storage mirror",
			Input: detections.Detection{
				Source: source.Source{
					Filename: "config/storage.yml",
					Language: "YAML",
					LanguageType: "config",
				},
				Value: rails.Storage{
					Name: "production",
					Service: "Mirror",
				},
				Type: detections.TypeFramework,
			},
			Want: nil,
			ShouldSucceed: true,
		},
	}

	classifier := frameworks.NewDefault()

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
