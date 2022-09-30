package classsification

import (
	"encoding/json"

	reportinterfaces "github.com/bearer/curio/pkg/report/interfaces"
	reportschema "github.com/bearer/curio/pkg/report/schema"

	"github.com/bearer/curio/pkg/classification/interfaces"
	"github.com/bearer/curio/pkg/classification/schema"
	"github.com/bearer/curio/pkg/report"
	"github.com/bearer/curio/pkg/util/jsonlinesreader"
)

type Classifier struct {
	config     Config
	interfaces interfaces.Classifier
	schema     schema.Classifier
}

type Config struct {
}

func NewClassifier(config *Config) *Classifier {
	// todo: config setup
	return &Classifier{}
}

func (classifer *Classifier) Classify(reportPath string) (classifyReportPath string, err error) {
	// get tmpdir
	// create a new file
	reader, err := jsonlinesreader.New(reportPath)
	if err != nil {
		return "", err
	}
	defer reader.Close()

	for reader.Next() {
		line := reader.Data()
		linePropertiesMap := line.(map[string]interface{})

		lineType := linePropertiesMap["type"].(string)

		switch lineType {
		case string(report.TypeInterface):
			var data reportinterfaces.Interface
			err := json.Unmarshal([]byte(reader.Text()), &data)
			if err != nil {
				// todo: we probabbly want to write error to report and continue
				return "", err
			}
			classifiedData, err := classifer.interfaces.Classify(data)
			if err != nil {
				// todo: we probabbly want to write error to report and continue
				return "", err
			}
			// write classifiedData to our newly created file as jsonline
		case string(report.TypeSchema):
			var data reportschema.Schema
			err := json.Unmarshal([]byte(reader.Text()), &data)
			if err != nil {
				// todo: we probabbly want to write error to report and continue
				return "", err
			}
			classifiedData, err := classifer.schema.Classify(data)
			if err != nil {
				// todo: we probabbly want to write error to report and continue
				return "", err
			}
			// write classifiedData to our newly created file as jsonline
		default:
			// write unclassifable data to a report file as jsonline
		}
	}

	return "", nil
}
