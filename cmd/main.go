package main

import (
	"github.com/mark-burns-0/vk-photo-parser/internal/app"
	"github.com/mark-burns-0/vk-photo-parser/internal/config"
)

func main() {
	config := config.MustLoadConfig()
	app.Run(config)
}
