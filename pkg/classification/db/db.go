package db

import (
	_ "embed"
	"encoding/json"
	"log"
)

//go:embed recipes.json
var Raw []byte

type Recipe struct {
	URLS     []string  `json:"urls"`
	Name     string    `json:"name"`
	Type     string    `json:"size"`
	Packages []Package `json:"packages"`
}

type Package struct {
	Name           string `json:"name"`
	PackageManager string `json:"package_manager"`
}

type RecipeType string

var RecipeTypeDataStore = RecipeType("data_store")
var RecipeTypeService = RecipeType("service")

func Default() []Recipe {
	return Unmarshal(Raw)
}

func Unmarshal(rawBytes []byte) []Recipe {
	var db []Recipe

	err := json.Unmarshal(rawBytes, &db)
	if err != nil {
		log.Fatalf("failed to unmarshall db %e", err)
	}

	return db
}
