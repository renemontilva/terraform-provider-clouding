package clouding

import (
	"bytes"
	"fmt"
	"net/http"
)

const (
	// Server
	ENDPOINT = "https://api.clouding.io"
	VERSION  = "v1"
)

type API struct {
	Endpoint string
	Token    string
}

type ErrorResponse struct {
	Type     string   `json:"type"`
	Title    string   `json:"title"`
	Status   int      `json:"status"`
	Detail   string   `json:"detail,omitempty"`
	Instance string   `json:"instance,omitempty"`
	TraceID  string   `json:"traceId,omitempty"`
	Errors   []string `json:"errors,omitempty"`
}

type option func(*API) error

func NewAPI(token string, options ...option) (*API, error) {
	api := API{
		Endpoint: ENDPOINT,
		Token:    token,
	}

	for _, option := range options {
		err := option(&api)
		if err != nil {
			return nil, err
		}
	}
	return &api, nil
}

func WithEndpoint(endpoint string) option {
	return func(a *API) error {
		a.Endpoint = endpoint
		return nil
	}
}

func (a *API) sendRequest(method string, path string, body []byte) (*http.Response, error) {
	request, err := http.NewRequest(method, fmt.Sprintf("%s/%s/%s", a.Endpoint, VERSION, path), bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	// Set headers
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-API-KEY", a.Token)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	return response, nil
}
