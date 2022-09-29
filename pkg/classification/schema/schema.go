package schema

import (
	"github.com/bearer/curio/pkg/report/schema"
)

type Classifier struct {
	config Config
}

type Config struct {
}

func (classifier *Classifier) Classify(data schema.Schema) (ClassifiedSchema, error) {
	// todo: implement interface classification (bigbear etc...)
	return ClassifiedSchema{}, nil
}
