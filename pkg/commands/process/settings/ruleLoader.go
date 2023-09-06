package settings

import (
	"archive/tar"
	"compress/gzip"
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

const BASE_RULE_FOLDER = "/"

func LoadRuleDefinitionsFromUrls(ruleDefinitions map[string]RuleDefinition, languageDownloads []string) (err error) {

	bearerRulesDir := bearerRulesDir()
	if _, err := os.Stat(bearerRulesDir); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(bearerRulesDir, os.ModePerm)
		if err != nil {
			return fmt.Errorf("could not create bearer-rules directory: %s", err)
		}
	}

	for _, languagePackageUrl := range languageDownloads {
		// Prepare filepath
		urlHash := md5.Sum([]byte(languagePackageUrl))
		filepath, err := filepath.Abs(filepath.Join(bearerRulesDir, fmt.Sprintf("%x.tar.gz", urlHash)))

		if err != nil {
			return err
		}

		if _, err := os.Stat(filepath); err == nil {
			log.Trace().Msgf("Using local cache for rule package: %s", languagePackageUrl)
			file, err := os.Open(filepath)
			if err != nil {
				return err
			}
			defer file.Close()

			if _, err = ReadRuleDefinitions(ruleDefinitions, file); err != nil {
				return err
			}
		} else {
			log.Trace().Msgf("Downloading rule package: %s", languagePackageUrl)
			httpClient := &http.Client{Timeout: 60 * time.Second}
			resp, err := httpClient.Get(languagePackageUrl)
			if err != nil {
				return err
			}
			defer resp.Body.Close()

			// Create file in rules dir
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

			if _, err = ReadRuleDefinitions(ruleDefinitions, file); err != nil {
				return err
			}
		}
	}

	return nil
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
