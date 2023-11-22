package git

import (
	"context"
	"strings"
)

func GetDefaultBranch(dir string) (string, error) {
	name, err := getRevParseAbbrevRef(dir, "origin/HEAD")
	if err != nil {
		if strings.Contains(err.Error(), "unknown revision") {
			return "", nil
		}

		return "", err
	}

	return strings.TrimPrefix(name, "origin/"), nil
}

// GetCurrentBranch gets the branch name. It is blank when detached.
func GetCurrentBranch(dir string) (string, error) {
	name, err := getRevParseAbbrevRef(dir, "HEAD")
	if err != nil {
		if strings.Contains(err.Error(), "unknown revision") {
			return "", nil
		}

		return "", err
	}

	if name == "HEAD" {
		return "", nil
	}

	return name, nil
}

func getRevParseAbbrevRef(dir string, ref string) (name string, err error) {
	output, err := captureCommandBasic(
		context.TODO(),
		dir,
		"rev-parse",
		"--abbrev-ref",
		ref,
	)

	return strings.TrimSpace(output), err
}
