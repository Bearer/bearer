package beego

import (
	"strings"

	"github.com/smacker/go-tree-sitter/golang"

	"github.com/bearer/bearer/internal/detectors/types"
	"github.com/bearer/bearer/internal/parser"
	"github.com/bearer/bearer/internal/parser/golang_util"
	"github.com/bearer/bearer/internal/report"
	"github.com/bearer/bearer/internal/report/detectors"
	"github.com/bearer/bearer/internal/report/frameworks/beego"
	"github.com/bearer/bearer/internal/util/file"
)

const ormPackage = "github.com/beego/beego/v2/client/orm"

var (
	language = golang.GetLanguage()

	ormDatabaseQuery = parser.QueryMustCompile(language, `
		(call_expression
			function: (_) @function
			arguments: (argument_list . (_) @name . (_) @driver))
	`)

	ormDriverQuery = parser.QueryMustCompile(language, `
		(call_expression
			function: (_) @function
			arguments: (argument_list . (_) @name . (selector_expression) @constant))
	`)
)

type detector struct{}

type fileInfo struct {
	imports       map[string]string
	ormImportName string
	tree          *parser.Tree
}

type databaseDriver struct {
	name          string
	targetPackage string
	typeConstant  string
}

func New() types.Detector {
	return &detector{}
}

func (detector *detector) AcceptDir(dir *file.Path) (bool, error) {
	if isGoProject := dir.Join("go.mod").Exists(); !isGoProject {
		return false, nil
	}

	return true, nil
}

func (detector *detector) ProcessFile(file *file.FileInfo, dir *file.Path, report report.Report) (bool, error) {
	if file.Extension != ".go" {
		return false, nil
	}

	fileInfo, err := loadFile(file)
	if err != nil {
		return false, err
	}
	defer fileInfo.tree.Close()

	drivers, err := extractDatabaseDrivers(fileInfo, report)
	if err != nil {
		return false, err
	}

	if err := extractDatabases(drivers, fileInfo, report); err != nil {
		return false, err
	}

	// Allow "golang" detector to process file
	return false, nil
}

func loadFile(file *file.FileInfo) (*fileInfo, error) {
	tree, err := parser.ParseFile(file, file.Path, language)
	if err != nil {
		return nil, err
	}

	imports, err := golang_util.GetImports(tree)
	if err != nil {
		tree.Close()
		return nil, err
	}

	ormImportName := ""
	for importName, packagePath := range imports {
		if packagePath == ormPackage {
			ormImportName = "orm"
			if importName != "" {
				ormImportName = importName
			}

			break
		}
	}

	return &fileInfo{imports: imports, ormImportName: ormImportName, tree: tree}, nil
}

func extractDatabases(drivers []*databaseDriver, fileInfo *fileInfo, report report.Report) error {
	expectedFunction := fileInfo.ormImportName + ".RegisterDataBase"

	return fileInfo.tree.Query(ormDatabaseQuery, func(captures parser.Captures) error {
		functionNode := captures["function"]
		function := functionNode.Content()
		name := stripQuotes(captures["name"].Content())
		targetdriver := stripQuotes(captures["driver"].Content())

		driverPackage := ""
		driverTypeConstant := ""

		if function == expectedFunction {
			for _, driver := range drivers {
				if driver.name == targetdriver {
					driverPackage = driver.targetPackage
					driverTypeConstant = driver.typeConstant
					break
				}
			}

			report.AddFramework(detectors.DetectorBeego, beego.TypeDatabase, beego.Database{
				Name:         name,
				Package:      driverPackage,
				DriverName:   targetdriver,
				TypeConstant: driverTypeConstant,
			}, functionNode.Source(true))

		}

		return nil
	})
}

func extractDatabaseDrivers(fileInfo *fileInfo, report report.Report) (drivers []*databaseDriver, err error) {
	expectedFunction := fileInfo.ormImportName + ".RegisterDriver"

	err = fileInfo.tree.Query(ormDriverQuery, func(captures parser.Captures) error {
		functionNode := captures["function"]
		function := functionNode.Content()

		if function == expectedFunction {
			name := stripQuotes(captures["name"].Content())
			fullConstant := captures["constant"].Content()

			typeConstant := fullConstant
			packagePath := ""
			splitConstant := strings.Split(fullConstant, ".")
			if len(splitConstant) == 2 {
				packagePath = fileInfo.imports[splitConstant[0]]
				typeConstant = splitConstant[1]
			}

			drivers = append(drivers, &databaseDriver{
				name:          name,
				targetPackage: packagePath,
				typeConstant:  typeConstant,
			})
		}

		return nil
	})

	return drivers, err
}

func stripQuotes(value string) string {
	return strings.Trim(value, "\"`")
}
