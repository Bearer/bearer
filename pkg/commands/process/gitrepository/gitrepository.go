package gitrepository

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/format/diff"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/hhatto/gocloc"
	"github.com/rs/zerolog/log"

	"github.com/bearer/bearer/pkg/commands/process/filelist/files"
	"github.com/bearer/bearer/pkg/commands/process/filelist/ignore"
	"github.com/bearer/bearer/pkg/commands/process/filelist/timeout"
	"github.com/bearer/bearer/pkg/commands/process/settings"
	"github.com/bearer/bearer/pkg/report/basebranchfindings"
	bbftypes "github.com/bearer/bearer/pkg/report/basebranchfindings/types"
)

type Repository struct {
	ctx        context.Context
	config     settings.Config
	repository *git.Repository
	baseBranch string
	baseLocalRefName,
	baseRemoteRefName plumbing.ReferenceName
	headRef *plumbing.Reference
}

func New(
	ctx context.Context,
	config settings.Config,
	path string,
	baseBranch string,
) (*Repository, error) {
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
		config:            config,
		repository:        repository,
		baseBranch:        baseBranch,
		baseLocalRefName:  baseLocalRefName,
		baseRemoteRefName: baseRemoteRefName,
		headRef:           headRef,
	}, nil
}

func (repository *Repository) ListFiles(
	ignore *ignore.FileIgnore,
	goclocResult *gocloc.Result,
	projectPath string,
) (*files.List, error) {
	if repository == nil {
		return nil, nil
	}

	headTree, err := repository.treeForRef(repository.headRef)
	if err != nil {
		return nil, err
	}

	if repository.baseBranch == "" {
		return repository.getTreeFiles(ignore, goclocResult, projectPath, headTree)
	}

	filePatches, err := repository.getDiffPatch(headTree)
	if err != nil {
		return nil, err
	}

	return repository.getDiffFiles(
		ignore,
		goclocResult,
		projectPath,
		filePatches,
	), nil
}

func (repository *Repository) treeForRef(ref *plumbing.Reference) (*object.Tree, error) {
	commit, err := repository.repository.CommitObject(ref.Hash())
	if err != nil {
		return nil, err
	}

	return commit.Tree()
}

func (repository *Repository) getTreeFiles(
	ignore *ignore.FileIgnore,
	goclocResult *gocloc.Result,
	projectPath string,
	tree *object.Tree,
) (*files.List, error) {
	var headFiles []files.File

	if err := tree.Files().ForEach(func(f *object.File) error {
		file := repository.fileFor(
			ignore,
			goclocResult,
			projectPath,
			f.Name,
		)

		if file != nil {
			headFiles = append(headFiles, *file)
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return &files.List{Files: headFiles}, nil
}

func (repository *Repository) getDiffPatch(headTree *object.Tree) ([]diff.FilePatch, error) {
	baseRef, err := repository.repository.Reference(repository.baseLocalRefName, true)
	if err != nil {
		return nil, err
	}

	baseTree, err := repository.treeForRef(baseRef)
	if err != nil {
		return nil, err
	}

	changes, err := object.DiffTreeWithOptions(
		repository.ctx,
		baseTree,
		headTree,
		object.DefaultDiffTreeOptions,
	)
	if err != nil {
		return nil, fmt.Errorf("error diffing tree: %w", err)
	}

	patch, err := changes.PatchContext(repository.ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting diff patch: %w", err)
	}

	return patch.FilePatches(), nil
}

func (repository *Repository) translateDiffChunks(gitChunks []diff.Chunk) bbftypes.Chunks {
	chunks := basebranchfindings.NewChunks()
	for _, chunk := range gitChunks {
		var changeType bbftypes.ChangeType
		switch chunk.Type() {
		case diff.Delete:
			changeType = bbftypes.ChunkRemove
		case diff.Add:
			changeType = bbftypes.ChunkAdd
		case diff.Equal:
			changeType = bbftypes.ChunkEqual
		default:
			panic(fmt.Sprintf("unexpected git chunk type %d", chunk.Type()))
		}

		chunks.Add(changeType, strings.Count(chunk.Content(), "\n"))
	}

	return chunks
}

func (repository *Repository) getDiffFiles(
	ignore *ignore.FileIgnore,
	goclocResult *gocloc.Result,
	projectPath string,
	filePatches []diff.FilePatch,
) *files.List {
	var baseFiles, headFiles []files.File
	renames := make(map[string]string)
	chunks := make(map[string]bbftypes.Chunks)

	for _, filePatch := range filePatches {
		fromFile, toFile := filePatch.Files()

		// we're not interested in removals
		if toFile == nil {
			continue
		}

		toPath := toFile.Path()
		headFile := repository.fileFor(ignore, goclocResult, projectPath, toPath)
		if headFile == nil {
			continue
		}

		headFiles = append(headFiles, *headFile)

		if fromFile == nil {
			continue
		}

		fromPath := fromFile.Path()
		baseFiles = append(baseFiles, files.File{
			Timeout:  headFile.Timeout,
			FilePath: "/" + fromPath,
		})

		if fromPath != toPath {
			renames[fromPath] = toPath
		}

		chunks[toPath] = repository.translateDiffChunks(filePatch.Chunks())
	}

	return &files.List{
		Files:     headFiles,
		BaseFiles: baseFiles,
		Renames:   renames,
		Chunks:    chunks,
	}
}

func (repository *Repository) FetchBaseIfNotPresent() error {
	if repository == nil || repository.baseBranch == "" {
		return nil
	}

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
	if repository.baseBranch == "" {
		return nil
	}

	branchName := plumbing.NewBranchReferenceName(repository.baseBranch)

	worktree, err := repository.repository.Worktree()
	if err != nil {
		return fmt.Errorf("error getting git worktree: %w", err)
	}

	defer repository.restoreHead(worktree)

	if err := worktree.Checkout(&git.CheckoutOptions{
		Branch: branchName,
	}); err != nil {
		return fmt.Errorf("error checking out base branch: %w", err)
	}

	return body()
}

func (repository *Repository) restoreHead(worktree *git.Worktree) {
	checkoutOptions := &git.CheckoutOptions{}
	if repository.headRef.Name().IsBranch() {
		checkoutOptions.Branch = repository.headRef.Name()
	} else {
		checkoutOptions.Hash = repository.headRef.Hash()
	}

	if err := worktree.Checkout(checkoutOptions); err != nil {
		log.Error().Msgf(
			"error restoring git worktree. your worktree may not have been restored to it's original state! %s",
			err,
		)
	}
}

func (repository *Repository) fileFor(
	ignore *ignore.FileIgnore,
	goclocResult *gocloc.Result,
	projectPath,
	relativePath string,
) *files.File {
	fullPath := filepath.Join(projectPath, relativePath)

	fileInfo, err := os.Stat(fullPath)
	if err != nil {
		log.Debug().Msgf("error getting file stat: %s, %s", fullPath, err)
		return nil
	}

	if ignore != nil && ignore.Ignore(projectPath, fullPath, goclocResult, fileInfo) {
		return nil
	}

	return &files.File{
		Timeout:  timeout.Assign(fileInfo, repository.config),
		FilePath: "/" + relativePath,
	}
}
