package db

import (
	"embed"
	_ "embed"
	"encoding/json"
	"fmt"

	"github.com/rs/zerolog/log"
)

//go:embed github
var githubDir embed.FS

type Category struct {
	Language string `json:"language" yaml:"language"`
	Items    []Item `json:"items" yaml:"items"`
}

type Item struct {
	FullName string `json:"full_name" yaml:"full_name"`
	HtmlUrl  string `json:"html_url" yaml:"html_url"`
}

func UnmarshalRaw() []Item {
	items := []Item{}

	files, err := githubDir.ReadDir("github")
	if err != nil {
		log.Fatal().Err(fmt.Errorf("unable to open directory %e", err)).Send()
	}

	for _, file := range files {
		val, err := githubDir.ReadFile("github/" + file.Name())
		if err != nil {
			log.Fatal().Err(fmt.Errorf("unable to open file %e", err)).Send()
		}

		var category Category
		rawBytes := []byte(val)
		err = json.Unmarshal(rawBytes, &category)
		if err != nil {
			log.Fatal().Err(fmt.Errorf("unable to unmarshal %e", err)).Send()
		}

		items = append(items, category.Items...)
	}

	return items
}

func Unmarshal(rawBytes []byte) Category {
	var category Category

	err := json.Unmarshal(rawBytes, &category)
	if err != nil {
		log.Fatal().Err(fmt.Errorf("failed to unmarshal category file %e", err)).Send()
	}

	return category
}
