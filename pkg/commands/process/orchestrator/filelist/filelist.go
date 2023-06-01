package filelist

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/bearer/bearer/pkg/commands/process/orchestrator/fileignore"
	"github.com/bearer/bearer/pkg/commands/process/orchestrator/timeout"
	"github.com/bearer/bearer/pkg/commands/process/settings"
	"github.com/bearer/bearer/pkg/commands/process/worker/work"
	"github.com/bearer/bearer/pkg/util/gitutil"
	"github.com/hhatto/gocloc"
	"github.com/rs/zerolog/log"
)

// Discover searches directory for files to scan, skipping the ones specified by skip config and assigning timeout speficfied by timeout config
func Discover(projectPath string, goclocResult *gocloc.Result, config settings.Config) ([]work.File, error) {
	var files []work.File

	haveDir, statErr := isDir(projectPath)
	if statErr != nil {
		return files, statErr
	}

	if haveDir {
		projectPath = strings.TrimPrefix(projectPath, "./")
		projectPath = strings.TrimSuffix(projectPath, "/")
		if projectPath != "." {
			projectPath += "/"
		}
	}

	ignore := fileignore.New(projectPath, config)

	pathsFromGit, err := gitutil.DiscoverFromGit(projectPath)
	if err != nil {
		log.Error().Msg("Git discovery failed")
		return nil, err
	}

	if len(pathsFromGit) != 0 {
		log.Debug().Msg("Files found from Git")

		for _, pathFromGit := range pathsFromGit {
			fullPath := projectPath + "/" + pathFromGit

			// Check if the file path itself should be ignored
			fileInfo, err := os.Stat(fullPath)
			if err != nil {
				log.Debug().Msgf("Skipping directory: %s, %s", fullPath, err.Error())
				continue
			}
			fileEntry := fs.FileInfoToDirEntry(fileInfo)
			if ignore.Ignore(projectPath, pathFromGit, goclocResult, fileEntry) {
				log.Debug().Msgf("Skipping file: %s", pathFromGit)
				continue
			}

			// Check if the parent directory should be ignored
			dirFullPath := filepath.Dir(fullPath)
			dirInfo, err := os.Stat(dirFullPath)
			if err != nil {
				return nil, err
			}
			dirEntry := fs.FileInfoToDirEntry(dirInfo)
			dirRelativePath := filepath.Dir(pathFromGit)
			if ignore.Ignore(projectPath, dirRelativePath, goclocResult, dirEntry) {
				log.Debug().Msgf("Skipping parent directory: %s", dirRelativePath)
				continue
			}

			relativePath := strings.TrimPrefix(pathFromGit, projectPath)
			relativePath = "/" + relativePath

			file := work.File{
				FilePath: relativePath,
				Timeout:  timeout.Assign(fileEntry, config),
			}

			files = append(files, file)
		}

		return files, nil
	}

	log.Debug().Msg("No files found from Git")

	err = filepath.WalkDir(projectPath, func(filePath string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			if ignore.Ignore(projectPath, filePath, goclocResult, d) {
				return filepath.SkipDir
			}

			return nil
		}

		relativePath := strings.TrimPrefix(filePath, projectPath)
		relativePath = "/" + relativePath

		if ignore.Ignore(projectPath, filePath, goclocResult, d) {
			log.Debug().Msgf("skipping file due to file skip rules: %s", relativePath)

			return nil
		}

		file := work.File{
			FilePath: relativePath,
			Timeout:  timeout.Assign(d, config),
		}

		files = append(files, file)

		return nil
	})

	return files, err
}

func isDir(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	return fileInfo.IsDir(), nil
}
