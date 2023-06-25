package client

import "net/http"

type impl struct {
	baseURL    string
	httpclient *http.Client
}

// New creates a new instance of Utility.
func New(baseURL string, httpclient *http.Client) Utility {
	return &impl{
		baseURL:    baseURL,
		httpclient: httpclient,
	}
}
