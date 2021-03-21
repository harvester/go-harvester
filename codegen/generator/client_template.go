package generator

var clientTemplate = `package client

import (
	"net/http"
	"net/url"
)

type Client struct {
	HTTPClient *http.Client
	BaseURL    *url.URL

	{{range .schemas}}
    {{.CodeName}}s *{{.CodeName}}Client
{{- end}}
}

func New(baseURL *url.URL, httpClient *http.Client) *Client {

	c := &Client{
		HTTPClient: httpClient,
		BaseURL:    baseURL,
	}

	{{range .schemas}}
    c.{{.CodeName}}s = new{{.CodeName}}Client(c)
{{- end}}

	return c
}
`
