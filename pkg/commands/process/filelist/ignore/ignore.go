package ignore

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/bearer/bearer/pkg/commands/process/settings"
	"github.com/hhatto/gocloc"
	"github.com/rs/zerolog/log"
	ignore "github.com/sabhiram/go-gitignore"
)

type FileIgnore struct {
	ignorer *ignore.GitIgnore

	config settings.Config
}

func New(projectPath string, config settings.Config) *FileIgnore {
	return &FileIgnore{
		ignorer: ignorerFromStrings(config.Scan.SkipPath),

		config: config,
	}
}

func (fileignore *FileIgnore) Ignore(
	projectPath string,
	filePath string,
	goclocResult *gocloc.Result,
	fileInfo fs.FileInfo,
) bool {
	relativePath := strings.TrimPrefix(filePath, projectPath)
	trimmedPath := strings.TrimPrefix(relativePath, "/")

	symlink, _ := isSymlink(projectPath + relativePath)
	if symlink {
		log.Debug().Msgf("skipping symlink: %s %s", projectPath, relativePath)
		return true
	}

	if fileignore.ignorer.MatchesPath(trimmedPath) {
		log.Debug().Msgf("file ignore match err: %s %s", projectPath, relativePath)
		return true
	}

	if !fileInfo.IsDir() {
		if fileInfo.Size() > int64(fileignore.config.Worker.FileSizeMaximum) {
			log.Debug().Msgf("skipping file due to size: %s %s", projectPath, relativePath)
			return true
		}
		if isMinified(fmt.Sprintf("%s%s", projectPath, filePath), fileInfo.Size(), goclocResult) {
			log.Debug().Msgf("skipping file (suspected minified JS): %s%s", projectPath, filePath)
			return true
		}
	}

	dirTrimmedPath := filepath.Dir(trimmedPath)
	dirPath := filepath.Join(projectPath, dirTrimmedPath)

	// No parent directory, allow
	if dirTrimmedPath == "." || dirTrimmedPath[len(dirTrimmedPath)-1] == filepath.Separator {
		return false
	}

	dirInfo, err := os.Stat(dirPath)
	if err != nil {
		log.Debug().Msgf("error getting dir stat %s: %s", dirPath, err)
		return true
	}

	return fileignore.Ignore(
		projectPath,
		dirPath,
		goclocResult,
		dirInfo,
	)
}

func isMinified(fullPath string, size int64, goclocResult *gocloc.Result) bool {
	if strings.HasSuffix(fullPath, ".min.js") {
		return true
	}

	if strings.HasSuffix(fullPath, "-min.js") {
		return true
	}

	if strings.HasSuffix(fullPath, ".js") {
		goclocFileResult := goclocResult.Files[fullPath]

		if goclocFileResult == nil {
			// couldn't find file
			return false
		}

		if goclocFileResult.Blanks == 0 && goclocFileResult.Comments == 0 && size > int64(5000) {
			// > 5KB JS file with no blank lines or comments -> assume minified
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

func ignorerFromStrings(paths []string) *ignore.GitIgnore {
	return ignore.CompileIgnoreLines(paths...)
}
