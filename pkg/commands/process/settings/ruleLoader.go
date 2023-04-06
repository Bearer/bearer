package settings

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

const LATEST_RELEASE_URL = "https://api.github.com/repos/bearer/bearer-rules/releases/latest"
const BASE_RULE_FOLDER = "/"

func LoadRuleDefinitionsFromGitHub(ruleDefinitions map[string]RuleDefinition, foundLanguages []string) error {
	resp, err := http.Get(LATEST_RELEASE_URL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Decode the response JSON to get the URL of the asset we want to download
	type Asset struct {
		BrowserDownloadUrl string `json:"browser_download_url"`
		Name               string `json:"name"`
	}
	var release struct {
		Id     int     `json:"id"`
		Assets []Asset `json:"assets"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return err
	}

	bearerRulesDir := bearerRulesDir()
	if _, err := os.Stat(bearerRulesDir); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(bearerRulesDir, os.ModePerm)
		if err != nil {
			return fmt.Errorf("could not create bearer-rules directory: %s", err)
		}
	}

	// loop assets and download tarballs for found languages
	for _, asset := range release.Assets {
		// we aren't expecting many found languages per repo
		for _, language := range foundLanguages {
			if asset.Name != language+".tar.gz" {
				continue
			}

			resp, err = http.Get(asset.BrowserDownloadUrl)
			if err != nil {
				return err
			}
			defer resp.Body.Close()

			// Create file in rules dir
			filepath, err := filepath.Abs(filepath.Join(bearerRulesDir, "source-"+language+".tar.gz"))
			if err != nil {
				return err
			}
			file, err := os.Create(filepath)
			if err != nil {
				return err
			}
			defer file.Close()

			// Copy the contents of the downloaded archive to the file
			if _, err := io.Copy(file, resp.Body); err != nil {
				return err
			}

			// reset file pointer to start of file
			_, err = file.Seek(0, 0)
			if err != nil {
				return err
			}

			if err = ReadRuleDefinitions(ruleDefinitions, file); err != nil {
				return err
			}
		}
	}

	return nil
}

func ReadRuleDefinitions(ruleDefinitions map[string]RuleDefinition, file *os.File) error {
	gzr, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		if !isRuleFile(header.Name) {
			continue
		}

		data := make([]byte, header.Size)
		_, err = io.ReadFull(tr, data)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %w", header.Name, err)
		}

		var ruleDefinition RuleDefinition
		err = yaml.Unmarshal(data, &ruleDefinition)
		if err != nil {
			return fmt.Errorf("failed to unmarshal rule %s: %w", header.Name, err)
		}

		id := ruleDefinition.Metadata.ID
		_, ruleExists := ruleDefinitions[id]
		if ruleExists {
			return fmt.Errorf("duplicate built-in rule ID %s", id)
		}

		ruleDefinitions[id] = ruleDefinition
	}

	return nil
}

func isRuleFile(headerName string) bool {
	if strings.Contains(headerName, ".snapshots") {
		return false
	}

	ext := filepath.Ext(headerName)
	if ext != ".yaml" && ext != ".yml" {
		return false
	}

	return strings.Contains(headerName, BASE_RULE_FOLDER)
}
