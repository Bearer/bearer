package db

import (
	"embed"
	"encoding/json"
	"fmt"
	"log"

	"github.com/google/uuid"
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
	Id                        int       `json:"id"`
	DataTypeUUID              uuid.UUID `json:"data_type_uuid"`
	IncludeRegexp             string    `json:"include_regexp,omitempty"`
	ExcludeRegexp             string    `json:"exclude_regexp,omitempty"`
	IncludeTypes              []string  `json:"include_types,omitempty"`
	ExcludeTypes              []string  `json:"exclude_types,omitempty"`
	FriendlyName              string    `json:"friendly_name"`
	HealthContextDataTypeUUID string    `json:"health_context_data_type_uuid,omitempty"`
	MatchColumn               bool      `json:"match_column"`
	MatchObject               bool      `json:"match_object"`
	ObjectType                []string  `json:"object_type"`
}

type KnownPersonObjectPattern struct {
	Id              int    `json:"id"`
	IncludeRegexp   string `json:"include_regexp,omitempty"`
	Category        string `json:"category"`
	ActAsIdentifier bool   `json:"act_as_identifier"`
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
		log.Fatalln(err)
	}

	for _, file := range files {
		val, err := recipesDir.ReadFile("recipes/" + file.Name())
		if err != nil {
			fmt.Println(err)
			continue
		}

		var recipe Recipe
		rawBytes := []byte(val)
		err = json.Unmarshal(rawBytes, &recipe)
		if err != nil {
			log.Fatalln(err)
		}

		recipes = append(recipes, recipe)
	}

	return recipes
}

func defaultDataTypes() []DataType {
	dataTypes := []DataType{}

	files, err := dataTypesDir.ReadDir("data_types")
	if err != nil {
		log.Fatalln(err)
	}

	for _, file := range files {
		val, err := dataTypesDir.ReadFile("data_types/" + file.Name())
		if err != nil {
			fmt.Println(err)
			continue
		}

		var dataType DataType
		rawBytes := []byte(val)
		err = json.Unmarshal(rawBytes, &dataType)
		if err != nil {
			log.Fatalln(err)
		}

		dataTypes = append(dataTypes, dataType)
	}

	return dataTypes
}

func defaultDataTypeClassificationPatterns() []DataTypeClassificationPattern {
	dataTypeClassificationPatterns := []DataTypeClassificationPattern{}

	files, err := dataTypeClassificationPatternsDir.ReadDir("data_type_classification_patterns")
	if err != nil {
		log.Fatalln(err)
	}

	for _, file := range files {
		val, err := dataTypeClassificationPatternsDir.ReadFile("data_type_classification_patterns/" + file.Name())
		if err != nil {
			fmt.Println(err)
			continue
		}

		var dataTypeClassificationPattern DataTypeClassificationPattern
		rawBytes := []byte(val)
		err = json.Unmarshal(rawBytes, &dataTypeClassificationPattern)
		if err != nil {
			log.Fatalln(err)
		}

		dataTypeClassificationPatterns = append(dataTypeClassificationPatterns, dataTypeClassificationPattern)
	}

	return dataTypeClassificationPatterns
}

func defaultKnownPersonObjectPatterns() []KnownPersonObjectPattern {
	knownPersonObjectPatterns := []KnownPersonObjectPattern{}

	files, err := knownPersonObjectPatternsDir.ReadDir("known_person_object_patterns")
	if err != nil {
		log.Fatalln(err)
	}

	for _, file := range files {
		val, err := knownPersonObjectPatternsDir.ReadFile("known_person_object_patterns/" + file.Name())
		if err != nil {
			fmt.Println(err)
			continue
		}

		var knownPersonObjectPattern KnownPersonObjectPattern
		rawBytes := []byte(val)
		err = json.Unmarshal(rawBytes, &knownPersonObjectPattern)
		if err != nil {
			log.Fatalln(err)
		}

		knownPersonObjectPatterns = append(knownPersonObjectPatterns, knownPersonObjectPattern)
	}

	return knownPersonObjectPatterns
}
