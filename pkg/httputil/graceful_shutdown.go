package httputil

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func GracefulShutdown(httpServer *http.Server, timeout time.Duration) error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return httpServer.Shutdown(ctx)
}
