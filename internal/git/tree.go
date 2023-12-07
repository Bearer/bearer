package git

import (
	"bufio"
	"context"
	"io"
	"strings"
)

type TreeFile struct {
	Filename string `json:"filename" yaml:"filename"`
	SHA      string `json:"sha" yaml:"sha"`
	Other    bool   `json:"other" yaml:"other"`
}

func HasUncommittedChanges(rootDir string) (bool, error) {
	output, err := captureCommandBasic(
		context.TODO(),
		rootDir,
		"status",
		"--porcelain=v1",
		"--no-renames",
	)

	return strings.TrimSpace(output) != "", err
}

func ListTree(rootDir, commitSHA string) ([]TreeFile, error) {
	result := []TreeFile{}

	err := captureCommand(
		context.TODO(),
		rootDir,
		[]string{"ls-tree", "-r", "-z", commitSHA},
		func(stdout io.Reader) error {
			stdoutBuf := bufio.NewReader(stdout)
			for {
				metadata, err := stdoutBuf.ReadString('\t')
				if err == io.EOF {
					break
				}
				if err != nil {
					return err
				}

				splitMeta := strings.Split(metadata[:len(metadata)-1], " ")
				if len(splitMeta) != 3 {
					continue
				}
				sha := splitMeta[2]

				filename, err := stdoutBuf.ReadString(0)
				if err != nil && err != io.EOF {
					return err
				}

				if len(filename) > 1 {
					result = append(result, TreeFile{Filename: filename[:len(filename)-1], SHA: sha})
				}
			}

			return nil
		},
	)

	if err != nil {
		return result, nil
	}

	err = captureCommand(
		context.TODO(),
		rootDir,
		[]string{"ls-files", "--others", "--exclude-standard", "-z"},
		func(stdout io.Reader) error {
			stdoutBuf := bufio.NewReader(stdout)
			for {
				filename, err := stdoutBuf.ReadString(0)
				if err == io.EOF {
					break
				}
				if err != nil {
					return err
				}

				if len(filename) > 1 {
					result = append(result, TreeFile{Filename: filename[:len(filename)-1], Other: true})
				}
			}

			return nil
		},
	)

	return result, err
}
