package api_v1

import (
	"net/http"
)

type Client struct {
	apiKey     string
	baseUrl    string
	httpClient *http.Client
}

func ensureHttpClient(httpClient *http.Client) *http.Client {
	if httpClient == nil {
		return &http.Client{}
	}
	return httpClient
}

func NewClient(apiKey string, baseUrl string, httpClient *http.Client) *Client {
	return &Client{
		apiKey:     apiKey,
		baseUrl:    baseUrl,
		httpClient: ensureHttpClient(httpClient),
	}
}
