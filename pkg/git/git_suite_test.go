package git_test

import (
	"os"
	"os/exec"
	"path"
	"strings"
	"testing"
)

func runGit(t *testing.T, dir string, args ...string) {
	t.Helper()

	command := exec.Command("git", args...)
	command.Dir = dir

	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("failed to run git command [%s]: %s\n%s", strings.Join(args, " "), err, output)
	}
}

func addAndCommit(t *testing.T, dir string) {
	t.Helper()

	runGit(t, dir, "add", ".")
	runGit(t, dir,
		"-c", "user.name=Bearer CI",
		"-c", "user.email=ci@bearer.com",
		"-c", "commit.gpgSign=false",
		"commit",
		"--allow-empty-message",
		"--message=",
	)
}

func writeFile(t *testing.T, tempDir, filename, content string) {
	t.Helper()

	if err := os.WriteFile(path.Join(tempDir, filename), []byte(content), 0600); err != nil {
		t.Fatalf("failed to write file: %s", err)
	}
}
