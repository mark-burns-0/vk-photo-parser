package parser

import "net/http"

type Parser struct {
	urls   []string
	client *ParserClient
	cfg    Configer
}

type ParserClient struct {
	client  *http.Client
	baseURL string
}

type VKPhotosResponse struct {
	Response struct {
		Count int           `json:"count"`
		Items []VKPhotoItem `json:"items"`
	} `json:"response"`
}

type VKPhotoItem struct {
	Sizes []struct {
		Type string `json:"type"`
		URL  string `json:"url"`
	} `json:"sizes"`
}
