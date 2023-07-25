package gitutil

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/rs/zerolog/log"
)

type seenFile struct {
	count int
	sha   plumbing.Hash
}

func DiscoverFromGit(path string, diffBaseBranch string) ([]string, error) {
	repository, err := open(path)
	if err != nil {
		return nil, err
	}
	if repository == nil {
		log.Debug().Msg("No .git directory found")
		return nil, nil
	}

	return getDiffPaths(repository, diffBaseBranch)
}

func open(path string) (*git.Repository, error) {
	gitDir := filepath.Join(path, git.GitDirName)
	if _, err := os.Stat(gitDir); os.IsNotExist(err) {
		gitDir = path
	}

	repository, err := git.PlainOpen(gitDir)
	if err != nil {
		if err == git.ErrRepositoryNotExists {
			return nil, nil
		}

		return nil, err
	}

	return repository, nil
}

func getDiffPaths(repository *git.Repository, baseBranch string) ([]string, error) {
	seenFiles := make(map[string]seenFile)

	if err := addBaseSeen(seenFiles, repository, baseBranch); err != nil {
		return nil, fmt.Errorf("error with diff base: %w", err)
	}

	headRef, err := repository.Head()
	if err != nil {
		return nil, err
	}

	headTree, err := treeForRef(repository, headRef)
	if err != nil {
		return nil, err
	}

	if err := headTree.Files().ForEach(func(f *object.File) error {
		seen := seenFiles[f.Name]

		if seen.sha == f.Hash {
			seen.count++
		}

		seenFiles[f.Name] = seen

		return nil
	}); err != nil {
		return nil, err
	}

	var result []string
	for name, seen := range seenFiles {
		if seen.count < 2 {
			result = append(result, name)
		}
	}

	return result, nil
}

func addBaseSeen(seenFiles map[string]seenFile, repository *git.Repository, baseBranch string) error {
	if baseBranch == "" {
		return nil
	}

	baseRef, err := repository.Reference(plumbing.NewBranchReferenceName(baseBranch), true)
	if err != nil {
		return err
	}

	baseTree, err := treeForRef(repository, baseRef)
	if err != nil {
		return err
	}

	return baseTree.Files().ForEach(func(f *object.File) error {
		seenFiles[f.Name] = seenFile{
			count: 1,
			sha:   f.Hash,
		}
		return nil
	})
}

func treeForRef(repository *git.Repository, ref *plumbing.Reference) (*object.Tree, error) {
	commit, err := repository.CommitObject(ref.Hash())
	if err != nil {
		return nil, err
	}

	return commit.Tree()
}
