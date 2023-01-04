package settings

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

func LoadExternalPolicies(directories []string) (map[string]*Policy, error) {
	policies := make(map[string]*Policy)
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

			policyName := strings.TrimSuffix(fileName, ext)

			var policy *Policy
			err = yaml.Unmarshal(fileContent, &policy)
			if err != nil {
				return fmt.Errorf("failed to unmarshal yaml file: %s %s", filePath, err)
			}

			for _, module := range policy.Modules {
				if module.Path != "" {
					dirPath := strings.TrimSuffix(filePath, fileName)
					modulePath := dirPath + "/" + module.Path
					moduleContent, err := os.ReadFile(modulePath)
					if err != nil {
						return fmt.Errorf("failed to read module at path %s %s", modulePath, err)
					}
					module.Content = string(moduleContent)
				}
			}

			policies[policyName] = policy

			return nil
		})

		if err != nil {
			return nil, err
		}

		if err != nil {
			return nil, err
		}
	}

	return policies, nil
}
