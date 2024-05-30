package git

import "context"

func FetchRef(ctx context.Context, rootDir string, ref string) error {
	return basicCommand(
		ctx,
		rootDir,
		"fetch",
		"--no-tags",
		"--no-recurse-submodules",
		"--depth=1",
		"origin",
		ref,
	)
}
