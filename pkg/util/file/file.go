package file

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	pathlib "path"

	"github.com/go-enry/go-enry/v2"
	"github.com/rs/zerolog/log"

	"github.com/bearer/bearer/pkg/util/linescanner"
	"github.com/bearer/bearer/pkg/util/regex"

	ignore "github.com/sabhiram/go-gitignore"
)

const guessByteCount = 16 * 1024

type Line struct {
	Extract    string
	LineNumber int
	Strip      bool
}

var ignoreTestFiles = []*regexp.Regexp{
	regexp.MustCompile(`(^|/)(?i:_*tests?_*)/`),
	regexp.MustCompile(`(^|/)specs?/`),
	regexp.MustCompile(`(^|/|[_-])(spec|test)s?\.`),
	regexp.MustCompile(`(?i:unit[-_]?tests?)`),
}

// IgnoredFilenames is the built-in list of path patterns that are skipped
// during a scan. It is exported so it can be overridden or extended at
// startup, before scanning begins.
var IgnoredFilenames = []*regexp.Regexp{
	regexp.MustCompile(`(^|/)\.git/`),
	regexp.MustCompile(`(^|/)testing/`),
	regexp.MustCompile(`(^|/)_*mocks?_*`),
	regexp.MustCompile(`(^|/)fixtures/`),
	regexp.MustCompile(`\.log$`),
	regexp.MustCompile(`(^|/)_*examples?_*(\.|/)`),
	regexp.MustCompile(`(^|/)_*samples?_*/`),
	regexp.MustCompile(`(^|/)node_modules/`),
	regexp.MustCompile(`(^|/)tmp/`),
	regexp.MustCompile(`[.-]min\.js$`),
	regexp.MustCompile(`\.map\.js$`),
}

type AllowDirFunction func(dir *Path) (bool, error)
type VisitFileFunction func(file *FileInfo) error

type Path struct {
	AbsolutePath string
	RelativePath string
	scope        *InputScope
	input        *BoundInput
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
		scope:        path.scope,
	}
}

func (fileInfo *FileInfo) isGlobalIgnored() bool {
	return fileInfo.isBinary || fileInfo.isGitIgnored || fileInfo.isImage
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
	if path.scope != nil {
		return path.scope.contextExists(path.AbsolutePath)
	}
	_, err := os.Stat(path.AbsolutePath)
	if os.IsNotExist(err) {
		return false
	}

	if err != nil {
		return false
	}

	return true
}

// ConfiguredExclusion applies the file eligibility rules that are shared by
// discovery and processing. filename is relative to rootDir. A non-empty
// reason is an intentional exclusion; a non-nil error means that the input
// could not be inspected and must be treated as a discovery failure.
func ConfiguredExclusion(rootDir, filename string, skipTest, skipGitIgnore bool) (reason string, retErr error) {
	scope, err := OpenInputScope(rootDir)
	if err != nil {
		return "", err
	}
	defer func() { retErr = errors.Join(retErr, scope.Close()) }()
	info, err := scope.statRelative(filename)
	if err != nil {
		return "", err
	}
	relative := filepath.Clean(filename)
	if scope.single != "" {
		relative = scope.single
	}
	if info.IsDir() {
		return configuredPathExclusion(relative+"/", skipTest), nil
	}
	bound, err := scope.Acquire(filename, nil)
	if err != nil {
		return "", err
	}
	if reason := configuredPathExclusion(relative, skipTest); reason != "" {
		return reason, nil
	}
	var gitIgnore *ignore.GitIgnore
	if !skipGitIgnore {
		gitIgnore, err = getScopedGitIgnore(scope)
		if err != nil {
			return "", err
		}
	}
	pathObject := &Path{AbsolutePath: filepath.Join(scope.logicalRoot(), relative), RelativePath: relative, scope: scope, input: bound}
	fi, err := newFileInfo(gitIgnore, pathObject, bound.Info())
	if err != nil {
		return "", err
	}
	if fi.isGlobalIgnored() {
		return "file-type-or-gitignore", nil
	}
	if skipTest && fi.isTest {
		return "test-type", nil
	}
	return "", nil
}

func configuredPathExclusion(relativePath string, skipTest bool) string {
	if regex.AnyMatch(IgnoredFilenames, relativePath) {
		return "built-in-path"
	}
	if skipTest && regex.AnyMatch(ignoreTestFiles, relativePath) {
		return "test-name"
	}
	return ""
}

