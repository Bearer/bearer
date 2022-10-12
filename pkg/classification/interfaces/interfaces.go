package interfaces

import (
	"strings"

	"github.com/bearer/curio/pkg/classification/db"
	"github.com/bearer/curio/pkg/report"
	"github.com/bearer/curio/pkg/report/interfaces"
	"github.com/bearer/curio/pkg/util/url_matcher"
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
	var classifiedInterface *ClassifiedInterface
	var classificationDecision ClassificationDecision

	// detected url, with unknown parts replaced with * wildcards
	value := data.Value.(interfaces.Interface).Value.ToString()
	recipeMatch, err := classifier.FindMatchingRecipeUrl(value)
	if err != nil {
		return classifiedInterface, err
	}

	if recipeMatch != nil {
		if strings.Contains(recipeMatch.DetectionURLPart, "*") {
			classificationDecision = ClassificationDecision{
				State:  Potential,
				Reason: "recipe_match_with_wildcard",
			}
		} else {
			classificationDecision = ClassificationDecision{
				State:  Valid,
				Reason: "recipe_match",
			}
		}

		return &ClassifiedInterface{
			Detection: &data,
			Classification: &Classification{
				URL:         recipeMatch.DetectionURLPart,
				RecipeMatch: true,
				RecipeName:  recipeMatch.RecipeName,
				Decision:    classificationDecision,
			},
		}, nil
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
			match, err := url_matcher.UrlMatcher(
				url_matcher.ComparableUrls{
					DetectionURL: detectionURL,
					RecipeURL:    recipeURL,
				},
			)
			if err != nil {
				return recipeURLMatch, err
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
