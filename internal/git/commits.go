package git

import (
	"context"
	"strings"
)

// GetCurrentCommit gets a current commit from a HEAD for a local directory
func GetCurrentCommit(dir string) (string, error) {
	output, err := captureCommandBasic(
		context.TODO(),
		dir,
		"rev-parse",
		"HEAD",
	)

	return strings.TrimSpace(output), err
}

func GetMergeBase(dir string, ref1, ref2 string) (sha string, err error) {
	output, err := captureCommandBasic(
		context.TODO(),
		dir,
		"merge-base",
		ref1,
		ref2,
	)

	return strings.TrimSpace(output), err
}

func CommitPresent(rootDir, sha string) bool {
	cmd := logAndBuildCommand(context.TODO(), "cat-file", "-t", sha)
	cmd.Dir = rootDir

	output, _ := cmd.Output()
	return string(output) == "commit\n"
}
