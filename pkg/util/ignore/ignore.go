package ignore

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/fatih/color"
)

type IgnoredFingerprint struct {
	Author    string
	Comment   string
	IgnoredAt string
}

type DuplicateIgnoredFingerprintError struct {
	Err error
}

func (f *DuplicateIgnoredFingerprintError) Error() string {
	return f.Err.Error()
}

var bold = color.New(color.Bold).SprintFunc()
var morePrefix = color.HiBlackString("├─ ")
var lastPrefix = color.HiBlackString("└─ ")

func DisplayIgnoredEntryTextString(fingerprintId string, entry IgnoredFingerprint) string {
	prefix := morePrefix
	result := fmt.Sprintf(bold(color.HiBlueString("%s \n")), fingerprintId)

	if entry.Author == "" && entry.Comment == "" {
		prefix = lastPrefix
	}
	result += fmt.Sprintf("%sIgnored At: %s\n", prefix, bold(entry.IgnoredAt))

	if entry.Author != "" {
		if entry.Comment == "" {
			prefix = lastPrefix
		}

		result += fmt.Sprintf("%sAuthor: %s\n", prefix, bold(entry.Author))
	}

	if entry.Comment != "" {
		result += fmt.Sprintf("%sComment: %s\n", lastPrefix, bold(entry.Comment))
	}

	return result
}

func GetIgnoredFingerprints(bearerIgnoreFilePath string) (ignoredFingerprints map[string]IgnoredFingerprint, err error) {
	fingerprints, err := readIgnoreFile(bearerIgnoreFilePath)
	if err != nil {
		return map[string]IgnoredFingerprint{}, err
	}

	return fingerprints, nil
}

func GetExistingIgnoredFingerprints(bearerIgnoreFilePath string) (existingIgnoredFingerprints map[string]IgnoredFingerprint, fileExists bool, err error) {
	if _, noFileErr := os.Stat(bearerIgnoreFilePath); noFileErr != nil {
		existingIgnoredFingerprints = make(map[string]IgnoredFingerprint)
	} else {
		fileExists = true
		existingIgnoredFingerprints, err = readIgnoreFile(bearerIgnoreFilePath)
	}

	return existingIgnoredFingerprints, fileExists, err
}

func AddToIgnoreFile(fingerprintsToIgnore map[string]IgnoredFingerprint, existingIgnoredFingerprints map[string]IgnoredFingerprint, bearerIgnoreFilePath string, force bool) error {
	for key, value := range fingerprintsToIgnore {
		if !force {
			if _, ok := existingIgnoredFingerprints[key]; ok {
				error := fmt.Errorf(
					"fingerprint '%s' already exists in bearer.ignore file. To view this entry run:\n\n$ bearer ignore show %s\n\nTo overwrite this entry, use --force",
					key,
					key,
				)
				return &DuplicateIgnoredFingerprintError{
					Err: error,
				}
			}
		}
		ignoredAt := time.Now().UTC()
		value.IgnoredAt = ignoredAt.Format(time.RFC3339)
		existingIgnoredFingerprints[key] = value
	}

	data, err := json.MarshalIndent(existingIgnoredFingerprints, "", "  ")
	if err != nil {
		// failed to marshall data
		return err
	}

	return os.WriteFile(bearerIgnoreFilePath, data, 0644)
}

func readIgnoreFile(bearerIgnoreFilePath string) (payload map[string]IgnoredFingerprint, err error) {
	if _, err := os.Stat(bearerIgnoreFilePath); err != nil {
		// bearer.ignore file does not exist
		// return blank payload
		return map[string]IgnoredFingerprint{}, nil
	}

	content, err := os.ReadFile(bearerIgnoreFilePath)
	if err != nil {
		return payload, err
	}

	err = json.Unmarshal(content, &payload)
	return payload, err
}
