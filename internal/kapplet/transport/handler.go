package transport

import (
	"context"
	"github.com/KirillMironov/kappa/internal/kappa/domain"
	"github.com/KirillMironov/kappa/pkg/logger"
	"github.com/KirillMironov/kappa/pkg/tcp"
	"time"
)

type Handler struct {
	pods              chan<- []domain.Pod
	port              string
	reconnectInterval time.Duration
	reconnectAttempts int
	logger            logger.Logger
}

func NewHandler(pods chan<- []domain.Pod, port string, logger logger.Logger) *Handler {
	return &Handler{
		pods:              pods,
		port:              port,
		reconnectInterval: time.Second * 3,
		reconnectAttempts: 3,
		logger:            logger,
	}
}

func (h Handler) Start(ctx context.Context) error {
	client, err := tcp.NewClient[[]domain.Pod]("", h.port)
	if err != nil {
		return err
	}
	defer client.Close()

	client.SetReconnectInterval(h.reconnectInterval)
	client.SetReconnectAttempts(h.reconnectAttempts)

	podsCh, errorCh := client.Receive()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case pods := <-podsCh:
			h.pods <- pods
		case err = <-errorCh:
			h.logger.Error(err)
			return err
		}
	}
}
