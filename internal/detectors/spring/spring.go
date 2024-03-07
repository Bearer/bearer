package spring

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"

	"github.com/bearer/bearer/internal/detectors/types"
	"github.com/bearer/bearer/internal/report"
	"github.com/bearer/bearer/internal/report/detectors"
	"github.com/bearer/bearer/internal/report/frameworks/spring"
	"github.com/bearer/bearer/internal/report/source"
	"github.com/bearer/bearer/internal/util/file"
	"github.com/bearer/bearer/internal/util/linescanner"
	"github.com/bearer/bearer/internal/util/pointers"
)

var (
	projectFiles = []string{
		"pom.xml",
		"build.gradle",
		"build.gradle.kts",
	}
)

type yamlConfig struct {
	Spring *yamlSpring `yaml:"spring"`
}

type yamlSpring struct {
	DataSource yaml.Node `yaml:"datasource"`
}

type yamlDataSource struct {
	URL          string `yaml:"url"`
	DriverClass1 string `yaml:"driver-class-name"`
	DriverClass2 string `yaml:"driverClassName"`
}

type detector struct{}

func New() types.Detector {
	return &detector{}
}

func (detector *detector) AcceptDir(dir *file.Path) (bool, error) {
	if !isJava(dir) {
		return false, nil
	}

	return true, nil
}

func (detector *detector) ProcessFile(file *file.FileInfo, dir *file.Path, report report.Report) (bool, error) {

	if file.Base == "application.properties" {
		if err := extractDataStoresFromProperties(file, report); err != nil {
			return false, fmt.Errorf("failed to load spring properties: %s", err)
		}

		return true, nil
	}

	if file.Base == "application.yml" {
		if err := extractDataStoresFromYAML(file, report); err != nil {
			return false, fmt.Errorf("failed to load spring YAML: %s", err)
		}

		return true, nil
	}

	return false, nil
}

type Property struct {
	value      string
	lineNumber int
	text       string
}

type AppConfigProperties map[string]Property

func ReadPropertiesFile(file *file.FileInfo) (AppConfigProperties, error) {
	config := AppConfigProperties{}

	fileBytes, err := os.ReadFile(file.AbsolutePath)
	if err != nil {
		log.Error().Msgf("There was an error while opening the file: %s", err.Error())
		return nil, err
	}

	scanner := linescanner.New(bytes.NewBuffer(fileBytes))
	for scanner.Scan() {
		line := scanner.Text()
		if equal := strings.Index(line, "="); equal >= 0 {
			if key := strings.TrimSpace(line[:equal]); len(key) > 0 {
				value := ""
				if len(line) > equal {
					value = strings.TrimSpace(line[equal+1:])
				}
				config[key] = Property{
					value:      value,
					lineNumber: scanner.LineNumber(),
					text:       scanner.Text(),
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Error().Msgf("There was an error processing properties: %s", err.Error())
		return nil, err
	}

	return config, nil
}

func extractDataStoresFromProperties(file *file.FileInfo, report report.Report) error {
	config, err := ReadPropertiesFile(file)
	if err != nil {
		return err
	}

	url := config["spring.datasource.url"]
	driverClass1 := config["spring.datasource.driver-class-name"]
	driverClass2 := config["spring.datasource.driverClassName"]

	property := getProperty(driverClass1, driverClass2, url)

	if property.value == "" {
		return nil
	}

	report.AddFramework(detectors.DetectorSpring, spring.TypeDatabase, spring.DataStore{
		Driver: getDriver(driverClass1.value, driverClass2.value, url.value),
	}, source.Source{
		Language:        file.Language,
		LanguageType:    file.LanguageTypeString(),
		Filename:        file.RelativePath,
		StartLineNumber: &property.lineNumber,
	})

	return nil
}

func getProperty(driverClass1 Property, driverClass2 Property, url Property) Property {
	if driverClass1.value != "" {
		return driverClass1
	}

	if driverClass2.value != "" {
		return driverClass2
	}

	return url
}

func extractDataStoresFromYAML(file *file.FileInfo, report report.Report) error {
	input, err := os.ReadFile(file.AbsolutePath)
	if err != nil {
		return err
	}

	var config yamlConfig
	if err = yaml.Unmarshal(input, &config); err != nil {
		return err
	}

	if config.Spring == nil {
		return nil
	}

	dataSourceNode := config.Spring.DataSource

	var dataSource yamlDataSource
	if err = dataSourceNode.Decode(&dataSource); err != nil {
		return fmt.Errorf("failed to decode datasource: %s", err)
	}

	if dataSource.URL == "" && dataSource.DriverClass1 == "" && dataSource.DriverClass2 == "" {
		return nil
	}

	report.AddFramework(detectors.DetectorSpring, spring.TypeDatabase, spring.DataStore{
		Driver: getDriver(dataSource.DriverClass1, dataSource.DriverClass2, dataSource.URL),
	}, source.Source{
		Language:        file.Language,
		LanguageType:    file.LanguageTypeString(),
		Filename:        file.RelativePath,
		StartLineNumber: pointers.Int(dataSourceNode.Line - 1),
	})

	return nil
}

func getDriver(driverClass1 string, driverClass2 string, url string) string {
	if driverClass1 != "" {
		return driverClass1
	}

	if driverClass2 != "" {
		return driverClass2
	}

	if strings.HasPrefix(url, "jdbc:") {
		pieces := strings.Split(url, ":")
		if len(pieces) >= 2 {
			return pieces[1]
		}
	}

	return ""
}

func isJava(dir *file.Path) bool {
	for _, filename := range projectFiles {
		if dir.Join(filename).Exists() {
			return true
		}
	}

	return false
}
