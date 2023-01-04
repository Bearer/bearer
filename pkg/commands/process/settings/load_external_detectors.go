package settings

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

func LoadExternalDetectors(directories []string) (map[string]Rule, error) {
	rules := make(map[string]Rule)

	for _, dirPath := range directories {

		err := filepath.WalkDir(dirPath, func(filePath string, d fs.DirEntry, errReading error) error {
			if errReading != nil {
				return errReading
			}

			fileName := d.Name()
			ext := filepath.Ext(fileName)

			if d.IsDir() {
				return nil
			}

			if ext != ".yaml" && ext != ".yml" {
				return nil
			}

			fileContent, err := os.ReadFile(filePath)
			if err != nil {
				return fmt.Errorf("error reading file: %s %s", filePath, err)
			}

			ruleName := strings.TrimSuffix(fileName, ext)

			var rule Rule
			err = yaml.Unmarshal(fileContent, &rule)
			if err != nil {
				return fmt.Errorf("failed to unmarshal yaml file: %s %s", filePath, err)
			}

			rules[ruleName] = rule

			return nil
		})

		if err != nil {
			return nil, err
		}
	}

	return rules, nil
}
