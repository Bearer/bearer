package db

import (
	"embed"
	_ "embed"
	"encoding/json"
	"fmt"

	"github.com/bearer/curio/battle_tests/build"
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

type ItemWithLanguage struct {
	FullName string `json:"full_name" yaml:"full_name"`
	HtmlUrl  string `json:"html_url" yaml:"html_url"`
	Language string `json:"language" yaml:"language"`
}

func UnmarshalRaw() []ItemWithLanguage {
	items := []ItemWithLanguage{}

	files, err := githubDir.ReadDir("github")
	if err != nil {
		log.Fatal().Err(fmt.Errorf("unable to open directory %e", err)).Send()
	}

	for _, file := range files {
		if build.Language != "all" && file.Name() != build.Language {
			continue
		}

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

		for _, item := range category.Items {
			newItem := ItemWithLanguage{
				HtmlUrl:  item.HtmlUrl,
				FullName: item.FullName,
				Language: category.Language,
			}

			items = append(items, newItem)
		}
	}

	return items
}
