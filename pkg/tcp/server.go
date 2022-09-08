package tcp

import (
	"encoding/gob"
	"net"
	"sync"
)

type Server[T any] struct {
	host string
	port string

	listener net.Listener
	clients  map[net.Conn]*gob.Encoder
	mu       sync.Mutex
}

func NewServer[T any](host string, port string) (*Server[T], error) {
	listener, err := net.Listen("tcp", net.JoinHostPort(host, port))
	if err != nil {
		return nil, err
	}

	server := &Server[T]{
		host:     host,
		port:     port,
		listener: listener,
		clients:  make(map[net.Conn]*gob.Encoder),
	}

	go server.listen()

	return server, nil
}

func (s *Server[T]) Send(obj T) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for conn, encoder := range s.clients {
		err := encoder.Encode(obj)
		if err != nil {
			delete(s.clients, conn)
			conn.Close()
			continue
		}
	}
}

func (s *Server[T]) Close() error {
	if s.listener != nil {
		return s.listener.Close()
	}
	return nil
}

func (s *Server[T]) listen() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			continue
		}

		s.mu.Lock()
		s.clients[conn] = gob.NewEncoder(conn)
		s.mu.Unlock()
	}
}
