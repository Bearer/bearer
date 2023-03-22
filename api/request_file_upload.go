package api

import (
	"encoding/json"
)

type RequestFileUpload struct {
	Checksum    string `json:"checksum"`
	ByteSize    int    `json:"byte_size"`
	UUID        string `json:"uuid"`
	Prefix      string `json:"prefix"`
	ContentType string `json:"content_type"`
}

type ActiveStorageDirectUpload struct {
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers"`
}
type FileUploadOffer struct {
	SignedID     string                    `json:"signed_id"`
	UUID         string                    `json:"uuid"`
	DirectUpload ActiveStorageDirectUpload `json:"direct_upload"`
}

func (api *API) RequestFileUpload(fileUpload RequestFileUpload, messageUuid MessageUuid) (*FileUploadOffer, error) {
	endpoint := Endpoints.RequestFileUpload
	bytes, err := api.makeRequest(endpoint.Route, endpoint.HttpMethod, fileUpload)
	if err != nil {
		return nil, err
	}

	var fileUPloadOffer FileUploadOffer

	err = json.Unmarshal(bytes, &fileUPloadOffer)
	if err != nil {
		return nil, err
	}

	return &fileUPloadOffer, nil
}
