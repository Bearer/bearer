package blamer

import (
	"github.com/bearer/curio/pkg/git"
	"github.com/rs/zerolog/log"
)

type Blamer struct {
	repositoryDir      string
	commitSHAsFilename string
	previousCommitSHA  string
	fileBlames         map[string]git.BlameResult
}

func New(repositoryDir, commitSHAsFilename string, previousCommitSHA string) *Blamer {
	return &Blamer{
		repositoryDir:      repositoryDir,
		commitSHAsFilename: commitSHAsFilename,
		previousCommitSHA:  previousCommitSHA,
		fileBlames:         make(map[string]git.BlameResult),
	}
}

func (blamer *Blamer) SHAForLine(filename string, lineNumber int) string {
	if blamer.previousCommitSHA == "" {
		return ""
	}

	blameResult, ok := blamer.fileBlames[filename]
	if !ok {
		var err error
		blameResult, err = git.Blame(blamer.repositoryDir, blamer.commitSHAsFilename, filename)
		if err != nil {
			log.Err(err).Msg("failed to get blame info for %s")
		}

		blamer.fileBlames[filename] = blameResult
	}

	sha := blameResult.SHAForLine(lineNumber)

	// If the commit is at the boundary of what we've checked out then it
	// could have actually been a prior commit, so don't report it
	if sha == blamer.previousCommitSHA {
		return ""
	}

	return sha
}
