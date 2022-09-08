package transport

import (
	"encoding/gob"
	"github.com/KirillMironov/kappa/internal/kappa/domain"
	"github.com/KirillMironov/kappa/pkg/logger"
	"net"
)

type Handler struct {
	port   string
	logger logger.Logger
}

func NewHandler(port string, logger logger.Logger) *Handler {
	return &Handler{
		port:   port,
		logger: logger,
	}
}

func (h Handler) Start() error {
	conn, err := net.Dial("tcp", ":"+h.port)
	if err != nil {
		return err
	}
	defer conn.Close()

	var decoder = gob.NewDecoder(conn)

	for {
		var pods []domain.Pod

		err = decoder.Decode(&pods)
		if err != nil {
			h.logger.Errorf("failed to decode data: %v", err)
			return err
		}

		h.logger.Infof("received pods: %v", pods)
	}
}
