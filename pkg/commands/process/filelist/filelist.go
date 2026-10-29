package filelist

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/hhatto/gocloc"
	"github.com/rs/zerolog/log"

	flfiles "github.com/bearer/bearer/pkg/commands/process/filelist/files"
	"github.com/bearer/bearer/pkg/commands/process/filelist/ignore"
	"github.com/bearer/bearer/pkg/commands/process/filelist/timeout"
	"github.com/bearer/bearer/pkg/commands/process/gitrepository"
	"github.com/bearer/bearer/pkg/commands/process/settings"
	"github.com/bearer/bearer/pkg/report/coverage"
)

// Discover searches directory for files to scan, skipping the ones specified by skip config and assigning timeout speficfied by timeout config
func Discover(repository *gitrepository.Repository, targetPath string, goclocResult *gocloc.Result, config settings.Config) (*flfiles.List, error) {
	ignore := ignore.New(targetPath, config)

	if !config.IgnoreGit {
		fileList, err := repository.ListFiles(ignore, goclocResult)
		if err != nil {
			log.Error().Msg("Git discovery failed")
			return nil, err
		}

		if fileList != nil {
			log.Debug().Msg("Files found from Git")
			return fileList, nil
		}

		log.Debug().Msg("No files found from Git")
	}

	var files []flfiles.File
	discovery := coverage.Discovery{
		Files:    []string{},
		Excluded: []coverage.Exclusion{},
		Errors:   []string{},
	}
	err := filepath.WalkDir(targetPath, func(filePath string, d fs.DirEntry, err error) error {
		if err != nil {
			discovery.Errors = append(discovery.Errors, fmt.Sprintf("%s: %v", filePath, err))
			return err
		}

		relativePath := strings.TrimPrefix(strings.TrimPrefix(filePath, targetPath), "/")
		if relativePath == "" {
			relativePath = "."
		}

		fileInfo, err := d.Info()
		if err != nil {
			discovery.Errors = append(discovery.Errors, fmt.Sprintf("%s: %v", relativePath, err))
			return err
		}

		if d.IsDir() {
			reason, decisionErr := ignore.Decision(targetPath, filePath, goclocResult, fileInfo)
			if decisionErr != nil {
				discovery.Errors = append(discovery.Errors, fmt.Sprintf("%s: %v", relativePath, decisionErr))
				return decisionErr
			}
			if reason != "" {
				discovery.Excluded = append(discovery.Excluded, coverage.Exclusion{Path: relativePath, Reason: reason})
				return filepath.SkipDir
			}

			return nil
		}

		reason, decisionErr := ignore.Decision(targetPath, filePath, goclocResult, fileInfo)
		if decisionErr != nil {
			discovery.Errors = append(discovery.Errors, fmt.Sprintf("%s: %v", relativePath, decisionErr))
			return decisionErr
		}
		if reason != "" {
			log.Debug().Msgf("skipping file due to file skip rules: %s", relativePath)
			discovery.Excluded = append(discovery.Excluded, coverage.Exclusion{Path: relativePath, Reason: reason})
			return nil
		}

		files = append(files, flfiles.File{
			FilePath: relativePath,
			Timeout:  timeout.Assign(fileInfo, config),
		})
		discovery.Files = append(discovery.Files, relativePath)

		return nil
	})

	discovery.Complete = err == nil
	return &flfiles.List{Files: files, Coverage: discovery}, err
}
