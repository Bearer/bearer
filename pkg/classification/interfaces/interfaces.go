package interfaces

import "github.com/bearer/curio/pkg/report/interfaces"

type Classifier struct {
	config Config
}

type Config struct {
}

func New(config Config) *Classifier {
	return &Classifier{}
}

func (classifier *Classifier) Classify(data interfaces.Interface) (ClassifiedInterface, error) {
	// todo: implement interface classification (bigbear etc...)
	return ClassifiedInterface{}, nil
}
