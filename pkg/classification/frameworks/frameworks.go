package frameworks

import (
	"errors"

	"github.com/bearer/curio/pkg/classification/db"
	"github.com/bearer/curio/pkg/report/detections"
	"github.com/bearer/curio/pkg/util/classify"
)

type classifiableFramework interface{
	GetTechnologyKey() string
}

type ClassifiedFramework struct {
	*detections.Detection
	Classification *Classification `json:"classification" yaml:"classification"`
}

type Classification struct {
	RecipeMatch   bool                            `json:"recipe_match" yaml:"recipe_match"`
	RecipeName    string                          `json:"recipe_name,omitempty"`
	RecipeUUID    string                          `json:"recipe_uuid,omitempty"`
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

	value, ok := data.Value.(classifiableFramework)
	if !ok {
		return &ClassifiedFramework{
			Detection:      &data,
			Classification: classification,
		}, errors.New("detection is not for a framework")
	}

	technologyKey = value.GetTechnologyKey()

	if technologyKey != "" {
		for _, recipe := range classifier.config.Recipes {
			if isRecipeMatch(recipe, technologyKey) {
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

	return &ClassifiedFramework{
		Detection:      &data,
		Classification: classification,
	}, nil
}

func isRecipeMatch(recipe db.Recipe, technologyKey string) bool {
	return recipe.UUID == technologyKey
}
