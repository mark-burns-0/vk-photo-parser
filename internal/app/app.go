package app

import (
	"log/slog"

	"github.com/mark-burns-0/vk-photo-parser/internal/parser"
)

func Run(config parser.Configer) {
	pr := parser.New(config, slog.LevelDebug)
	pr.Parse().
		Download()
}
