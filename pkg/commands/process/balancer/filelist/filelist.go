package filelist

import (
	"io/fs"
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

	ignore, err := fileignore.New(config)
	if err != nil {
		return nil, err
	}

	err = filepath.WalkDir(projectPath, func(filePath string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		relativePath := strings.TrimPrefix(filePath, projectPath)

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
