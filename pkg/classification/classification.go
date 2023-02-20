package classification

import (
	"github.com/bearer/bearer/pkg/classification/db"
	"github.com/bearer/bearer/pkg/classification/dependencies"
	"github.com/bearer/bearer/pkg/classification/frameworks"
	"github.com/bearer/bearer/pkg/classification/interfaces"
	"github.com/bearer/bearer/pkg/classification/schema"
	config "github.com/bearer/bearer/pkg/commands/process/settings"
	"github.com/bearer/bearer/pkg/util/url"
)

type Classifier struct {
	config Config

	Interfaces   *interfaces.Classifier
	Schema       *schema.Classifier
	Dependencies *dependencies.Classifier
	Frameworks   *frameworks.Classifier
}

type Config struct {
	Config config.Config
}

func NewClassifier(config *Config) (*Classifier, error) {
	interfacesClassifier, err := interfaces.New(
		interfaces.Config{
			Recipes:         db.Default().Recipes,
			InternalDomains: config.Config.Scan.InternalDomains,
			DomainResolver: url.NewDomainResolver(
				!config.Config.Scan.DisableDomainResolution,
				config.Config.Scan.DomainResolutionTimeout,
			),
		},
	)
	if err != nil {
		return nil, err
	}

	// apply subject mapping override, if present
	var knownPersonObjectPatterns []db.KnownPersonObjectPattern
	if config.Config.Scan.DataSubjectMapping != "" {
		knownPersonObjectPatterns = db.DefaultWithMapping(config.Config.Scan.DataSubjectMapping).KnownPersonObjectPatterns
	} else {
		knownPersonObjectPatterns = db.Default().KnownPersonObjectPatterns
	}

	schemaClassifier := schema.New(
		schema.Config{
			DataTypes:                      db.Default().DataTypes,
			DataTypeClassificationPatterns: db.Default().DataTypeClassificationPatterns,
			KnownPersonObjectPatterns:      knownPersonObjectPatterns,
			Context:                        config.Config.Scan.Context,
		},
	)

	dependenciesClassifier := dependencies.New(
		dependencies.Config{
			Recipes: db.Default().Recipes,
		},
	)

	frameworksClassifier := frameworks.New(
		frameworks.Config{
			Recipes: db.Default().Recipes,
		},
	)

	return &Classifier{
		config:       *config,
		Dependencies: dependenciesClassifier,
		Interfaces:   interfacesClassifier,
		Schema:       schemaClassifier,
		Frameworks:   frameworksClassifier,
	}, nil
}
