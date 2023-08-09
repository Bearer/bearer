package gitrepository

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/format/diff"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/client"
	githttp "github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/google/go-github/github"
	"github.com/hhatto/gocloc"
	"github.com/rs/zerolog/log"
	"golang.org/x/oauth2"

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
	baseRemoteRefName plumbing.ReferenceName
	headRef           *plumbing.Reference
	headCommit,
	mergeBaseCommit *object.Commit
	githubToken string
}

func New(
	ctx context.Context,
	config settings.Config,
	targetPath string,
	baseBranch string,
) (*Repository, error) {
	gitRepository, err := git.PlainOpenWithOptions(targetPath, &git.PlainOpenOptions{
		DetectDotGit: true,
	})
	if err != nil {
		return nil, translateOpenError(baseBranch, err)
	}

	worktree, err := gitRepository.Worktree()
	if err != nil {
		return nil, fmt.Errorf("failed to get worktree: %w", err)
	}

	rootPath := worktree.Filesystem.Root()
	gitTargetPath, err := filepath.Rel(rootPath, targetPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get relative target: %w", err)
	}

	log.Debug().Msgf("git target: [%s/]%s", rootPath, gitTargetPath)

	headRef, err := gitRepository.Head()
	if err != nil {
		return nil, fmt.Errorf("error getting head ref: %w", err)
	}

	headCommit, err := gitRepository.CommitObject(headRef.Hash())
	if err != nil {
		return nil, fmt.Errorf("failed to get head commit: %w", err)
	}

	repository := &Repository{
		ctx:               ctx,
		config:            config,
		git:               gitRepository,
		rootPath:          rootPath,
		targetPath:        targetPath,
		gitTargetPath:     gitTargetPath,
		baseRemoteRefName: plumbing.NewRemoteReferenceName("origin", baseBranch),
		headRef:           headRef,
		headCommit:        headCommit,
		githubToken:       config.Scan.GithubToken,
	}

	repository.mergeBaseCommit, err = repository.fetchMergeBaseCommit(baseBranch)
	if err != nil {
		return nil, err
	}

	return repository, nil
}

func (repository *Repository) ListFiles(
	ignore *ignore.FileIgnore,
	goclocResult *gocloc.Result,
) (*files.List, error) {
	if repository == nil {
		return nil, nil
	}

	headTree, err := repository.headCommit.Tree()
	if err != nil {
		return nil, fmt.Errorf("failed to get head tree: %w", err)
	}

	if repository.mergeBaseCommit == nil {
		return repository.getTreeFiles(ignore, goclocResult, headTree)
	}

	filePatches, err := repository.getDiffPatch(headTree)
	if err != nil {
		return nil, err
	}

	return repository.getDiffFiles(ignore, goclocResult, filePatches)
}

