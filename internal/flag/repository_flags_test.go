package flag

import (
	"testing"
)

func Test_getRepositoryURLFlag(t *testing.T) {
	testCases := []TestCase{
		{
			name:      "Repository URL. Default",
			flag:      RepositoryURLFlag,
			flagValue: "",
			want:      nil,
		},
		{
			name: "Repository URL. ORIGIN_URL env",
			flag: RepositoryURLFlag,
			env: Env{
				key:   "ORIGIN_URL",
				value: "https://example.com",
			},
			want: []string{
				string("https://example.com"),
			},
		},
		{
			name: "Repository URL. CI_REPOSITORY_URL env",
			flag: RepositoryURLFlag,
			env: Env{
				key:   "CI_REPOSITORY_URL",
				value: "https://example.com",
			},
			want: []string{
				string("https://example.com"),
			},
		},
	}

	RunFlagTests(testCases, t)
}

func Test_getRepositoryBranchFlag(t *testing.T) {
	testCases := []TestCase{
		{
			name:      "Repository Branch. Default",
			flag:      BranchFlag,
			flagValue: "",
			want:      nil,
		},
		{
			name: "Repository Branch. CURRENT_BRANCH env",
			flag: BranchFlag,
			env: Env{
				key:   "CURRENT_BRANCH",
				value: "main",
			},
			want: []string{
				string("main"),
			},
		},
		{
			name: "Repository Branch. CI_COMMIT_REF_NAME env",
			flag: BranchFlag,
			env: Env{
				key:   "CI_COMMIT_REF_NAME",
				value: "main",
			},
			want: []string{
				string("main"),
			},
		},
	}

	RunFlagTests(testCases, t)
}

func Test_getRepositoryCommitFlag(t *testing.T) {
	testCases := []TestCase{
		{
			name:      "Repository Commit. Default",
			flag:      CommitFlag,
			flagValue: "",
			want:      nil,
		},
		{
			name: "Repository Commit. SHA env",
			flag: CommitFlag,
			env: Env{
				key:   "SHA",
				value: "abc123",
			},
			want: []string{
				string("abc123"),
			},
		},
		{
			name: "Repository Commit. CI_COMMIT_SHA env",
			flag: CommitFlag,
			env: Env{
				key:   "CI_COMMIT_SHA",
				value: "abc123",
			},
			want: []string{
				string("abc123"),
			},
		},
	}

	RunFlagTests(testCases, t)
}

func Test_getRepositoryDefaultBranchFlag(t *testing.T) {
	testCases := []TestCase{
		{
			name:      "Repository Default Branch. Default",
			flag:      DefaultBranchFlag,
			flagValue: "",
			want:      nil,
		},
		{
			name: "Repository Default Branch. DEFAULT_BRANCH env",
			flag: DefaultBranchFlag,
			env: Env{
				key:   "DEFAULT_BRANCH",
				value: "main",
			},
			want: []string{
				string("main"),
			},
		},
		{
			name: "Repository Default Branch. CI_DEFAULT_BRANCH env",
			flag: DefaultBranchFlag,
			env: Env{
				key:   "CI_DEFAULT_BRANCH",
				value: "main",
			},
			want: []string{
				string("main"),
			},
		},
	}

	RunFlagTests(testCases, t)
}

func Test_getRepositoryDiffBaseBranchFlag(t *testing.T) {
	testCases := []TestCase{
		{
			name:      "Repository Diff Base Branch. Default",
			flag:      DiffBaseBranchFlag,
			flagValue: "",
			want:      nil,
		},
		{
			name: "Repository Diff Base Branch. DIFF_BASE_BRANCH env",
			flag: DiffBaseBranchFlag,
			env: Env{
				key:   "DIFF_BASE_BRANCH",
				value: "main",
			},
			want: []string{
				string("main"),
			},
		},
		{
			name: "Repository Diff Base Branch. CI_MERGE_REQUEST_TARGET_BRANCH_NAME env",
			flag: DiffBaseBranchFlag,
			env: Env{
				key:   "CI_MERGE_REQUEST_TARGET_BRANCH_NAME",
				value: "main",
			},
			want: []string{
				string("main"),
			},
		},
	}

	RunFlagTests(testCases, t)
}

