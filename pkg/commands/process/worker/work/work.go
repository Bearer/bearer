package work

import (
	"time"

	config "github.com/bearer/curio/pkg/commands/process/settings"
)

type ProcessResponse struct {
	Error error
}

type Repository struct {
	Dir               string
	PreviousCommitSHA string
	CommitSHA         string
}

type ProcessRequest struct {
	Repository
	Files                  []File
	FilePath               string
	BlameRevisionsFilePath string
	CustomDetectorConfig   *config.RulesConfig
}

type File struct {
	Timeout  time.Duration
	FilePath string
}

var RouteStatus = "/status"
var RouteProcess = "/process"
var RouteCustomDetector = "/custom_detector"
