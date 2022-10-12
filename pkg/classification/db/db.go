package db

import (
	"embed"
	"encoding/json"
	"fmt"
	"log"
)

//go:embed recipes
var recipesDir embed.FS

type Recipe struct {
	URLS     []string  `json:"urls"`
	Name     string    `json:"name"`
	Type     string    `json:"size"`
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

		recipe, err := Unmarshal(val)
		if err != nil {
			fmt.Println(err)
			continue
		}

		fmt.Printf("%+v\n", recipe)

		recipes = append(recipes, *recipe)
	}

	return recipes
}

func Unmarshal(rawBytes []byte) (*Recipe, error) {
	var db Recipe

	err := json.Unmarshal(rawBytes, &db)

	if err != nil {
		log.Fatalf("failed to unmarshal db %e", err)
		return nil, err
	}

	return &db, nil
}
