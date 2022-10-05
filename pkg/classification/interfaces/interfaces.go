package interfaces

import (
	"github.com/bearer/curio/pkg/classification/db"
	"github.com/bearer/curio/pkg/report"
)

type ClassifiedInterface struct {
	*report.Detection
	Classification *Classification `json:"classification"`
}

type Classification struct {
	RecipeName string
}

type Classifier struct {
	config Config
}

type Config struct {
	recipes []db.Recipe
}

func New(config Config) *Classifier {
	return &Classifier{config}
}

func (classifier *Classifier) Classify(data report.Detection) (ClassifiedInterface, error) {
	// todo: implement interface classification (bigbear etc...)
	for _, recipe := range classifier.config.recipes { //nolint:all,unused

	}

	return ClassifiedInterface{
		Detection: &data,
		Classification: &Classification{
			RecipeName: "stripe",
		},
	}, nil
}
