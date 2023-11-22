package flag

type repositoryFlagGroup struct{ flagGroupBase }

var RepositoryFlagGroup = &repositoryFlagGroup{flagGroupBase{name: "Repository"}}

var (
	RepositoryURLFlag = RepositoryFlagGroup.add(Flag{
		Name:       "repository-url",
		ConfigName: "repository.url",
		Value:      "",
		Usage:      "The remote URL of the repository.",
		EnvironmentVariables: []string{
			"ORIGIN_URL",        // legacy
			"CI_REPOSITORY_URL", // gitlab
		},
		DisableInConfig: true,
		Hide:            true,
	})
	BranchFlag = RepositoryFlagGroup.add(Flag{
		Name:       "branch",
		ConfigName: "repository.branch",
		Value:      "",
		Usage:      "The name of the branch being scanned.",
		EnvironmentVariables: []string{
			"CURRENT_BRANCH",     // legacy
			"CI_COMMIT_REF_NAME", // gitlab
		},
		DisableInConfig: true,
		Hide:            true,
	})
	CommitFlag = RepositoryFlagGroup.add(Flag{
		Name:       "commit",
		ConfigName: "repository.commit",
		Value:      "",
		Usage:      "The hash of the commit being scanned.",
		EnvironmentVariables: []string{
			"SHA",           // legacy
			"CI_COMMIT_SHA", // gitlab
		},
		DisableInConfig: true,
		Hide:            true,
	})
	DefaultBranchFlag = RepositoryFlagGroup.add(Flag{
		Name:       "default-branch",
		ConfigName: "repository.default-branch",
		Value:      "",
		Usage:      "The name of the default branch.",
		EnvironmentVariables: []string{
			"DEFAULT_BRANCH",    // legacy
			"CI_DEFAULT_BRANCH", // gitlab
		},
		DisableInConfig: true,
		Hide:            true,
	})
	DiffBaseBranchFlag = RepositoryFlagGroup.add(Flag{
		Name:       "diff-base-branch",
		ConfigName: "repository.diff-base-branch",
		Value:      "",
		Usage:      "The name of the base branch to use for diff scanning.",
		EnvironmentVariables: []string{
			"DIFF_BASE_BRANCH",                    // legacy
			"CI_MERGE_REQUEST_TARGET_BRANCH_NAME", // gitlab
		},
		DisableInConfig: true,
		Hide:            true,
	})
	DiffBaseCommitFlag = RepositoryFlagGroup.add(Flag{
		Name:       "diff-base-commit",
		ConfigName: "repository.diff-base-commit",
		Value:      "",
		Usage:      "The hash of the base commit to use for diff scanning.",
		EnvironmentVariables: []string{
			"DIFF_BASE_COMMIT",               // legacy
			"CI_MERGE_REQUEST_DIFF_BASE_SHA", // gitlab
		},
		DisableInConfig: true,
		Hide:            true,
	})
	GithubTokenFlag = RepositoryFlagGroup.add(Flag{
		Name:       "github-token",
		ConfigName: "repository.github-token",
		Value:      "",
		Usage:      "An access token for the Github API.",
		EnvironmentVariables: []string{
			"GITHUB_TOKEN", // github
		},
		DisableInConfig: true,
		Hide:            true,
	})
	GithubRepositoryFlag = RepositoryFlagGroup.add(Flag{
		Name:       "github-repository",
		ConfigName: "repository.github-repository",
		Value:      "",
		Usage:      "The owner and name of the repository on Github. eg. Bearer/bearer",
		EnvironmentVariables: []string{
			"GITHUB_REPOSITORY", // github
		},
		DisableInConfig: true,
		Hide:            true,
	})
	GithubAPIURLFlag = RepositoryFlagGroup.add(Flag{
		Name:       "github-api-url",
		ConfigName: "repository.github-api-url",
		Value:      "",
		Usage:      "A non-standard URL to use for the Github API",
		EnvironmentVariables: []string{
			"GITHUB_API_URL", // github
		},
		DisableInConfig: true,
		Hide:            true,
	})
)

type RepositoryOptions struct {
	OriginURL        string
	Branch           string
	Commit           string
	DefaultBranch    string
	DiffBaseBranch   string
	DiffBaseCommit   string
	GithubToken      string
	GithubRepository string
	GithubAPIURL     string
}

func (repositoryFlagGroup) SetOptions(options *Options, args []string) error {
	options.RepositoryOptions = RepositoryOptions{
		OriginURL:        getString(RepositoryURLFlag),
		Branch:           getString(BranchFlag),
		Commit:           getString(CommitFlag),
		DefaultBranch:    getString(DefaultBranchFlag),
		DiffBaseBranch:   getString(DiffBaseBranchFlag),
		DiffBaseCommit:   getString(DiffBaseCommitFlag),
		GithubToken:      getString(GithubTokenFlag),
		GithubRepository: getString(GithubRepositoryFlag),
		GithubAPIURL:     getString(GithubAPIURLFlag),
	}

	return nil
}
