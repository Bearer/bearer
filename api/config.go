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
	Hello             Endpoint
}

var Endpoints = APIEndpoints{
	RequestFileUpload: Endpoint{
		HttpMethod: "POST",
		Route:      "/cloud/file_uploads",
	},
	ScanFinished: Endpoint{
		HttpMethod: "PUT",
		Route:      "/cloud/scans/{uuid}",
	},
	Hello: Endpoint{
		HttpMethod: "POST",
		Route:      "/cloud/hello",
	},
}
