package db

import (
	"embed"
	"encoding/json"
	"log"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"github.com/tangzero/inflector"
)

//go:embed recipes
var recipesDir embed.FS

//go:embed data_types
var dataTypesDir embed.FS

//go:embed data_type_classification_patterns
var dataTypeClassificationPatternsDir embed.FS

//go:embed known_person_object_patterns
var knownPersonObjectPatternsDir embed.FS

type DefaultDB struct {
	Recipes                        []Recipe
	DataTypes                      []DataType
	DataTypeClassificationPatterns []DataTypeClassificationPattern
	KnownPersonObjectPatterns      []KnownPersonObjectPattern
}

type Recipe struct {
	URLS     []string  `json:"urls"`
	Name     string    `json:"name"`
	Type     string    `json:"type"`
	Packages []Package `json:"packages"`
}

type Package struct {
	Name           string `json:"name"`
	PackageManager string `json:"package_manager"`
	Group          string `json:"group"`
}

type RecipeType string

var RecipeTypeDataStore = RecipeType("data_store")
var RecipeTypeService = RecipeType("service")

type DataType struct {
	DataCategoryName string    `json:"data_category_name"`
	DefaultCategory  string    `json:"default_category"`
	Id               int       `json:"id"`
	UUID             uuid.UUID `json:"uuid"`
}

type DataTypeClassificationPattern struct {
	Id                        int                 `json:"id"`
	DataTypeUUID              *uuid.UUID          `json:"data_type_uuid,omitempty"`
	IncludeRegexp             string              `json:"include_regexp"`
	IncludeRegexpMatcher      *regexp.Regexp      `json:"include_regexp_matcher"`
	ExcludeRegexp             string              `json:"exclude_regexp,omitempty"`
	ExcludeRegexpMatcher      *regexp.Regexp      `json:"exclude_regexp_matcher"`
	ExcludeTypes              []string            `json:"exclude_types"`
	ExcludeTypesMapping       map[string]struct{} `json:"exclude_types_mapping"`
	FriendlyName              string              `json:"friendly_name"`
	HealthContextDataTypeUUID string              `json:"health_context_data_type_uuid,omitempty"`
	MatchColumn               bool                `json:"match_column"`
	MatchObject               bool                `json:"match_object"`
	ObjectType                []string            `json:"object_type"`
	ObjectTypeMapping         map[string]struct{} `json:"object_types_mapping"`
}

type KnownPersonObjectPattern struct {
	Id                      int            `json:"id"`
	IncludeRegexp           string         `json:"include_regexp"`
	IncludeRegexpMatcher    *regexp.Regexp `json:"include_regexp_matcher"`
	ExcludeRegexp           string         `json:"exclude_regexp,omitempty"`
	ExcludeRegexpMatcher    *regexp.Regexp `json:"exclude_regexp_matcher"`
	Category                string         `json:"category"`
	ActAsIdentifier         bool           `json:"act_as_identifier"`
	IdentifierRegexpMatcher *regexp.Regexp `json:"identifier_regexp_matcher"`
}

func Default() DefaultDB {
	return DefaultDB{
		Recipes:                        defaultRecipes(),
		DataTypes:                      defaultDataTypes(),
		DataTypeClassificationPatterns: defaultDataTypeClassificationPatterns(),
		KnownPersonObjectPatterns:      defaultKnownPersonObjectPatterns(),
	}
}

func defaultRecipes() []Recipe {
	recipes := []Recipe{}

	files, err := recipesDir.ReadDir("recipes")
	if err != nil {
		handleError(err)
	}

	for _, file := range files {
		val, err := recipesDir.ReadFile("recipes/" + file.Name())
		if err != nil {
			handleError(err)
		}

		var recipe Recipe
		rawBytes := []byte(val)
		err = json.Unmarshal(rawBytes, &recipe)
		if err != nil {
			handleError(err)
		}

		recipes = append(recipes, recipe)
	}

	return recipes
}