func configuredExclusion(rootDir, relativePath string, fileStat os.FileInfo, skipTest bool, gitIgnore *ignore.GitIgnore) (string, error) {
	if fileStat.IsDir() {
		relativePath += "/"
	}
	if reason := configuredPathExclusion(relativePath, skipTest); reason != "" {
		return reason, nil
	}
	if fileStat.IsDir() {
		return "", nil
	}
	if fileStat.Mode()&(os.ModeSymlink|os.ModeSocket) != 0 {
		return "", nil
	}

	pathObject := &Path{
		AbsolutePath: filepath.Join(rootDir, filepath.FromSlash(relativePath)),
		RelativePath: relativePath,
	}
	fileInfo, err := newFileInfo(gitIgnore, pathObject, fileStat)
	if err != nil {
		return "", err
	}
	if fileInfo.isGlobalIgnored() {
		return "file-type-or-gitignore", nil
	}
	if skipTest && fileInfo.isTest {
		return "test-type", nil
	}
	return "", nil
}

func IterateFilesList(
	rootDir string,
	files []string,
	skipTest bool,
	skipGitIgnore bool,
	allowDir AllowDirFunction,
	visitFile VisitFileFunction,
) error {
	return IterateFilesListWithExclusions(rootDir, files, skipTest, skipGitIgnore, allowDir, visitFile, nil)
}

// IterateFilesListWithExclusions observes the existing filtering decisions. It
// does not predict processing from discovery or change which filters apply.
func IterateFilesListWithExclusions(rootDir string, files []string, skipTest, skipGitIgnore bool,
	allowDir AllowDirFunction, visitFile VisitFileFunction, excluded func(string)) error {
	return IterateFilesListWithInput(rootDir, files, skipTest, skipGitIgnore, allowDir, visitFile, excluded, nil)
}

// IterateFilesListWithInput binds all consumer reads, including context reads,
// to one acquired root. The expectation is required by worker entrypoints.
func IterateFilesListWithInput(rootDir string, files []string, skipTest, skipGitIgnore bool,
	allowDir AllowDirFunction, visitFile VisitFileFunction, excluded func(string), expectation *InputExpectation, completed ...func(*InputScope) error) (retErr error) {
	scope, err := OpenInputScope(rootDir)
	if err != nil {
		return err
	}
	defer func() { retErr = errors.Join(retErr, scope.Close()) }()
	if expectation != nil && len(files) != 1 {
		return errors.New("expectation requires one input")
	}
	// A worker request binds the root as well as the selected bytes. Enforce
	// that binding before acquiring any context from this newly opened root.
	// Reuse the acquisition below rather than validating by a pathname check.
	var expectedInput *BoundInput
	if expectation != nil {
		expectedInput, err = scope.Acquire(files[0], expectation)
		if err != nil {
			return err
		}
	}
	var gitIgnore *ignore.GitIgnore
	if !skipGitIgnore {
		gitIgnore, err = getScopedGitIgnore(scope)
		if err != nil {
			return err
		}
	}
	skip := func(reason string) {
		if excluded != nil {
			excluded(reason)
		}
	}
	for _, relative := range files {
		bound := expectedInput
		if bound == nil {
			bound, err = scope.Acquire(relative, nil)
			if err != nil {
				return err
			}
		}
		eligibility, err := scope.selector(relative)
		if err != nil {
			return err
		}
		absolute := filepath.Join(scope.logicalRoot(), eligibility)
		if reason := configuredPathExclusion(eligibility, skipTest); reason != "" {
			skip(reason)
			if err = bound.Close(); err != nil {
				return err
			}
			continue
		}
		directories := []string{"."}
		dir := filepath.Dir(eligibility)
		if dir != "." {
			acc := ""
			for _, part := range strings.Split(dir, string(filepath.Separator)) {
				acc = filepath.Join(acc, part)
				directories = append(directories, acc)
			}
		}
		allowed := true
		for _, dir := range directories {
			ok, e := allowDir(&Path{AbsolutePath: filepath.Join(scope.logicalRoot(), dir), RelativePath: strings.TrimPrefix(filepath.Join(scope.logicalRoot(), dir), scope.logicalRoot()), scope: scope})
			// A bool-only detector predicate must not turn a scoped I/O fault
			// into an ordinary directory exclusion.
			if e = errors.Join(e, scope.Err()); e != nil {
				return e
			}
			if !ok {
				allowed = false
				break
			}
		}
		if !allowed {
			skip("directory-filter")
			if err = bound.Close(); err != nil {
				return err
			}
			continue
		}
		pathObject := &Path{AbsolutePath: absolute, RelativePath: eligibility, scope: scope, input: bound}
		info, err := newFileInfo(gitIgnore, pathObject, bound.Info())
		if err != nil {
			return err
		}
		reason := ""
		if info.isGlobalIgnored() {
			reason = "file-type-or-gitignore"
		} else if skipTest && info.isTest {
			reason = "test-type"
		}
		if reason != "" {
			skip(reason)
			if err = bound.Close(); err != nil {
				return err
			}
			continue
		}
		info.RelativePath = relative
		if scope.single != "" {
			info.RelativePath = "."
		}
		err = visitFile(info)
		err = errors.Join(err, bound.Close(), scope.Err())
		if err != nil {
			return err
		}
	}
	if err := scope.Validate(); err != nil {
		return err
	}
	for _, transfer := range completed {
		if transfer == nil {
			return errors.New("context handoff unavailable")
		}
		if err := transfer(scope); err != nil {
			return err
		}
	}
	return nil
}

