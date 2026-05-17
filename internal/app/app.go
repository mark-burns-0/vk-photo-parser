package app

import (
	"log/slog"
	"os"

	"github.com/mark-burns-0/vk-photo-parser/internal/parser"
)

func Run(config parser.Configer) {
	slog.SetDefault(slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: config.GetLogLevel(),
		}),
	))

	pr := parser.New(config)
	pr.ParsePhoto().
		Download()
}
