package scanid

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/google/uuid"

	"github.com/bearer/bearer/cmd/bearer/build"
	"github.com/bearer/bearer/internal/commands/process/gitrepository"
	"github.com/bearer/bearer/internal/commands/process/settings"
	"github.com/bearer/bearer/internal/util/file"
)

func Build(scanSettings settings.Config, gitContext *gitrepository.Context) (string, error) {
	var sha string
	if gitContext != nil && !gitContext.HasUncommittedChanges {
		if gitContext.BaseCommitHash == "" {
			sha = gitContext.CurrentCommitHash
		} else {
			sha = gitContext.BaseCommitHash + "_" + gitContext.CurrentCommitHash
		}
	}

	if sha == "" {
		sha = uuid.NewString()
	}

	configHash, err := hashConfig(scanSettings)
	if err != nil {
		return "", err
	}

	// we want sha as it might change detections
	buildSHA := build.CommitSHA
	scanID := strings.TrimSuffix(string(sha), "\n") + "-" + buildSHA + "-" + configHash + ".jsonl"

	return scanID, nil
}

func hashConfig(scanSettings settings.Config) (string, error) {
	ruleHash, err := hashRules(scanSettings.Rules)
	if err != nil {
		return "", fmt.Errorf("error building rule hash: %w", err)
	}

	scannersHash, err := hashScanners(scanSettings.Scan.Scanner)
	if err != nil {
		return "", fmt.Errorf("error building scanners hash: %w", err)
	}

	absTarget, err := file.CanonicalPath(scanSettings.Scan.Target)
	if err != nil {
		return "", fmt.Errorf("error getting absolute path to target: %w", err)
	}

	targetHash := md5.Sum([]byte(absTarget))

	hashBuilder := md5.New()
	if _, err := hashBuilder.Write(targetHash[:]); err != nil {
		return "", err
	}
	if _, err := hashBuilder.Write(ruleHash); err != nil {
		return "", err
	}
	if _, err := hashBuilder.Write(scannersHash); err != nil {
		return "", err
	}

	return hex.EncodeToString(hashBuilder.Sum(nil)[:]), nil
}

func hashRules(rules map[string]*settings.Rule) ([]byte, error) {
	hashBuilder := md5.New()

	var ruleNames []string
	for key := range rules {
		ruleNames = append(ruleNames, key)
	}
	sort.Strings(ruleNames)

	for _, ruleName := range ruleNames {
		detectorContent, err := json.Marshal(rules[ruleName])
		if err != nil {
			return nil, err
		}

		if _, err := hashBuilder.Write([]byte(ruleName)); err != nil {
			return nil, err
		}
		if _, err = hashBuilder.Write(detectorContent); err != nil {
			return nil, err
		}
	}

	return hashBuilder.Sum(nil), nil
}

func hashScanners(scanners []string) ([]byte, error) {
	hashBuilder := md5.New()

	sortedScanners := make([]string, len(scanners))
	copy(sortedScanners, scanners)
	sort.Strings(sortedScanners)

	for _, scanner := range sortedScanners {
		_, err := hashBuilder.Write([]byte(scanner))
		if err != nil {
			return nil, err
		}
	}

	return hashBuilder.Sum(nil), nil
}
