package main

import (
	"github.com/KirillMironov/kappa/internal/kappctl/config"
	"github.com/KirillMironov/kappa/internal/kappctl/transport/cmd"
	"github.com/KirillMironov/kappa/internal/kappctl/transport/http"
	"log"
	"time"
)

func main() {
	// Config
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	// DI
	var (
		requester = http.NewRequester(cfg.Host, cfg.Port, time.Second*3)
		command   = cmd.NewCmd(requester)
	)

	// Cobra
	command.Execute()
}
