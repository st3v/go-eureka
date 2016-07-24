package eureka

import (
	"net/http"

	"github.com/st3v/go-eureka/retry"
)

type Option func(*Client)

func HttpClient(hc *http.Client) Option {
	return func(c *Client) {
		c.httpClient = hc
	}
}

func Retry(selector retry.Selector, allow retry.Allow, delay retry.Delay) Option {
	return func(c *Client) {
		c.retrySelector = selector
		c.retryAllow = allow
		c.retryDelay = delay
	}
}
