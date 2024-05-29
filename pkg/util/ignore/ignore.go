package ignore

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/fatih/color"

	types "github.com/bearer/bearer/pkg/util/ignore/types"
	pointer "github.com/bearer/bearer/pkg/util/pointers"
)

const DefaultIgnoreFilepath = "bearer.ignore"

func GetIgnoredFingerprints(filePath string, target *string) (ignoredFingerprints map[string]types.IgnoredFingerprint, ignoreFilePath string, fileExists bool, err error) {
	if filePath == "" {
		// nothing to do here
		return map[string]types.IgnoredFingerprint{}, filePath, false, nil
	}

	ignoreFilePath, isDefaultPath, fileExists, err := GetIgnoreFilePath(filePath, target)
	if err != nil {
		if isDefaultPath && !fileExists {
			// default bearer.ignore file does not exist: expected scenario
			return map[string]types.IgnoredFingerprint{}, ignoreFilePath, false, nil
		}

		return ignoredFingerprints, ignoreFilePath, fileExists, err
	}

	// file exists
	content, err := os.ReadFile(ignoreFilePath)
	if err != nil {
		return ignoredFingerprints, ignoreFilePath, true, err
	}

	err = json.Unmarshal(content, &ignoredFingerprints)
	if err != nil {
		err = fmt.Errorf("ignore file '%s' is invalid - %s", ignoreFilePath, err)
	}
	return ignoredFingerprints, ignoreFilePath, true, err
}

func MergeIgnoredFingerprints(fingerprintsToIgnore map[string]types.IgnoredFingerprint, ignoredFingerprints map[string]types.IgnoredFingerprint, force bool) error {
	for key, value := range fingerprintsToIgnore {
		if !force {
			if _, ok := ignoredFingerprints[key]; ok {
				return fmt.Errorf(
					"fingerprint '%s' already exists in your ignore file. To view this entry run:\n\n$ bearer ignore show %s\n\nTo overwrite this entry, use --force",
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

func GetIgnoreFilePath(ignoreFilePath string, target *string) (
	path string,
	isDefaultPath bool,
	fileExists bool,
	err error,
) {
	isDefaultPath = ignoreFilePath == DefaultIgnoreFilepath

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
		return ignoreFilePath, isDefaultPath, fileExists, targetErr
	}

	ignoreFilePath = filepath.Join(targetPath, ignoreFilePath)
	info, err := os.Stat(ignoreFilePath)
	if err != nil {
		return ignoreFilePath, isDefaultPath, fileExists, err
	}

	if info.IsDir() {
		return ignoreFilePath, isDefaultPath, fileExists, fmt.Errorf("ignore file path %s is a dir not a file", ignoreFilePath)
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
