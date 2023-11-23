package httpclient

import "time"

type Options func(*Client)

func Timeout(t time.Duration) Options {
	return func(c *Client) {
		c.timeout = t
	}
}
