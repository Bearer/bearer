package interfaces

import (
	"errors"
	"strings"

	"github.com/bearer/curio/pkg/classification/db"
	"github.com/bearer/curio/pkg/report"
	"github.com/bearer/curio/pkg/report/interfaces"
	"github.com/bearer/curio/pkg/util/url"
)

const (
	Valid ClassificationState = iota + 1
	Invalid
	Potential
)

type ClassifiedInterface struct {
	*report.Detection
	Classification *Classification `json:"classification"`
}

type ClassificationState int

type ClassificationDecision struct {
	State  ClassificationState
	Reason string
}

type Classification struct {
	URL         string
	RecipeMatch bool
	RecipeName  string
	Decision    ClassificationDecision
}

type Classifier struct {
	config Config
}

type Config struct {
	Recipes []db.Recipe
}

type RecipeURLMatch struct {
	DetectionURLPart string
	RecipeURL        string
	RecipeName       string
}

func New(config Config) *Classifier {
	return &Classifier{config}
}

func NewDefault() *Classifier {
	return &Classifier{
		config: Config{
			Recipes: db.Default(),
		},
	}
}

func (classifier *Classifier) Classify(data report.Detection) (*ClassifiedInterface, error) {
	detectedInterface, ok := data.Value.(interfaces.Interface)
	if !ok {
		return nil, errors.New("detectiosn is not an interface")
	}

	// detected url, with unknown parts replaced with * wildcards
	value := detectedInterface.Value.ToString()
	recipeMatch, err := classifier.FindMatchingRecipeUrl(value)
	if err != nil {
		return nil, err
	}

	if recipeMatch != nil {
		classifiedInterface := &ClassifiedInterface{
			Detection: &data,
			Classification: &Classification{
				URL:         recipeMatch.DetectionURLPart,
				RecipeMatch: true,
				RecipeName:  recipeMatch.RecipeName,
			},
		}

		if strings.Contains(recipeMatch.DetectionURLPart, "*") {
			classifiedInterface.Classification.Decision = ClassificationDecision{
				State:  Potential,
				Reason: "recipe_match_with_wildcard",
			}
			return classifiedInterface, nil
		}

		classifiedInterface.Classification.Decision = ClassificationDecision{
			State:  Valid,
			Reason: "recipe_match",
		}
		return classifiedInterface, nil
	}

	// todo: handle other URL cases (internal, invalid subdomains, etc)
	return &ClassifiedInterface{
		Detection:      &data,
		Classification: &Classification{},
	}, nil
}

func (classifier *Classifier) FindMatchingRecipeUrl(detectionURL string) (*RecipeURLMatch, error) {
	var recipeURLMatch *RecipeURLMatch

	matchSize := 0
	for _, recipe := range classifier.config.Recipes {
		for _, recipeURL := range recipe.URLS {
			match, err := url.Match(
				url.ComparableUrls{
					DetectionURL: detectionURL,
					RecipeURL:    recipeURL,
				},
			)
			if err != nil {
				return nil, err
			}

			if match == "" {
				// no match found; move to next recipe URL
				continue
			}

			candidateSize := len(strings.ReplaceAll(recipeURL, "*", ""))
			if candidateSize <= matchSize {
				// we have a more accurate match already; move to next recipe URL
				continue
			}

			matchSize = candidateSize
			recipeURLMatch = &RecipeURLMatch{
				DetectionURLPart: match,
				RecipeURL:        recipeURL,
				RecipeName:       recipe.Name,
			}
		}
	}

	return recipeURLMatch, nil
}
