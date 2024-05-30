package engine

import (
	"github.com/bearer/bearer/pkg/commands/process/filelist/files"
	"github.com/bearer/bearer/pkg/commands/process/settings"
	"github.com/bearer/bearer/pkg/scanner/language"
	"github.com/bearer/bearer/pkg/scanner/stats"
)

type Engine interface {
	GetLanguages() []language.Language
	GetLanguageById(id string) language.Language
	Initialize(logLevel string) error
	LoadRule(yamlDefinition string) error
	Scan(config *settings.Config, stats *stats.Stats, reportPath, targetPath string, files []files.File) error
	Close()
}
