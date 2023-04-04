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

	"github.com/bearer/bearer/pkg/util/output"
	"gopkg.in/yaml.v3"
)

const LATEST_RELEASE_URL = "https://api.github.com/repos/bearer/bearer-rules/releases/latest"
const BASE_RULE_FOLDER = "/"

func LoadRuleDefinitionsFromGitHub(ruleDefinitions map[string]RuleDefinition, force bool, quiet bool) error {
	resp, err := http.Get(LATEST_RELEASE_URL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Decode the response JSON to get the URL of the asset we want to download
	var release struct {
		Id         int    `json:"id"`
		TarballUrl string `json:"tarball_url"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return err
	}

	var ruleTarball *os.File

	if release.Id == 0 {
		return fmt.Errorf("could not find ID for latest release")
	}

	bearerRulesDir := os.TempDir() + "bearer-rules/"

	fileExists := fileExistsForReleaseId(release.Id, bearerRulesDir)
	if fileExists {
		// cached & up-to-date rule tarball found
		ruleTarball, err = os.Open(formatFileName(release.Id, bearerRulesDir))
		if err != nil {
			return fmt.Errorf("could not open file %s", formatFileName(release.Id, bearerRulesDir))
		}
		defer ruleTarball.Close()
		if !quiet {
			output.StdErrLogger().Msgf("Using cached rules")
		}
	}

	if ruleTarball == nil || force {
		if !quiet {
			output.StdErrLogger().Msgf("Downloading rules from source")
		}
		// either no cached rule tarballs found or more recent rules version available
		// therefore : download tarball from GitHub

		// create dir if it doesn't exist
		if _, err := os.Stat(bearerRulesDir); errors.Is(err, os.ErrNotExist) {
			err := os.Mkdir(bearerRulesDir, 0700)
			if err != nil {
				return fmt.Errorf("could not create bearer-rules directory: %s", err)
			}
		} else {
			// clean up any (now-outdated) files in the bearer-rules dir
			err = cleanupRuleDirFiles(bearerRulesDir)
			if err != nil {
				return fmt.Errorf("could not clean up bearer-rules dir files: %s", err)
			}
		}

		if release.TarballUrl == "" {
			return fmt.Errorf("could not find source.tar.gz asset in latest release")
		}

		resp, err = http.Get(release.TarballUrl)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		filepath, err := filepath.Abs(formatFileName(release.Id, bearerRulesDir))
		if err != nil {
			return err
		}
		ruleTarball, err = os.Create(filepath)
		if err != nil {
			return err
		}
		defer ruleTarball.Close()

		// Copy the contents of the downloaded archive to the temporary file
		if _, err := io.Copy(ruleTarball, resp.Body); err != nil {
			return err
		}

		// reset file pointer to start of file
		_, err = ruleTarball.Seek(0, 0)
		if err != nil {
			return err
		}
	}

	gzr, err := gzip.NewReader(ruleTarball)
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

		// TODO: only load rules for detected language(s)
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

func fileExistsForReleaseId(releaseId int, bearerRulesDir string) bool {
	_, err := os.Stat(formatFileName(releaseId, bearerRulesDir))

	return err == nil
}

func formatFileName(releaseId int, bearerRulesDir string) string {
	return bearerRulesDir + fmt.Sprintf("source-%d.tar.gz", releaseId)
}

func cleanupRuleDirFiles(bearerRulesDir string) error {
	files, err := filepath.Glob(filepath.Join(bearerRulesDir, "*"))
	if err != nil {
		return err
	}
	for _, file := range files {
		err = os.RemoveAll(file)
		if err != nil {
			return err
		}
	}
	return nil

}
