package git

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

// CloneLatest retrieves the latest commit for the given repository and
// extracts it into a subfolder of targetDir. It returns the subfolder name

func CloneLatest(targetDir string, url string) (string, error) {
	cmd := exec.Command("git", "clone", "--depth=1", url)
	cmd.Env = append(os.Environ(), "GIT_TERMINAL_PROMPT=0", "GIT_LFS_SKIP_SMUDGE=1", "GIT_CURL_VERBOSE=1")
	cmd.Dir = targetDir

	if output, err := cmd.CombinedOutput(); err != nil {
		return "", errors.New(fmt.Sprint("clone failed: ", err, "\n-- git output --\n", string(output), "-- end git output --"))
	}

	files, err := ioutil.ReadDir(targetDir)
	if err != nil {
		return "", err
	}

	return filepath.Join(targetDir, files[0].Name()), nil
}
