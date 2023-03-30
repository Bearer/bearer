package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/bearer/bearer/cmd/bearer/build"
)

type API struct {
	client *http.Client
	Host   string
	Token  string
}

type MessageType string

const MessageTypeSuccess MessageType = "success"
const MessageTypeError MessageType = "error"

type Message struct {
	Type MessageType `json:"type"`
	Data interface{} `json:"data"`
}

type ErrorData struct {
	Message string `json:"message"`
}

func New(config API) *API {
	return &API{
		client: &http.Client{},
		Token:  config.Token,
		Host:   config.Host,
	}
}

var ErrTokenInvalid = errors.New("bearer token is invalid")

func (api *API) makeRequest(route string, httpMethod string, data interface{}) ([]byte, error) {
	fullURL := fmt.Sprintf("https://%s%s", api.Host, route)

	var req *http.Request

	if data != nil {
		sendingData, err := json.Marshal(data)
		if err != nil {
			return nil, fmt.Errorf("fail marshaling data %e", err)
		}

		req, err = http.NewRequest(httpMethod, fullURL, bytes.NewBuffer(sendingData))
		if err != nil {
			return nil, fmt.Errorf("fail creating request %e %s", err, fullURL)
		}
		req.Header.Set("Content-Type", "application/json")
		defer req.Body.Close()
	} else {
		var err error
		req, err = http.NewRequest(httpMethod, fullURL, nil)
		if err != nil {
			return nil, fmt.Errorf("fail creating request %e %s", err, fullURL)
		}
	}

	req.Header.Set("Authorization", api.Token)
	req.Header.Set("X-Bearer-SHA", build.CommitSHA)
	req.Header.Set("X-Bearer-Version", build.Version)

	resp, err := api.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fail getting response %e %s", err, fullURL)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("fail reading response body %s %s", err, fullURL)
	}

	if resp.StatusCode == http.StatusUnauthorized {
		var BearerTokenInvalidMsgErr = "invalid token"
		type BearerTokenInvalidMsg struct {
			Error string `json:"error"`
		}

		var unathorizedErr BearerTokenInvalidMsg
		err := json.Unmarshal(body, &unathorizedErr)
		if err != nil {
			return nil, err
		}

		if unathorizedErr.Error == BearerTokenInvalidMsgErr {
			return nil, ErrTokenInvalid
		}
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != 204 {
		return nil, fmt.Errorf("didn't get response status 200/204 got %d %s", resp.StatusCode, string(body))
	}

	return body, nil
}
