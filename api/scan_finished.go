package api

import (
	"strings"
)

type ScanResult struct {
	SignedID string `json:"signed_id"`
}

func (api *API) makeScanFinishedRequest(scanUUID MessageUuid, message Message) error {
	endpoint := Endpoints.ScanFinished
	route := strings.ReplaceAll(endpoint.Route, "{uuid}", string(scanUUID))
	_, err := api.makeRequest(route, endpoint.HttpMethod, message)
	return err
}

func (api *API) ScanFinished(scanResult ScanResult, scanUUID MessageUuid) error {
	return api.makeScanFinishedRequest(scanUUID, Message{Type: MessageTypeSuccess, Data: scanResult})
}

func (api *API) ScanFinishedError(err error, scanUUID MessageUuid) error {
	return api.makeScanFinishedRequest(scanUUID, Message{Type: MessageTypeError, Data: ErrorData{Message: err.Error()}})
}
