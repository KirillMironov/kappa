package httputil

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient(baseURL string, timeout time.Duration) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

func (c *Client) Do(method string, path string, body io.Reader, expectedStatus int) error {
	joinedURL, err := url.JoinPath(c.baseURL, path)
	if err != nil {
		return &httpError{Internal: err}
	}

	req, err := http.NewRequest(method, joinedURL, body)
	if err != nil {
		return &httpError{Internal: err}
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return &httpError{Internal: err}
	}
	defer resp.Body.Close()

	if resp.StatusCode != expectedStatus {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return &httpError{Internal: err}
		}
		return &httpError{
			Code:    resp.StatusCode,
			Message: string(bodyBytes),
		}
	}

	return nil
}

type httpError struct {
	Code     int
	Message  string
	Internal error
}

func (he *httpError) Error() string {
	switch {
	case he.Internal != nil:
		return he.Internal.Error()
	case he.Message != "":
		return fmt.Sprintf("status: %s (%d), message: %s", http.StatusText(he.Code), he.Code, he.Message)
	default:
		return fmt.Sprintf("status: %s (%d)", http.StatusText(he.Code), he.Code)
	}
}
