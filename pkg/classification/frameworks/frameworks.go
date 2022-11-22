package frameworks

import (
	"errors"

	"github.com/bearer/curio/pkg/classification/db"
	"github.com/bearer/curio/pkg/report/detections"
	"github.com/bearer/curio/pkg/report/frameworks/rails"
	"github.com/bearer/curio/pkg/util/classify"
)

type ClassifiedFramework struct {
	*detections.Detection
	Classification *Classification `json:"classification"`
}

type Classification struct {
	RecipeMatch bool   `json:"recipe_match"`
	RecipeName  string `json:"recipe_name,omitempty"`
	RecipeUUID  string `json:"recipe_uuid,omitempty"`
	Decision    classify.ClassificationDecision `json:"decision"`
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

func (classifier *Classifier) Classify(data detections.Detection) (*ClassifiedFramework, error) {
	var classification *Classification

	if classify.IsVendored(data.Source.Filename) {
		return &ClassifiedFramework{
			Detection:      &data,
			Classification: classification,
		}, nil
	}

	if classify.IsPotentialDetector(data.DetectorType) {
		return &ClassifiedFramework{
			Detection:      &data,
			Classification: classification,
		}, nil
	}

	var technologyKey string
	switch value := data.Value.(type) {
	case rails.Cache:
		technologyKey = value.GetTechnologyKey()
	case rails.Database:
		technologyKey = value.GetTechnologyKey()
	case rails.Storage:
		technologyKey = value.GetTechnologyKey()
	default:
		return &ClassifiedFramework{
			Detection:      &data,
			Classification: classification,
		}, errors.New("detection is not for a framework")
	}

	if technologyKey != "" {
		for _, recipe := range classifier.config.Recipes {
			if isRecipeMatch(recipe, technologyKey) {
				classification = &Classification{
					RecipeUUID:  recipe.UUID,
					RecipeName:  recipe.Name,
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

	return &ClassifiedFramework{
		Detection:      &data,
		Classification: classification,
	}, nil
}

func isRecipeMatch(recipe db.Recipe, technologyKey string) bool {
	return recipe.Name == technologyKey
}
