package gitrepository

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/gitsight/go-vcsurl"
	"github.com/google/go-github/github"
	"github.com/rs/zerolog/log"
	"golang.org/x/oauth2"

	"github.com/bearer/bearer/internal/flag"
	"github.com/bearer/bearer/internal/git"
)

type Context struct {
	RootDir string
	CurrentBranch,
	DefaultBranch,
	BaseBranch string
	CurrentCommitHash,
	BaseCommitHash string
	OriginURL string
	ID,
	Host,
	Owner,
	Name,
	FullName string
	HasUncommittedChanges bool
}

func NewContext(options *flag.Options) (*Context, error) {
	rootDir := git.GetRoot(options.Target)
	if rootDir == "" {
		return nil, nil
	}

	currentBranch, err := getCurrentBranch(options, rootDir)
	if err != nil {
		return nil, fmt.Errorf("error getting current branch name: %w", err)
	}

	defaultBranch, err := getDefaultBranch(options, rootDir)
	if err != nil {
		return nil, fmt.Errorf("error getting default branch name: %w", err)
	}

	baseBranch, err := getBaseBranch(options)
	if err != nil {
		return nil, fmt.Errorf("error getting base branch name: %w", err)
	}

	currentCommitHash, err := getCurrentCommitHash(options, rootDir)
	if err != nil {
		return nil, fmt.Errorf("error getting current commit hash: %w", err)
	}

	baseCommitHash, err := getBaseCommitHash(options, rootDir, baseBranch, currentCommitHash)
	if err != nil {
		return nil, fmt.Errorf("error getting base commit hash: %w", err)
	}

	hasUncommittedChanges, err := git.HasUncommittedChanges(rootDir)
	if err != nil {
		return nil, fmt.Errorf("error checking for uncommitted changes: %w", err)
	}

	originURL, err := getOriginURL(options, rootDir)
	if err != nil {
		return nil, fmt.Errorf("error getting origin url: %w", err)
	}

	var id, host, owner, name, fullName string
	if originURL != "" {
		urlInfo, err := vcsurl.Parse(originURL)
		if err != nil {
			return nil, fmt.Errorf("couldn't parse origin url: %w", err)
		}

		id = urlInfo.ID
		host = string(urlInfo.Host)
		owner = urlInfo.Username
		name = urlInfo.Name
		fullName = urlInfo.FullName
	}

	return &Context{
		RootDir:               rootDir,
		CurrentBranch:         currentBranch,
		DefaultBranch:         defaultBranch,
		BaseBranch:            baseBranch,
		CurrentCommitHash:     currentCommitHash,
		BaseCommitHash:        baseCommitHash,
		OriginURL:             originURL,
		ID:                    id,
		Host:                  host,
		Owner:                 owner,
		Name:                  name,
		FullName:              fullName,
		HasUncommittedChanges: hasUncommittedChanges,
	}, nil
}

func getCurrentBranch(options *flag.Options, rootDir string) (string, error) {
	if options.CurrentBranch != "" {
		return options.CurrentBranch, nil
	}

	name, err := git.GetCurrentBranch(rootDir)
	if err != nil {
		return "", err
	}

	return name, err
}

func getDefaultBranch(options *flag.Options, rootDir string) (string, error) {
	if options.DefaultBranch != "" {
		return options.DefaultBranch, nil
	}

	return git.GetDefaultBranch(rootDir)
}

func getBaseBranch(options *flag.Options) (string, error) {
	if !options.Diff {
		return "", nil
	}

	if options.DiffBaseBranch != "" {
		return options.DiffBaseBranch, nil
	}

	return "", errors.New(
		"couldn't determine base branch for diff scanning. " +
			"please set the 'BEARER_DIFF_BASE_BRANCH' environment variable",
	)
}

func getCurrentCommitHash(options *flag.Options, rootDir string) (string, error) {
	if options.CurrentCommit != "" {
		return options.CurrentCommit, nil
	}

	return git.GetCurrentCommit(rootDir)
}

func getBaseCommitHash(
	options *flag.Options,
	rootDir string,
	baseBranch string,
	currentCommitHash string,
) (string, error) {
	if baseBranch == "" {
		return "", nil
	}

	if options.DiffBaseCommit != "" {
		return options.DiffBaseCommit, nil
	}

	if hash, err := lookupBaseCommitHashFromGithub(options, baseBranch, currentCommitHash); hash != "" || err != nil {
		return hash, err
	}

	log.Debug().Msg("finding merge base using local repository")
	hash, err := git.GetMergeBase(rootDir, "origin/"+baseBranch, currentCommitHash)
	if err != nil {
		if !strings.Contains(err.Error(), "Not a valid object name") {
			return "", fmt.Errorf("invalid ref: %w", err)
		}
	}

	if hash != "" {
		return hash, nil
	}

	log.Debug().Msg("remote ref not found, trying local ref")
	hash, err = git.GetMergeBase(rootDir, baseBranch, currentCommitHash)
	if err != nil {
		return "", fmt.Errorf("invalid ref: %w", err)
	}

	if hash != "" {
		return hash, nil
	}

	return "", fmt.Errorf(
		"could not find common ancestor between the current and %s branch. "+
			"please check that the base branch is correct, and that you have "+
			"fetched enough git history to include the latest common ancestor",
		baseBranch,
	)
}

func lookupBaseCommitHashFromGithub(options *flag.Options, baseBranch string, currentCommitHash string) (string, error) {
	if options.GithubToken == "" || options.GithubRepository == "" {
		return "", nil
	}

	log.Debug().Msg("finding merge base using github api")

	splitRepository := strings.SplitN(options.GithubRepository, "/", 2)
	if len(splitRepository) != 2 {
		return "", fmt.Errorf("invalid github repository name '%s'", options.GithubRepository)
	}

	client, err := newGithubClient(options)
	if err != nil {
		return "", err
	}

	comparison, _, err := client.Repositories.CompareCommits(
		context.Background(),
		splitRepository[0],
		splitRepository[1],
		baseBranch,
		currentCommitHash,
	)
	if err != nil {
		return "", fmt.Errorf("error calling github compare api: %w", err)
	}

	if comparison.MergeBaseCommit == nil {
		return "", nil
	}

	return *comparison.MergeBaseCommit.SHA, nil
}

func getOriginURL(options *flag.Options, rootDir string) (string, error) {
	if options.OriginURL != "" {
		return options.OriginURL, nil
	}

	return git.GetOriginURL(rootDir)
}

func newGithubClient(options *flag.Options) (*github.Client, error) {
	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: options.GithubToken})
	httpClient := oauth2.NewClient(context.Background(), tokenSource)
	client := github.NewClient(httpClient)

	if options.GithubAPIURL != "" {
		parsedURL, err := url.Parse(options.GithubAPIURL)
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
