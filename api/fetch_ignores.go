package api

import (
	"encoding/json"

	"github.com/rs/zerolog/log"
)

type CloudIgnoreData struct {
	ProjectFound bool     `json:"project_found"`
	CloudIgnores []string `json:"cloud_ignores"`
}

func (api *API) FetchIgnores(fullname string) (*CloudIgnoreData, error) {
	endpoint := Endpoints.FetchIgnores
	bytes, err := api.makeRequest(endpoint.Route, endpoint.HttpMethod,
		Message{
			Type: MessageTypeSuccess,
			Data: fullname,
		})
	if err != nil {
		return nil, err
	}

	var cloudIgnoreData CloudIgnoreData
	err = json.Unmarshal(bytes, &cloudIgnoreData)
	if err != nil {
		return nil, err
	}

	log.Error().Msgf("Ignores: %#v", cloudIgnoreData)

	return &cloudIgnoreData, err
}
