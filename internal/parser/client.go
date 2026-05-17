package parser

import (
	"io"
	"log/slog"
	"net/http"

	"github.com/mark-burns-0/vk-photo-parser/pkg/transport"
)

func newClient(log *slog.Logger) *ParserClient {
	tr := transport.NewTransportBuilder(nil).
		WithMiddlewares(
			transport.WithLogging(log),
			transport.WithRetry(3),
		).
		Build()

	return &ParserClient{
		client: &http.Client{
			Transport: tr,
		},
	}
}

func (c *ParserClient) Get(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	return c.client.Do(req)
}

func (c *ParserClient) Post(url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}

	return c.client.Do(req)
}
