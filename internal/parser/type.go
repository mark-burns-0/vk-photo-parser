package parser

import "net/http"

type Parser struct {
	urls   []string
	client *ParserClient
}

type ParserClient struct {
	client *http.Client
}
