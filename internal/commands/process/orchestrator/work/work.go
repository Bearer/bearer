package work

import (
	"github.com/bearer/bearer/internal/commands/process/filelist/files"
	"github.com/bearer/bearer/internal/scanner/stats"
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
	File       files.File
	ReportPath string
}

var RouteInitialize = "/initialize"
var RouteProcess = "/process"
var RouteReduceMemory = "/reduce_memory"
