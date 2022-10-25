package work

import (
	"time"
)

type StatusResponse struct {
	ClassifierError     string
	CustomDetectorError string
}

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
}

type File struct {
	Timeout  time.Duration
	FilePath string
}

var RouteStatus = "/status"
var RouteProcess = "/process"
var RouteCustomDetector = "/custom_detector"
