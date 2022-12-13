package fileignore

import (
	"bytes"
	"io/fs"
	"os"
	"strings"

	"github.com/bearer/curio/pkg/commands/process/settings"
	"github.com/monochromegane/go-gitignore"
	"github.com/rs/zerolog/log"
)

type FileIgnore struct {
	ignorer gitignore.IgnoreMatcher

	config settings.Config
}

func New(projectPath string, config settings.Config) *FileIgnore {
	return &FileIgnore{
		ignorer: ignorerFromStrings(config.Scan.SkipPath),

		config: config,
	}
}

func (fileignore *FileIgnore) Ignore(projectPath string, filePath string, d fs.DirEntry) bool {
	relativePath := strings.TrimPrefix(filePath, projectPath)
	trimmedPath := strings.TrimPrefix(relativePath, "/")

	fileInfo, err := d.Info()
	if err != nil {
		log.Error().Msgf("fileInfo err: %s %s", projectPath, relativePath)
		return true
	}

	symlink, _ := isSymlink(projectPath + relativePath)
	if symlink {
		log.Debug().Msgf("skipping symlink: %s %s", projectPath, relativePath)
		return true
	}

	if fileignore.ignorer.Match(trimmedPath, fileInfo.IsDir()) {
		log.Error().Msgf("file ignore match err: %s %s", projectPath, relativePath)
		return true
	}

	if !fileInfo.IsDir() {
		if fileInfo.Size() > int64(fileignore.config.Worker.FileSizeMaximum) {
			log.Debug().Msgf("skipping file due to size: %s %s", projectPath, relativePath)
			return true
		}
	}

	return false
}

func isSymlink(path string) (bool, error) {
	fileInfo, err := os.Lstat(path)
	if err != nil {
		return false, err
	}

	return fileInfo.Mode()&os.ModeSymlink == os.ModeSymlink, nil
}

func ignorerFromStrings(paths []string) gitignore.IgnoreMatcher {
	buffer := bytes.NewBuffer([]byte{})
	for _, path := range paths {
		buffer.Write([]byte(path + "\n"))
	}

	return gitignore.NewGitIgnoreFromReader(".", buffer)
}
