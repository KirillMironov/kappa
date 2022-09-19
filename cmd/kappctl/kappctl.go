package main

import (
	"github.com/KirillMironov/kappa/internal/kappctl/config"
	"github.com/KirillMironov/kappa/internal/kappctl/transport"
	"log"
)

func main() {
	// Config
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	// DI
	var command = transport.NewCmd(cfg.BaseApiURL)

	// Cobra
	command.Execute()
}
