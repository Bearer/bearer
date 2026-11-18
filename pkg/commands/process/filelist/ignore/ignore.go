package ignore

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/bearer/bearer/pkg/commands/process/settings"
	fileutil "github.com/bearer/bearer/pkg/util/file"
	"github.com/hhatto/gocloc"
	"github.com/rs/zerolog/log"
	ignore "github.com/sabhiram/go-gitignore"
)

type FileIgnore struct {
	ignorer *ignore.GitIgnore
	config  settings.Config
}

func New(projectPath string, config settings.Config) *FileIgnore {
	return &FileIgnore{
		ignorer: ignorerFromStrings(config.Scan.SkipPath),
		config:  config,
	}
}

func (fileignore *FileIgnore) Ignore(
	projectPath string,
	filePath string,
	goclocResult *gocloc.Result,
	fileInfo fs.FileInfo,
) bool {
	reason, err := fileignore.Decision(projectPath, filePath, goclocResult, fileInfo)
	return err != nil || reason != ""
}

// Decision returns a stable exclusion reason. A non-nil error is a discovery
// failure, not an exclusion, and callers must record the scan as incomplete.
func (fileignore *FileIgnore) Decision(
	projectPath string,
	filePath string,
	goclocResult *gocloc.Result,
	fileInfo fs.FileInfo,
) (string, error) {
	relativePath := strings.TrimPrefix(filePath, projectPath)
	trimmedPath := strings.TrimPrefix(relativePath, "/")

	symlink, err := isSymlink(filePath)
	if err != nil {
		return "", err
	}
	if symlink {
		log.Debug().Msgf("skipping symlink: %s %s", projectPath, relativePath)
		return "symlink", nil
	}

	if fileignore.ignorer.MatchesPath(trimmedPath) {
		log.Debug().Msgf("file ignore match: %s %s", projectPath, relativePath)
		return "skip-path", nil
	}

	if !fileInfo.IsDir() {
		maxSize := fileignore.config.Worker.FileSizeMaximum
		if maxSize > 0 && fileInfo.Size() > int64(maxSize) {
			log.Debug().Msgf("skipping file due to size: %s %s", projectPath, relativePath)
			return "file-size", nil
		}
		if isMinified(filePath, fileInfo.Size(), goclocResult) {
			log.Debug().Msgf("skipping file (suspected minified JS): %s %s", projectPath, relativePath)
			return "minified-js", nil
		}
		if fileignore.skipByLanguage(filePath, goclocResult) {
			log.Debug().Msgf("skipping file (language not in --language filter): %s %s", projectPath, relativePath)
			return "language-filter", nil
		}
	}

	dirTrimmedPath := filepath.Dir(trimmedPath)
	dirPath := filepath.Join(projectPath, dirTrimmedPath)

	// No parent directory, allow
	if dirTrimmedPath == "." || dirTrimmedPath[len(dirTrimmedPath)-1] == filepath.Separator {
		return fileignore.configuredDecision(projectPath, trimmedPath)
	}

	dirInfo, err := os.Stat(dirPath)
	if err != nil {
		log.Debug().Msgf("error getting dir stat %s: %s", dirPath, err)
		return "", err
	}

	reason, err := fileignore.Decision(
		projectPath,
		dirPath,
		goclocResult,
		dirInfo,
	)
	if err != nil || reason != "" {
		return reason, err
	}

	return fileignore.configuredDecision(projectPath, trimmedPath)
}

func (fileignore *FileIgnore) configuredDecision(projectPath, filename string) (string, error) {
	// Diff discovery keeps its historical head-side decision here so the
	// caller can preserve rename/chunk metadata. The diff runner evaluates
	// configured eligibility independently in each checked-out tree.
	if fileignore.config.Scan.Diff {
		return "", nil
	}

	return fileutil.ConfiguredExclusion(
		projectPath,
		filename,
		fileignore.config.Scan.SkipTest,
		fileignore.config.Scan.SkipGitIgnore,
	)
}

// skipByLanguage reports whether a file should be ignored because of the
// --language restriction. When no languages are requested it never skips.
// Otherwise the gocloc result has already been restricted to the requested
// languages upstream, so any file absent from it is not one of those languages
// (a different language or unrecognised) and is skipped.
func (fileignore *FileIgnore) skipByLanguage(fullPath string, goclocResult *gocloc.Result) bool {
	if len(fileignore.config.Scan.Language) == 0 {
		return false
	}

	if goclocResult == nil {
		return true
	}

	return goclocResult.Files[fullPath] == nil
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
