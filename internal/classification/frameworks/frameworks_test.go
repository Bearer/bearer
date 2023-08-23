package frameworks_test

import (
	"testing"

	"github.com/bearer/bearer/internal/classification/frameworks"
	"github.com/bearer/bearer/internal/report/detections"
	"github.com/bearer/bearer/internal/report/frameworks/beego"
	"github.com/bearer/bearer/internal/report/frameworks/django"
	"github.com/bearer/bearer/internal/report/frameworks/dotnet"
	"github.com/bearer/bearer/internal/report/frameworks/rails"
	"github.com/bearer/bearer/internal/report/frameworks/spring"
	"github.com/bearer/bearer/internal/report/frameworks/symfony"
	"github.com/bearer/bearer/internal/report/source"
	"github.com/bearer/bearer/internal/util/classify"

	"github.com/stretchr/testify/assert"
)

func TestFrameworks(t *testing.T) {
	tests := []struct {
		Name          string
		Input         detections.Detection
		Want          *frameworks.Classification
		ShouldSucceed bool
	}{
		{
			Name: "Framework match for Rails cache",
			Input: detections.Detection{
				Source: source.Source{
					Filename:     "config/application.rb",
					Language:     "Ruby",
					LanguageType: "programming",
				},
				Value: rails.Cache{
					Type: "redis_cache_store",
				},
				Type: detections.TypeFramework,
			},
			Want: &frameworks.Classification{
				RecipeMatch:   true,
				RecipeName:    "Redis",
				RecipeType:    "data_store",
				RecipeSubType: "key_value_cache",
				RecipeUUID:    "62c20409-c1bf-4be9-a859-6fe6be7b11e3",
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
					Filename:     "config/database.yml",
					Language:     "YAML",
					LanguageType: "config",
				},
				Value: rails.Database{
					Name:    "",
					Adapter: "postgresql",
				},
				Type: detections.TypeFramework,
			},
			Want: &frameworks.Classification{
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
			Name: "Framework match for Rails storage",
			Input: detections.Detection{
				Source: source.Source{
					Filename:     "config/storage.yml",
					Language:     "YAML",
					LanguageType: "config",
				},
				Value: rails.Storage{
					Name:       "amazon",
					Service:    "S3",
					Encryption: "AES256",
				},
				Type: detections.TypeFramework,
			},
			Want: &frameworks.Classification{
				RecipeMatch:   true,
				RecipeName:    "AWS S3",
				RecipeType:    "data_store",
				RecipeSubType: "object_storage",
				RecipeUUID:    "4e5a3a3a-47cd-4b0e-b0a6-fa30a0a62499",
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
					Name:    "",
					Adapter: "my-non-matching-adapter",
				},
				Type: detections.TypeFramework,
			},
			Want:          nil,
			ShouldSucceed: true,
		},
		{
			Name: "Invalid detection (vendored detection source)",
			Input: detections.Detection{
				Source: source.Source{
					Filename: "vendor/vendor.rb",
				},
				Value: rails.Database{
					Name:    "",
					Adapter: "postgresql",
				},
				Type: detections.TypeFramework,
			},
			Want:          nil,
			ShouldSucceed: true,
		},
		{
			Name: "Non-framework detection",
			Input: detections.Detection{
				Value: 12345,
				Type:  detections.TypeFramework,
			},
			Want:          nil,
			ShouldSucceed: false,
		},
		{
			Name: "Rails cache store in ignore list",
			Input: detections.Detection{
				Source: source.Source{
					Filename:     "config/application.rb",
					Language:     "Ruby",
					LanguageType: "programming",
				},
				Value: rails.Cache{
					Type: "memory_store",
				},
				Type: detections.TypeFramework,
			},
			Want:          nil,
			ShouldSucceed: true,
		},
		{
			Name: "Rails storage with name including 'test'",
			Input: detections.Detection{
				Source: source.Source{
					Filename:     "config/storage.yml",
					Language:     "YAML",
					LanguageType: "config",
				},
				Value: rails.Storage{
					Name:       "my_test",
					Service:    "S3",
					Encryption: "AES256",
				},
				Type: detections.TypeFramework,
			},
			Want:          nil,
			ShouldSucceed: true,
		},
		{
			Name: "Rails storage mirror",
			Input: detections.Detection{
				Source: source.Source{
					Filename:     "config/storage.yml",
					Language:     "YAML",
					LanguageType: "config",
				},
				Value: rails.Storage{
					Name:    "production",
					Service: "Mirror",
				},
				Type: detections.TypeFramework,
			},
			Want:          nil,
			ShouldSucceed: true,
		},
		{
			Name: "Beego: driver name defined",
			Input: detections.Detection{
				Source: source.Source{
					Filename:     "orm.go",
					Language:     "Go",
					LanguageType: "programming",
				},
				Value: beego.Database{
					Name:         "default",
					DriverName:   "mysql",
					Package:      "",
					TypeConstant: "",
				},
				Type: detections.TypeFramework,
			},
			Want: &frameworks.Classification{
				RecipeMatch:   true,
				RecipeName:    "MySQL",
				RecipeType:    "data_store",
				RecipeSubType: "database",
				RecipeUUID:    "ffa70264-2b19-445d-a5c9-be82b64fe750",
				Decision: classify.ClassificationDecision{
					State:  classify.Valid,
					Reason: "recipe_match",
				},
			},
			ShouldSucceed: true,
		},
		{
			Name: "Beego: package defined",
			Input: detections.Detection{
				Source: source.Source{
					Filename:     "orm.go",
					Language:     "Go",
					LanguageType: "programming",
				},
				Value: beego.Database{
					Name:         "default",
					DriverName:   "",
					Package:      "github.com/beego/beego/v2/client/orm",
					TypeConstant: "DRSqlite",
				},
				Type: detections.TypeFramework,
			},
			Want: &frameworks.Classification{
				RecipeMatch:   true,
				RecipeName:    "SQLite",
				RecipeType:    "data_store",
				RecipeSubType: "database",
				RecipeUUID:    "aa706b3c-0f6d-4a7b-a7a5-71ee0c5b6c00",
				Decision: classify.ClassificationDecision{
					State:  classify.Valid,
					Reason: "recipe_match",
				},
			},
			ShouldSucceed: true,
		},
		{
			Name: "Django: database match",
			Input: detections.Detection{
				Source: source.Source{
					Filename:     "orm.py",
					Language:     "Python",
					LanguageType: "programming",
				},
				Value: django.Database{
					Name:   "default",
					Engine: "django.db.backends.mysql",
				},
				Type: detections.TypeFramework,
			},
			Want: &frameworks.Classification{
				RecipeMatch:   true,
				RecipeName:    "MySQL",
				RecipeType:    "data_store",
				RecipeSubType: "database",
				RecipeUUID:    "ffa70264-2b19-445d-a5c9-be82b64fe750",
				Decision: classify.ClassificationDecision{
					State:  classify.Valid,
					Reason: "recipe_match",
				},
			},
			ShouldSucceed: true,
		},
		{
			Name: ".NET: database match",
			Input: detections.Detection{
				Source: source.Source{
					Filename:     "Startup.cs",
					Language:     "C#",
					LanguageType: "programming",
				},
				Value: dotnet.DBContext{
					UseDbMethodName: "UseSqlServer",
				},
				Type: detections.TypeFramework,
			},
			Want: &frameworks.Classification{
				RecipeMatch:   true,
				RecipeName:    "Microsoft SQL Server",
				RecipeType:    "data_store",
				RecipeSubType: "database",
				RecipeUUID:    "e4db4505-b837-4b76-9184-c3cec3b5e522",
				Decision: classify.ClassificationDecision{
					State:  classify.Valid,
					Reason: "recipe_match",
				},
			},
			ShouldSucceed: true,
		},
		{
			Name: "Spring: database match",
			Input: detections.Detection{
				Source: source.Source{
					Filename:     "src/main/application.properties",
					Language:     "Properties",
					LanguageType: "config",
				},
				Value: spring.DataStore{
					Driver: "com.mysql.jdbc.Driver",
				},
				Type: detections.TypeFramework,
			},
			Want: &frameworks.Classification{
				RecipeMatch:   true,
				RecipeName:    "MySQL",
				RecipeType:    "data_store",
				RecipeSubType: "database",
				RecipeUUID:    "ffa70264-2b19-445d-a5c9-be82b64fe750",
				Decision: classify.ClassificationDecision{
					State:  classify.Valid,
					Reason: "recipe_match",
				},
			},
			ShouldSucceed: true,
		},
		{
			Name: "Symfony: database match",
			Input: detections.Detection{
				Source: source.Source{
					Filename:     "config/packages/doctrine.yml",
					Language:     "YAML",
					LanguageType: "config",
				},
				Value: symfony.Database{
					Name:   "production",
					Driver: "oci8",
				},
				Type: detections.TypeFramework,
			},
			Want: &frameworks.Classification{
				RecipeMatch:   true,
				RecipeName:    "Oracle",
				RecipeType:    "data_store",
				RecipeSubType: "database",
				RecipeUUID:    "80886e2a-ee2c-423d-98bc-0a3d743787b4",
				Decision: classify.ClassificationDecision{
					State:  classify.Valid,
					Reason: "recipe_match",
				},
			},
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
