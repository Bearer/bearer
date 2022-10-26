package interfaces

import (
	"errors"
	"regexp"
	"strings"

	"github.com/bearer/curio/pkg/classification/db"
	"github.com/bearer/curio/pkg/report/detections"
	"github.com/bearer/curio/pkg/report/interfaces"
	"github.com/bearer/curio/pkg/util/url"
)

type ClassifiedInterface struct {
	*detections.Detection
	Classification *Classification `json:"classification"`
}

type ClassificationDecision struct {
	State  url.ValidationState `json:"state"`
	Reason string              `json:"reason"`
}

type Classification struct {
	URL         string                 `json:"url"`
	RecipeMatch bool                   `json:"recipe_match"`
	RecipeName  string                 `json:"recipe_name,omitempty"`
	Decision    ClassificationDecision `json:"decision"`
}

type Classifier struct {
	Recipes                []Recipe
	InternalDomainMatchers []*regexp.Regexp
	DomainResolver         *url.DomainResolver
}

type Config struct {
	Recipes         []db.Recipe
	InternalDomains []string
	DomainResolver  *url.DomainResolver
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

var ErrInvalidRecipes = errors.New("invalid interface recipe")
var ErrInvalidInternalDomainRegexp = errors.New("could not parse internal domains as regexp")

func New(config Config) (*Classifier, error) {
	// prepare regular expressions for recipes
	var preparedRecipes []Recipe
	for _, recipe := range config.Recipes {
		preparedRecipe := Recipe{
			Name: recipe.Name,
			Type: recipe.Type,
		}
		for _, recipeURL := range recipe.URLS {
			regexpMatcher, err := url.PrepareRegexpMatcher(recipeURL)
			if err != nil {
				return nil, ErrInvalidRecipes
			}

			preparedRecipeURL := RecipeURL{
				URL:           recipeURL,
				RegexpMatcher: regexpMatcher,
			}
			preparedRecipe.URLS = append(preparedRecipe.URLS, preparedRecipeURL)
		}
		preparedRecipes = append(preparedRecipes, preparedRecipe)
	}

	// parse internal domains as regular expressions
	var internalDomainMatchers []*regexp.Regexp
	for _, internalDomain := range config.InternalDomains {
		internalDomainMatcher, err := regexp.Compile(internalDomain)
		if err != nil {
			return nil, ErrInvalidInternalDomainRegexp
		}

		internalDomainMatchers = append(internalDomainMatchers, internalDomainMatcher)
	}

	return &Classifier{
		Recipes:                preparedRecipes,
		InternalDomainMatchers: internalDomainMatchers,
		DomainResolver:         config.DomainResolver,
	}, nil
}

func NewDefault() (*Classifier, error) {
	return New(
		Config{
			Recipes:         db.Default().Recipes,
			InternalDomains: []string{},
			DomainResolver:  url.NewDomainResolverDefault(),
		},
	)
}

func (classifier *Classifier) Classify(data detections.Detection) (*ClassifiedInterface, error) {
	detectedInterface, ok := data.Value.(interfaces.Interface)
	if !ok {
		return nil, errors.New("detection is not an interface")
	}

	value, err := url.PrepareURLValue(detectedInterface.Value.ToString())
	if err != nil {
		return nil, err
	}

	// check URL format
	formatValidityCheck, err := url.ValidateFormat(value, &data)
	if err != nil {
		return nil, err
	}
	if formatValidityCheck.State == url.Invalid {
		return &ClassifiedInterface{
			Detection: &data,
			Classification: &Classification{
				URL: value,
				Decision: ClassificationDecision{
					State:  formatValidityCheck.State,
					Reason: formatValidityCheck.Reason,
				},
			},
		}, nil
	}

	// check if URL is internal
	var isInternal = false
	for _, matcher := range classifier.InternalDomainMatchers {
		if matcher.MatchString(value) {
			isInternal = true
			break
		}
	}

	if isInternal {
		internalValidityCheck, err := url.ValidateInternal(value)
		if err != nil {
			return nil, err
		}

		return &ClassifiedInterface{
			Detection: &data,
			Classification: &Classification{
				URL: value,
				Decision: ClassificationDecision{
					State:  internalValidityCheck.State,
					Reason: internalValidityCheck.Reason,
				},
			},
		}, nil
	}

	// check for matching recipe
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
				State:  url.Potential,
				Reason: "recipe_match_with_wildcard",
			}
			return classifiedInterface, nil
		}

		classifiedInterface.Classification.Decision = ClassificationDecision{
			State:  url.Valid,
			Reason: "recipe_match",
		}
		return classifiedInterface, nil
	}

	// URL is not internal & no recipe found : validate URL and return result
	validityCheck, err := url.Validate(value, classifier.DomainResolver)
	if err != nil {
		return nil, err
	}

	return &ClassifiedInterface{
		Detection: &data,
		Classification: &Classification{
			URL: value,
			Decision: ClassificationDecision{
				State:  validityCheck.State,
				Reason: validityCheck.Reason,
			},
		},
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
