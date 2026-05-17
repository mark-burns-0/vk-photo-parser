package parser

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/mark-burns-0/vk-photo-parser/internal/pool"
)

const (
	photosGetMethod = "photos.get"
)

const (
	offsetStep      = 100
	countStep       = 100
	amountOfWorkers = 6
)

type Configer interface {
	GetOwnerID() string
	GetAlbumID() string
	GetVersion() string
	GetBaseURL() string
	GetToken() string
	GetLogLevel() slog.Level
	GetOutputFolder() string
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

		urls, err := extractBaseUrls(resp.Response.Items)
		if err != nil {
			slog.Error("Failed to extract base urls from items", "error", err, "operation", op)
		} else {
			p.urls = append(p.urls, urls...)
		}

		offset += count
		body["offset"] = strconv.Itoa(offset)
	}

	return p
}

func (p *Parser) Download() *Parser {
	pl := pool.NewPool(amountOfWorkers, func(num int, url string) error {
		err := os.MkdirAll(filepath.Join(
			p.cfg.GetOutputFolder(),
			time.Now().Format("2006-01-02")),
			0755,
		)
		if err != nil {
			return err
		}

		out, err := os.Create(filepath.Join(
			p.cfg.GetOutputFolder(),
			time.Now().Format("2006-01-02"),
			fmt.Sprintf(
				"%s_%s.jpg",
				time.Now().Format("2006-01-02_15-04-05"),
				uuid.New().String(),
			),
		))
		if err != nil {
			return err
		}
		defer out.Close()

		resp, err := p.client.Get(url)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		_, err = io.Copy(out, resp.Body)
		time.Sleep(time.Duration(300) * time.Millisecond)
		return nil
	})
	pl.Create()

	for _, url := range p.urls {
		pl.Handle(url)
	}
	pl.Wait()
	pl.Stats()

	return p
}

func extractBaseUrls(data []VKPhotoItem) ([]string, error) {
	var op = "parser.extractBaseUrls"

	urls := make([]string, 0, len(data))

	for _, d := range data {
		for _, size := range d.Sizes {
			if size.Type == "base" {
				urls = append(urls, size.URL)
			}
		}
	}

	if len(urls) == 0 {
		return nil, fmt.Errorf("%s: %s", op, "exctracted nothing")
	}

	return urls, nil
}
