package queries

import (
	"bytes"
	"encoding/json"
	"os"

	"github.com/bearer/bearer/pkg/report/operations"
	"github.com/bearer/bearer/pkg/util/file"
	"gopkg.in/yaml.v3"
)

type Document struct {
	Servers []Url `yaml:"servers" json:"servers"`
}

type Url struct {
	Url       string              `yaml:"url" json:"url"`
	Variables map[string]Variable `yaml:"variables" json:"variables"`
}

type Variable struct {
	Name    string   `yaml:"-" json:""`
	Values  []string `yaml:"enum" json:"enum"`
	Default string   `yaml:"default" json:"default"`
}

func FindUrls(file *file.FileInfo) (urls []operations.Url) {
	fileBytes, err := os.ReadFile(file.AbsolutePath)
	if err != nil {
		return nil
	}

	var document Document

	err = yaml.NewDecoder(bytes.NewBuffer(fileBytes)).Decode(&document)

	if err != nil {
		json.NewDecoder(bytes.NewBuffer(fileBytes)).Decode(&document) //nolint:all,errcheck
	}

	for _, url := range document.Servers {

		returnedUrl := operations.Url{
			Url: url.Url,
		}

		for key, variable := range url.Variables {
			variable.Name = key

			values := variable.Values

			defaultExists := false

			for _, value := range values {
				if value == variable.Default {
					defaultExists = true
					break
				}
			}

			if !defaultExists && variable.Default != "" {
				values = append(values, variable.Default)
			}

			returnedUrl.Variables = append(returnedUrl.Variables, operations.Variable{
				Name:   variable.Name,
				Values: values,
			})

		}

		urls = append(urls, returnedUrl)
	}

	return
}
