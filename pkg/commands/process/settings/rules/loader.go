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
	"github.com/bearer/bearer/pkg/report/customdetectors"
	"github.com/bearer/bearer/pkg/structurederrors"
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

// ruleDir is one rule directory: its path, for messages, and its files.
type ruleDir struct {
	path string
	fs   fs.FS
}

// loadDefinitions loads rule directories in order: a rule is validated against the rules of its own
// directory, and a rule ID in a later directory replaces the same ID from an earlier one. When there
// is more than one directory, a rule may also import a `type: shared` rule from any of the other
// directories, if that shared rule is valid in its own directory, so a custom rule can reuse a shared
// rule kept in another directory.
func loadDefinitions(
	definitions map[string]settings.RuleDefinition,
	dirs []ruleDir,
	isBuiltIn bool,
	languageIDs []string,
	onlyRule map[string]bool,
	errorCollector *structurederrors.Collector,
) (int, error) {
	count := 0
	dirDefinitions := make([]map[string]settings.RuleDefinition, 0, len(dirs))
	sharedDefinitions := make(map[string]settings.RuleDefinition)

	for _, dir := range dirs {
		loadedDefinitions, dirCount, err := readCustomDefinitions(isBuiltIn, dir, languageIDs, onlyRule, errorCollector)
		if err != nil {
			return 0, err
		}

		// With a single directory every shared rule is already in that directory, so there is nothing
		// to share and loading stays exactly as it was. A shared rule is only offered to other
		// directories if it is valid in its own: importing an invalid one would otherwise make the
		// engine fail to start, where before the importing rule was just dropped.
		if len(dirs) > 1 {
			for id, definition := range loadedDefinitions {
				if definition.Type == customdetectors.TypeShared &&
					len(validateCustomRuleDefinition(loadedDefinitions, &definition)) == 0 {
					sharedDefinitions[id] = definition
				}
			}
		}

		dirDefinitions = append(dirDefinitions, loadedDefinitions)
		count += dirCount
	}

	for _, loadedDefinitions := range dirDefinitions {
		if len(sharedDefinitions) == 0 {
			addValidDefinitions(definitions, loadedDefinitions, loadedDefinitions, errorCollector)
			continue
		}

		visibleDefinitions := make(map[string]settings.RuleDefinition, len(sharedDefinitions)+len(loadedDefinitions))
		for id, definition := range sharedDefinitions {
			visibleDefinitions[id] = definition
		}
		for id, definition := range loadedDefinitions {
			visibleDefinitions[id] = definition
		}

		addValidDefinitions(definitions, loadedDefinitions, visibleDefinitions, errorCollector)
	}

	return count, nil
}

// addValidDefinitions copies into definitions each loaded definition that is valid, checking its
// imports and filter references against visibleDefinitions. Invalid definitions are logged and left
// out.
func addValidDefinitions(
	definitions map[string]settings.RuleDefinition,
	loadedDefinitions map[string]settings.RuleDefinition,
	visibleDefinitions map[string]settings.RuleDefinition,
	errorCollector *structurederrors.Collector,
) {
	for id, definition := range loadedDefinitions {
		issues := validateCustomRuleDefinition(visibleDefinitions, &definition)
		if len(issues) == 0 {
			definitions[id] = definition
			continue
		}

		for _, issue := range issues {
			log.Debug().Msgf("%s: %s", id, issue)
			errorCollector.Add(structurederrors.Item{
				Severity:      structurederrors.SeverityError,
				Code:          issue.code,
				Message:       issue.message,
				RuleID:        id,
				RelatedRuleID: issue.relatedRuleID,
			})
		}

		log.Debug().Msgf("%s ignored due to validation errors", id)
	}
}

func readCustomDefinitions(
	isBuiltIn bool,
	dir ruleDir,
	languageIDs []string,
	onlyRule map[string]bool,
	errorCollector *structurederrors.Collector,
) (map[string]settings.RuleDefinition, int, error) {
	count := 0
	loadedDefinitions := make(map[string]settings.RuleDefinition)
	if err := fs.WalkDir(dir.fs, ".", func(path string, dirEntry fs.DirEntry, err error) error {
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

		rulePath := filepath.Join(dir.path, path)
		ruleError := func(code structurederrors.Code, ruleID string, err error) {
			errorCollector.Add(structurederrors.Item{
				Severity: structurederrors.SeverityError,
				Code:     code,
				Message:  err.Error(),
				RuleID:   ruleID,
				RuleFile: rulePath,
			})
		}

		entry, err := fs.ReadFile(dir.fs, path)
		if err != nil {
			ruleError(structurederrors.CodeRuleReadFailed, "", err)
			return fmt.Errorf("failed to read file %s: %w", path, err)
		}

		var ruleDefinition settings.RuleDefinition
		err = yaml.Unmarshal(entry, &ruleDefinition)
		if err != nil {
			ruleError(structurederrors.CodeRuleYAMLInvalid, "", err)
			output.StdErrLog(validateCustomRuleSchema(entry, filename))
			return fmt.Errorf("rule file was invalid: %w", err)
		}

		if ruleDefinition.Metadata == nil {
			log.Debug().Msgf("rule file has invalid metadata %s", path)
			ruleError(structurederrors.CodeRuleMissingMetadata, "", errors.New("rule file has no metadata"))
			return nil
		}

		id := ruleDefinition.Metadata.ID
		if id == "" {
			log.Debug().Msgf("rule file missing metadata.id %s", path)
			ruleError(structurederrors.CodeRuleMissingID, "", errors.New("metadata.id must be specified"))
			return nil
		}

		// The only-rule option already names the rules to load, so the language
		// filter must not remove them (or the shared rules they import) first.
		supported := isBuiltIn || len(onlyRule) > 0
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
			errorCollector.Add(structurederrors.Item{
				Severity: structurederrors.SeverityWarning,
				Code:     structurederrors.CodeRuleLanguageNotFound,
				Message: fmt.Sprintf(
					"rule languages [%s] were not found in the scan target (found: [%s]); the rule was not run",
					strings.Join(ruleDefinition.Languages, ", "),
					strings.Join(languageIDs, ", "),
				),
				RuleID:   id,
				RuleFile: rulePath,
			})
			return nil
		}

		if _, exists := loadedDefinitions[id]; exists {
			err := fmt.Errorf("duplicate rule ID %s", id)
			ruleError(structurederrors.CodeRuleDuplicateID, id, err)
			return err
		}

		loadedDefinitions[id] = ruleDefinition

		errorCollector.RegisterRule(id, rulePath, "")
		for _, auxiliaryDefinition := range ruleDefinition.Auxiliary {
			errorCollector.RegisterRule(auxiliaryDefinition.Id, rulePath, id)
		}

		return nil
	}); err != nil {
		return nil, 0, err
	}

	return loadedDefinitions, count, nil
}

func bearerRulesDir() string {
	return filepath.Join(os.TempDir(), "bearer-rules")
}