func defaultDataTypes() []DataType {
	dataTypes := []DataType{}

	files, err := dataTypesDir.ReadDir("data_types")
	if err != nil {
		handleError(err)
	}

	for _, file := range files {
		val, err := dataTypesDir.ReadFile("data_types/" + file.Name())
		if err != nil {
			handleError(err)
		}

		var dataType DataType
		rawBytes := []byte(val)
		err = json.Unmarshal(rawBytes, &dataType)
		if err != nil {
			handleError(err)
		}

		dataTypes = append(dataTypes, dataType)
	}

	return dataTypes
}

func defaultDataTypeClassificationPatterns() []DataTypeClassificationPattern {
	dataTypeClassificationPatterns := []DataTypeClassificationPattern{}

	files, err := dataTypeClassificationPatternsDir.ReadDir("data_type_classification_patterns")
	if err != nil {
		handleError(err)
	}

	for _, file := range files {
		val, err := dataTypeClassificationPatternsDir.ReadFile("data_type_classification_patterns/" + file.Name())
		if err != nil {
			handleError(err)
		}

		var dataTypeClassificationPattern DataTypeClassificationPattern
		rawBytes := []byte(val)
		err = json.Unmarshal(rawBytes, &dataTypeClassificationPattern)
		if err != nil {
			handleError(err)
		}

		// compile regexp matchers
		dataTypeClassificationPattern.IncludeRegexpMatcher, err = regexp.Compile(dataTypeClassificationPattern.IncludeRegexp)
		if err != nil {
			handleError(err)
		}
		if dataTypeClassificationPattern.ExcludeRegexp != "" {
			dataTypeClassificationPattern.ExcludeRegexpMatcher, err = regexp.Compile(dataTypeClassificationPattern.ExcludeRegexp)
			if err != nil {
				handleError(err)
			}
		}

		// add mappings for performant inclusion checks
		dataTypeClassificationPattern.ExcludeTypesMapping = map[string]struct{}{}
		for _, excludeType := range dataTypeClassificationPattern.ExcludeTypes {
			dataTypeClassificationPattern.ExcludeTypesMapping[excludeType] = struct{}{}
		}
		dataTypeClassificationPattern.ObjectTypeMapping = map[string]struct{}{}
		for _, objectType := range dataTypeClassificationPattern.ObjectType {
			dataTypeClassificationPattern.ObjectTypeMapping[objectType] = struct{}{}
		}

		dataTypeClassificationPatterns = append(dataTypeClassificationPatterns, dataTypeClassificationPattern)
	}

	return dataTypeClassificationPatterns
}

func defaultKnownPersonObjectPatterns() []KnownPersonObjectPattern {
	knownPersonObjectPatterns := []KnownPersonObjectPattern{}

	files, err := knownPersonObjectPatternsDir.ReadDir("known_person_object_patterns")
	if err != nil {
		handleError(err)
	}

	for _, file := range files {
		val, err := knownPersonObjectPatternsDir.ReadFile("known_person_object_patterns/" + file.Name())
		if err != nil {
			handleError(err)
		}

		var knownPersonObjectPattern KnownPersonObjectPattern
		rawBytes := []byte(val)
		err = json.Unmarshal(rawBytes, &knownPersonObjectPattern)
		if err != nil {
			handleError(err)
		}

		// compile regexp matchers
		knownPersonObjectPattern.IncludeRegexpMatcher, err = regexp.Compile(knownPersonObjectPattern.IncludeRegexp)
		if err != nil {
			handleError(err)
		}
		if knownPersonObjectPattern.ExcludeRegexp != "" {
			knownPersonObjectPattern.ExcludeRegexpMatcher, err = regexp.Compile(knownPersonObjectPattern.ExcludeRegexp)
			if err != nil {
				handleError(err)
			}
		}
		if knownPersonObjectPattern.ActAsIdentifier {
			category := strings.ToLower(knownPersonObjectPattern.Category)
			pluralCategory := inflector.Pluralize(category)

			knownPersonObjectPattern.IdentifierRegexpMatcher, err = regexp.Compile("(?i)^[\\S]*(" + category + "|" + pluralCategory + ")\\s?(uu)?id")

			if err != nil {
				handleError(err)
			}
		}

		knownPersonObjectPatterns = append(knownPersonObjectPatterns, knownPersonObjectPattern)
	}

	return knownPersonObjectPatterns
}

func handleError(err error) {
	log.Fatalln(err)
}
