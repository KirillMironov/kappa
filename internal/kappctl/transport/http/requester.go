package http

import (
	"io"
	"net"
	"net/http"
	"net/url"
	"time"
)

type Requester struct {
	client  *http.Client
	baseURL string
}

func NewRequester(host, port string, timeout time.Duration) Requester {
	return Requester{
		client: &http.Client{
			Timeout: timeout,
		},
		baseURL: net.JoinHostPort(host, port),
	}
}

func (r Requester) Do(method string, path string, body io.Reader) (respBody string, err error) {
	joinedURL, err := url.JoinPath("http://", r.baseURL, path)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(method, joinedURL, body)
	if err != nil {
		return "", err
	}

	resp, err := r.client.Do(req)
	if err != nil {
		return "", err
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(bodyBytes), nil
}
