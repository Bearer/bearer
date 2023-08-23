package flag

var (
	FetchBranchFlag = Flag{
		Name:       "branch",
		ConfigName: "repository.branch",
		Value:      "",
		Usage:      "Specify the branch name to be scanned.",
	}
	FetchCommitFlag = Flag{
		Name:       "commit",
		ConfigName: "repository.commit",
		Value:      "",
		Usage:      "Specify the commit hash to be scanned.",
	}
	FetchTagFlag = Flag{
		Name:       "tag",
		ConfigName: "repository.tag",
		Value:      "",
		Usage:      "Specify the tag name to be scanned.",
	}
)

type RepoFlagGroup struct {
	Branch *Flag
	Commit *Flag
	Tag    *Flag
}

type RepoOptions struct {
	RepoBranch string `mapstructure:"branch" json:"branch" yaml:"branch"`
	RepoCommit string `mapstructure:"commit" json:"commit" yaml:"commit"`
	RepoTag    string `mapstructure:"tag" json:"tag" yaml:"tag"`
}

func NewRepoFlagGroup() *RepoFlagGroup {
	return &RepoFlagGroup{
		Branch: &FetchBranchFlag,
		Commit: &FetchCommitFlag,
		Tag:    &FetchTagFlag,
	}
}

func (f *RepoFlagGroup) Name() string {
	return "Repository"
}

func (f *RepoFlagGroup) Flags() []*Flag {
	return []*Flag{f.Branch, f.Commit, f.Tag}
}

func (f *RepoFlagGroup) ToOptions() RepoOptions {
	return RepoOptions{
		RepoBranch: getString(f.Branch),
		RepoCommit: getString(f.Commit),
		RepoTag:    getString(f.Tag),
	}
}
