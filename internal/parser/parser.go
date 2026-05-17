package parser

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"strconv"

	"github.com/k0kubun/pp"
)

type Configer interface {
	GetOwnerID() string
	GetAlbumID() string
	GetVersion() string
	GetBaseURL() string
	GetToken() string
}

func New(cfg Configer, logLvl slog.Level) *Parser {
	client := newClient(
		slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
				Level: logLvl,
			}),
		),
		cfg.GetBaseURL(),
	)

	return &Parser{
		client: client,
		cfg:    cfg,
	}
}

func (p *Parser) ParsePhoto() *Parser {
	op := "parser.Parse"

	var err error
	var r *http.Response
	offset := 0
	count := 100

	body := map[string]string{
		"owner_id":     p.cfg.GetOwnerID(),
		"album_id":     p.cfg.GetAlbumID(),
		"access_token": p.cfg.GetToken(),
		"offset":       strconv.Itoa(offset),
		"count":        strconv.Itoa(count),
		"v":            p.cfg.GetVersion(),
	}
	resp := &VKPhotosResponse{}
	allResponses := []VKPhotosResponse{}

	for {
		r, err = p.client.Post("photos.get", body)
		if err != nil {
			slog.Error("Failed to get photos", "error", err, "operation", op)
			break
		}
		json.NewDecoder(r.Body).Decode(&resp)
		r.Body.Close()

		if len(resp.Response.Items) == 0 {
			break
		}
		allResponses = append(allResponses, *resp)

		offset += count
		body["offset"] = strconv.Itoa(offset)
	}
	pp.Print(len(allResponses))
	return p
}

func (p *Parser) Download() *Parser {
	return p
}
