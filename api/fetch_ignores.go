package api

import (
	"encoding/json"

	ignoretypes "github.com/bearer/bearer/internal/util/ignore/types"
)

type CloudIgnoreData struct {
	ProjectFound             bool                                      `json:"project_found"`
	Ignores                  []string                                  `json:"ignores"`
	StaleIgnores             []string                                  `json:"stale_local_ignores"`
	CloudIgnoredFingerprints map[string]ignoretypes.IgnoredFingerprint `json:"detailed_cloud_ignores"`
}

type CloudIgnorePayload struct {
	Project           string   `json:"project"`
	LocalIgnores      []string `json:"local_ignores"`
	PullRequestNumber string   `json:"pull_request_number,omitempty"`
}

func (api *API) FetchIgnores(fullname string, pullRequestNumber string, localIgnores []string) (*CloudIgnoreData, error) {
	endpoint := Endpoints.FetchIgnores

	bytes, err := api.makeRequest(endpoint.Route, endpoint.HttpMethod,
		Message{
			Type: MessageTypeSuccess,
			Data: CloudIgnorePayload{
				Project:           fullname,
				LocalIgnores:      localIgnores,
				PullRequestNumber: pullRequestNumber,
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
