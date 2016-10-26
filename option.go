package eureka

import (
	"crypto/tls"
	"net"
	"net/http"
	"time"

	"golang.org/x/oauth2/clientcredentials"

	"github.com/st3v/go-eureka/retry"
)

var (
	DefaultRetrySelector retry.Selector = retry.RoundRobin
	DefaultRetryLimit    retry.Allow    = retry.MaxRetries(3)
	DefaultRetryDelay    retry.Delay    = retry.ConstantDelay(1 * time.Second)

	DefaultTransport = &http.Transport{
		Dial: (&net.Dialer{
			Timeout:   5 * time.Second,
			KeepAlive: 60 * time.Second,
		}).Dial,
		TLSHandshakeTimeout:   5 * time.Second,
		ResponseHeaderTimeout: 5 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		MaxIdleConnsPerHost:   1,
	}

	DefaultTimeout = 10 * time.Second
)

type Option func(*Client)

// HTTPTimeout sets the timeout for the internal HTTP client.
func HTTPTimeout(t time.Duration) Option {
	return func(c *Client) {
		c.timeout = t
	}
}

// HTTPTransport sets the transport for the internal HTTP client.
func HTTPTransport(t *http.Transport) Option {
	return func(c *Client) {
		c.transport = t
	}
}

// TLSConfig sets the TLS config for the internal HTTP client.
func TLSConfig(config *tls.Config) Option {
	return func(c *Client) {
		c.tlsConfig = config
	}
}

// Oauth2ClientCredentials instructs the internal http client to use the
// Oauth2 Client Credential flow to authenticate with the Eureka server.
func Oauth2ClientCredentials(clientID, clientSecret, tokenURI string, scopes ...string) Option {
	return func(c *Client) {
		c.oauth2Config = &clientcredentials.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			TokenURL:     tokenURI,
			Scopes:       scopes,
		}
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
