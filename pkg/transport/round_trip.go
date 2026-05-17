package transport

import (
	"net/http"
)

type Transport struct {
	base http.RoundTripper
}

func NewTransport(base http.RoundTripper) http.RoundTripper {
	if base == nil {
		base = http.DefaultTransport
	}

	return &Transport{
		base: http.DefaultTransport,
	}
}

func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	return t.base.RoundTrip(req)
}

func (t *Transport) WithToken(token string) *Transport {
	return &Transport{
		base: t.base,
	}
}

func (t *Transport) WithLogging(token string) *Transport {
	return &Transport{
		base: t.base,
	}
}
