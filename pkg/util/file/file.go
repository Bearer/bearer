package file

import (
	"bufio"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	pathlib "path"

	"github.com/go-enry/go-enry/v2"
	"github.com/rs/zerolog/log"

	"github.com/bearer/curio/pkg/util/regex"

	ignore "github.com/sabhiram/go-gitignore"
)

const guessByteCount = 16 * 1024

var ignoredFilenames = []*regexp.Regexp{
	regexp.MustCompile(`(^|/)\.git/`),
	regexp.MustCompile(`(^|/)(?i:_*tests?_*)/`),
	regexp.MustCompile(`(^|/)specs?/`),
	regexp.MustCompile(`(^|/)testing/`),
	regexp.MustCompile(`(^|/|[_-])(spec|test)s?\.`),
	regexp.MustCompile(`(?i:unit[-_]?tests?)`),
	regexp.MustCompile(`(^|/)_*mocks?_*`),
	regexp.MustCompile(`(^|/)fixtures/`),
	regexp.MustCompile(`\.log$`),
	regexp.MustCompile(`(^|/)_*examples?_*(\.|/)`),
	regexp.MustCompile(`(^|/)_*samples?_*/`),
	regexp.MustCompile(`(^|/)node_modules/`),
	regexp.MustCompile(`(^|/)tmp/`),
	regexp.MustCompile(`\.min\.js$`),
	regexp.MustCompile(`\.map\.js$`),
}

type AllowDirFunction func(dir *Path) (bool, error)
type VisitFileFunction func(file *FileInfo) error

type Path struct {
	AbsolutePath string
	RelativePath string
}

type FileInfo struct {
	*Path
	os.FileInfo
	isBinary        bool
	isDocumentation bool
	IsGenerated     bool
	isImage         bool
	isTest          bool
	isGitIgnored    bool
	Extension       string
	Base            string
	IsConfiguration bool
	IsDotFile       bool
	IsVendor        bool
	Language        string
	LanguageType    enry.Type
}

func (path *Path) Join(elements ...string) *Path {
	return &Path{
		AbsolutePath: filepath.Join(append([]string{path.AbsolutePath}, elements...)...),
		RelativePath: filepath.Join(append([]string{path.RelativePath}, elements...)...),
	}
}

func (fileInfo *FileInfo) isGlobalIgnored() bool {
	return fileInfo.isBinary || fileInfo.isGitIgnored || fileInfo.isImage || fileInfo.isTest
}

func (fileInfo *FileInfo) LanguageTypeString() string {
	switch fileInfo.LanguageType {
	case enry.Data:
		return "data"
	case enry.Programming:
		return "programming"
	case enry.Markup:
		return "markup"
	case enry.Prose:
		return "prose"
	default:
		return ""
	}
}

func (path *Path) Exists() bool {
	_, err := os.Stat(path.AbsolutePath)
	if os.IsNotExist(err) {
		return false
	}

	if err != nil {
		log.Debug().Msgf("file '%s' doesn't exist", path.AbsolutePath)
		return false
	}

	return true
}

func IterateFilesList(rootDir string, files []string, allowDir AllowDirFunction, visitFile VisitFileFunction) error {
	gitIgnore := getGitIgnore(rootDir)

	rootDir, err := filepath.Abs(rootDir)
	if err != nil {
		return err
	}

	for _, fileToScan := range files {
		path := rootDir + "/" + fileToScan
		path, err = filepath.Abs(path)
		if err != nil {
			return err
		}

		fileStat, err := os.Stat(path)
		if _, ok := err.(*os.PathError); ok {
			log.Debug().Msgf("%s: skipping due to err: %s", path, err)
			return nil
		}

		if err != nil {
			return err
		}

		relativePath, err := filepath.Rel(rootDir, path)
		if err != nil {
			return err
		}

		if fileStat.IsDir() {
			relativePath += "/"
		}

		if regex.AnyMatch(ignoredFilenames, relativePath) {
			continue
		}

		// skip file based on parent directory
		var parentDirectories = []string{}
		possibleDirectories := strings.Split(relativePath, "/")
		startPath := rootDir
		for _, dir := range possibleDirectories {
			parentDirectories = append(parentDirectories, startPath)

			startPath = startPath + "/" + dir
		}
		isAllowed := true
		for _, parentDirectory := range parentDirectories {
			allowed, err := allowDir(&Path{
				AbsolutePath: parentDirectory,
				RelativePath: strings.TrimPrefix(parentDirectory, rootDir),
			})

			if err != nil {
				return err
			}

			if !allowed {
				isAllowed = false
				break
			}
		}
		if !isAllowed {
			continue
		}

		pathObject := &Path{
			AbsolutePath: path,
			RelativePath: relativePath,
		}

		if (fileStat.Mode() & (os.ModeSymlink | os.ModeSocket)) != 0 {
			continue
		}

		if fileStat.IsDir() {
			continue
		}

		fileInfo, err := newFileInfo(gitIgnore, pathObject, fileStat)
		if err != nil {
			log.Debug().Msgf("skipping %s due to err: %s", path, err)
		}

		if fileInfo.isGlobalIgnored() {
			log.Debug().Msgf("ignored due to file type %s", path)
			continue
		}

		if err := visitFile(fileInfo); err != nil {
			return err
		}
	}
	return nil
}

