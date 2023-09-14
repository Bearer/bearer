package ignore

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/exp/maps"

	"github.com/fatih/color"

	"github.com/bearer/bearer/api"
	types "github.com/bearer/bearer/internal/util/ignore/types"
	pointer "github.com/bearer/bearer/internal/util/pointers"
)

const DefaultIgnoreFilepath = "bearer.ignore"

func GetIgnoredFingerprints(ignoreFilePath string, target *string) (ignoredFingerprints map[string]types.IgnoredFingerprint, fileExists bool, err error) {
	ignorePath, isDefaultPath, fileExists, err := getIgnoreFilePath(ignoreFilePath, target)
	if err != nil {
		if isDefaultPath && !fileExists {
			// bearer.ignore file does not exist: expected scenario
			return map[string]types.IgnoredFingerprint{}, false, nil
		}

		return ignoredFingerprints, fileExists, err
	}

	// file exists
	content, err := os.ReadFile(ignorePath)
	if err != nil {
		return ignoredFingerprints, true, err
	}

	err = json.Unmarshal(content, &ignoredFingerprints)
	return ignoredFingerprints, true, err
}

func GetIgnoredFingerprintsFromCloud(
	client *api.API,
	fullname string,
	localIgnores map[string]types.IgnoredFingerprint,
) (
	useCloudIgnores bool,
	ignoredFingerprints map[string]types.IgnoredFingerprint,
	staleIgnoredFingerprintIds []string,
	err error,
) {
	data, err := client.FetchIgnores(fullname, maps.Keys(localIgnores))
	if err != nil {
		return useCloudIgnores, ignoredFingerprints, staleIgnoredFingerprintIds, err
	}

	ignoredFingerprints = make(map[string]types.IgnoredFingerprint)
	for _, fingerprint := range data.Ignores {
		item := types.IgnoredFingerprint{}

		_, persistedInCloud := data.CloudIgnoredFingerprints[fingerprint]
		if !persistedInCloud {
			// it is a new addition; use information from ignore file
			item = localIgnores[fingerprint]
		}

		ignoredFingerprints[fingerprint] = item
	}

	return data.ProjectFound, ignoredFingerprints, data.StaleIgnores, nil
}

func MergeIgnoredFingerprints(fingerprintsToIgnore map[string]types.IgnoredFingerprint, ignoredFingerprints map[string]types.IgnoredFingerprint, force bool) error {
	for key, value := range fingerprintsToIgnore {
		if !force {
			if _, ok := ignoredFingerprints[key]; ok {
				return fmt.Errorf(
					"fingerprint '%s' already exists in the bearer.ignore file. To view this entry run:\n\n$ bearer ignore show %s\n\nTo overwrite this entry, use --force",
					key,
					key,
				)
			}
		}
		ignoredAt := time.Now().UTC()
		value.IgnoredAt = ignoredAt.Format(time.RFC3339)
		ignoredFingerprints[key] = value
	}
	return nil
}

var bold = color.New(color.Bold).SprintFunc()
var morePrefix = color.HiBlackString("├─ ")
var lastPrefix = color.HiBlackString("└─ ")

func DisplayIgnoredEntryTextString(fingerprintId string, entry types.IgnoredFingerprint, noColor bool) string {
	initialColorSetting := color.NoColor
	if noColor && !initialColorSetting {
		color.NoColor = true
	}
	prefix := morePrefix
	result := fmt.Sprintf(bold(color.HiBlueString("%s \n")), fingerprintId)

	if entry.Author == nil && entry.Comment == nil {
		prefix = lastPrefix
	}
	result += fmt.Sprintf("%sIgnored At: %s", prefix, bold(entry.IgnoredAt))

	if entry.Author != nil {
		result += fmt.Sprintf("\n%sAuthor: %s", prefix, bold(*entry.Author))
	}

	if entry.Comment == nil {
		prefix = lastPrefix
	}
	var falsePositiveStr string
	if entry.FalsePositive {
		falsePositiveStr = "Yes"
	} else {
		falsePositiveStr = "No"
	}
	result += fmt.Sprintf("\n%sFalse positive? %s", prefix, bold(falsePositiveStr))

	if entry.Comment != nil {
		result += fmt.Sprintf("\n%sComment: %s", lastPrefix, bold(*entry.Comment))
	}

	color.NoColor = initialColorSetting

	return result
}

func GetAuthor() (*string, error) {
	nameBytes, err := exec.Command("git", "config", "user.name").Output()
	if err != nil {
		return nil, err
	}

	return pointer.String(strings.TrimSuffix(string(nameBytes), "\n")), nil
}

func getIgnoreFilePath(ignoreFilePath string, target *string) (
	ignorePath string,
	isDefaultPath bool,
	fileExists bool,
	err error,
) {
	if ignoreFilePath == "" {
		// use default ignore file path
		isDefaultPath = true
		ignoreFilePath = DefaultIgnoreFilepath
	}

	_, err = os.Stat(ignoreFilePath)
	if err == nil {
		// file is found (all good)
		return ignoreFilePath, isDefaultPath, true, err
	}
	fileNotFoundErr := os.IsNotExist(err)
	if !isDefaultPath || !fileNotFoundErr {
		// custom ignore file is not found (fail early)
		// or unexpected error has occurred
		return ignoreFilePath, isDefaultPath, fileExists, err
	}

	// file not found
	fileExists = false

	// append default path to target path and try again
	targetPath, targetErr := targetPath(target)
	if targetErr != nil {
		return "", isDefaultPath, fileExists, targetErr
	}

	ignoreFilePath = filepath.Join(targetPath, ignoreFilePath)
	info, err := os.Stat(ignoreFilePath)
	if err != nil {
		return "", isDefaultPath, fileExists, err
	}

	if info.IsDir() {
		return "", isDefaultPath, fileExists, fmt.Errorf("ignore file path %s is a dir not a file", ignoreFilePath)
	}

	return ignoreFilePath, isDefaultPath, fileExists, nil
}

// returns target directory from target
func targetPath(target *string) (string, error) {
	if target == nil {
		return "", nil
	}

	targetPath, err := filepath.Abs(*target)
	if err != nil {
		return "", err
	}

	info, err := os.Stat(targetPath)
	if err != nil {
		return "", err
	}

	if info.IsDir() {
		return targetPath, nil
	}

	// not a directory
	return filepath.Dir(targetPath), nil
}
