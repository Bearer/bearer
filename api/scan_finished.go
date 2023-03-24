package api

type ScanResult struct {
	SignedID string `json:"signed_id"`
}

func (api *API) ScanFinished(meta interface{}) error {
	endpoint := Endpoints.ScanFinished
	_, err := api.makeRequest(
		endpoint.Route,
		endpoint.HttpMethod,
		Message{
			Type: MessageTypeSuccess,
			Data: meta,
		})

	return err
}
