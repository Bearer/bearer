package git

import (
	"context"
	"io"
	"strings"
)

// GetCurrentCommit gets a current commit from a HEAD for a local directory
func GetCurrentCommit(dir string) (sha string, err error) {
	err = captureCommand(context.TODO(), dir, []string{"rev-parse", "HEAD"}, func(stdout io.Reader) error {
		data, err := io.ReadAll(stdout)
		if err != nil {
			return err
		}

		sha = strings.TrimSpace(string(data))
		return nil
	})

	return
}

// GetCurrentBranch gets the branch name. It is blank when detached.
func GetCurrentBranch(dir string) (name string, err error) {
	err = captureCommand(
		context.TODO(),
		dir,
		[]string{"rev-parse", "--abbrev-ref", "HEAD"},
		func(stdout io.Reader) error {
			data, err := io.ReadAll(stdout)
			if err != nil {
				return err
			}

			name = strings.TrimSpace(string(data))
			if name == "HEAD" {
				name = ""
			}

			return nil
		},
	)

	return
}

func GetMergeBase(dir string, ref1, ref2 string) (sha string, err error) {
	err = captureCommand(context.TODO(), dir, []string{"merge-base", ref1, ref2}, func(stdout io.Reader) error {
		data, err := io.ReadAll(stdout)
		if err != nil {
			return err
		}

		sha = strings.TrimSpace(string(data))
		return nil
	})

	return
}

func CommitPresent(rootDir, sha string) bool {
	cmd := logAndBuildCommand(context.TODO(), "cat-file", "-t", sha)
	cmd.Dir = rootDir

	output, _ := cmd.Output()
	return string(output) == "commit\n"
}
