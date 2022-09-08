package transport

import (
	"context"
	"encoding/gob"
	"github.com/KirillMironov/kappa/internal/kappa/domain"
	"github.com/KirillMironov/kappa/pkg/logger"
	"net"
	"sync"
)

type Handler struct {
	pods        <-chan []domain.Pod
	connections map[net.Conn]*gob.Encoder
	port        string
	logger      logger.Logger
	sync.Mutex
}

func NewHandler(pods <-chan []domain.Pod, port string, logger logger.Logger) *Handler {
	return &Handler{
		pods:        pods,
		connections: make(map[net.Conn]*gob.Encoder),
		port:        port,
		logger:      logger,
	}
}

func (h *Handler) Start(ctx context.Context) error {
	listener, err := net.Listen("tcp", ":"+h.port)
	if err != nil {
		return err
	}
	defer listener.Close()

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				h.logger.Errorf("failed to accept connection: %v", err)
				continue
			}

			h.Lock()
			h.connections[conn] = gob.NewEncoder(conn)
			h.Unlock()
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case pods := <-h.pods:
			h.Lock()

			for conn, encoder := range h.connections {
				err = encoder.Encode(pods)
				if err != nil {
					h.logger.Errorf("failed to encode data: %v", err)
					delete(h.connections, conn)
					conn.Close()
					continue
				}
			}

			h.Unlock()
		}
	}
}
