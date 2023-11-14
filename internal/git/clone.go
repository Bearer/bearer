package git

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"
	"slices"
	"time"

	"github.com/rs/zerolog/log"
)

func CloneAndGetTree(token string, url *url.URL, branchName string) (*Tree, error) {
	tempDir, err := os.MkdirTemp("", "tree")
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

	treeFiles, err := ListTree(targetDir, commit.SHA)
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

	if err := checkoutFiles(targetDir, commit.SHA, filenames); err != nil {
		return false, err
	}

	return isRange, nil
}

func appendMailmap(filenames []string, treeFiles []TreeFile) []string {
	if slices.Contains(filenames, mailmapFilename) {
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
		context.TODO(),
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

	if !CommitPresent(targetDir, commit.SHA) {
		return false, errors.New("target commit not found")
	}

	if previousCommit != nil && !CommitPresent(targetDir, previousCommit.SHA) {
		log.Debug().Msg("previous commit is missing, possible re-written history")
		// Fallback to non range clone
		return false, nil
	}

	return previousCommit != nil, nil
}

func cloneTreeSince(url *url.URL, branchName string, commit CommitIdentifier, targetDir string) error {
	cmd := logAndBuildCommand(
		context.TODO(),
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
		cmd.Cancel() //nolint:errcheck
		return newError(err, logWriter.AllOutput())
	}

	return nil
}

func removeRemote(rootDir string) error {
	return basicCommand(context.TODO(), rootDir, "remote", "remove", "origin")
}
