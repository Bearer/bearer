package filelist

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/hhatto/gocloc"
	"github.com/rs/zerolog/log"

	flfiles "github.com/bearer/bearer/pkg/commands/process/filelist/files"
	"github.com/bearer/bearer/pkg/commands/process/filelist/ignore"
	"github.com/bearer/bearer/pkg/commands/process/filelist/timeout"
	"github.com/bearer/bearer/pkg/commands/process/gitrepository"
	"github.com/bearer/bearer/pkg/commands/process/settings"
)

// Discover searches directory for files to scan, skipping the ones specified by skip config and assigning timeout speficfied by timeout config
func Discover(repository *gitrepository.Repository, projectPath string, goclocResult *gocloc.Result, config settings.Config) (*flfiles.List, error) {
	haveDir, statErr := isDir(projectPath)
	if statErr != nil {
		return nil, statErr
	}

	if haveDir {
		projectPath = strings.TrimPrefix(projectPath, "./")
		projectPath = strings.TrimSuffix(projectPath, "/")
		if projectPath != "." {
			projectPath += "/"
		}
	}

	ignore := ignore.New(projectPath, config)

	fileList, err := repository.ListFiles(ignore, goclocResult, projectPath)
	if err != nil {
		log.Error().Msg("Git discovery failed")
		return nil, err
	}

	if fileList != nil {
		log.Debug().Msg("Files found from Git")
		return fileList, nil
	}

	log.Debug().Msg("No files found from Git")

	var files []flfiles.File
	err = filepath.WalkDir(projectPath, func(filePath string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		relativePath := strings.TrimPrefix(filePath, projectPath)
		relativePath = "/" + relativePath

		fileInfo, err := d.Info()
		if err != nil {
			log.Debug().Msgf("skipping due to info error %s: %s", relativePath, err)
			return nil
		}

		if d.IsDir() {
			if ignore.Ignore(projectPath, filePath, goclocResult, fileInfo) {
				return filepath.SkipDir
			}

			return nil
		}

		if ignore.Ignore(projectPath, filePath, goclocResult, fileInfo) {
			log.Debug().Msgf("skipping file due to file skip rules: %s", relativePath)

			return nil
		}

		files = append(files, flfiles.File{
			FilePath: relativePath,
			Timeout:  timeout.Assign(fileInfo, config),
		})

		return nil
	})

	return &flfiles.List{Files: files}, err
}

func isDir(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	return fileInfo.IsDir(), nil
}
