package parser

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"

	"github.com/k0kubun/pp"
	"github.com/mark-burns-0/vk-photo-parser/internal/pool"
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

	pl := pool.NewPool(4, func(num int, data VKPhotosResponse) error {
		fmt.Println(num, len(data.Response.Items))

		return nil
	})
	pl.Create()

	for _, response := range allResponses {
		pl.Handle(response)
	}
	pl.Wait()
	pl.Stats()

	return p
}

func (p *Parser) Download() *Parser {
	return p
}
