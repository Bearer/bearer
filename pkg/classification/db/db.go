package db

import (
	"embed"
	"encoding/json"
	"log"
	"regexp"
	"strings"

	"github.com/bearer/curio/pkg/flag"
	"github.com/tangzero/inflector"
)

var PHIDataCategoryGroupUUID = "247fa503-115b-490a-96e5-bcd357bd5686"

//go:embed recipes
var recipesDir embed.FS

//go:embed data_types
var dataTypesDir embed.FS

//go:embed data_categories
var dataCategoriesDir embed.FS

//go:embed data_type_classification_patterns
var dataTypeClassificationPatternsDir embed.FS

//go:embed known_person_object_patterns
var knownPersonObjectPatternsDir embed.FS

//go:embed category_grouping.json
var categoryGroupingFile embed.FS

type DefaultDB struct {
	Recipes                        []Recipe
	DataTypes                      []DataType
	DataCategories                 []DataCategory
	DataTypeClassificationPatterns []DataTypeClassificationPattern
	KnownPersonObjectPatterns      []KnownPersonObjectPattern
}

type Recipe struct {
	URLS     []string  `json:"urls" yaml:"urls"`
	Name     string    `json:"name" yaml:"name"`
	Type     string    `json:"type" yaml:"type"`
	Packages []Package `json:"packages" yaml:"packages"`
	UUID     string    `json:"uuid" yaml:"uuid"`
}

type Package struct {
	Name           string `json:"name" yaml:"name"`
	PackageManager string `json:"package_manager" yaml:"package_manager"`
	Group          string `json:"group" yaml:"group"`
}

type RecipeType string

var RecipeTypeDataStore = RecipeType("data_store")
var RecipeTypeService = RecipeType("service")

type DataType struct {
	Name         string `json:"name" yaml:"name"`
	UUID         string `json:"uuid" yaml:"uuid"`
	CategoryUUID string `json:"category_uuid" yaml:"category_uuid"`
}

type DataCategory struct {
	Name   string                       `json:"name" yaml:"name"`
	UUID   string                       `json:"uuid" yaml:"uuid"`
	Groups map[string]DataCategoryGroup `json:"groups" yaml:"groups"`
}

type DataCategoryGroup struct {
	Name string `json:"name" yaml:"name"`
	UUID string `json:"uuid,omitempty" yaml:"uuid,omitempty"`
}

type DataCategoryGrouping struct {
	Groups map[string]struct {
		Name        string   `json:"name" yaml:"name"`
		ParentUUIDs []string `json:"parent_uuids,omitempty" yaml:"parent_uuids,omitempty"`
	} `json:"groups"`
	CategoryMapping map[string]struct {
		Name       string   `json:"name" yaml:"name"`
		GroupUUIDs []string `json:"group_uuids" yaml:"group_uuids"`
	} `json:"category_mapping"`
}

type ObjectType string

var KnownObject ObjectType = "known"
var ExtendedUnknownObject ObjectType = "unknown_extended"
var UnknownObject ObjectType = "unknown"
var AssociatedObject ObjectType = "associated"
var KnownDataObject ObjectType = "known_data_object"

type DataTypeClassificationPattern struct {
	Id                        int                 `json:"id" yaml:"id"`
	DataTypeUUID              string              `json:"data_type_uuid,omitempty"`
	DataType                  DataType            `json:"data_type" yaml:"data_type"`
	IncludeRegexp             string              `json:"include_regexp" yaml:"include_regexp"`
	IncludeRegexpMatcher      *regexp.Regexp      `json:"include_regexp_matcher" yaml:"include_regexp_matcher"`
	ExcludeRegexp             string              `json:"exclude_regexp,omitempty"`
	ExcludeRegexpMatcher      *regexp.Regexp      `json:"exclude_regexp_matcher" yaml:"exclude_regexp_matcher"`
	ExcludeTypes              []string            `json:"exclude_types" yaml:"exclude_types"`
	ExcludeTypesMapping       map[string]struct{} `json:"exclude_types_mapping" yaml:"exclude_types_mapping"`
	FriendlyName              string              `json:"friendly_name" yaml:"friendly_name"`
	HealthContextDataTypeUUID string              `json:"health_context_data_type_uuid,omitempty"`
	HealthContextDataType     DataType            `json:"health_context_data_type" yaml:"health_context_data_type"`
	MatchColumn               bool                `json:"match_column" yaml:"match_column"`
	MatchObject               bool                `json:"match_object" yaml:"match_object"`
	ObjectType                []string            `json:"object_type" yaml:"object_type"`
	ObjectTypeMapping         map[string]struct{} `json:"object_types_mapping" yaml:"object_types_mapping"`
}

