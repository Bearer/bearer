package classsification

import (
	"github.com/bearer/curio/pkg/classification/db"
	"github.com/bearer/curio/pkg/classification/dependencies"
	"github.com/bearer/curio/pkg/classification/interfaces"
	"github.com/bearer/curio/pkg/classification/schema"
)

type Classifier struct {
	config Config

	Interfaces   interfaces.Classifier
	Schema       schema.Classifier
	Dependencies *dependencies.Classifier
}

type Config struct {
}

func NewClassifier(config *Config) *Classifier {
	// todo: config setup
	dependenciesClassifier := dependencies.New(
		dependencies.Config{
			Recipes: db.Default(),
		},
	)

	return &Classifier{
		config:       *config,
		Dependencies: dependenciesClassifier,
	}
}
