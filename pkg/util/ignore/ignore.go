package ignore

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type IgnoredFingerprint struct {
	Author    *string `json:"author,omitempty"`
	Comment   *string `json:"comment,omitempty"`
	IgnoredAt string  `json:"ignored_at"`
}

func GetIgnoredFingerprints(bearerIgnoreFilePath string) (ignoredFingerprints map[string]IgnoredFingerprint, fileExists bool, err error) {
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
