package interfaces

import (
	"errors"
	"regexp"
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
	Recipes []Recipe
}

type Config struct {
	Recipes []db.Recipe
}

type Recipe struct {
	Name string
	Type string
	URLS []RecipeURL
}

type RecipeURL struct {
	URL           string
	RegexpMatcher *regexp.Regexp
}

type RecipeURLMatch struct {
	DetectionURLPart string
	RecipeURL        string
	RecipeName       string
}

func New(config Config) *Classifier {
	var preparedRecipes []Recipe
	for _, recipe := range config.Recipes {
		preparedRecipe := Recipe{
			Name: recipe.Name,
			Type: recipe.Type,
		}
		for _, recipeURL := range recipe.URLS {
			regexpMatcher, err := url.PrepareRegexpMatcher(recipeURL)
			if err != nil {
				panic(err) // todo: how to handle error in New()?
			}

			preparedRecipeURL := RecipeURL{
				URL:           recipeURL,
				RegexpMatcher: regexpMatcher,
			}
			preparedRecipe.URLS = append(preparedRecipe.URLS, preparedRecipeURL)
		}
		preparedRecipes = append(preparedRecipes, preparedRecipe)
	}

	return &Classifier{Recipes: preparedRecipes}
}

func NewDefault() *Classifier {
	return New(Config{Recipes: db.Default()})
}

func (classifier *Classifier) Classify(data report.Detection) (*ClassifiedInterface, error) {
	detectedInterface, ok := data.Value.(interfaces.Interface)
	if !ok {
		return nil, errors.New("detection is not an interface")
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
	for _, recipe := range classifier.Recipes {
		for _, recipeURL := range recipe.URLS {
			match, err := url.Match(detectionURL, recipeURL.RegexpMatcher)
			if err != nil {
				return nil, err
			}

			if match == "" {
				// no match found; move to next recipe URL
				continue
			}

			candidateSize := len(strings.ReplaceAll(recipeURL.URL, "*", ""))
			if candidateSize <= matchSize {
				// we have a more accurate match already; move to next recipe URL
				continue
			}

			matchSize = candidateSize
			recipeURLMatch = &RecipeURLMatch{
				DetectionURLPart: match,
				RecipeURL:        recipeURL.URL,
				RecipeName:       recipe.Name,
			}
		}
	}

	return recipeURLMatch, nil
}