func (repository *Repository) fetchMergeBaseCommit(baseBranch string) (*object.Commit, error) {
	if baseBranch == "" {
		return nil, nil
	}

	hash, err := repository.lookupMergeBaseHash(baseBranch)
	if err != nil {
		return nil, fmt.Errorf("error looking up hash: %w", err)
	}

	if hash == nil {
		return nil, fmt.Errorf(
			"could not find common ancestor between the current and %s branch. "+
				"please check that the base branch is correct, and that you have "+
				"fetched enough git history to include the latest common ancestor",
			baseBranch,
		)
	}

	log.Debug().Msgf("merge base commit: %s", hash)

	commit, err := repository.git.CommitObject(*hash)
	if err == nil {
		return commit, nil
	}

	if err != plumbing.ErrObjectNotFound {
		return nil, fmt.Errorf("error looking up commit: %w", err)
	}

	log.Debug().Msgf("merge base commit not present, fetching")

	if err := repository.git.FetchContext(repository.ctx, &git.FetchOptions{
		RefSpecs: []config.RefSpec{config.RefSpec(
			fmt.Sprintf("+%s:%s", hash.String(), repository.baseRemoteRefName),
		)},
		Depth: 1,
		Tags:  git.NoTags,
	}); err != nil {
		return nil, fmt.Errorf("error fetching: %w", err)
	}

	log.Debug().Msgf("merge base commit fetched")

	return repository.git.CommitObject(*hash)
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
	baseTree, err := repository.mergeBaseCommit.Tree()
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

func (repository *Repository) WithBaseBranch(body func() error) error {
	if repository == nil || repository.mergeBaseCommit == nil {
		return nil
	}

	worktree, err := repository.git.Worktree()
	if err != nil {
		return fmt.Errorf("error getting git worktree: %w", err)
	}

	defer repository.restoreHead(worktree)

	if err := worktree.Checkout(&git.CheckoutOptions{
		Hash: repository.mergeBaseCommit.Hash,
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

func (repository *Repository) lookupMergeBaseHash(baseBranch string) (*plumbing.Hash, error) {
	if hash := repository.lookupMergeBaseRefFromVariable(); hash != nil {
		return hash, nil
	}

	if hash, err := repository.lookupMergeBaseRefFromGithub(baseBranch); hash != nil || err != nil {
		return hash, err
	}

	log.Debug().Msg("finding merge base using local repository")

	ref, err := repository.git.Reference(repository.baseRemoteRefName, true)
	if err != nil {
		if err == plumbing.ErrReferenceNotFound {
			return nil, nil
		}

		return nil, fmt.Errorf("invalid ref: %w", err)
	}

	baseCommit, err := repository.git.CommitObject(ref.Hash())
	if err != nil {
		if err == plumbing.ErrObjectNotFound {
			return nil, nil
		}

		return nil, fmt.Errorf("error looking up base commit: %w", err)
	}

	commonAncestors, err := repository.headCommit.MergeBase(baseCommit)
	if err != nil {
		if err == plumbing.ErrObjectNotFound {
			return nil, nil
		}

		return nil, fmt.Errorf("error computing merge base: %w", err)
	}
	if len(commonAncestors) == 0 {
		return nil, nil
	}

	return &commonAncestors[0].Hash, nil
}

func (repository *Repository) lookupMergeBaseRefFromVariable() *plumbing.Hash {
	sha := os.Getenv("DIFF_BASE_COMMIT")
	if sha == "" {
		return nil
	}

	hash := plumbing.NewHash(sha)
	return &hash
}

func (repository *Repository) lookupMergeBaseRefFromGithub(baseBranch string) (*plumbing.Hash, error) {
	if repository.githubToken == "" {
		return nil, nil
	}

	githubRepository := os.Getenv("GITHUB_REPOSITORY")
	if githubRepository == "" {
		return nil, nil
	}

	log.Debug().Msg("finding merge base using github api")

	splitRepository := strings.SplitN(githubRepository, "/", 2)
	if len(splitRepository) != 2 {
		return nil, fmt.Errorf("invalid github repository name '%s'", githubRepository)
	}

	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: repository.githubToken})
	httpClient := oauth2.NewClient(context.Background(), tokenSource)
	client := github.NewClient(httpClient)

	if githubAPIURL := os.Getenv("GITHUB_API_URL"); githubAPIURL != "" {
		parsedURL, err := url.Parse(githubAPIURL)
		if err != nil {
			return nil, fmt.Errorf("failed to parse github api url: %w", err)
		}

		client.BaseURL = parsedURL
	}

	comparison, _, err := client.Repositories.CompareCommits(
		context.Background(),
		splitRepository[0],
		splitRepository[1],
		baseBranch,
		repository.headRef.Hash().String(),
	)
	if err != nil {
		return nil, fmt.Errorf("error calling github compare api: %w", err)
	}

	if comparison.MergeBaseCommit == nil {
		return nil, nil
	}

	hash := plumbing.NewHash(*comparison.MergeBaseCommit.SHA)
	return &hash, nil
}

type GithubTransport struct {
	githubToken string
}

func (transport *GithubTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.SetBasicAuth("x-access-token", transport.githubToken)

	return http.DefaultTransport.RoundTrip(req)
}

func ConfigureGithubAuth(githubToken string) {
	if githubToken == "" {
		return
	}

	githubClient := githttp.NewClient(&http.Client{
		Transport: &GithubTransport{githubToken: githubToken},
	})

	client.InstallProtocol("http", githubClient)
	client.InstallProtocol("https", githubClient)
}
