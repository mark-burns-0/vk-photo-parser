package parser

import (
	"log/slog"
	"net/http"
	url_ "net/url"
	"strings"

	"github.com/mark-burns-0/vk-photo-parser/pkg/transport"
)

func newClient(log *slog.Logger, baseURL string) *ParserClient {
	tr := transport.NewTransportBuilder(nil).
		WithMiddlewares(
			transport.WithLogging(log),
			transport.WithRetry(3),
			transport.WithHeader(map[string]string{
				"Accept":       "application/json",
				"Content-Type": "application/json",
			}),
		).
		Build()

	return &ParserClient{
		client: &http.Client{
			Transport: tr,
		},
		baseURL: baseURL,
	}
}

func (c *ParserClient) Get(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", c.baseURL+url, nil)
	if err != nil {
		return nil, err
	}

	return c.client.Do(req)
}

func (c *ParserClient) Post(url string, params map[string]string) (*http.Response, error) {
	form := url_.Values{}

	for key, value := range params {
		form.Set(key, value)
	}

	req, err := http.NewRequest("POST", c.baseURL+url, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}

	return c.client.Do(req)
}
