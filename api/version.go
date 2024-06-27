package api

import (
	"fmt"
	"net/url"
)

func (api *API) Version(languages []string) ([]byte, error) {
	endpoint := Endpoints.Version
	queryString := "/?"
	for _, lang := range languages {
		queryString += fmt.Sprintf("_json[]=%s&", url.QueryEscape(lang))
	}
	response, err := api.makeRequest(endpoint.Route+queryString, endpoint.HttpMethod, nil)

	return response, err
}
