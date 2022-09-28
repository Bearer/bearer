package rails

import (
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"

	"github.com/bearer/curio/pkg/detectors/rails/cache"
	"github.com/bearer/curio/pkg/detectors/rails/personal_data"
	"github.com/bearer/curio/pkg/detectors/types"
	"github.com/bearer/curio/pkg/parser/nodeid"
	"github.com/bearer/curio/pkg/report"
	"github.com/bearer/curio/pkg/report/detectors"
	"github.com/bearer/curio/pkg/report/frameworks/rails"
	"github.com/bearer/curio/pkg/report/source"
	"github.com/bearer/curio/pkg/util/file"
	"github.com/bearer/curio/pkg/util/maputil"
	"github.com/bearer/curio/pkg/util/pointers"
)

var (
	databasePath           = filepath.Join("config", "database.yml")
	storagePath            = filepath.Join("config", "storage.yml")
	applicationConfigPath  = filepath.Join("config", "application.rb")
	seedsPath              = filepath.Join("db", "seeds.rb")
	productionConfigPath   = filepath.Join("config", "environments", "production.rb")
	rubyDatabaseSchemaPath = filepath.Join("db", "schema.rb")
)

type storageUploadInfo struct {
	ServerSideEncryption string `yaml:"server_side_encryption"`
}

type storageConfig struct {
	Service string             `yaml:"service"`
	Upload  *storageUploadInfo `yaml:"upload"`
}

type detector struct {
	idGenerator nodeid.Generator
}

func New(idGenerator nodeid.Generator) types.Detector {
	return &detector{
		idGenerator: idGenerator,
	}
}

func (detector *detector) AcceptDir(dir *file.Path) (bool, error) {
	return dir.Join("bin", "rails").Exists() || dir.Join("script", "rails").Exists(), nil
}

func (detector *detector) ProcessFile(file *file.FileInfo, dir *file.Path, report report.Report) (bool, error) {
	projectPath := strings.TrimPrefix("/"+file.RelativePath, dir.RelativePath)
	projectPath = strings.TrimPrefix(projectPath, "/")

	if projectPath == seedsPath {
		return true, nil
	}

	switch projectPath {
	case databasePath:
		if err := extractDatabases(file, report); err != nil {
			return false, err
		}
	case storagePath:
		if err := extractStorage(file, report); err != nil {
			return false, err
		}
	case applicationConfigPath, productionConfigPath:
		if err := cache.ExtractCaches(file, report); err != nil {
			return false, err
		}

		// Allow "ruby" detector to process file
		return false, nil
	case rubyDatabaseSchemaPath:
		if err := personal_data.ExtractFromDatabaseSchema(detector.idGenerator, file, report); err != nil {
			return false, err
		}
	default:
		return false, nil
	}

	return true, nil
}

func extractDatabases(file *file.FileInfo, report report.Report) error {
	bytes, err := ioutil.ReadFile(file.AbsolutePath)
	if len(bytes) == 0 || err != nil {
		return err
	}

	config := make(map[string]yaml.Node)
	if err := yaml.Unmarshal(bytes, &config); err != nil {
		return err
	}

	productionNode, exists := config["production"]
	if !exists {
		return err
	}

	database, source, err := getAdapterDatabase(file, productionNode, "")
	if err != nil {
		return err
	}

	if database != nil {
		report.AddFramework(detectors.DetectorRails, rails.TypeDatabase, database, *source)
	}

	productionConfig := make(map[string]yaml.Node)
	if err := productionNode.Decode(&productionConfig); err != nil {
		return err
	}

	for _, name := range maputil.SortedStringKeys(productionConfig) {
		databaseNode := productionConfig[name]
		if databaseNode.Kind == yaml.MappingNode {
			database, source, err := getAdapterDatabase(file, databaseNode, name)
			if err != nil {
				return err
			}
			if database != nil {
				report.AddFramework(detectors.DetectorRails, rails.TypeDatabase, database, *source)
			}
		}
	}

	return nil
}

func getAdapterDatabase(file *file.FileInfo, node yaml.Node, name string) (*rails.Database, *source.Source, error) {
	config := make(map[string]yaml.Node)
	if err := node.Decode(&config); err != nil {
		return nil, nil, err
	}

	if adapterNode, exists := config["adapter"]; exists {
		var adapter string
		if adapterNode.Decode(&adapter) == nil {
			return &rails.Database{
					Name:    name,
					Adapter: adapter,
				}, &source.Source{
					Language:     file.Language,
					LanguageType: file.LanguageTypeString(),
					Filename:     file.RelativePath,
					LineNumber:   &node.Line,
				}, nil
		}
	}

	return nil, nil, nil
}

func extractStorage(file *file.FileInfo, report report.Report) error {
	bytes, err := ioutil.ReadFile(file.AbsolutePath)
	if len(bytes) == 0 || err != nil {
		return err
	}

	config := make(map[string]yaml.Node)
	if err := yaml.Unmarshal(bytes, &config); err != nil {
		return err
	}

	for _, name := range maputil.SortedStringKeys(config) {
		storageNode := config[name]
		var storageConfig storageConfig
		if err := storageNode.Decode(&storageConfig); err != nil {
			log.Error().Msgf("failed to decode ActiveStorage config for '%s': %s", name, err)
		}

		encryption := ""
		if storageConfig.Upload != nil {
			encryption = strings.TrimSpace(storageConfig.Upload.ServerSideEncryption)
		}

		report.AddFramework(detectors.DetectorRails,
			rails.TypeStorage,
			rails.Storage{
				Name:       name,
				Service:    storageConfig.Service,
				Encryption: encryption,
			},
			source.Source{
				Language:     file.Language,
				LanguageType: file.LanguageTypeString(),
				Filename:     file.RelativePath,
				LineNumber:   pointers.Int(storageNode.Line - 1),
			})
	}

	return nil
}
