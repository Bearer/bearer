package gitrepository

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/go-github/github"
	"github.com/hhatto/gocloc"
	"github.com/rs/zerolog/log"
	"golang.org/x/oauth2"

	"github.com/bearer/bearer/internal/commands/process/filelist/files"
	"github.com/bearer/bearer/internal/commands/process/filelist/ignore"
	"github.com/bearer/bearer/internal/commands/process/filelist/timeout"
	"github.com/bearer/bearer/internal/commands/process/settings"
	"github.com/bearer/bearer/internal/git"
	"github.com/bearer/bearer/internal/util/output"
)

type Repository struct {
	ctx    context.Context
	config settings.Config
	rootPath,
	targetPath,
	gitTargetPath string
	headBranch string
	headCommitHash,
	mergeBaseCommitHash string
	githubToken string
}

func New(
	ctx context.Context,
	config settings.Config,
	targetPath string,
	baseBranch string,
) (*Repository, error) {
	rootPath := git.GetRoot(targetPath)
	if rootPath == "" {
		log.Debug().Msg("no git repository found")

		if baseBranch != "" {
			return nil, errors.New("base branch specified but no git repository found")
		}

		return nil, nil
	}

	gitTargetPath, err := filepath.Rel(rootPath, targetPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get relative target: %w", err)
	}

	log.Debug().Msgf("git target: [%s/]%s", rootPath, gitTargetPath)

	headCommitHash, err := git.GetCurrentCommit(rootPath)
	if err != nil {
		return nil, fmt.Errorf("error getting head ref: %w", err)
	}

	headBranch, err := git.GetCurrentBranch(rootPath)
	if err != nil {
		return nil, fmt.Errorf("error getting head ref: %w", err)
	}

	repository := &Repository{
		ctx:            ctx,
		config:         config,
		rootPath:       rootPath,
		targetPath:     targetPath,
		gitTargetPath:  gitTargetPath,
		headBranch:     headBranch,
		headCommitHash: headCommitHash,
		githubToken:    config.Scan.GithubToken,
	}

	repository.mergeBaseCommitHash, err = repository.fetchMergeBaseCommit(baseBranch)
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

	if repository.mergeBaseCommitHash == "" {
		return repository.getTreeFiles(ignore, goclocResult, repository.headCommitHash)
	}

	return repository.getDiffFiles(ignore, goclocResult)
}

func (repository *Repository) fetchMergeBaseCommit(baseBranch string) (string, error) {
	if baseBranch == "" {
		return "", nil
	}

	hash, err := repository.lookupMergeBaseHash(baseBranch)
	if err != nil {
		return "", fmt.Errorf("error looking up hash: %w", err)
	}

	if hash == "" {
		return "", fmt.Errorf(
			"could not find common ancestor between the current and %s branch. "+
				"please check that the base branch is correct, and that you have "+
				"fetched enough git history to include the latest common ancestor",
			baseBranch,
		)
	}

	log.Debug().Msgf("merge base commit: %s", hash)

	if git.CommitPresent(repository.rootPath, hash) {
		return hash, nil
	}

	log.Debug().Msgf("merge base commit not present, fetching")

	if err := git.FetchRef(repository.ctx, repository.rootPath, hash); err != nil {
		return "", err
	}

	log.Debug().Msgf("merge base commit fetched")

	return hash, nil
}

func (repository *Repository) getTreeFiles(
	ignore *ignore.FileIgnore,
	goclocResult *gocloc.Result,
	commitSHA string,
) (*files.List, error) {
	var headFiles []files.File

	gitFiles, err := git.ListTree(repository.rootPath, commitSHA)
	if err != nil {
		return nil, err
	}

	for _, file := range gitFiles {
		if file := repository.fileFor(
			ignore,
			goclocResult,
			file.Filename,
		); file != nil {
			headFiles = append(headFiles, *file)
		}
	}

	return &files.List{Files: headFiles}, nil
}

