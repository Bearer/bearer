package interfaces

import (
	"errors"
	"regexp"
	"strings"

	"github.com/bearer/bearer/pkg/classification/db"
	"github.com/bearer/bearer/pkg/report/detections"
	"github.com/bearer/bearer/pkg/report/interfaces"
	"github.com/bearer/bearer/pkg/util/classify"
	"github.com/bearer/bearer/pkg/util/url"
)

type ClassifiedInterface struct {
	*detections.Detection
	Classification *Classification `json:"classification" yaml:"classification"`
}

type Classification struct {
	URL           string                          `json:"url" yaml:"url"`
	RecipeMatch   bool                            `json:"recipe_match" yaml:"recipe_match"`
	RecipeName    string                          `json:"recipe_name,omitempty"`
	RecipeUUID    string                          `json:"recipe_uuid,omitempty"`
	RecipeSubType string                          `json:"recipe_sub_type,omitempty"`
	Decision      classify.ClassificationDecision `json:"decision" yaml:"decision"`
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
	UUID        string
	Name        string
	Type        string
	SubType     string
	URLS        []RecipeURL
	ExcludeURLS []RecipeURL
}

type RecipeURL struct {
	URL           string
	RegexpMatcher *regexp.Regexp
}

type RecipeURLMatch struct {
	DetectionURLPart string
	RecipeURL        string
	RecipeUUID       string
	RecipeName       string
	RecipeSubType    string
	ExcludedURL      bool
}

var ErrInvalidRecipes = errors.New("invalid interface recipe")
var ErrInvalidInternalDomainRegexp = errors.New("could not parse internal domains as regexp")

func (classification *Classification) Name() string {
	if classification.RecipeMatch {
		return classification.RecipeName
	} else {
		return classification.URL
	}
}

func New(config Config) (*Classifier, error) {
	// prepare regular expressions for recipes
	var preparedRecipes []Recipe
	for _, recipe := range config.Recipes {
		preparedRecipe := Recipe{
			UUID:    recipe.UUID,
			Name:    recipe.Name,
			Type:    recipe.Type,
			SubType: recipe.SubType,
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

		for _, excludedRecipeURL := range recipe.ExcludeURLS {
			regexpMatcher, err := url.PrepareRegexpMatcher(excludedRecipeURL)
			if err != nil {
				return nil, ErrInvalidRecipes
			}

			preparedRecipeURL := RecipeURL{
				URL:           excludedRecipeURL,
				RegexpMatcher: regexpMatcher,
			}
			preparedRecipe.ExcludeURLS = append(preparedRecipe.ExcludeURLS, preparedRecipeURL)
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
	if formatValidityCheck.State == classify.Invalid {
		return &ClassifiedInterface{
			Detection: &data,
			Classification: &Classification{
				URL: value,
				Decision: classify.ClassificationDecision{
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
				Decision: classify.ClassificationDecision{
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
		if recipeMatch.ExcludedURL {
			return &ClassifiedInterface{
				Detection: &data,
				Classification: &Classification{
					URL: value,
					Decision: classify.ClassificationDecision{
						State:  classify.Invalid,
						Reason: "ignored_url_in_recipe",
					},
				},
			}, nil
		}

		classifiedInterface := &ClassifiedInterface{
			Detection: &data,
			Classification: &Classification{
				URL:           recipeMatch.DetectionURLPart,
				RecipeMatch:   true,
				RecipeUUID:    recipeMatch.RecipeUUID,
				RecipeName:    recipeMatch.RecipeName,
				RecipeSubType: recipeMatch.RecipeSubType,
			},
		}
		if strings.Contains(recipeMatch.DetectionURLPart, "*") {
			classifiedInterface.Classification.Decision = classify.ClassificationDecision{
				State:  classify.Potential,
				Reason: "recipe_match_with_wildcard",
			}
			return classifiedInterface, nil
		}

		classifiedInterface.Classification.Decision = classify.ClassificationDecision{
			State:  classify.Valid,
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
			Decision: classify.ClassificationDecision{
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
		for _, recipeURL := range recipe.ExcludeURLS {
			match, err := url.Match(detectionURL, recipeURL.RegexpMatcher)
			if err != nil {
				return nil, err
			}

			if match != "" {
				return &RecipeURLMatch{
					ExcludedURL: true,
				}, nil
			}
		}

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
				RecipeUUID:       recipe.UUID,
				RecipeName:       recipe.Name,
				RecipeSubType:    recipe.SubType,
			}
		}
	}

	return recipeURLMatch, nil
}
