package app

import (
	"github.com/mark-burns-0/vk-photo-parser/internal/parser"
)

type Configer interface {
	GetOwnerID() string
	GetAlbumID() string
	GetToken() string
}

func Run(config Configer) {
	pr := parser.New(config)
	pr.Parse()
}
