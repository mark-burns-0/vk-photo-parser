package config

type Config struct {
	Token   string `env:"VK_TOKEN"`
	AlbumID string `env:"VK_ALBUM_ID"`
	OwnerID string `env:"VK_OWNER_ID"`
	Version string `env:"VK_API_VERSION" env-default:"5.131"`
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
