package git_test

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestGit(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Git Suite")
}

func runGit(dir string, args ...string) {
	command := exec.Command("git", args...)
	command.Dir = dir

	output, err := command.CombinedOutput()
	if err != nil {
		Fail(fmt.Sprintf("failed to run git command: %s\n%s", err, output))
	}
}

func addAndCommit(dir string) {
	runGit(dir, "add", ".")
	runGit(dir, "commit", "--allow-empty-message", "--message=")
}

func writeFile(tempDir, filename, content string) {
	Expect(os.WriteFile(path.Join(tempDir, filename), []byte(content), 0600)).To(Succeed())
}
