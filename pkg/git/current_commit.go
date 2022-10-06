package git

import "bufio"

// GetCurrentCommit gets a current commit from a HEAD for a local directory
func GetCurrentCommit(dir string) (hash string, err error) {
	cmd := logAndBuildCommand(
		"rev-parse",
		"HEAD",
	)
	cmd.Dir = dir

	logWriter := &debugLogWriter{}
	cmd.Stderr = logWriter

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		killProcess(cmd)
		return "", err
	}

	if err := cmd.Start(); err != nil {
		killProcess(cmd)
		return "", err
	}

	if err := cmd.Wait(); err != nil {
		killProcess(cmd)
		return "", newError(err, logWriter.AllOutput())
	}

	scanner := bufio.NewScanner(stdout)
	hashB := scanner.Bytes()

	return string(hashB), nil
}