type KnownPersonObjectPattern struct {
	Id                      int            `json:"id" yaml:"id"`
	DataType                DataType       `json:"data_type" yaml:"data_type"`
	IncludeRegexp           string         `json:"include_regexp" yaml:"include_regexp"`
	IncludeRegexpMatcher    *regexp.Regexp `json:"include_regexp_matcher" yaml:"include_regexp_matcher"`
	ExcludeRegexp           string         `json:"exclude_regexp,omitempty"`
	ExcludeRegexpMatcher    *regexp.Regexp `json:"exclude_regexp_matcher" yaml:"exclude_regexp_matcher"`
	Category                string         `json:"category" yaml:"category"`
	ActAsIdentifier         bool           `json:"act_as_identifier" yaml:"act_as_identifier"`
	IdentifierRegexpMatcher *regexp.Regexp `json:"identifier_regexp_matcher" yaml:"identifier_regexp_matcher"`
}

func Default() DefaultDB {
	return defaultDB("")
}

func DefaultWithContext(context flag.Context) DefaultDB {
	return defaultDB(context)
}

func defaultDB(context flag.Context) DefaultDB {
	dataTypes := defaultDataTypes()
	return DefaultDB{
		Recipes:                        defaultRecipes(),
		DataTypes:                      dataTypes,
		DataCategories:                 defaultDataCategories(context),
		DataTypeClassificationPatterns: defaultDataTypeClassificationPatterns(dataTypes),
		KnownPersonObjectPatterns:      defaultKnownPersonObjectPatterns(dataTypes),
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

func defaultDataCategories(context flag.Context) []DataCategory {
	skipHealthContext := true
	if context == flag.Health {
		skipHealthContext = false
	}

	dataCategories := []DataCategory{}

	categoryGroupingJson, err := categoryGroupingFile.ReadFile("category_grouping.json")
	if err != nil {
		handleError(err)
	}

	var dataCategoryGrouping DataCategoryGrouping
	rawBytes := []byte(categoryGroupingJson)
	err = json.Unmarshal(rawBytes, &dataCategoryGrouping)
	if err != nil {
		handleError(err)
	}

	files, err := dataCategoriesDir.ReadDir("data_categories")
	if err != nil {
		handleError(err)
	}

	for _, file := range files {
		val, err := dataCategoriesDir.ReadFile("data_categories/" + file.Name())
		if err != nil {
			handleError(err)
		}

		var dataCategory DataCategory
		rawBytes := []byte(val)
		err = json.Unmarshal(rawBytes, &dataCategory)
		if err != nil {
			handleError(err)
		}

		// Add all category groups
		dataCategory.Groups = make(map[string]DataCategoryGroup)
		categoryFromMapping := dataCategoryGrouping.CategoryMapping[dataCategory.UUID]
		for _, groupUUID := range categoryFromMapping.GroupUUIDs {
			if skipHealthContext && groupUUID == PHIDataCategoryGroupUUID {
				continue // skip health context
			}
			group := dataCategoryGrouping.Groups[groupUUID]
			dataCategory.Groups[groupUUID] = DataCategoryGroup{
				Name: group.Name,
				UUID: groupUUID,
			}
			// add parent group if present
			for _, parentUUID := range group.ParentUUIDs {
				dataCategory.Groups[parentUUID] = DataCategoryGroup{
					Name: group.Name,
					UUID: parentUUID,
				}
			}
		}

		dataCategories = append(dataCategories, dataCategory)
	}

	return dataCategories
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

func defaultDataTypeClassificationPatterns(dataTypes []DataType) []DataTypeClassificationPattern {
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

		// add data type and health context data type
		for _, dataType := range dataTypes {
			if dataType.UUID == dataTypeClassificationPattern.DataTypeUUID {
				dataTypeClassificationPattern.DataType = dataType
				break
			}
		}

		for _, dataType := range dataTypes {
			if dataType.UUID == dataTypeClassificationPattern.HealthContextDataTypeUUID {
				dataTypeClassificationPattern.HealthContextDataType = dataType
				break
			}
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

func defaultKnownPersonObjectPatterns(dataTypes []DataType) []KnownPersonObjectPattern {
	knownPersonObjectPatterns := []KnownPersonObjectPattern{}

	// "Identification" > "Unique Identifier" data type
	// Applies to all known person object patterns e.g.
	// "profile", "user", "supplier", etc
	uniqueIdentifierDataTypeUUID := "12d44ae0-1df7-4faf-9fb1-b46cc4b4dce9"
	var uniqueIdentifierDataType DataType
	for _, dataType := range dataTypes {
		if dataType.UUID == uniqueIdentifierDataTypeUUID {
			uniqueIdentifierDataType = dataType
			break
		}
	}
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

		// add data type UUID and data type
		knownPersonObjectPattern.DataType = uniqueIdentifierDataType

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
