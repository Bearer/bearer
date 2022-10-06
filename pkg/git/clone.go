package git

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/ssoroka/slice"
)

func CloneAndGetTree(token string, url *url.URL, branchName string) (*Tree, error) {
	tempDir, err := ioutil.TempDir("", "tree")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp dir: %s", err)
	}
	defer os.RemoveAll(tempDir)

	if err := cloneTree(tempDir, urlWithCredentials(url, token), branchName); err != nil {
		return nil, fmt.Errorf("failed to clone: %s", err)
	}

	return GetTree(tempDir)
}

func CloneRangeAndCheckoutFiles(
	token string,
	url *url.URL,
	branchName string,
	previousCommit *CommitIdentifier,
	commit CommitIdentifier,
	filenames []string,
	targetDir string,
) (bool, error) {
	isRange, err := cloneTreeRange(urlWithCredentials(url, token), branchName, previousCommit, commit, targetDir)
	if err != nil {
		return false, err
	}

	firstCommitSHA := commit.SHA
	if isRange {
		firstCommitSHA = previousCommit.SHA
	}

	treeFiles, err := listTree(targetDir, commit.SHA)
	if err != nil {
		return false, err
	}

	filenames = appendMailmap(filenames, treeFiles)

	if err := fetchBlobsForRange(targetDir, firstCommitSHA, commit.SHA, filenames); err != nil {
		return false, err
	}

	// Remove remote to avoid accidentally fetching more data
	if err := removeRemote(targetDir); err != nil {
		return false, err
	}

	if err := checkout(targetDir, commit.SHA, filenames); err != nil {
		return false, err
	}

	return isRange, nil
}

func appendMailmap(filenames []string, treeFiles []TreeFile) []string {
	if slice.Contains(filenames, mailmapFilename) {
		return filenames
	}

	for _, treeFile := range treeFiles {
		if treeFile.Filename == mailmapFilename {
			return append(filenames, mailmapFilename)
		}
	}

	return filenames
}

func cloneTree(targetDir string, url *url.URL, branchName string) error {
	return basicCommand(
		"",
		"clone",
		"--depth=1",
		"--branch",
		branchName,
		"--no-checkout",
		"--no-tags",
		"--filter=blob:none",
		"--progress",
		url.String(),
		targetDir,
	)
}

func cloneTreeRange(
	url *url.URL,
	branchName string,
	previousCommit *CommitIdentifier,
	commit CommitIdentifier,
	targetDir string,
) (bool, error) {
	if previousCommit == nil {
		if err := cloneTreeSince(url, branchName, commit, targetDir); err != nil {
			return false, err
		}
	} else {
		if err := cloneTreeSince(url, branchName, *previousCommit, targetDir); err != nil {
			return false, err
		}
	}

	if !commitPresent(targetDir, commit.SHA) {
		return false, errors.New("target commit not found")
	}

	if previousCommit != nil && !commitPresent(targetDir, previousCommit.SHA) {
		log.Debug().Msg("previous commit is missing, possible re-written history")
		// Fallback to non range clone
		return false, nil
	}

	return previousCommit != nil, nil
}

func cloneTreeSince(url *url.URL, branchName string, commit CommitIdentifier, targetDir string) error {
	cmd := logAndBuildCommand(
		"clone",
		"--shallow-since="+commit.Timestamp.Format(time.RFC3339),
		"--branch",
		branchName,
		"--single-branch",
		"--no-checkout",
		"--no-tags",
		"--filter=blob:none",
		"--progress",
		url.String(),
		targetDir,
	)

	logWriter := &debugLogWriter{}
	cmd.Stdout = logWriter
	cmd.Stderr = logWriter

	if err := cmd.Run(); err != nil {
		killProcess(cmd)
		return newError(err, logWriter.AllOutput())
	}

	return nil
}

func removeRemote(rootDir string) error {
	return basicCommand(rootDir, "remote", "remove", "origin")
}
