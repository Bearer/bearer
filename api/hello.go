package api

func (api *API) Hello() ([]byte, error) {
	endpoint := Endpoints.Hello
	response, err := api.makeRequest(endpoint.Route, endpoint.HttpMethod, nil)

	return response, err
}
