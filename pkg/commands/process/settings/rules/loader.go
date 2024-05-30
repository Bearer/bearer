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
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"

	"github.com/bearer/bearer/pkg/commands/process/settings"
	"github.com/bearer/bearer/pkg/engine"
	flagtypes "github.com/bearer/bearer/pkg/flag/types"
	"github.com/bearer/bearer/pkg/util/output"
	"github.com/bearer/bearer/pkg/version_check"
)

func loadDefinitionsFromRemote(
	definitions map[string]settings.RuleDefinition,
	options flagtypes.RuleOptions,
	versionMeta *version_check.VersionMeta,
) error {
	if options.DisableDefaultRules {
		return nil
	}

	if versionMeta.Rules.Version == nil {
		log.Debug().Msg("No rule packages found")
		return nil
	}

	urls := make([]string, 0, len(versionMeta.Rules.Packages))
	for _, value := range versionMeta.Rules.Packages {
		log.Debug().Msgf("Added rule package URL %s", value)
		urls = append(urls, value)
	}

	if err := readDefinitionsFromUrls(definitions, urls); err != nil {
		return fmt.Errorf("loading rules failed: %s", err)
	}

	return nil
}

func readDefinitionsFromUrls(ruleDefinitions map[string]settings.RuleDefinition, languageDownloads []string) (err error) {
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

			if err = readRuleDefinitionZip(ruleDefinitions, file); err != nil {
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

			if err = readRuleDefinitionZip(ruleDefinitions, file); err != nil {
				return err
			}
		}
	}

	return nil
}

func readRuleDefinitionZip(ruleDefinitions map[string]settings.RuleDefinition, file *os.File) error {
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

		var ruleDefinition settings.RuleDefinition
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

func loadCustomDefinitions(engine engine.Engine, definitions map[string]settings.RuleDefinition, isBuiltIn bool, dir fs.FS) error {
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
			return fmt.Errorf("rule file was invalid")
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
			language := engine.GetLanguageById(languageID)
			if language != nil {
				supported = true
			}
		}

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
		return err
	}

	for id, definition := range loadedDefinitions {
		if validateCustomRuleDefinition(loadedDefinitions, &definition) {
			definitions[id] = definition
		}
	}

	return nil
}

func bearerRulesDir() string {
	return filepath.Join(os.TempDir(), "bearer-rules")
}
