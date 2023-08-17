package api

import (
	"net/http"
	"time"
)

var UploadClient = &http.Client{Timeout: 60 * time.Second}

type Endpoint struct {
	HttpMethod string
	Route      string
}

type APIEndpoints struct {
	RequestFileUpload Endpoint
	ScanFinished      Endpoint
	FetchIgnores      Endpoint
	Hello             Endpoint
}

var Endpoints = APIEndpoints{
	RequestFileUpload: Endpoint{
		HttpMethod: "POST",
		Route:      "/cloud/file_uploads",
	},
	ScanFinished: Endpoint{
		HttpMethod: "POST",
		Route:      "/cloud/scans",
	},
	FetchIgnores: Endpoint{
		HttpMethod: "GET",
		Route:      "/cloud/ignores",
	},
	Hello: Endpoint{
		HttpMethod: "POST",
		Route:      "/cloud/hello",
	},
}
