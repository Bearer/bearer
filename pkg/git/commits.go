package git

import (
	"context"
	"strings"
)

// GetCurrentCommit gets a current commit from a HEAD for a local directory
func GetCurrentCommit(rootDir string) (string, error) {
	output, err := captureCommandBasic(
		context.TODO(),
		rootDir,
		"rev-parse",
		"HEAD",
	)

	if err != nil && strings.Contains(err.Error(), "unknown revision") {
		return "", nil
	}

	return strings.TrimSpace(output), err
}

func GetMergeBase(rootDir string, ref1, ref2 string) (string, error) {
	output, err := captureCommandBasic(
		context.TODO(),
		rootDir,
		"merge-base",
		ref1,
		ref2,
	)

	return strings.TrimSpace(output), err
}

func CommitPresent(rootDir, hash string) (bool, error) {
	output, err := captureCommandBasic(
		context.TODO(),
		rootDir,
		"cat-file",
		"-t",
		hash,
	)
	if err != nil {
		if strings.Contains(err.Error(), "could not get object info") {
			return false, nil
		}

		return false, err
	}

	return output == "commit\n", nil
}
