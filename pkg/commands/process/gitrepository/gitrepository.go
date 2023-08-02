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
	ctx    context.Context
	config settings.Config
	git    *git.Repository
	rootPath,
	targetPath,
	gitTargetPath string
	baseBranch        string
	baseRemoteRefName plumbing.ReferenceName
	headRef           *plumbing.Reference
}

func New(
	ctx context.Context,
	config settings.Config,
	targetPath string,
	baseBranch string,
) (*Repository, error) {
	repository, err := git.PlainOpenWithOptions(targetPath, &git.PlainOpenOptions{
		DetectDotGit: true,
	})
	if err != nil {
		return nil, translateOpenError(baseBranch, err)
	}

	worktree, err := repository.Worktree()
	if err != nil {
		return nil, fmt.Errorf("failed to get worktree: %w", err)
	}

	rootPath := worktree.Filesystem.Root()
	gitTargetPath, err := filepath.Rel(rootPath, targetPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get git relative target: %w", err)
	}

	log.Debug().Msgf("git target: [%s/]%s", rootPath, gitTargetPath)

	headRef, err := repository.Head()
	if err != nil {
		return nil, fmt.Errorf("error getting git head: %w", err)
	}

	return &Repository{
		ctx:               ctx,
		config:            config,
		git:               repository,
		rootPath:          rootPath,
		targetPath:        targetPath,
		gitTargetPath:     gitTargetPath,
		baseBranch:        baseBranch,
		baseRemoteRefName: plumbing.NewRemoteReferenceName("origin", baseBranch),
		headRef:           headRef,
	}, nil
}

func (repository *Repository) ListFiles(
	ignore *ignore.FileIgnore,
	goclocResult *gocloc.Result,
) (*files.List, error) {
	if repository == nil {
		return nil, nil
	}

	headTree, err := repository.treeForRef(repository.headRef)
	if err != nil {
		return nil, fmt.Errorf("failed to get head tree: %w", err)
	}

	if repository.baseBranch == "" {
		return repository.getTreeFiles(ignore, goclocResult, headTree)
	}

	filePatches, err := repository.getDiffPatch(headTree)
	if err != nil {
		return nil, err
	}

	return repository.getDiffFiles(ignore, goclocResult, filePatches)
}

func (repository *Repository) treeForRef(ref *plumbing.Reference) (*object.Tree, error) {
	commit, err := repository.git.CommitObject(ref.Hash())
	if err != nil {
		return nil, fmt.Errorf("failed to get commit for ref: %w", err)
	}

	return commit.Tree()
}

func (repository *Repository) getTreeFiles(
	ignore *ignore.FileIgnore,
	goclocResult *gocloc.Result,
	tree *object.Tree,
) (*files.List, error) {
	var headFiles []files.File

	if err := tree.Files().ForEach(func(f *object.File) error {
		if file := repository.fileFor(
			ignore,
			goclocResult,
			f.Name,
		); file != nil {
			headFiles = append(headFiles, *file)
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return &files.List{Files: headFiles}, nil
}

func (repository *Repository) getDiffPatch(headTree *object.Tree) ([]diff.FilePatch, error) {
	baseRef, err := repository.git.Reference(repository.baseRemoteRefName, true)
	if err != nil {
		return nil, fmt.Errorf("failed to get base ref %s: %w", repository.baseRemoteRefName, err)
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
	filePatches []diff.FilePatch,
) (*files.List, error) {
	var baseFiles, headFiles []files.File
	renames := make(map[string]string)
	chunks := make(map[string]bbftypes.Chunks)

	for _, filePatch := range filePatches {
		fromFile, toFile := filePatch.Files()

		// we're not interested in removals
		if toFile == nil {
			continue
		}

		headFile := repository.fileFor(ignore, goclocResult, toFile.Path())
		if headFile == nil {
			continue
		}

		headFiles = append(headFiles, *headFile)

		if fromFile == nil {
			continue
		}

		relativeFromPath, err := filepath.Rel(repository.gitTargetPath, fromFile.Path())
		if err != nil {
			return nil, err
		}
		baseFiles = append(baseFiles, files.File{
			Timeout:  headFile.Timeout,
			FilePath: relativeFromPath,
		})

		if relativeFromPath != headFile.FilePath {
			renames[relativeFromPath] = headFile.FilePath
		}

		chunks[headFile.FilePath] = repository.translateDiffChunks(filePatch.Chunks())
	}

	return &files.List{
		Files:     headFiles,
		BaseFiles: baseFiles,
		Renames:   renames,
		Chunks:    chunks,
	}, nil
}

func (repository *Repository) FetchBaseIfNotPresent() error {
	if repository == nil || repository.baseBranch == "" {
		return nil
	}

	ref, err := repository.git.Reference(repository.baseRemoteRefName, true)
	if err != nil && err != plumbing.ErrReferenceNotFound {
		return fmt.Errorf("invalid branch %s: %w", repository.baseBranch, err)
	}

	// Already exists
	if ref != nil {
		return nil
	}

	localRefName := plumbing.NewBranchReferenceName(repository.baseBranch)
	if err := repository.git.FetchContext(repository.ctx, &git.FetchOptions{
		RefSpecs: []config.RefSpec{config.RefSpec(
			fmt.Sprintf("+%s:%s", localRefName, repository.baseRemoteRefName),
		)},
		Depth: 1,
		Tags:  git.NoTags,
	}); err != nil {
		return fmt.Errorf("error fetching branch %s: %w", repository.baseBranch, err)
	}

	return nil
}

func (repository *Repository) WithBaseBranch(body func() error) error {
	if repository == nil || repository.baseBranch == "" {
		return nil
	}

	worktree, err := repository.git.Worktree()
	if err != nil {
		return fmt.Errorf("error getting git worktree: %w", err)
	}

	defer repository.restoreHead(worktree)

	if err := worktree.Checkout(&git.CheckoutOptions{
		Branch: repository.baseRemoteRefName,
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
	gitRelativePath string,
) *files.File {
	relativePath, err := filepath.Rel(repository.gitTargetPath, gitRelativePath)
	if err != nil {
		log.Error().Msg(err.Error())
		return nil
	}

	if strings.Contains(relativePath, "..") {
		return nil
	}

	fullPath := filepath.Join(repository.targetPath, relativePath)

	fileInfo, err := os.Stat(fullPath)
	if err != nil {
		log.Debug().Msgf("error getting file stat: %s, %s", fullPath, err)
		return nil
	}

	if ignore != nil && ignore.Ignore(repository.targetPath, fullPath, goclocResult, fileInfo) {
		return nil
	}

	return &files.File{
		Timeout:  timeout.Assign(fileInfo, repository.config),
		FilePath: relativePath,
	}
}

func translateOpenError(baseBranch string, err error) error {
	if err != git.ErrRepositoryNotExists {
		return err
	}

	log.Debug().Msg("no git repository found")

	if baseBranch != "" {
		return errors.New("base branch specified but no git repository found")
	}

	return nil
}