func newFileInfo(gitIgnore *ignore.GitIgnore, path *Path, file os.FileInfo) (*FileInfo, error) {
	guessBytes, err := guessBytes(path.AbsolutePath)
	if err != nil {
		return nil, err
	}

	language := enry.GetLanguage(pathlib.Base(path.AbsolutePath), guessBytes)
	languageType := enry.GetLanguageType(language)

	isGitIgnored := false
	if gitIgnore != nil {
		isGitIgnored = gitIgnore.MatchesPath(path.RelativePath)
	}

	fileInfo := &FileInfo{
		FileInfo:        file,
		Path:            path,
		isBinary:        enry.IsBinary(guessBytes),
		isDocumentation: enry.IsDocumentation(path.RelativePath),
		IsGenerated:     enry.IsGenerated(path.RelativePath, guessBytes),
		isImage:         enry.IsImage(path.RelativePath),
		isTest:          enry.IsTest(path.RelativePath),
		isGitIgnored:    isGitIgnored,
		Base:            filepath.Base(path.RelativePath),
		Extension:       strings.ToLower(filepath.Ext(path.RelativePath)),
		IsConfiguration: enry.IsConfiguration(path.RelativePath),
		IsDotFile:       enry.IsDotFile(path.RelativePath),
		IsVendor:        enry.IsVendor(path.RelativePath),
		Language:        language,
		LanguageType:    languageType,
	}

	return fileInfo, nil
}

func FileInfoFromPath(filePath string) (*FileInfo, error) {
	file, err := os.Stat(filePath)
	if err != nil {
		return nil, err
	}

	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return nil, err
	}

	path := &Path{
		AbsolutePath: filepath.Clean(absPath),
		RelativePath: filePath,
	}
	return newFileInfo(nil, path, file)
}

// returns the first guessByteCount of the specified file
func guessBytes(filename string) ([]byte, error) {
	result := make([]byte, guessByteCount)

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	n, err := file.Read(result)
	file.Close()
	if err != nil && err != io.EOF {
		return nil, err
	}

	return result[:n], nil
}

func getGitIgnore(rootPath string) *ignore.GitIgnore {
	gitIgnoreFilename := filepath.Join(rootPath, ".gitignore")

	var err error
	if _, err = os.Stat(gitIgnoreFilename); err == nil {
		var gitIgnore *ignore.GitIgnore
		gitIgnore, err = ignore.CompileIgnoreFile(gitIgnoreFilename)
		if err == nil {
			return gitIgnore
		}
	}

	// log.Debug().Msgf("no .gitignore, or error reading it: %s", err)
	return nil
}

func EnsureFileExists(filePath string) *os.File {
	// file exists
	if _, err := os.Stat(filePath); !os.IsNotExist(err) {
		err := os.Truncate(filePath, 0)
		if err != nil {
			log.Panic().Msgf("Failed to truncate existing file %s %e", filePath, err)
		}
	}

	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		log.Printf("Error creating output file %e", err)
	}

	return file
}

func ReadFileSingleLine(filePath string, lineNumber int) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineCounter := 1
	for scanner.Scan() {
		if lineCounter == lineNumber {
			return scanner.Text(), nil
		}
		lineCounter++
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return "", nil
}
