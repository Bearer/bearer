package ignore

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/fatih/color"
)

type IgnoredFingerprint struct {
	Author    *string `json:"author,omitempty"`
	Comment   *string `json:"comment,omitempty"`
	IgnoredAt string  `json:"ignored_at"`
}

func GetIgnoredFingerprints(bearerIgnoreFilePath string, target *string) (ignoredFingerprints map[string]IgnoredFingerprint, fileExists bool, err error) {
	if target != nil {
		targetPath := ""
		if targetPath, err = filepath.Abs(*target); err != nil {
			return ignoredFingerprints, fileExists, err
		}
		bearerIgnoreFilePath = filepath.Join(targetPath, bearerIgnoreFilePath)
	}

	if _, noFileErr := os.Stat(bearerIgnoreFilePath); noFileErr != nil {
		ignoredFingerprints = make(map[string]IgnoredFingerprint)
		return ignoredFingerprints, false, err
	}

	// file exists
	content, err := os.ReadFile(bearerIgnoreFilePath)
	if err != nil {
		return ignoredFingerprints, true, err
	}

	err = json.Unmarshal(content, &ignoredFingerprints)
	return ignoredFingerprints, true, err
}

func MergeIgnoredFingerprints(fingerprintsToIgnore map[string]IgnoredFingerprint, ignoredFingerprints map[string]IgnoredFingerprint, force bool) error {
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

func DisplayIgnoredEntryTextString(fingerprintId string, entry IgnoredFingerprint, noColor bool) string {
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
		if entry.Comment == nil {
			prefix = lastPrefix
		}

		result += fmt.Sprintf("\n%sAuthor: %s", prefix, bold(*entry.Author))
	}

	if entry.Comment != nil {
		result += fmt.Sprintf("\n%sComment: %s", lastPrefix, bold(*entry.Comment))
	}

	color.NoColor = initialColorSetting

	return result
}
