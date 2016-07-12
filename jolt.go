package jolt

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	endpoints  []string
	httpClient *http.Client
}

func init() {
	rand.Seed(time.Now().Unix())
}

type Option func(*Client)

func HttpClient(c *http.Client) Option {
	return func(client *Client) {
		client.httpClient = c
	}
}

func NewClient(endpoints []string, options ...Option) *Client {
	c := &Client{
		endpoints:  endpoints,
		httpClient: http.DefaultClient,
	}

	for _, opt := range options {
		opt(c)
	}

	return c
}

func (c *Client) endpoint() string {
	return strings.TrimRight(c.endpoints[rand.Intn(len(c.endpoints))], " /")
}

func (c *Client) appURI(instance Instance) string {
	return fmt.Sprintf("%s/apps/%s", c.endpoint(), instance.App)
}

func (c *Client) instanceURI(instance Instance) string {
	return fmt.Sprintf("%s/%s", c.appURI(instance), instance.Id)
}

func (c *Client) Register(instance Instance) error {
	data, err := xml.Marshal(instance)
	if err != nil {
		return err
	}

	return c.request("POST", c.appURI(instance), data, http.StatusNoContent)
}

func (c *Client) Deregister(instance Instance) error {
	return c.request("DELETE", c.instanceURI(instance), nil, http.StatusOK)
}

func (c *Client) Heartbeat(instance Instance) error {
	return c.request("PUT", c.instanceURI(instance), nil, http.StatusOK)
}

func (c *Client) request(method, uri string, body []byte, respCode int) error {
	req, err := http.NewRequest(method, uri, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/xml")
	req.Header.Add("Accept", "application/xml")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != respCode {
		return fmt.Errorf("Unexpected response code %d", resp.StatusCode)
	}

	return nil
}
