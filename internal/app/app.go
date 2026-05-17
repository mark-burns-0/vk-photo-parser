package app

import (
	"github.com/mark-burns-0/vk-photo-parser/internal/parser"
)

func Run(config parser.Configer) {
	pr := parser.New(config, config.GetLogLevel())
	pr.ParsePhoto().
		Download()
}
