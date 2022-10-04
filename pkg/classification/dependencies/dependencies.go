package interfaces

import (
	"github.com/bearer/curio/pkg/report/dependencies"
)

type Classifier struct {
	config Config
}

type Config struct {
}

func New(config Config) *Classifier {
	return &Classifier{}
}

func (classifier *Classifier) Classify(data dependencies.Dependency) (ClassifiedDependency, error) {
	// todo: implement interface classification (bigbear etc...)
	return ClassifiedDependency{}, nil
}
