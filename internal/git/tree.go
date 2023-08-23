package git

import (
	"bufio"
	"fmt"
	"io"
	"path"
	"strings"
	"time"

	"golang.org/x/exp/slices"
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

func GetTree(rootDir string) (*Tree, error) {
	commit, err := getHeadCommitIdentifier(rootDir)
	if err != nil {
		return nil, fmt.Errorf("failed to get commit identifier: %s", err)
	}

	files, err := listTree(rootDir, commit.SHA)
	if err != nil {
		return nil, fmt.Errorf("failed to list tree: %s", err)
	}

	return &Tree{Commit: *commit, Files: files}, nil
}

func listTree(rootDir, commitSHA string) ([]TreeFile, error) {
	result := []TreeFile{}

	cmd := logAndBuildCommand("ls-tree", "-r", "-z", commitSHA)
	cmd.Dir = rootDir

	logWriter := &debugLogWriter{}
	cmd.Stderr = logWriter

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		killProcess(cmd)
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		killProcess(cmd)
		return nil, err
	}

	stdoutBuf := bufio.NewReader(stdout)
	for {
		metadata, err := stdoutBuf.ReadString('\t')
		if err == io.EOF {
			break
		}
		if err != nil {
			killProcess(cmd)
			return nil, err
		}

		splitMeta := strings.Split(metadata[:len(metadata)-1], " ")
		if len(splitMeta) != 3 {
			continue
		}
		sha := splitMeta[2]

		filename, err := stdoutBuf.ReadString(0)
		if err != nil && err != io.EOF {
			killProcess(cmd)
			return nil, err
		}

		if len(filename) > 1 {
			result = append(result, TreeFile{Filename: filename[:len(filename)-1], SHA: sha})
		}
	}

	stdout.Close()

	if err := cmd.Wait(); err != nil {
		killProcess(cmd)
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
		killProcess(cmd)
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		killProcess(cmd)
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
		killProcess(cmd)
		return nil, err
	}

	stdout.Close()

	if err := cmd.Wait(); err != nil {
		killProcess(cmd)
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

	treeFiles, err := listTree(rootDir, commitSHA)
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
	cmd := logAndBuildCommand("log", "-1", "--format=%H %cI")
	cmd.Dir = rootDir

	logWriter := &debugLogWriter{}
	cmd.Stderr = logWriter

	output, err := cmd.Output()
	if err != nil {
		killProcess(cmd)
		return nil, newError(err, logWriter.AllOutput())
	}

	splitOutput := strings.SplitN(strings.TrimSpace(string(output)), " ", 2)

	parsedTimestamp, err := time.Parse(time.RFC3339, splitOutput[1])
	if err != nil {
		killProcess(cmd)
		return nil, err
	}

	return &CommitIdentifier{SHA: splitOutput[0], Timestamp: parsedTimestamp.UTC()}, nil
}

func commitPresent(rootDir, sha string) bool {
	cmd := logAndBuildCommand("cat-file", "-t", sha)
	cmd.Dir = rootDir

	output, _ := cmd.Output()
	return string(output) == "commit\n"
}
