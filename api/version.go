package api

func (api *API) Version(languages []string) ([]byte, error) {
	endpoint := Endpoints.Version
	response, err := api.makeRequest(endpoint.Route, endpoint.HttpMethod, languages)

	return response, err
}
