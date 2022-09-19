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
		return &Error{Internal: err}
	}

	req, err := http.NewRequest(method, joinedURL, body)
	if err != nil {
		return &Error{Internal: err}
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return &Error{Internal: err}
	}
	defer resp.Body.Close()

	if resp.StatusCode != expectedStatus {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return &Error{Internal: err}
		}
		return &Error{
			Code:    resp.StatusCode,
			Message: string(bodyBytes),
		}
	}

	return nil
}

type Error struct {
	Code     int
	Message  string
	Internal error
}

func (e *Error) Error() string {
	switch {
	case e.Internal != nil:
		return e.Internal.Error()
	case e.Message != "":
		return fmt.Sprintf("status: %s (%d), message: %s", http.StatusText(e.Code), e.Code, e.Message)
	default:
		return fmt.Sprintf("status: %s (%d)", http.StatusText(e.Code), e.Code)
	}
}
