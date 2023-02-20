package repo_info

import (
	"fmt"
	"io"

	"github.com/wlredeye/jsonlines"

	"github.com/bearer/bearer/pkg/commands/process/worker/work"
	"github.com/bearer/bearer/pkg/git"
)

type renamedFileReport struct {
	Type  string            `json:"type" yaml:"type"`
	Files []git.RenamedFile `json:"files" yaml:"files"`
}

type commitInfoReport struct {
	Type    string           `json:"type" yaml:"type"`
	Commits []git.CommitInfo `json:"commits" yaml:"commits"`
}

func ReportRepositoryInfo(reportWriter io.Writer, repository work.Repository, commitList []git.CommitInfo) error {
	if err := reportRenamedFiles(reportWriter, repository); err != nil {
		return fmt.Errorf("failed to add renamed files to report: %s", err)
	}

	if err := reportCommitInfo(reportWriter, commitList); err != nil {
		return fmt.Errorf("failed to add commit info to report: %s", err)
	}

	return nil
}

func reportRenamedFiles(reportWriter io.Writer, repository work.Repository) error {
	if repository.PreviousCommitSHA == "" {
		return nil
	}

	renamedFiles, err := git.GetRenames(repository.Dir, repository.PreviousCommitSHA, repository.CommitSHA)
	if err != nil {
		return err
	}

	if err := jsonlines.Encode(
		reportWriter,
		&[]renamedFileReport{{Type: "renamed_files", Files: renamedFiles}},
	); err != nil {
		return fmt.Errorf("encoding error: %s", err)
	}

	return nil
}

func reportCommitInfo(reportWriter io.Writer, commitList []git.CommitInfo) error {
	if len(commitList) <= 1 {
		return nil
	}

	commitListWithoutPrevCommit := commitList[:len(commitList)-1]

	if err := jsonlines.Encode(
		reportWriter,
		&[]commitInfoReport{{Type: "commit_info", Commits: commitListWithoutPrevCommit}},
	); err != nil {
		return fmt.Errorf("encoding error: %s", err)
	}

	return nil
}
