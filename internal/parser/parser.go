package parser

import (
	"log/slog"
	"net/http"

	"github.com/mark-burns-0/vk-photo-parser/pkg/transport"
)

type Configer interface {
	GetOwnerID() string
	GetAlbumID() string
	GetToken() string
	GetVersion() string
}

func New(cfg Configer) *Parser {

	return &Parser{}
}

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

func (p *Parser) Parse() *Parser {
	return p
}

func (p *Parser) Download() *Parser {
	return p
}
