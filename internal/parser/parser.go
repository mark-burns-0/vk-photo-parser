package parser

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"strconv"

	"github.com/k0kubun/pp"
)

const (
	photosGetMethod = "photos.get"
)

const (
	offsetStep = 100
	countStep  = 100
)

type Configer interface {
	GetOwnerID() string
	GetAlbumID() string
	GetVersion() string
	GetBaseURL() string
	GetToken() string
	GetLogLevel() slog.Level
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
	count := countStep

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
		r, err = p.client.Post(photosGetMethod, body)
		if err != nil {
			slog.Error("Failed to get photos", "error", err, "operation", op)
			break
		}
		err = json.NewDecoder(r.Body).Decode(&resp)
		if err != nil {
			r.Body.Close()
			slog.Error("Failed to decode body", "error", err, "operation", op)
			break
		}
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
