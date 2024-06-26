package api

import "net/url"

func (api *API) Version(languages []string) ([]byte, error) {
	endpoint := Endpoints.Version
	languageQuery := url.Values{}
	for _, lang := range languages {
		languageQuery.Add("_json[]", lang)
	}
	response, err := api.makeRequest(endpoint.Route+languageQuery.Encode(), endpoint.HttpMethod, nil)

	return response, err
}
