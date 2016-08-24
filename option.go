package eureka

import (
	"crypto/tls"
	"net"
	"net/http"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
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

type httpClientOptions struct {
	timeout      time.Duration
	transport    *http.Transport
	oauth2Config *clientcredentials.Config
	tlsConfig    *tls.Config
}

func newDefaultHTTPClientOptions() *httpClientOptions {
	return &httpClientOptions{
		timeout:   DefaultTimeout,
		transport: DefaultTransport,
	}
}

func newHTTPClient(opts *httpClientOptions) *http.Client {
	if opts.tlsConfig != nil {
		opts.transport.TLSClientConfig = opts.tlsConfig
	}

	c := &http.Client{
		Timeout:   opts.timeout,
		Transport: opts.transport,
	}

	if opts.oauth2Config != nil {
		ctx := context.WithValue(context.Background(), oauth2.HTTPClient, c)
		c = opts.oauth2Config.Client(ctx)
	}

	return c
}

// HTTPTimeout sets the timeout for the internal HTTP client.
func HTTPTimeout(t time.Duration) Option {
	return func(c *Client) {
		c.httpClientOptions.timeout = t
	}
}

// HTTPTransport sets the transport for the internal HTTP client.
func HTTPTransport(t *http.Transport) Option {
	return func(c *Client) {
		c.httpClientOptions.transport = t
	}
}

// TLSConfig sets the TLS config for the internal HTTP client.
func TLSConfig(config *tls.Config) Option {
	return func(c *Client) {
		c.httpClientOptions.tlsConfig = config
	}
}

// Oauth2ClientCredentials instructs the internal http client to use the
// Oauth2 Client Credential flow to authenticate with the Eureka server.
func Oauth2ClientCredentials(clientID, clientSecret, tokenURI string, scopes ...string) Option {
	return func(c *Client) {
		c.httpClientOptions.oauth2Config = &clientcredentials.Config{
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
