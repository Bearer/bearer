package work

import (
	"time"

	"github.com/bearer/bearer/new/detector/evaluator/stats"
)

type InitializeResponse struct {
	Error string
}

type ProcessResponse struct {
	FileStats *stats.FileStats
	Error     string
}

type Repository struct {
	Dir               string
	PreviousCommitSHA string
	CommitSHA         string
}

type ProcessRequest struct {
	Repository
	File       File
	ReportPath string
}

type File struct {
	Timeout  time.Duration
	FilePath string
}

var RouteInitialize = "/initialize"
var RouteProcess = "/process"
