package gitutil

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

type seenFile struct {
	count int
	sha   plumbing.Hash
}

func Open(path string) (*git.Repository, error) {
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

func ListFiles(repository *git.Repository, path string, diffBaseBranch string) ([]string, error) {
	if repository == nil {
		return nil, nil
	}

	seenFiles := make(map[string]seenFile)

	if err := addBaseSeen(seenFiles, repository, diffBaseBranch); err != nil {
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

func FetchIfNotPresent(ctx context.Context, repository *git.Repository, branchName string) error {
	localRefName := plumbing.NewBranchReferenceName(branchName)
	remoteRefName := plumbing.NewRemoteReferenceName("origin", branchName)

	ref, err := repository.Reference(localRefName, true)
	if err != nil && err != plumbing.ErrReferenceNotFound {
		return fmt.Errorf("invalid branch %s: %w", branchName, err)
	}

	// Already exists
	if ref != nil {
		return nil
	}

	if err := repository.FetchContext(ctx, &git.FetchOptions{
		RefSpecs: []config.RefSpec{config.RefSpec(
			fmt.Sprintf("+%s:%s", localRefName, remoteRefName),
		)},
		Depth: 1,
		Tags:  git.NoTags,
	}); err != nil {
		return fmt.Errorf("error fetching branch %s: %w", branchName, err)
	}

	return nil
}