func Test_getRepositoryDiffBaseCommitFlag(t *testing.T) {
	testCases := []TestCase{
		{
			name:      "Repository Diff Base Commit. Default",
			flag:      DiffBaseCommitFlag,
			flagValue: "",
			want:      nil,
		},
		{
			name: "Repository Diff Base Commit. DIFF_BASE_COMMIT env",
			flag: DiffBaseCommitFlag,
			env: Env{
				key:   "DIFF_BASE_COMMIT",
				value: "abc123",
			},
			want: []string{
				string("abc123"),
			},
		},
		{
			name: "Repository Diff Base Commit. CI_MERGE_REQUEST_DIFF_BASE_SHA env",
			flag: DiffBaseCommitFlag,
			env: Env{
				key:   "CI_MERGE_REQUEST_DIFF_BASE_SHA",
				value: "abc123",
			},
			want: []string{
				string("abc123"),
			},
		},
	}

	RunFlagTests(testCases, t)
}

func Test_getRepositoryGithubTokenFlag(t *testing.T) {
	testCases := []TestCase{
		{
			name:      "Repository GithubTokenFlag. Default",
			flag:      GithubTokenFlag,
			flagValue: "",
			want:      nil,
		},
		{
			name: "Repository GithubTokenFlag. GITHUB_TOKEN env",
			flag: GithubTokenFlag,
			env: Env{
				key:   "GITHUB_TOKEN",
				value: "abc123",
			},
			want: []string{
				string("abc123"),
			},
		},
	}

	RunFlagTests(testCases, t)
}

func Test_getRepositoryGithubRepositoryFlag(t *testing.T) {
	testCases := []TestCase{
		{
			name:      "Repository GithubRepositoryFlag. Default",
			flag:      GithubRepositoryFlag,
			flagValue: "",
			want:      nil,
		},
		{
			name: "Repository GithubRepositoryFlag. GITHUB_REPOSITORY env",
			flag: GithubRepositoryFlag,
			env: Env{
				key:   "GITHUB_REPOSITORY",
				value: "Bearer/bearer",
			},
			want: []string{
				string("Bearer/bearer"),
			},
		},
	}

	RunFlagTests(testCases, t)
}

func Test_getRepositoryGithubAPIURLFlag(t *testing.T) {
	testCases := []TestCase{
		{
			name:      "Repository GithubAPIURLFlag. Default",
			flag:      GithubAPIURLFlag,
			flagValue: "",
			want:      nil,
		},
		{
			name: "Repository GithubAPIURLFlag. GITHUB_API_URL env",
			flag: GithubAPIURLFlag,
			env: Env{
				key:   "GITHUB_API_URL",
				value: "https://github.com/bearer/bearer",
			},
			want: []string{
				string("https://github.com/bearer/bearer"),
			},
		},
	}

	RunFlagTests(testCases, t)
}

func Test_getPullRequestNumberFlag(t *testing.T) {
	testCases := []TestCase{
		{
			name:      "Repository PullRequestNumber. Default",
			flag:      PullRequestNumberFlag,
			flagValue: "",
			want:      nil,
		},
		{
			name: "Repository PullRequestNumber. PR_NUMBER env",
			flag: PullRequestNumberFlag,
			env: Env{
				key:   "PR_NUMBER",
				value: "42",
			},
			want: []string{
				string("42"),
			},
		},
		{
			name: "Repository PullRequestNumber. CI_MERGE_REQUEST_ID env",
			flag: PullRequestNumberFlag,
			env: Env{
				key:   "CI_MERGE_REQUEST_ID",
				value: "24",
			},
			want: []string{
				string("24"),
			},
		},
	}

	RunFlagTests(testCases, t)
}
