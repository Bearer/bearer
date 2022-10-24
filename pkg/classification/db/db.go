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

func DefaultDataTypes() []DataType {
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

		dataType, err := UnmarshalDataType(val)
		if err != nil {
			fmt.Println(err)
			continue
		}

		dataTypes = append(dataTypes, *dataType)
	}

	return dataTypes
}

func DefaultDataTypeClassificationPatterns() []DataTypeClassificationPattern {
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

		dataTypeClassificationPattern, err := UnmarshalDataTypeClassificationPattern(val)
		if err != nil {
			log.Fatalf("failed to unmarshal data type classification #%v", dataTypeClassificationPatterns)
			fmt.Println(err)
			continue
		}

		dataTypeClassificationPatterns = append(dataTypeClassificationPatterns, *dataTypeClassificationPattern)
	}

	return dataTypeClassificationPatterns
}

func DefaultKnownPersonObjectPatterns() []KnownPersonObjectPattern {
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

		knownPersonObjectPattern, err := UnmarshalKnownPersonObjectPattern(val)
		if err != nil {
			fmt.Println(err)
			continue
		}

		knownPersonObjectPatterns = append(knownPersonObjectPatterns, *knownPersonObjectPattern)
	}

	return knownPersonObjectPatterns
}

func Default() []Recipe {
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

		recipe, err := UnmarshalRecipe(val)
		if err != nil {
			fmt.Println(err)
			continue
		}

		recipes = append(recipes, *recipe)
	}

	return recipes
}

func UnmarshalRecipe(rawBytes []byte) (*Recipe, error) {
	var db Recipe

	err := json.Unmarshal(rawBytes, &db)

	if err != nil {
		log.Fatalf("failed to unmarshal db %e", err)
		return nil, err
	}

	return &db, nil
}

func UnmarshalDataType(rawBytes []byte) (*DataType, error) {
	var db DataType

	err := json.Unmarshal(rawBytes, &db)

	if err != nil {
		log.Fatalf("failed to unmarshal db %e", err)
		return nil, err
	}

	return &db, nil
}

func UnmarshalDataTypeClassificationPattern(rawBytes []byte) (*DataTypeClassificationPattern, error) {
	var db DataTypeClassificationPattern

	err := json.Unmarshal(rawBytes, &db)

	if err != nil {
		log.Fatalf("failed to unmarshal db %e", err)
		return nil, err
	}

	return &db, nil
}

func UnmarshalKnownPersonObjectPattern(rawBytes []byte) (*KnownPersonObjectPattern, error) {
	var db KnownPersonObjectPattern

	err := json.Unmarshal(rawBytes, &db)

	if err != nil {
		log.Fatalf("failed to unmarshal db %e", err)
		return nil, err
	}

	return &db, nil
}
