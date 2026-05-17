package transport

import (
	"log/slog"
	"net/http"
)

type Middleware func(http.RoundTripper) http.RoundTripper
type RoundTripperFunc func(*http.Request) (*http.Response, error)

func (f RoundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

type TransportBuilder struct {
	base        http.RoundTripper
	middlewares []Middleware
}

func NewTransportBuilder(base http.RoundTripper) *TransportBuilder {
	if base == nil {
		base = http.DefaultTransport
	}
	return &TransportBuilder{
		base:        base,
		middlewares: []Middleware{},
	}
}

func (tb *TransportBuilder) WithMiddlewares(mw ...Middleware) *TransportBuilder {
	tb.middlewares = append(tb.middlewares, mw...)
	return tb
}

func (tb *TransportBuilder) Build() http.RoundTripper {
	rt := tb.base
	for i := len(tb.middlewares) - 1; i >= 0; i-- {
		rt = tb.middlewares[i](rt)
	}
	return rt
}

func WithHeader(headers map[string]string) Middleware {
	return func(next http.RoundTripper) http.RoundTripper {
		return RoundTripperFunc(func(req *http.Request) (*http.Response, error) {
			reqClone := req.Clone(req.Context())

			for key, value := range headers {
				reqClone.Header.Set(key, value)
			}

			return next.RoundTrip(reqClone)
		})
	}
}

func WithLogging(logger *slog.Logger) Middleware {
	return func(next http.RoundTripper) http.RoundTripper {
		return RoundTripperFunc(func(req *http.Request) (*http.Response, error) {
			logger.Info("Request", "url", req.URL.String(), "method", req.Method)

			resp, err := next.RoundTrip(req)
			if err != nil {
				logger.Error("Request failed", "error", err)
				return nil, err
			}

			logger.Info("Response received", "status", resp.StatusCode, "url", req.URL.String())
			return resp, err
		})
	}
}

func WithRetry(retries int) Middleware {
	return func(next http.RoundTripper) http.RoundTripper {
		return RoundTripperFunc(func(req *http.Request) (*http.Response, error) {
			var err error
			for range retries {
				resp, err := next.RoundTrip(req)
				if err == nil {
					return resp, nil
				}
			}
			return nil, err
		})
	}
}
