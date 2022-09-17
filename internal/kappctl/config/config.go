package config

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Host string `default:"localhost" envconfig:"HOST"`
	Port string `default:"20501" envconfig:"PORT"`
}

func Load() (cfg Config, _ error) {
	return cfg, envconfig.Process("", &cfg)
}
