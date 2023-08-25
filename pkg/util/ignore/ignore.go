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

func GetIgnoredFingerprints(target *string) (ignoredFingerprints map[string]IgnoredFingerprint, err error) {
	fingerprints, err := readIgnoreFile(target)
	if err != nil {
		return map[string]IgnoredFingerprint{}, err
	}

	return fingerprints, nil
}

func AddToIgnoreFile(fingerprintsToIgnore map[string]IgnoredFingerprint, force bool) error {
	var existingIgnoredFingerprints map[string]IgnoredFingerprint
	if _, err := os.Stat("./bearer.ignore"); err != nil {
		existingIgnoredFingerprints = make(map[string]IgnoredFingerprint)
	} else {
		if existingIgnoredFingerprints, err = readIgnoreFile(nil); err != nil {
			return err
		}
	}

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

	return os.WriteFile("./bearer.ignore", data, 0644)
}

func readIgnoreFile(target *string) (payload map[string]IgnoredFingerprint, err error) {
	targetPath := ""
	if target != nil {
		if targetPath, err = filepath.Abs(*target); err != nil {
			return payload, err
		}
	}

	path := filepath.Join(targetPath, "bearer.ignore")

	if _, err := os.Stat(path); err != nil {
		// bearer.ignore file does not exist
		// return blank payload
		return map[string]IgnoredFingerprint{}, nil
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return payload, err
	}

	err = json.Unmarshal(content, &payload)
	return payload, err
}
