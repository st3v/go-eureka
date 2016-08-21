package eureka

import (
	"net"
	"net/http"
	"time"

	"github.com/st3v/go-eureka/retry"
)

var (
	DefaultHTTPClient = &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout:   5 * time.Second,
				KeepAlive: 60 * time.Second,
			}).Dial,
			TLSHandshakeTimeout:   5 * time.Second,
			ResponseHeaderTimeout: 5 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			MaxIdleConnsPerHost:   1,
		},
	}

	DefaultRetrySelector retry.Selector = retry.RoundRobin
	DefaultRetryLimit    retry.Allow    = retry.MaxRetries(3)
	DefaultRetryDelay    retry.Delay    = retry.ConstantDelay(1 * time.Second)
)

type Option func(*Client)

func HTTPClient(hc *http.Client) Option {
	return func(c *Client) {
		c.httpClient = hc
	}
}

func RetryLimit(limit retry.Allow) Option {
	return func(c *Client) {
		c.retryLimit = limit
	}
}

func RetrySelector(selector retry.Selector) Option {
	return func(c *Client) {
		c.retrySelector = selector
	}
}

func RetryDelay(delay retry.Delay) Option {
	return func(c *Client) {
		c.retryDelay = delay
	}
}
