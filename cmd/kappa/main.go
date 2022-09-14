package main

import (
	"context"
	"errors"
	"github.com/KirillMironov/kappa/internal/kappa/config"
	"github.com/KirillMironov/kappa/internal/kappa/service"
	"github.com/KirillMironov/kappa/internal/kappa/transport"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: "01|02 15:04:05.000",
	})

	cfg, err := config.Load()
	if err != nil {
		logger.Fatal(err)
	}

	var (
		deployer = service.NewDeployer(logger)
		handler  = transport.NewHandler(deployer)
	)

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: handler.InitRoutes(),
	}

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
		<-quit

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		err := srv.Shutdown(ctx)
		if err != nil {
			logger.Fatal(err)
		}
	}()

	err = srv.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Fatal(err)
	}
}
