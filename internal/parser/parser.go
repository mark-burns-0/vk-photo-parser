package parser

import (
	"log/slog"
	"os"
)

type Configer interface {
	GetOwnerID() string
	GetAlbumID() string
	GetToken() string
	GetVersion() string
}

func New(cfg Configer, logLvl slog.Level) *Parser {
	client := newClient(
		slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
				Level: logLvl,
			}),
		),
	)

	return &Parser{
		client: client,
	}
}

func (p *Parser) Parse() *Parser {
	p.client.Post("https://api.vk.com/method/photos.get", nil)

	return p
}

func (p *Parser) Download() *Parser {
	return p
}
