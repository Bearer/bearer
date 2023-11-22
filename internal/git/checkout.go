package git

import (
	"context"
)

func Switch(rootDir, ref string, detach bool) error {
	args := []string{"switch"}

	if detach {
		args = append(args, "--detach")
	}

	return basicCommand(
		context.TODO(),
		rootDir,
		append(args, ref)...,
	)
}
