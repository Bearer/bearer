package ignore

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type IgnoredFingerprint struct {
	Author    string
	Comment   string
	IgnoredAt string
}

func GetIgnoredFingerprints() (ignoredFingerprints map[string]IgnoredFingerprint, err error) {
	fingerprints, err := readIgnoreFile()
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
		if existingIgnoredFingerprints, err = readIgnoreFile(); err != nil {
			return err
		} else {
			for key, value := range fingerprintsToIgnore {
				if !force {
					if existingIgnoredFingerprint, ok := existingIgnoredFingerprints[key]; ok {
						return fmt.Errorf(
							"fingerprint %s already exists in bearer.ignore file, with author %s and comment %s. To overwrite this entry, use --force",
							key,
							existingIgnoredFingerprint.Author,
							existingIgnoredFingerprint.Comment,
						)
					}
				}

				ignoredAt := time.Now().UTC()
				value.IgnoredAt = ignoredAt.Format(time.RFC3339)
				existingIgnoredFingerprints[key] = value
			}
		}
	}
	data, err := json.Marshal(existingIgnoredFingerprints)
	if err != nil {
		// failed to marshall data
		return err
	}

	return os.WriteFile("./bearer.ignore", data, 0644)
}

func readIgnoreFile() (payload map[string]IgnoredFingerprint, err error) {
	content, err := os.ReadFile("./bearer.ignore")
	if err != nil {
		return payload, err
	}

	err = json.Unmarshal(content, &payload)
	return payload, err
}
