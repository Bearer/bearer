package git

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"path"
	"slices"
	"strings"
	"time"

	"github.com/bearer/bearer/internal/util/file"
)

const blankID = "0000000000000000000000000000000000000000"

type Tree struct {
	Commit CommitIdentifier `json:"commit" yaml:"commit"`
	Files  []TreeFile       `json:"files" yaml:"files"`
}

type TreeFile struct {
	Filename string `json:"filename" yaml:"filename"`
	SHA      string `json:"sha" yaml:"sha"`
}

func GetRoot(targetPath string) string {
	dir := targetPath
	if !file.IsDir(dir) {
		dir = path.Dir(dir)
	}

	command := logAndBuildCommand(context.TODO(), "rev-parse", "--show-toplevel")
	command.Dir = dir

	output, err := command.Output()
	if err != nil {
		return ""
	}

	path := strings.TrimSpace(string(output))
	if path == "" {
		return ""
	}

	canonicalPath, _ := file.CanonicalPath(path)
	return canonicalPath
}

func HasUncommittedChanges(rootDir string) (bool, error) {
	output, err := captureCommandBasic(
		context.TODO(),
		rootDir,
		"status",
		" --porcelain=v1",
		" --no-renames",
	)

	return strings.Count(output, "\n") > 0, err
}

func GetTree(rootDir string) (*Tree, error) {
	commit, err := getHeadCommitIdentifier(rootDir)
	if err != nil {
		return nil, fmt.Errorf("failed to get commit identifier: %s", err)
	}

	files, err := ListTree(rootDir, commit.SHA)
	if err != nil {
		return nil, fmt.Errorf("failed to list tree: %s", err)
	}

	return &Tree{Commit: *commit, Files: files}, nil
}

func ListTree(rootDir, commitSHA string) ([]TreeFile, error) {
	result := []TreeFile{}

	cmd := logAndBuildCommand(context.TODO(), "ls-tree", "-r", "-z", commitSHA)
	cmd.Dir = rootDir

	logWriter := &debugLogWriter{}
	cmd.Stderr = logWriter

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		cmd.Cancel() //nolint:errcheck
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		cmd.Cancel() //nolint:errcheck
		return nil, err
	}

	stdoutBuf := bufio.NewReader(stdout)
	for {
		metadata, err := stdoutBuf.ReadString('\t')
		if err == io.EOF {
			break
		}
		if err != nil {
			cmd.Cancel() //nolint:errcheck
			return nil, err
		}

		splitMeta := strings.Split(metadata[:len(metadata)-1], " ")
		if len(splitMeta) != 3 {
			continue
		}
		sha := splitMeta[2]

		filename, err := stdoutBuf.ReadString(0)
		if err != nil && err != io.EOF {
			cmd.Cancel() //nolint:errcheck
			return nil, err
		}

		if len(filename) > 1 {
			result = append(result, TreeFile{Filename: filename[:len(filename)-1], SHA: sha})
		}
	}

	stdout.Close()

	if err := cmd.Wait(); err != nil {
		cmd.Cancel() //nolint:errcheck
		return nil, newError(err, logWriter.AllOutput())
	}

	return result, nil
}

func getObjectIDsForRangeFiles(rootDir, firstCommitSHA, lastCommitSHA string, filenames []string) ([]string, error) {
	firstCommitFileObjectIDs, err := getObjectIDsForFiles(rootDir, firstCommitSHA, filenames)
	if err != nil {
		return nil, err
	}

	rangeUsedObjectIDs, err := getObjectIDsUsedByRange(rootDir, firstCommitSHA, lastCommitSHA)
	if err != nil {
		return nil, err
	}

	ids := append(firstCommitFileObjectIDs, rangeUsedObjectIDs...)
	slices.Sort(ids)
	return slices.Compact(ids), nil
}

// Returns all the object ids of files touched by the given range of commits.
func getObjectIDsUsedByRange(rootDir, firstCommitSHA, lastCommitSHA string) ([]string, error) {
	result := []string{}

	cmd := logAndBuildCommand(
		context.TODO(),
		"log",
		"--no-renames",
		"--first-parent",
		"--raw",
		"--no-abbrev",
		"--format=%H",
		firstCommitSHA+".."+lastCommitSHA,
	)

	cmd.Dir = rootDir

	logWriter := &debugLogWriter{}
	cmd.Stderr = logWriter

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		cmd.Cancel() //nolint:errcheck
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		cmd.Cancel() //nolint:errcheck
		return nil, err
	}

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" || line[0] != ':' {
			continue
		}

		splitLine := strings.Split(line, " ")
		if splitLine[2] != blankID {
			result = append(result, splitLine[2])
		}
		if splitLine[3] != blankID {
			result = append(result, splitLine[3])
		}
	}

	if err := scanner.Err(); err != nil {
		cmd.Cancel() //nolint:errcheck
		return nil, err
	}

	stdout.Close()

	if err := cmd.Wait(); err != nil {
		cmd.Cancel() //nolint:errcheck
		return nil, newError(err, logWriter.AllOutput())
	}

	return result, nil
}

// Get's the object ids of the given filenames in the specified commit.
// Also includes special git metadata files
func getObjectIDsForFiles(rootDir, commitSHA string, filenames []string) ([]string, error) {
	wantedFilenames := make(map[string]struct{})
	for _, filename := range filenames {
		wantedFilenames[filename] = struct{}{}
	}

	treeFiles, err := ListTree(rootDir, commitSHA)
	if err != nil {
		return nil, err
	}

	objectIDs := []string{}

	for _, treeFile := range treeFiles {
		if _, ok := wantedFilenames[treeFile.Filename]; ok {
			objectIDs = append(objectIDs, treeFile.SHA)
			continue
		}

		if slices.Contains(specialFiles, path.Base(treeFile.Filename)) {
			objectIDs = append(objectIDs, treeFile.SHA)
		}
	}

	return objectIDs, nil
}

func getHeadCommitIdentifier(rootDir string) (*CommitIdentifier, error) {
	cmd := logAndBuildCommand(context.TODO(), "log", "-1", "--format=%H %cI")
	cmd.Dir = rootDir

	logWriter := &debugLogWriter{}
	cmd.Stderr = logWriter

	output, err := cmd.Output()
	if err != nil {
		cmd.Cancel() //nolint:errcheck
		return nil, newError(err, logWriter.AllOutput())
	}

	splitOutput := strings.SplitN(strings.TrimSpace(string(output)), " ", 2)

	parsedTimestamp, err := time.Parse(time.RFC3339, splitOutput[1])
	if err != nil {
		cmd.Cancel() //nolint:errcheck
		return nil, err
	}

	return &CommitIdentifier{SHA: splitOutput[0], Timestamp: parsedTimestamp.UTC()}, nil
}
