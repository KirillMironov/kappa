package tcp

import (
	"encoding/gob"
	"errors"
	"io"
	"net"
	"time"
)

var ErrReconnectAttemptsExceeded = errors.New("reconnect attempts exceeded")

type Client[T any] struct {
	host string
	port string

	reconnectInterval time.Duration
	reconnectAttempts int

	conn    net.Conn
	decoder *gob.Decoder
}

func NewClient[T any](host string, port string) (*Client[T], error) {
	client := &Client[T]{
		host: host,
		port: port,
	}
	return client, client.connect()
}

func (c *Client[T]) Receive() (<-chan T, <-chan error) {
	resultCh := make(chan T)
	errorCh := make(chan error)

	go func() {
		defer close(resultCh)
		defer close(errorCh)

		for {
			var objects T

			err := c.decoder.Decode(&objects)
			if err != nil {
				if errors.Is(err, io.EOF) {
					err = c.reconnect()
					if err != nil {
						errorCh <- err
						return
					}
					continue
				}
				errorCh <- err
				return
			}

			resultCh <- objects
		}
	}()

	return resultCh, errorCh
}

func (c *Client[T]) SetReconnectInterval(d time.Duration) {
	c.reconnectInterval = d
}

func (c *Client[T]) SetReconnectAttempts(n int) {
	c.reconnectAttempts = n
}

func (c *Client[T]) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

func (c *Client[T]) connect() error {
	conn, err := net.Dial("tcp", net.JoinHostPort(c.host, c.port))
	if err != nil {
		return err
	}

	c.conn = conn
	c.decoder = gob.NewDecoder(conn)

	return nil
}

func (c *Client[T]) reconnect() error {
	c.Close()

	for i := 0; i < c.reconnectAttempts; i++ {
		time.Sleep(c.reconnectInterval)

		err := c.connect()
		if err == nil {
			return nil
		}
	}

	return ErrReconnectAttemptsExceeded
}
