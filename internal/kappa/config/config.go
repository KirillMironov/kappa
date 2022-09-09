package config

import "github.com/kelseyhightower/envconfig"

type Config struct {
	Port    string `default:"20501" envconfig:"PORT"`
	PodsDir string `default:"." envconfig:"PODS_DIR"`
}

func Load() (cfg Config, _ error) {
	return cfg, envconfig.Process("", &cfg)
}
