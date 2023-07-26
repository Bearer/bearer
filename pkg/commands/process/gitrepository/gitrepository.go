package gitrepository

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/rs/zerolog/log"
)

type Repository struct {
	ctx        context.Context
	repository *git.Repository
	baseBranch string
	baseLocalRefName,
	baseRemoteRefName plumbing.ReferenceName
	headRef *plumbing.Reference
}

type seenFile struct {
	count int
	sha   plumbing.Hash
}

func New(ctx context.Context, path string, baseBranch string) (*Repository, error) {
	repository, err := git.PlainOpen(path)
	if err != nil {
		if err == git.ErrRepositoryNotExists {
			log.Debug().Msg("no git repository found")

			if baseBranch != "" {
				return nil, errors.New("base branch specified but no git repository found")
			}

			return nil, nil
		}

		return nil, err
	}

	headRef, err := repository.Head()
	if err != nil {
		return nil, fmt.Errorf("error getting git head: %w", err)
	}

	baseLocalRefName := plumbing.NewBranchReferenceName(baseBranch)
	baseRemoteRefName := plumbing.NewRemoteReferenceName("origin", baseBranch)

	return &Repository{
		ctx:               ctx,
		repository:        repository,
		baseBranch:        baseBranch,
		baseLocalRefName:  baseLocalRefName,
		baseRemoteRefName: baseRemoteRefName,
		headRef:           headRef,
	}, nil
}

func (repository *Repository) ListFiles() ([]string, error) {
	if repository == nil {
		return nil, nil
	}

	seenFiles := make(map[string]seenFile)

	if err := repository.addBaseSeen(seenFiles); err != nil {
		return nil, fmt.Errorf("error with diff base: %w", err)
	}

	headTree, err := repository.treeForRef(repository.headRef)
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

func (repository *Repository) addBaseSeen(seenFiles map[string]seenFile) error {
	if repository.baseBranch == "" {
		return nil
	}

	baseRef, err := repository.repository.Reference(repository.baseLocalRefName, true)
	if err != nil {
		return err
	}

	baseTree, err := repository.treeForRef(baseRef)
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

func (repository *Repository) treeForRef(ref *plumbing.Reference) (*object.Tree, error) {
	commit, err := repository.repository.CommitObject(ref.Hash())
	if err != nil {
		return nil, err
	}

	return commit.Tree()
}

func (repository *Repository) FetchBaseIfNotPresent() error {
	ref, err := repository.repository.Reference(repository.baseLocalRefName, true)
	if err != nil && err != plumbing.ErrReferenceNotFound {
		return fmt.Errorf("invalid branch %s: %w", repository.baseBranch, err)
	}

	// Already exists
	if ref != nil {
		return nil
	}

	if err := repository.repository.FetchContext(repository.ctx, &git.FetchOptions{
		RefSpecs: []config.RefSpec{config.RefSpec(
			fmt.Sprintf("+%s:%s", repository.baseLocalRefName, repository.baseRemoteRefName),
		)},
		Depth: 1,
		Tags:  git.NoTags,
	}); err != nil {
		return fmt.Errorf("error fetching branch %s: %w", repository.baseBranch, err)
	}

	return nil
}

func (repository *Repository) WithBaseBranch(body func() error) error {
	branchName := plumbing.NewBranchReferenceName(repository.baseBranch)

	worktree, err := repository.repository.Worktree()
	if err != nil {
		return fmt.Errorf("error getting git worktree: %w", err)
	}

	if err := worktree.Checkout(&git.CheckoutOptions{
		Branch: branchName,
		Keep:   true,
	}); err != nil {
		repository.restoreHead(worktree)
		return fmt.Errorf("error checking out base branch: %w", err)
	}

	err = body()
	repository.restoreHead(worktree)
	return err
}

func (repository *Repository) restoreHead(worktree *git.Worktree) {
	checkoutOptions := &git.CheckoutOptions{
		Keep: true,
	}
	if repository.headRef.Name().IsBranch() {
		checkoutOptions.Branch = repository.headRef.Name()
	} else {
		checkoutOptions.Hash = repository.headRef.Hash()
	}

	if err := worktree.Checkout(checkoutOptions); err != nil {
		log.Error().Msgf("error restoring git worktree. your worktree may not have been restored to it's original state! %s", err)
	}
}
