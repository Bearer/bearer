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
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

const LATEST_RELEASE_URL = "https://api.github.com/repos/bearer/bearer-rules/releases/latest"
const BASE_RULE_FOLDER = "/"

func LoadRuleDefinitionsFromGitHub(ruleDefinitions map[string]RuleDefinition, foundLanguages []string) (tagName string, err error) {
	resp, err := http.Get(LATEST_RELEASE_URL)
	if err != nil {
		return tagName, err
	}
	defer resp.Body.Close()
	headers := resp.Header

	if headers.Get("x-ratelimit-remaining") == "0" {
		resetString := headers.Get("x-ratelimit-reset")
		unixTime, err := strconv.ParseInt(resetString, 10, 64)
		if err != nil {
			return tagName, fmt.Errorf("rules download is rate limited. Please wait until: %s", resetString)
		}
		tm := time.Unix(unixTime, 0)
		return tagName, fmt.Errorf("rules download is rate limited. Please wait until: %s", tm.Format("2006-01-02 15:04:05"))
	}

	if resp.StatusCode != 200 {
		return tagName, errors.New("rules download returned non 200 status code - could not download rules")
	}

	// Decode the response JSON to get the URL of the asset we want to download
	type Asset struct {
		BrowserDownloadUrl string `json:"browser_download_url"`
		Name               string `json:"name"`
	}
	var release struct {
		Id      int     `json:"id"`
		TagName string  `json:"tag_name"`
		Assets  []Asset `json:"assets"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return tagName, err
	}

	if release.TagName == "" {
		return tagName, errors.New("could not find valid release for rules")
	}

	bearerRulesDir := bearerRulesDir()
	if _, err := os.Stat(bearerRulesDir); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(bearerRulesDir, os.ModePerm)
		if err != nil {
			return tagName, fmt.Errorf("could not create bearer-rules directory: %s", err)
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
				return tagName, err
			}
			defer resp.Body.Close()

			// Create file in rules dir
			filepath, err := filepath.Abs(filepath.Join(bearerRulesDir, "source-"+language+".tar.gz"))
			if err != nil {
				return tagName, err
			}
			file, err := os.Create(filepath)
			if err != nil {
				return tagName, err
			}
			defer file.Close()

			// Copy the contents of the downloaded archive to the file
			if _, err := io.Copy(file, resp.Body); err != nil {
				return tagName, err
			}

			// reset file pointer to start of file
			_, err = file.Seek(0, 0)
			if err != nil {
				return tagName, err
			}

			if _, err = ReadRuleDefinitions(ruleDefinitions, file); err != nil {
				return tagName, err
			}
		}
	}

	return release.TagName, nil
}

func ReadRuleDefinitions(ruleDefinitions map[string]RuleDefinition, file *os.File) (map[string]bool, error) {
	ruleLanguages := make(map[string]bool)

	gzr, err := gzip.NewReader(file)
	if err != nil {
		return ruleLanguages, err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return ruleLanguages, err
		}

		if !isRuleFile(header.Name) {
			continue
		}

		data := make([]byte, header.Size)
		_, err = io.ReadFull(tr, data)
		if err != nil {
			return ruleLanguages, fmt.Errorf("failed to read file %s: %w", header.Name, err)
		}

		var ruleDefinition RuleDefinition
		err = yaml.Unmarshal(data, &ruleDefinition)
		if err != nil {
			return ruleLanguages, fmt.Errorf("failed to unmarshal rule %s: %w", header.Name, err)
		}

		id := ruleDefinition.Metadata.ID
		_, ruleExists := ruleDefinitions[id]
		if ruleExists {
			return ruleLanguages, fmt.Errorf("duplicate built-in rule ID %s", id)
		}

		for _, lang := range ruleDefinition.Languages {
			ruleLanguages[lang] = true
		}

		ruleDefinitions[id] = ruleDefinition
	}

	return ruleLanguages, nil
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
