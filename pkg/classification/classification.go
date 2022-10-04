package classsification

import (
	"github.com/bearer/curio/pkg/classification/frameworks"
	"github.com/bearer/curio/pkg/classification/interfaces"
	"github.com/bearer/curio/pkg/classification/schema"
)

type Classifier struct {
	config Config

	interfaces interfaces.Classifier
	schema     schema.Classifier
	frameworks frameworks.Classifier
}

type Config struct {
}

func NewClassifier(config *Config) *Classifier {
	// todo: config setup
	return &Classifier{}
}
