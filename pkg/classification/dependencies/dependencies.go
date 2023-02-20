package dependencies

import (
	"errors"

	"github.com/bearer/bearer/pkg/classification/db"
	"github.com/bearer/bearer/pkg/report/dependencies"
	"github.com/bearer/bearer/pkg/report/detections"
	"github.com/bearer/bearer/pkg/util/classify"
)

type ClassifiedDependency struct {
	*detections.Detection
	Classification *Classification `json:"classification" yaml:"classification"`
}

type Classification struct {
	RecipeMatch   bool                            `json:"recipe_match" yaml:"recipe_match"`
	RecipeUUID    string                          `json:"recipe_uuid,omitempty"`
	RecipeName    string                          `json:"recipe_name,omitempty"`
	RecipeSubType string                          `json:"recipe_sub_type,omitempty"`
	Decision      classify.ClassificationDecision `json:"decision" yaml:"decision"`
}

type Classifier struct {
	config Config
}

type Config struct {
	Recipes []db.Recipe
}

func New(config Config) *Classifier {
	return &Classifier{config: config}
}

func NewDefault() *Classifier {
	return &Classifier{
		config: Config{
			Recipes: db.Default().Recipes,
		},
	}
}

func (classifier *Classifier) Classify(data detections.Detection) (*ClassifiedDependency, error) {
	var classification *Classification
	value, ok := data.Value.(dependencies.Dependency)
	if !ok {
		return &ClassifiedDependency{
			Detection:      &data,
			Classification: classification,
		}, errors.New("detection is not an dependency")
	}

	if classify.IsVendored(data.Source.Filename) {
		return &ClassifiedDependency{
			Detection:      &data,
			Classification: classification,
		}, nil
	}

	if classify.IsPotentialDetector(data.DetectorType) {
		return &ClassifiedDependency{
			Detection:      &data,
			Classification: classification,
		}, nil
	}

	for _, recipe := range classifier.config.Recipes {
		for _, recipePackage := range recipe.Packages {
			if isRecipeMatch(recipePackage, value) {
				classification = &Classification{
					RecipeUUID:    recipe.UUID,
					RecipeName:    recipe.Name,
					RecipeSubType: recipe.SubType,
					RecipeMatch: true,
					Decision: classify.ClassificationDecision{
						State:  classify.Valid,
						Reason: "recipe_match",
					},
				}
				break
			}
		}
	}

	return &ClassifiedDependency{
		Detection:      &data,
		Classification: classification,
	}, nil
}

func isRecipeMatch(recipePackage db.Package, value dependencies.Dependency) bool {
	if isJavaPackage(recipePackage.PackageManager) {
		return recipePackage.PackageManager == value.PackageManager &&
			recipePackage.Name == value.Name &&
			recipePackage.Group == value.Group
	} else {
		return recipePackage.PackageManager == value.PackageManager &&
			recipePackage.Name == value.Name
	}
}

func isJavaPackage(packageManager string) bool {
	return packageManager == "maven"
}
