package transport

import (
	"context"
	"github.com/KirillMironov/kappa/internal/kappa/domain"
	"github.com/KirillMironov/kappa/pkg/logger"
	"github.com/KirillMironov/kappa/pkg/tcp"
)

type Handler struct {
	pods   <-chan []domain.Pod
	port   string
	logger logger.Logger
}

func NewHandler(pods <-chan []domain.Pod, port string, logger logger.Logger) *Handler {
	return &Handler{
		pods:   pods,
		port:   port,
		logger: logger,
	}
}

func (h Handler) Start(ctx context.Context) error {
	server, err := tcp.NewServer[[]domain.Pod]("", h.port)
	if err != nil {
		return err
	}
	defer server.Close()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case pods := <-h.pods:
			server.Send(pods)
		}
	}
}
