package main

import (
	"context"
	"github.com/KirillMironov/kappa/internal/kappa/domain"
	"github.com/KirillMironov/kappa/internal/kappa/service"
	"github.com/KirillMironov/kappa/internal/kappa/transport"
	"github.com/sirupsen/logrus"
	"time"
)

func main() {
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: "01|02 15:04:05.000",
	})

	var (
		pods = make(chan []domain.Pod)

		parser  = service.Parser{}
		loader  = service.NewLoader(pods, ".", time.Second, parser, logger)
		handler = transport.NewHandler(pods, "20501", logger)
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go loader.Start(ctx)

	err := handler.Start(ctx)
	if err != nil {
		logger.Fatal(err)
	}
}
