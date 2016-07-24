package eureka

import "net/http"

type Option func(*Client)

func HttpClient(hc *http.Client) Option {
	return func(c *Client) {
		c.httpClient = hc
	}
}
