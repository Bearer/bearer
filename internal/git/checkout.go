package git

import "fmt"

func checkout(rootDir, ref string, filenames []string) error {
	cmd := logAndBuildCommand(
		"-c",
		"advice.detachedHead=false",
		"checkout",
		ref,
		"--pathspec-from-file=-",
		"--pathspec-file-nul",
	)
	cmd.Dir = rootDir

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}

	logWriter := &debugLogWriter{}
	cmd.Stdout = logWriter
	cmd.Stderr = logWriter

	if err := cmd.Start(); err != nil {
		killProcess(cmd)
		return err
	}

	for _, filename := range filenames {
		_, err := stdin.Write([]byte(filename))
		if err != nil {
			killProcess(cmd)
			return err
		}

		_, err = stdin.Write([]byte{0})
		if err != nil {
			killProcess(cmd)
			return err
		}
	}

	if err := stdin.Close(); err != nil {
		killProcess(cmd)
		return err
	}

	if err := cmd.Wait(); err != nil {
		killProcess(cmd)
		return newError(err, logWriter.AllOutput())
	}

	// Using pathspec with checkout doesn't update the HEAD ref so do it manually
	return basicCommand(rootDir, "update-ref", "HEAD", ref)
}

func fetchBlobsForRange(rootDir, firstCommitSHA, lastCommitSHA string, filenames []string) error {
	objectIDs, err := getObjectIDsForRangeFiles(rootDir, firstCommitSHA, lastCommitSHA, filenames)
	if err != nil {
		return err
	}

	return fetchBlobs(rootDir, objectIDs)
}

// Fetches the given list of objects/blobs.
//
// There's no command in git that does this directly but we can get the desired
// behaviour by creating a pack and throwing it away
func fetchBlobs(rootDir string, objectIDs []string) error {
	cmd := logAndBuildCommand("pack-objects", "--progress", "--stdout")
	cmd.Dir = rootDir

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}

	logWriter := &debugLogWriter{}
	cmd.Stderr = logWriter

	if err := cmd.Start(); err != nil {
		killProcess(cmd)
		return err
	}

	for _, objectID := range objectIDs {
		fmt.Fprintln(stdin, objectID)
	}

	if err := stdin.Close(); err != nil {
		killProcess(cmd)
		return err
	}

	if err := cmd.Wait(); err != nil {
		killProcess(cmd)
		return newError(err, logWriter.AllOutput())
	}

	return nil
}
