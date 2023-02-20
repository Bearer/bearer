package simple

import (
	"os"
	"regexp"

	"github.com/bearer/bearer/pkg/detectors/types"
	"github.com/bearer/bearer/pkg/parser/interfaces"
	"github.com/bearer/bearer/pkg/report"
	"github.com/bearer/bearer/pkg/report/detectors"
	interfacestype "github.com/bearer/bearer/pkg/report/interfaces"
	"github.com/bearer/bearer/pkg/report/source"
	"github.com/bearer/bearer/pkg/report/values"
	"github.com/bearer/bearer/pkg/util/file"
	"github.com/bearer/bearer/pkg/util/linescanner"
	"github.com/bearer/bearer/pkg/util/pointers"

	"github.com/go-enry/go-enry/v2"
)

var (

	// Match lines for block comments. eg
	//  /* ignore me
	//   * ignore me
	//   ignore me */
	blockCommentPattern = regexp.MustCompile(`^\s*/?\*|\*/\s*$`)

	// Looks for a prefix string followed by a URL candidate
	lineURLPattern = regexp.MustCompile(`(.*?)(https?://[a-zA-Z0-9?/\-._~%=:+]+)`)

	// Match comments like:
	//  # ignore me
	//  // ignore me
	lineCommentPattern = regexp.MustCompile(`#|//|;`)
)

type detector struct {
}

func New() types.Detector {
	return &detector{}
}

func (detector *detector) AcceptDir(dir *file.Path) (bool, error) {
	return true, nil
}

func (detector *detector) ProcessFile(fileInfo *file.FileInfo, dir *file.Path, report report.Report) (bool, error) {
	if fileInfo.Language == "CSS" || fileInfo.Language == "SCSS" {
		return false, nil
	}

	if fileInfo.LanguageType != enry.Programming && fileInfo.LanguageType != enry.Markup &&
		fileInfo.Language != "SQL" && fileInfo.Language != "GraphQL" {

		return false, nil
	}

	file, err := os.Open(fileInfo.AbsolutePath)
	if err != nil {
		return false, err
	}
	defer file.Close()

	scanner := linescanner.New(file)
	for scanner.Scan() {
		line := scanner.Text()

		if blockCommentPattern.MatchString(line) {
			continue
		}

		extractURLs(fileInfo, line, scanner.LineNumber(), report)
	}

	return true, scanner.Err()
}

func extractURLs(fileInfo *file.FileInfo, line string, lineNumber int, report report.Report) {
	globalOffset := 0

	for {
		match := lineURLPattern.FindStringSubmatchIndex(line)
		if match == nil {
			break
		}

		startOffset := match[4]
		endOffset := match[5]

		prefix := line[:startOffset]
		if lineCommentPattern.MatchString(prefix) {
			break
		}

		url := line[startOffset:endOffset]

		value := values.New()
		value.AppendString(url)

		if interfaceType, isInterface := interfaces.GetType(value, false); isInterface {

			report.AddInterface(detectors.DetectorSimple, interfacestype.Interface{
				Value: value,
				Type:  interfaceType,
			}, source.Source{
				Filename:     fileInfo.RelativePath,
				Language:     fileInfo.Language,
				LanguageType: fileInfo.LanguageTypeString(),
				LineNumber:   &lineNumber,
				ColumnNumber: pointers.Int(globalOffset + startOffset),
				Text:         &url,
			})
		}

		line = line[endOffset:]
		globalOffset += endOffset
	}
}
