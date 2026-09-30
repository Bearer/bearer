package git

import (
	"bufio"
	"context"
	"encoding/hex"
	"fmt"
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
	err := captureCommand(context.TODO(), rootDir, []string{"ls-tree", "-r", "-z", commitSHA}, func(stdout io.Reader) error {
		return readTreeRecords(stdout, false, &result)
	})
	if err != nil {
		return result, err
	}
	err = captureCommand(context.TODO(), rootDir, []string{"ls-files", "--others", "--exclude-standard", "-z"}, func(stdout io.Reader) error {
		return readTreeRecords(stdout, true, &result)
	})
	return result, err
}

// Git's -z record terminator is mandatory; EOF inside a record is partial input.
func readTreeRecords(input io.Reader, other bool, result *[]TreeFile) error {
	reader := bufio.NewReader(input)
	for {
		record, err := reader.ReadString(0)
		if err == io.EOF && len(record) == 0 {
			return nil
		}
		if err != nil {
			if err == io.EOF {
				return io.ErrUnexpectedEOF
			}
			return err
		}
		record = strings.TrimSuffix(record, "\x00")
		if record == "" {
			return fmt.Errorf("empty Git inventory record")
		}
		entry := TreeFile{Other: other}
		if other {
			entry.Filename = record
		} else {
			parts := strings.SplitN(record, "\t", 2)
			if len(parts) != 2 || parts[1] == "" {
				return fmt.Errorf("invalid Git tree record")
			}
			metadata := strings.Split(parts[0], " ")
			if len(metadata) != 3 || metadata[0] == "" || (metadata[1] != "blob" && metadata[1] != "tree" && metadata[1] != "commit") {
				return fmt.Errorf("invalid Git tree metadata")
			}
			for _, ch := range metadata[0] {
				if ch < '0' || ch > '7' {
					return fmt.Errorf("invalid Git file mode")
				}
			}
			hash, e := hex.DecodeString(metadata[2])
			if e != nil || (len(hash) != 20 && len(hash) != 32) {
				return fmt.Errorf("invalid Git object id")
			}
			entry.Filename = parts[1]
			entry.SHA = metadata[2]
		}
		*result = append(*result, entry)
	}
}
