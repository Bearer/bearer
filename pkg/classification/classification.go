package classification

import (
	"regexp"

	"github.com/bearer/curio/pkg/classification/db"
	"github.com/bearer/curio/pkg/classification/dependencies"
	"github.com/bearer/curio/pkg/classification/interfaces"
	"github.com/bearer/curio/pkg/classification/schema"
	config "github.com/bearer/curio/pkg/commands/process/settings"
	"github.com/bearer/curio/pkg/util/url"
)

type Classifier struct {
	config Config

	Interfaces   *interfaces.Classifier
	Schema       schema.Classifier
	Dependencies *dependencies.Classifier
}

type Config struct {
	Config config.Config
}

func NewClassifier(config *Config) (*Classifier, error) {
	interfacesClassifier, err := interfaces.New(
		interfaces.Config{
			Recipes:                db.Default(),
			InternalDomainMatchers: []*regexp.Regexp{},
			DomainResolver: url.NewDomainResolver(
				!config.Config.Scan.DisableDomainResolution,
				config.Config.Scan.DomainResolutionTimeout,
			),
		},
	)
	if err != nil {
		return nil, err
	}

	dependenciesClassifier := dependencies.New(
		dependencies.Config{
			Recipes: db.Default(),
		},
	)

	return &Classifier{
		config:       *config,
		Dependencies: dependenciesClassifier,
		Interfaces:   interfacesClassifier,
	}, nil
}
