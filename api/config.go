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
	Version Endpoint
}

var Endpoints = APIEndpoints{
	Version: Endpoint{
		HttpMethod: "GET",
		Route:      "/api/v1/bearerpublic/version",
	},
}
