package fileignore

import (
	"bytes"
	"io/fs"
	"strings"

	"github.com/bearer/curio/pkg/commands/process/settings"
	"github.com/monochromegane/go-gitignore"
)

type FileIgnore struct {
	ignorer gitignore.IgnoreMatcher

	config settings.Config
}

func New(projectPath string, config settings.Config) *FileIgnore {
	return &FileIgnore{
		ignorer: ignorerFromStrings(config.Scan.Skip),

		config: config,
	}
}

func (fileignore *FileIgnore) Ignore(projectPath string, filePath string, d fs.DirEntry) bool {
	relativePath := strings.TrimPrefix(filePath, projectPath)
	trimmedPath := strings.TrimPrefix(relativePath, "/")

	fileInfo, err := d.Info()
	if err != nil {
		return true
	}

	if fileignore.ignorer.Match(trimmedPath, fileInfo.IsDir()) {
		return true
	}

	if !fileInfo.IsDir() {
		if fileInfo.Size() > int64(fileignore.config.Worker.FileSizeMaximum) {
			return true
		}
	}

	return false
}

func ignorerFromStrings(paths []string) gitignore.IgnoreMatcher {
	buffer := bytes.NewBuffer([]byte{})
	for _, path := range paths {
		buffer.Write([]byte(path + "\n"))
	}

	return gitignore.NewGitIgnoreFromReader(".", buffer)
}
