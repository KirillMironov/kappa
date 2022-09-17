package main

import (
	"errors"
	"github.com/KirillMironov/kappa/internal/kappa/config"
	"github.com/KirillMironov/kappa/internal/kappa/service"
	"github.com/KirillMironov/kappa/internal/kappa/transport"
	"github.com/KirillMironov/kappa/pkg/httputil"
	"github.com/sirupsen/logrus"
	"net/http"
)

func main() {
	// Logger
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: "01|02 15:04:05.000",
	})

	// Config
	cfg, err := config.Load()
	if err != nil {
		logger.Fatal(err)
	}

	// DI
	var (
		deployer = service.NewDeployer(logger)
		handler  = transport.NewHandler(deployer)
	)

	// HTTP server
	httpServer := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: handler.Router(),
	}

	// Graceful shutdown
	go func() {
		err := httputil.GracefulShutdown(httpServer, cfg.ShutdownTimeout)
		if err != nil {
			logger.Fatal(err)
		}
	}()

	// Start HTTP server
	err = httpServer.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Fatal(err)
	}
}
