package work

import (
	"time"
)

type InitializeResponse struct {
	Error string
}

type ProcessResponse struct {
	Error string
}

type Repository struct {
	Dir               string
	PreviousCommitSHA string
	CommitSHA         string
}

type ProcessRequest struct {
	Repository
	File                   File
	ReportPath             string
	BlameRevisionsFilePath string
}

type File struct {
	Timeout  time.Duration
	FilePath string
}

var RouteInitialize = "/initialize"
var RouteProcess = "/process"
