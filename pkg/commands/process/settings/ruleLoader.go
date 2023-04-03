package settings

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
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

func LoadRuleDefinitionsFromGitHub(ruleDefinitions map[string]RuleDefinition) error {
	resp, err := http.Get(LATEST_RELEASE_URL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Decode the response JSON to get the URL of the asset we want to download
	var release struct {
		TarballUrl string `json:"tarball_url"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return err
	}

	if release.TarballUrl == "" {
		return fmt.Errorf("could not find source.tar.gz asset in latest release")
	}

	resp, err = http.Get(release.TarballUrl)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create a temporary file
	tmpfile, err := os.CreateTemp("", "source-*.tar.gz")
	if err != nil {
		return err
	}
	defer os.Remove(tmpfile.Name())
	defer tmpfile.Close()

	// Copy the contents of the downloaded archive to the temporary file
	if _, err := io.Copy(tmpfile, resp.Body); err != nil {
		return err
	}

	// reset file pointer to start of file
	_, err = tmpfile.Seek(0, 0)
	if err != nil {
		return err
	}

	gzr, err := gzip.NewReader(tmpfile)
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
