package config

import (
	"log/slog"
	"strconv"
)

type Config struct {
	Token    string `env:"VK_TOKEN"`
	AlbumID  string `env:"VK_ALBUM_ID"`
	OwnerID  string `env:"VK_OWNER_ID"`
	Version  string `env:"VK_API_VERSION" env-default:"5.131"`
	BaseURL  string `env:"VK_API_BASE_URL" env-default:"https://api.vk.com/method/"`
	LogLevel string `env:"LOG_LEVEL" env-default:"8"` // default to error level
}

func (c *Config) GetOwnerID() string {
	return c.OwnerID
}

func (c *Config) GetAlbumID() string {
	return c.AlbumID
}

func (c *Config) GetToken() string {
	return c.Token
}

func (c *Config) GetVersion() string {
	return c.Version
}

func (c *Config) GetBaseURL() string {
	return c.BaseURL
}

func (c *Config) GetLogLevel() slog.Level {
	rawString, err := strconv.Atoi(c.LogLevel)
	if err != nil {
		return slog.LevelInfo
	}
	return slog.Level(rawString)
}
