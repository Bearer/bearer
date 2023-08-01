package filelist

import (
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
)

// Discover searches directory for files to scan, skipping the ones specified by skip config and assigning timeout speficfied by timeout config
func Discover(repository *gitrepository.Repository, targetPath string, goclocResult *gocloc.Result, config settings.Config) (*flfiles.List, error) {
	ignore := ignore.New(targetPath, config)

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

	var files []flfiles.File
	err = filepath.WalkDir(targetPath, func(filePath string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		relativePath := strings.TrimPrefix(strings.TrimPrefix(filePath, targetPath), "/")
		if relativePath == "" {
			relativePath = "."
		}

		fileInfo, err := d.Info()
		if err != nil {
			log.Debug().Msgf("skipping due to info error %s: %s", relativePath, err)
			return nil
		}

		if d.IsDir() {
			if ignore.Ignore(targetPath, filePath, goclocResult, fileInfo) {
				return filepath.SkipDir
			}

			return nil
		}

		if ignore.Ignore(targetPath, filePath, goclocResult, fileInfo) {
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
