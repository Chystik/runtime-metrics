package httpclient

import (
	"net/http"
	"time"
)

const (
	defaultTimeout = 20 * time.Second
)

type Client struct {
	client  *http.Client
	timeout time.Duration
}

func NewClient(opts ...Options) *Client {
	httpClietn := &http.Client{
		Timeout: defaultTimeout,
	}

	client := &Client{client: httpClietn}

	for _, opt := range opts {
		opt(client)
	}

	return client
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	return c.client.Do(req)
}
