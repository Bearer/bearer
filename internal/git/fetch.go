package git

func FetchBranchLatest(rootDir string, branchName string) error {
	return basicCommand(
		rootDir,
		"git",
		"fetch",
		"--no-tags",
		"--no-recurse-submodules",
		"--depth=1",
		"origin",
		branchName,
	)
}
