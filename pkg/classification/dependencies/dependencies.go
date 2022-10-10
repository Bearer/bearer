package dependencies

import (
	"github.com/bearer/curio/pkg/classification/db"
	"github.com/bearer/curio/pkg/report"
	"github.com/bearer/curio/pkg/report/dependencies"
	"github.com/rs/zerolog/log"
)

type ClassifiedDependency struct {
	*report.Detection
	Classification *Classification `json:"classification"`
}

type Classification struct {
	RecipeMatch bool
	RecipeName  string
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

func (classifier *Classifier) Classify(data report.Detection) (ClassifiedDependency, error) {
	var classification *Classification
	for _, recipe := range classifier.config.Recipes {
		for _, recipePackage := range recipe.Packages {
			log.Debug().Msgf(
				"Matching recipe package manager %s against data detector type %s",
				recipePackage.PackageManager,
				data.DetectorType,
			)

			if isRecipeMatch(recipePackage, data) {
				classification = &Classification{
					RecipeName:  recipe.Name,
					RecipeMatch: true,
				}
			}
		}
	}

	return ClassifiedDependency{
		Detection:      &data,
		Classification: classification,
	}, nil
}

func isRecipeMatch(recipePackage db.Package, data report.Detection) bool {
	value := data.Value.(dependencies.Dependency)

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
