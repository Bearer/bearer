package dependencies

import (
	"github.com/bearer/curio/pkg/classification/db"
	"github.com/bearer/curio/pkg/report"
)

type ClassifiedDependency struct {
	*report.Detection
	Classification Classification `json:"classification"`
}

type Classification struct {
	RecipeMatch bool
	RecipeName  string
}

type Classifier struct {
	config Config
}

type Config struct {
	recipes []db.Recipe
}

func New(config Config) *Classifier {
	return &Classifier{config: config}
}

func (classifier *Classifier) Classify(data report.Detection) (ClassifiedDependency, error) {
	// todo: implement interface classification (bigbear etc...)
	for _, recipe := range classifier.config.recipes { //nolint:all,unused

	}

	return ClassifiedDependency{
		Detection: &data,
		Classification: Classification{
			RecipeMatch: true,
			RecipeName:  "stripe",
		},
	}, nil
}
