package api

import (
	"encoding/json"
)

type CloudIgnoreData struct {
	ProjectFound bool     `json:"project_found"`
	CloudIgnores []string `json:"cloud_ignores"`
}

type CloudIgnorePayload struct {
	Project      string   `json:"project"`
	LocalIgnores []string `json:"local_ignores"`
}

func (api *API) FetchIgnores(fullname string, localIgnores []string) (*CloudIgnoreData, error) {
	endpoint := Endpoints.FetchIgnores
	bytes, err := api.makeRequest(endpoint.Route, endpoint.HttpMethod,
		Message{
			Type: MessageTypeSuccess,
			Data: CloudIgnorePayload{
				Project:      fullname,
				LocalIgnores: localIgnores,
			},
		})
	if err != nil {
		return nil, err
	}

	var cloudIgnoreData CloudIgnoreData
	err = json.Unmarshal(bytes, &cloudIgnoreData)
	if err != nil {
		return nil, err
	}

	return &cloudIgnoreData, err
}
