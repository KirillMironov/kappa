package config

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	BaseApiURL string `default:"http://localhost:20501" envconfig:"BASE_API_URL"`
}

func Load() (cfg Config, _ error) {
	return cfg, envconfig.Process("", &cfg)
}
