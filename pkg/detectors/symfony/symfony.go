package symfony

import (
	"io/ioutil"
	"net/url"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/bearer/curio/pkg/detectors/types"
	"github.com/bearer/curio/pkg/report"
	"github.com/bearer/curio/pkg/report/detectors"
	"github.com/bearer/curio/pkg/report/frameworks/symfony"
	"github.com/bearer/curio/pkg/report/source"
	"github.com/bearer/curio/pkg/util/file"
	"github.com/bearer/curio/pkg/util/maputil"
	"github.com/bearer/curio/pkg/util/pointers"
)

var (
	appConfigPattern      = regexp.MustCompile(`app/config/config.ya?ml$`)
	doctrineConfigPattern = regexp.MustCompile(`config/packages(/.*)?/doctrine.ya?ml$`)
)

type appConfig struct {
	DoctrineConfig *doctrineConfig `yaml:"doctrine"`
}

type doctrineConfig struct {
	DBAccessLayer yaml.Node `yaml:"dbal"`
}

type databaseConfig struct {
	Driver string `yaml:"driver"`
	URL    string `yaml:"url"`
}

type detector struct{}

func New() types.Detector {
	return &detector{}
}

func (detector *detector) AcceptDir(dir *file.Path) (bool, error) {
	if isComposerProject := isComposer(dir); !isComposerProject {
		return false, nil
	}

	return true, nil
}

func (detector *detector) ProcessFile(file *file.FileInfo, dir *file.Path, report report.Report) (bool, error) {
	projectPath := strings.TrimPrefix("/"+file.RelativePath, dir.RelativePath)
	projectPath = strings.TrimPrefix(projectPath, "/")

	if appConfigPattern.MatchString(projectPath) || doctrineConfigPattern.MatchString(projectPath) {
		if err := extractDatabasesFromConfig(file, report); err != nil {
			return false, err
		}

		return true, nil
	}

	return false, nil
}

func extractDatabasesFromConfig(file *file.FileInfo, report report.Report) error {
	input, err := ioutil.ReadFile(file.AbsolutePath)
	if err != nil {
		return err
	}

	var appConfig appConfig
	if err := yaml.Unmarshal(input, &appConfig); err != nil {
		return err
	}

	if appConfig.DoctrineConfig == nil {
		return nil
	}

	dbAccessLayerNode := appConfig.DoctrineConfig.DBAccessLayer
	dbAccessLayer := make(map[string]yaml.Node)
	if err := dbAccessLayerNode.Decode(&dbAccessLayer); err != nil {
		return err
	}

	connectionsNode := dbAccessLayer["connections"]
	connections := make(map[string]yaml.Node)
	if connectionsNode.Kind == 0 {
		connections = map[string]yaml.Node{"": dbAccessLayerNode}
	} else {
		if err := connectionsNode.Decode(&connections); err != nil {
			return err
		}
	}

	for _, name := range maputil.SortedStringKeys(connections) {
		connectionNode := connections[name]
		var connection databaseConfig
		if err := connectionNode.Decode(&connection); err != nil {
			return err
		}

		driver := connection.Driver
		if driver == "" {
			connectionURL, err := url.Parse(connection.URL)
			if err == nil {
				driver = connectionURL.Scheme
			}
		}

		report.AddFramework(detectors.DetectorSymfony, symfony.TypeDatabase, symfony.Database{
			Name:   name,
			Driver: driver,
		}, source.Source{
			Language:     file.Language,
			LanguageType: file.LanguageTypeString(),
			Filename:     file.RelativePath,
			LineNumber:   pointers.Int(connectionNode.Line - 1),
		})
	}

	return nil
}

func isComposer(dir *file.Path) bool {
	return dir.Join("composer.json").Exists()
}
