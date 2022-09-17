package config

import (
	"github.com/kelseyhightower/envconfig"
	"time"
)

type Config struct {
	Port            string        `default:"20501" envconfig:"PORT"`
	ShutdownTimeout time.Duration `default:"5s" envconfig:"SHUTDOWN_TIMEOUT"`
}

func Load() (cfg Config, _ error) {
	return cfg, envconfig.Process("", &cfg)
}
