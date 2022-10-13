package fileignore

import (
	"io/fs"
	"strings"

	"github.com/bearer/curio/pkg/commands/process/settings"
	"github.com/monochromegane/go-gitignore"
)

type FileIgnore struct {
	ignorer    gitignore.IgnoreMatcher
	hasIgnorer bool

	config settings.Config
}

func New(projectPath string, config settings.Config) (*FileIgnore, error) {
	var ignorer gitignore.IgnoreMatcher
	var err error

	hasIgnorer := config.Scan.SkipConfig != ""

	if hasIgnorer {
		ignorer, err = gitignore.NewGitIgnore(config.Scan.SkipConfig, ".")
		if err != nil {
			return nil, err
		}
	}

	return &FileIgnore{
		ignorer:    ignorer,
		hasIgnorer: hasIgnorer,

		config: config,
	}, nil
}

func (fileignore *FileIgnore) Ignore(projectPath string, filePath string, d fs.DirEntry) bool {
	relativePath := strings.TrimPrefix(filePath, projectPath)

	if fileignore.hasIgnorer {
		trimmedPath := strings.TrimPrefix(relativePath, "/")
		if fileignore.ignorer.Match(trimmedPath, false) {
			return true
		}
	}

	fileInfo, err := d.Info()
	if err != nil {
		return true
	}

	if fileInfo.Size() > int64(fileignore.config.Worker.FileSizeMaximum) {
		return true
	}

	return false
}
