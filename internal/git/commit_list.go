package git

import (
	"bufio"
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"
)

var coAuthorValidPattern = regexp.MustCompile(`<.*>`)

type CommitInfo struct {
	CommitIdentifier
	Committer string   `json:"committer" yaml:"committer"`
	Author    string   `json:"author" yaml:"author"`
	CoAuthors []string `json:"co_authors" yaml:"co_authors"`
}

func GetCommitList(rootDir, firstCommitSHA, lastCommitSHA string) ([]CommitInfo, error) {
	separator := "---"
	result := []CommitInfo{}

	cmd := logAndBuildCommand(
		context.TODO(),
		"log",
		"--first-parent",
		"--format=%H %cI%n%cN <%cE>%n%aN <%aE>%n%(trailers:unfold,valueonly,key=Co-authored-by)"+separator,
		lastCommitSHA,
		"--",
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
	info := CommitInfo{}
	n := 0

	for scanner.Scan() {
		line := scanner.Text()
		if line == separator {
			result = append(result, info)

			if info.SHA == firstCommitSHA {
				break
			}

			info = CommitInfo{}
			n = 0
			continue
		}

		switch n {
		case 0:
			splitLine := strings.SplitN(line, " ", 2)

			parsedTimestamp, err := time.Parse(time.RFC3339, splitLine[1])
			if err != nil {
				return nil, err
			}

			info.SHA = splitLine[0]
			info.Timestamp = parsedTimestamp.UTC()
		case 1:
			info.Committer = line
		case 2:
			info.Author = line
		default:
			info.CoAuthors = append(info.CoAuthors, line)
		}

		n++
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

	if err := translateCoAuthors(rootDir, result); err != nil {
		return nil, err
	}

	return result, nil
}

func translateCoAuthors(rootDir string, commitList []CommitInfo) error {
	translatedCoAuthors := make(map[string]string)
	toTranslate := []string{}

	for _, commitInfo := range commitList {
		for _, author := range commitInfo.CoAuthors {
			translated := author
			if coAuthorValidPattern.MatchString(author) {
				translated = ""
				toTranslate = append(toTranslate, author)
			}

			translatedCoAuthors[author] = translated
		}
	}

	mailmapAuthors, err := checkMailmap(rootDir, toTranslate)
	if err != nil {
		return err
	}

	for i, author := range toTranslate {
		translatedCoAuthors[author] = mailmapAuthors[i]
	}

	for i := range commitList {
		commitInfo := &commitList[i]

		for j := range commitInfo.CoAuthors {
			commitInfo.CoAuthors[j] = translatedCoAuthors[commitInfo.CoAuthors[j]]
		}
	}

	return nil
}

func checkMailmap(rootDir string, authors []string) ([]string, error) {
	cmd := logAndBuildCommand(context.TODO(), "check-mailmap", "--stdin")
	cmd.Dir = rootDir

	stdin, err := cmd.StdinPipe()
	if err != nil {
		cmd.Cancel() //nolint:errcheck
		return nil, err
	}

	go func() {
		for _, author := range authors {
			fmt.Fprintln(stdin, author)
		}

		stdin.Close()
	}()

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

	result := []string{}
	scanner := bufio.NewScanner(stdout)

	for scanner.Scan() {
		result = append(result, scanner.Text())
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
