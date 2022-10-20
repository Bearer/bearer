package dependencies

import (
	"errors"

	"github.com/bearer/curio/pkg/classification/db"
	"github.com/bearer/curio/pkg/report"
	"github.com/bearer/curio/pkg/report/dependencies"
)

type ClassifiedDependency struct {
	*report.Detection
	Classification *Classification `json:"classification"`
}

type Classification struct {
	RecipeMatch bool   `json:"recipe_match"`
	RecipeName  string `json:"recipe_name,omitempty"`
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
			Recipes: db.Default(),
		},
	}
}

func (classifier *Classifier) Classify(data report.Detection) (*ClassifiedDependency, error) {
	var classification *Classification
	value, ok := data.Value.(dependencies.Dependency)
	if !ok {
		return nil, errors.New("detection is not an dependency")
	}

	for _, recipe := range classifier.config.Recipes {
		for _, recipePackage := range recipe.Packages {
			if isRecipeMatch(recipePackage, value) {
				classification = &Classification{
					RecipeName:  recipe.Name,
					RecipeMatch: true,
				}
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