func newFileInfo(gitIgnore *ignore.GitIgnore, path *Path, file os.FileInfo) (*FileInfo, error) {
	guessBytes, err := path.ReadBytes()
	if err != nil {
		return nil, err
	}

	if len(guessBytes) > guessByteCount {
		guessBytes = guessBytes[:guessByteCount]
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

// Discovery metadata is acquired safely and does not certify terminal processing.
func guessBytes(filename string) ([]byte, error) {
	data, err := (&Path{AbsolutePath: filename}).ReadBytes()
	if len(data) > guessByteCount {
		data = data[:guessByteCount]
	}
	return data, err
}
func getGitIgnore(rootPath string) (*ignore.GitIgnore, error) {
	scope, err := OpenInputScope(rootPath)
	if err != nil {
		return nil, err
	}
	result, err := getScopedGitIgnore(scope)
	return result, errors.Join(err, scope.Close())
}

func EnsureFileExists(filePath string) *os.File {
	// file exists
	if _, err := os.Stat(filePath); !os.IsNotExist(err) {
		err := os.Truncate(filePath, 0)
		if err != nil {
			log.Panic().Msgf("Failed to truncate existing file %s %s", filePath, err)
		}
	}

	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		log.Printf("Error creating output file %s", err)
	}

	return file
}

func ReadFileSinkLines(
	filePath string,
	sinkStartLine int,
	sinkEndLine int,
	sourceStartLine int,
	sourceEndLine int,
	buffer int,
) ([]Line, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return []Line{}, err
	}

	defer file.Close()

	return ReadSinkLines(file, sinkStartLine, sinkEndLine, sourceStartLine, sourceEndLine, buffer)
}

// ReadSinkLines consumes an already acquired reader without reopening a path.
func ReadSinkLines(reader io.Reader, sinkStartLine, sinkEndLine, sourceStartLine, sourceEndLine, buffer int) ([]Line, error) {
	scanner := linescanner.New(reader)
	lineCounter := 1
	var extract []Line

	for scanner.Scan() {
		if lineCounter >= sinkStartLine && lineCounter <= sinkEndLine {
			if lineCounter >= sourceStartLine+buffer && lineCounter <= sourceEndLine-buffer {
				if lineCounter == sourceStartLine+buffer {
					extract = append(extract, Line{
						Extract:    fmt.Sprintf("%s...omitted (buffer value %d)", strings.Repeat(" ", countLeadingSpaces(scanner.Text())), buffer),
						LineNumber: lineCounter,
						Strip:      true,
					})
				}
			} else {
				extract = append(extract, Line{
					Extract:    scanner.Text(),
					LineNumber: lineCounter,
					Strip:      false,
				})
			}
		}

		lineCounter++
	}

	if err := scanner.Err(); err != nil {
		return []Line{}, err
	}

	return extract, nil
}

func countLeadingSpaces(line string) int {
	return len(line) - len(strings.TrimLeft(line, " "))
}

func ReadFileSingleLine(filePath string, lineNumber int) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := linescanner.New(file)
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

func GetFullFilename(path string, filename string) string {
	path = strings.TrimSuffix(path, "/")
	filename = strings.TrimPrefix(filename, "/")

	if filename == "." {
		return path
	}

	if path == "" || path == "." {
		return filename
	}

	return path + "/" + filename
}

func IsDir(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}

	return fileInfo.IsDir()
}

func CanonicalPath(path string) (string, error) {
	absolutePath, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	resolvedPath, err := filepath.EvalSymlinks(absolutePath)
	if err != nil {
		return "", err
	}

	return filepath.Abs(resolvedPath)
}
