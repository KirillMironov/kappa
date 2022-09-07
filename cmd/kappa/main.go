package main

import (
	"context"
	"github.com/KirillMironov/kappa/internal/kappa/core"
	"github.com/KirillMironov/kappa/internal/kappa/service"
	"github.com/KirillMironov/kappa/pkg/log"
	"time"
)

func main() {
	var (
		logger = log.New()

		pods = make(chan []core.Pod)

		parser = service.Parser{}
		loader = service.NewLoader(pods, ".", time.Second, parser, logger)
	)

	go func() {
		for v := range pods {
			logger.Info(v)
		}
	}()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	loader.Start(ctx)
}