func (repository *Repository) getDiffFiles(
	ignore *ignore.FileIgnore,
	goclocResult *gocloc.Result,
) (*files.List, error) {
	var baseFiles, headFiles []files.File
	renames := make(map[string]string)
	chunks := make(map[string]git.Chunks)

	filePatches, err := git.Diff(repository.rootPath, repository.mergeBaseCommitHash)
	if err != nil {
		return nil, err
	}

	for _, patch := range filePatches {
		// we're not interested in removals
		if patch.ToPath == "" {
			continue
		}

		headFile := repository.fileFor(ignore, goclocResult, patch.ToPath)
		if headFile == nil {
			continue
		}

		headFiles = append(headFiles, *headFile)

		if patch.FromPath == "" {
			continue
		}

		relativeFromPath, err := filepath.Rel(repository.gitTargetPath, patch.FromPath)
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

		chunks[headFile.FilePath] = patch.Chunks
	}

	return &files.List{
		Files:     headFiles,
		BaseFiles: baseFiles,
		Renames:   renames,
		Chunks:    chunks,
	}, nil
}

func (repository *Repository) WithBaseBranch(body func() error) error {
	if repository == nil || repository.mergeBaseCommitHash == "" {
		return nil
	}

	if err := git.Switch(repository.rootPath, repository.mergeBaseCommitHash, true); err != nil {
		return fmt.Errorf("error switching to base branch: %w", err)
	}

	err := body()

	if restoreErr := repository.restoreHead(); restoreErr != nil {
		wrappedErr := fmt.Errorf("error restoring to current commit: %w", restoreErr)
		if err == nil {
			return wrappedErr
		}

		output.StdErrLog(wrappedErr.Error())
	}

	return err
}

func (repository *Repository) restoreHead() error {
	if repository.headBranch == "" {
		return git.Switch(repository.rootPath, repository.headCommitHash, true)
	}

	return git.Switch(repository.rootPath, repository.headBranch, false)
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

func (repository *Repository) lookupMergeBaseHash(baseBranch string) (string, error) {
	if sha := os.Getenv("DIFF_BASE_COMMIT"); sha != "" {
		return sha, nil
	}

	if sha, err := repository.lookupMergeBaseHashFromGithub(baseBranch); sha != "" || err != nil {
		return sha, err
	}

	log.Debug().Msg("finding merge base using local repository")

	sha, err := git.GetMergeBase(repository.rootPath, "origin/"+baseBranch, repository.headCommitHash)
	if err != nil {
		if !strings.Contains(err.Error(), "Not a valid object name") {
			return "", fmt.Errorf("invalid ref: %w", err)
		}
	}

	if sha != "" {
		return sha, nil
	}

	log.Debug().Msg("remote ref not found, trying local ref")
	sha, err = git.GetMergeBase(repository.rootPath, baseBranch, repository.headCommitHash)
	if err != nil {
		return "", fmt.Errorf("invalid ref: %w", err)
	}

	return sha, nil
}

func (repository *Repository) lookupMergeBaseHashFromGithub(baseBranch string) (string, error) {
	if repository.githubToken == "" {
		return "", nil
	}

	githubRepository := os.Getenv("GITHUB_REPOSITORY")
	if githubRepository == "" {
		return "", nil
	}

	log.Debug().Msg("finding merge base using github api")

	splitRepository := strings.SplitN(githubRepository, "/", 2)
	if len(splitRepository) != 2 {
		return "", fmt.Errorf("invalid github repository name '%s'", githubRepository)
	}

	client, err := repository.newGithubClient()
	if err != nil {
		return "", err
	}

	comparison, _, err := client.Repositories.CompareCommits(
		context.Background(),
		splitRepository[0],
		splitRepository[1],
		baseBranch,
		repository.headCommitHash,
	)
	if err != nil {
		return "", fmt.Errorf("error calling github compare api: %w", err)
	}

	if comparison.MergeBaseCommit == nil {
		return "", nil
	}

	return *comparison.MergeBaseCommit.SHA, nil
}

func (repository *Repository) newGithubClient() (*github.Client, error) {
	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: repository.githubToken})
	httpClient := oauth2.NewClient(context.Background(), tokenSource)
	client := github.NewClient(httpClient)

	if githubAPIURL := os.Getenv("GITHUB_API_URL"); githubAPIURL != "" {
		parsedURL, err := url.Parse(githubAPIURL)
		if err != nil {
			return nil, fmt.Errorf("failed to parse github api url: %w", err)
		}

		if !strings.HasSuffix(parsedURL.Path, "/") {
			parsedURL.Path += "/"
		}

		client.BaseURL = parsedURL
	}

	return client, nil
}
