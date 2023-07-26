package work

import (
	"github.com/bearer/bearer/new/detector/evaluator/stats"
	"github.com/bearer/bearer/pkg/commands/process/filelist"
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
	File       filelist.File
	ReportPath string
}

var RouteInitialize = "/initialize"
var RouteProcess = "/process"
