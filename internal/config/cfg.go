package config

import "github.com/ilyakaznacheev/cleanenv"

func MustLoadConfig() *Config {
	config := &Config{}

	if err := cleanenv.ReadConfig(".env", config); err != nil {
		panic(err)
	}

	return config
}
