package implementation

import (
	"github.com/bearer/bearer/pkg/commands/process/filelist/files"
	"github.com/bearer/bearer/pkg/commands/process/orchestrator"
	"github.com/bearer/bearer/pkg/commands/process/orchestrator/work"
	"github.com/bearer/bearer/pkg/commands/process/settings"
	"github.com/bearer/bearer/pkg/engine"
	"github.com/bearer/bearer/pkg/scanner/language"
	"github.com/bearer/bearer/pkg/scanner/stats"
)

type implementation struct {
	languages    []language.Language
	orchestrator *orchestrator.Orchestrator
}

func New(languages []language.Language) engine.Engine {
	return &implementation{languages: languages}
}

func (engine *implementation) GetLanguages() []language.Language {
	return engine.languages
}

func (engine *implementation) GetLanguageById(id string) language.Language {
	for _, language := range engine.languages {
		if language.ID() == id {
			return language
		}
	}

	return nil
}

func (engine *implementation) Initialize(logLevel string) error {
	return nil
}

func (engine *implementation) LoadRule(yamlDefinition string) error {
	return nil
}

func (engine *implementation) Scan(
	config *settings.Config,
	stats *stats.Stats,
	reportPath,
	targetPath string,
	files []files.File,
) error {
	if engine.orchestrator == nil {
		var err error
		engine.orchestrator, err = orchestrator.New(work.Repository{Dir: targetPath}, config, stats, len(files))
		if err != nil {
			return err
		}
	}

	return engine.orchestrator.Scan(reportPath, files)
}

func (engine *implementation) Close() {
	if engine.orchestrator != nil {
		engine.orchestrator.Close()
	}
}
