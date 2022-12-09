package filelist

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/bearer/curio/pkg/commands/process/balancer/fileignore"
	"github.com/bearer/curio/pkg/commands/process/balancer/timeout"
	"github.com/bearer/curio/pkg/commands/process/settings"
	"github.com/bearer/curio/pkg/commands/process/worker/work"
	"github.com/rs/zerolog/log"
)

// Discover searches directory for files to scan, skipping the ones specified by skip config and assigning timeout speficfied by timeout config
func Discover(projectPath string, config settings.Config) ([]work.File, error) {
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

	err := filepath.WalkDir(projectPath, func(filePath string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			if ignore.Ignore(projectPath, filePath, d) {
				return filepath.SkipDir
			}

			return nil
		}

		relativePath := strings.TrimPrefix(filePath, projectPath)
		relativePath = "/" + relativePath

		if ignore.Ignore(projectPath, filePath, d) {
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
