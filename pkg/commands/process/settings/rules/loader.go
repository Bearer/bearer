package rules

import (
	"archive/tar"
	"compress/gzip"
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"

	"github.com/bearer/bearer/pkg/commands/process/settings"
	flagtypes "github.com/bearer/bearer/pkg/flag/types"
	"github.com/bearer/bearer/pkg/util/output"
	"github.com/bearer/bearer/pkg/version_check"
)

func loadDefinitionsFromRemote(
	definitions map[string]settings.RuleDefinition,
	options flagtypes.RuleOptions,
	versionMeta *version_check.VersionMeta,
) (int, error) {
	if options.DisableDefaultRules {
		return 0, nil
	}

	if versionMeta.Rules.Version == nil {
		log.Debug().Msg("No rule packages found")
		return 0, nil
	}

	urls := make([]string, 0, len(versionMeta.Rules.Packages))
	for _, value := range versionMeta.Rules.Packages {
		log.Debug().Msgf("Added rule package URL %s", value)
		urls = append(urls, value)
	}

	count, err := readDefinitionsFromUrls(definitions, urls)
	if err != nil {
		return 0, fmt.Errorf("loading rules failed: %s", err)
	}

	return count, nil
}

func readDefinitionsFromUrls(ruleDefinitions map[string]settings.RuleDefinition, languageDownloads []string) (int, error) {
	bearerRulesDir := bearerRulesDir()
	if _, err := os.Stat(bearerRulesDir); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(bearerRulesDir, os.ModePerm)
		if err != nil {
			return 0, fmt.Errorf("could not create bearer-rules directory: %s", err)
		}
	}

	count := 0
	for _, languagePackageUrl := range languageDownloads {
		// Prepare filepath
		urlHash := md5.Sum([]byte(languagePackageUrl))
		filepath, err := filepath.Abs(filepath.Join(bearerRulesDir, fmt.Sprintf("%x.tar.gz", urlHash)))
		if err != nil {
			return 0, err
		}

		var languageCount int
		if _, err := os.Stat(filepath); err == nil {
			log.Trace().Msgf("Using local cache for rule package: %s", languagePackageUrl)
			file, err := os.Open(filepath)
			if err != nil {
				return 0, err
			}
			defer file.Close()

			languageCount, err = readRuleDefinitionZip(ruleDefinitions, file)
			if err != nil {
				return 0, err
			}
		} else {
			log.Trace().Msgf("Downloading rule package: %s", languagePackageUrl)
			httpClient := &http.Client{Timeout: 60 * time.Second}
			resp, err := httpClient.Get(languagePackageUrl)
			if err != nil {
				return 0, err
			}
			defer resp.Body.Close()

			// Create file in rules dir
			file, err := os.Create(filepath)
			if err != nil {
				return 0, err
			}
			defer file.Close()

			// Copy the contents of the downloaded archive to the file
			if _, err := io.Copy(file, resp.Body); err != nil {
				return 0, err
			}
			// reset file pointer to start of file
			_, err = file.Seek(0, 0)
			if err != nil {
				return 0, err
			}

			languageCount, err = readRuleDefinitionZip(ruleDefinitions, file)
			if err != nil {
				return 0, err
			}
		}

		count += languageCount
	}

	return count, nil
}

func readRuleDefinitionZip(ruleDefinitions map[string]settings.RuleDefinition, file *os.File) (int, error) {
	gzr, err := gzip.NewReader(file)
	if err != nil {
		return 0, err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)
	count := 0
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return 0, err
		}

		if !isRuleFile(header.Name) {
			continue
		}

		data := make([]byte, header.Size)
		_, err = io.ReadFull(tr, data)
		if err != nil {
			return 0, fmt.Errorf("failed to read file %s: %w", header.Name, err)
		}

		var ruleDefinition settings.RuleDefinition
		err = yaml.Unmarshal(data, &ruleDefinition)
		if err != nil {
			return 0, fmt.Errorf("failed to unmarshal rule %s: %w", header.Name, err)
		}

		id := ruleDefinition.Metadata.ID
		_, ruleExists := ruleDefinitions[id]
		if ruleExists {
			return 0, fmt.Errorf("duplicate built-in rule ID %s", id)
		}

		ruleDefinitions[id] = ruleDefinition
		count += 1
	}

	return count, nil
}

func loadCustomDefinitions(
	definitions map[string]settings.RuleDefinition,
	isBuiltIn bool,
	dir fs.FS,
	languageIDs []string,
) (int, error) {
	count := 0
	loadedDefinitions := make(map[string]settings.RuleDefinition)
	if err := fs.WalkDir(dir, ".", func(path string, dirEntry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if dirEntry.IsDir() {
			return nil
		}

		filename := dirEntry.Name()
		ext := filepath.Ext(filename)

		if ext != ".yaml" && ext != ".yml" {
			return nil
		}

		entry, err := fs.ReadFile(dir, path)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %w", path, err)
		}

		var ruleDefinition settings.RuleDefinition
		err = yaml.Unmarshal(entry, &ruleDefinition)
		if err != nil {
			output.StdErrLog(validateCustomRuleSchema(entry, filename))
			return fmt.Errorf("rule file was invalid: %w", err)
		}

		if ruleDefinition.Metadata == nil {
			log.Debug().Msgf("rule file has invalid metadata %s", path)
			return nil
		}

		id := ruleDefinition.Metadata.ID
		if id == "" {
			log.Debug().Msgf("rule file missing metadata.id %s", path)
			return nil
		}

		supported := isBuiltIn
		for _, languageID := range ruleDefinition.Languages {
			if slices.Contains(languageIDs, languageID) {
				supported = true
			}
		}

		count += 1

		if !supported {
			log.Debug().Msgf(
				"rule file has no supported languages[%s] %s",
				strings.Join(ruleDefinition.Languages, ", "),
				path,
			)
			return nil
		}

		if _, exists := loadedDefinitions[id]; exists {
			return fmt.Errorf("duplicate rule ID %s", id)
		}

		loadedDefinitions[id] = ruleDefinition

		return nil
	}); err != nil {
		return 0, err
	}

	for id, definition := range loadedDefinitions {
		if validateCustomRuleDefinition(loadedDefinitions, &definition) {
			definitions[id] = definition
		}
	}

	return count, nil
}

func bearerRulesDir() string {
	return filepath.Join(os.TempDir(), "bearer-rules")
}
