package flag

import flagtypes "github.com/bearer/bearer/internal/flag/types"

type repositoryFlagGroup struct{ flagGroupBase }

var RepositoryFlagGroup = &repositoryFlagGroup{flagGroupBase{name: "Repository"}}

var (
	RepositoryURLFlag = RepositoryFlagGroup.add(flagtypes.Flag{
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
	BranchFlag = RepositoryFlagGroup.add(flagtypes.Flag{
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
	CommitFlag = RepositoryFlagGroup.add(flagtypes.Flag{
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
	DefaultBranchFlag = RepositoryFlagGroup.add(flagtypes.Flag{
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
	DiffBaseBranchFlag = RepositoryFlagGroup.add(flagtypes.Flag{
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
	DiffBaseCommitFlag = RepositoryFlagGroup.add(flagtypes.Flag{
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
	GithubTokenFlag = RepositoryFlagGroup.add(flagtypes.Flag{
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
	GithubRepositoryFlag = RepositoryFlagGroup.add(flagtypes.Flag{
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
	GithubAPIURLFlag = RepositoryFlagGroup.add(flagtypes.Flag{
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
	PullRequestNumberFlag = RepositoryFlagGroup.add(flagtypes.Flag{
		Name:       "pull-request-number",
		ConfigName: "repository.pull-request-number",
		Value:      "",
		Usage:      "Used when fetching branch level ignores for a PR/MR",
		EnvironmentVariables: []string{
			"PR_NUMBER",           // github
			"CI_MERGE_REQUEST_ID", //gitlab
		},
		DisableInConfig: true,
		Hide:            true,
	})
)

func (repositoryFlagGroup) SetOptions(options *flagtypes.Options, args []string) error {
	options.RepositoryOptions = flagtypes.RepositoryOptions{
		OriginURL:         getString(RepositoryURLFlag),
		Branch:            getString(BranchFlag),
		Commit:            getString(CommitFlag),
		DefaultBranch:     getString(DefaultBranchFlag),
		DiffBaseBranch:    getString(DiffBaseBranchFlag),
		DiffBaseCommit:    getString(DiffBaseCommitFlag),
		GithubToken:       getString(GithubTokenFlag),
		GithubRepository:  getString(GithubRepositoryFlag),
		GithubAPIURL:      getString(GithubAPIURLFlag),
		PullRequestNumber: getString(PullRequestNumberFlag),
	}

	return nil
}
