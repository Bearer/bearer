package dotnet

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/bearer/curio/pkg/detectors/types"
	"github.com/bearer/curio/pkg/parser"
	"github.com/bearer/curio/pkg/report"
	"github.com/bearer/curio/pkg/report/detectors"
	"github.com/bearer/curio/pkg/report/frameworks/dotnet"
	"github.com/bearer/curio/pkg/util/file"

	"github.com/smacker/go-tree-sitter/csharp"
)

var (
	fileProjectExt  = ".csproj"
	startupFileName = "Startup.cs"

	language = csharp.GetLanguage()

	query = parser.QueryMustCompile(language, `
	(invocation_expression
		(member_access_expression
			name: (generic_name (identifier) @name
			         (type_argument_list
						(identifier) @typeName)))
		(argument_list
			(argument
				(lambda_expression
					(invocation_expression
						(member_access_expression
							name: (identifier) @useDbMethodName)))))
		(#match? @name "^AddDbContext$"))`)
)

type detector struct{}

func New() types.Detector {
	return &detector{}
}

func (detector *detector) AcceptDir(dir *file.Path) (bool, error) {
	if isDotnetProject, err := isProject(dir.AbsolutePath); err != nil || !isDotnetProject {
		return false, err
	}

	return true, nil
}

func (detector *detector) ProcessFile(file *file.FileInfo, dir *file.Path, report report.Report) (bool, error) {
	if file.Base != startupFileName {
		return false, nil
	}

	tree, err := parser.ParseFile(file, file.Path, language)
	if err != nil {
		return false, err
	}
	defer tree.Close()

	err = tree.Query(query, func(captures parser.Captures) error {
		nameNode := captures["name"]
		typeNameNode := captures["typeName"]
		typeName := typeNameNode.Content()
		useDbMethodNameNode := captures["useDbMethodName"]
		useDbMethodName := useDbMethodNameNode.Content()

		report.AddFramework(detectors.DetectorDotnet, dotnet.TypeDatabase, dotnet.DBContext{
			TypeName:        typeName,
			UseDbMethodName: useDbMethodName,
		}, nameNode.Source(false))

		return nil
	})

	// Allow "csharp" detector to process file
	return false, err
}

func isProject(path string) (bool, error) {
	handleDir, err := isDir(path)
	if err != nil {
		return false, err
	}

	if handleDir {
		fileInfos, err := ioutil.ReadDir(path)
		if err != nil {
			return false, err
		}

		for _, file := range fileInfos {
			testMatch, errMatch := matchFilepath(file.Name())
			if errMatch != nil {
				return false, errMatch
			}

			if !testMatch {
				return true, nil
			}
		}
	}

	// we have a file and not a directory
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	testMatch, errMatch := matchFilepath(fileInfo.Name())
	if errMatch != nil {
		return false, errMatch
	}
	if !testMatch {
		return true, nil
	}

	return false, nil
}

func matchFilepath(filename string) (bool, error) {
	if filepath.Ext(filename) == fileProjectExt {
		return filepath.Match("*Test*", filename)
	}

	return false, nil
}

func isDir(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	return fileInfo.IsDir(), nil
}
