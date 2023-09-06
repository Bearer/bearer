package api

import "fmt"

func (api *API) Version(languages []string) ([]byte, error) {
	endpoint := Endpoints.Version

	queryString := "?"
	for _, language := range languages {
		queryString += fmt.Sprintf("language[]=%s&", language)
	}

	response, err := api.makeRequest(endpoint.Route+queryString, endpoint.HttpMethod, nil)

	return response, err
}
