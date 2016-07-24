package eureka

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/st3v/go-eureka/retry"
)

type Client struct {
	endpoints     []string
	httpClient    *http.Client
	retrySelector retry.Selector
	retryAllow    retry.Allow
	retryDelay    retry.Delay
}

var defaultHttpClient = &http.Client{
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

func NewClient(endpoints []string, options ...Option) *Client {
	for i, e := range endpoints {
		endpoints[i] = strings.TrimRight(e, " /")
	}

	c := &Client{
		endpoints:     endpoints,
		httpClient:    defaultHttpClient,
		retrySelector: retry.RoundRobin,
		retryAllow:    retry.Limit(3),
		retryDelay:    retry.Linear(1 * time.Second),
	}

	for _, opt := range options {
		opt(c)
	}

	return c
}

func (c *Client) Register(instance *Instance) error {
	data, err := xml.Marshal(instance)
	if err != nil {
		return err
	}

	return c.retry(c.do("POST", c.appPath(instance.AppName), data, http.StatusNoContent))
}

func (c *Client) Deregister(instance *Instance) error {
	return c.retry(c.do("DELETE", c.appInstancePath(instance.AppName, instance.Id), nil, http.StatusOK))
}

func (c *Client) Heartbeat(instance *Instance) error {
	return c.retry(c.do("PUT", c.appInstancePath(instance.AppName, instance.Id), nil, http.StatusOK))
}

func (c *Client) Apps() ([]*App, error) {
	result := new(Registry)
	if err := c.retry(c.get(c.appsPath(), result)); err != nil {
		return nil, err
	}

	return result.Apps, nil
}

func (c *Client) App(appName string) (*App, error) {
	app := new(App)
	err := c.retry(c.get(c.appPath(appName), app))
	return app, err
}

func (c *Client) AppInstance(appName, instanceId string) (*Instance, error) {
	instance := new(Instance)
	err := c.retry(c.get(c.appInstancePath(appName, instanceId), instance))
	return instance, err
}

func (c *Client) Instance(instanceId string) (*Instance, error) {
	instance := new(Instance)
	err := c.retry(c.get(c.instancePath(instanceId), instance))
	return instance, err
}

func (c *Client) StatusOverride(instance *Instance, status Status) error {
	return c.retry(c.do("PUT", c.appInstanceStatusPath(instance.AppName, instance.Id, status), nil, http.StatusOK))
}

func (c *Client) RemoveStatusOverride(instance *Instance, fallback Status) error {
	return c.retry(c.do("DELETE", c.appInstanceStatusPath(instance.AppName, instance.Id, fallback), nil, http.StatusOK))
}

func (c *Client) retry(action retry.Action) error {
	return retry.NewStrategy(
		c.retrySelector(c.endpoints),
		c.retryAllow,
		c.retryDelay,
	).Apply(action)
}

func (c *Client) do(method, path string, body []byte, respCode int) retry.Action {
	return func(endpoint string) error {
		req, err := http.NewRequest(method, fmt.Sprintf("%s/%s", endpoint, path), bytes.NewBuffer(body))
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
}

func (c *Client) get(path string, result interface{}) retry.Action {
	return func(endpoint string) error {
		req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s", endpoint, path), nil)
		if err != nil {
			return err
		}

		req.Header.Add("Accept", "application/xml")

		resp, err := c.httpClient.Do(req)
		if err != nil {
			return err
		}

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("Unexpected response code %d", resp.StatusCode)
		}

		defer resp.Body.Close()
		if err := xml.NewDecoder(resp.Body).Decode(result); err != nil {
			return err
		}

		return nil
	}
}

func (c *Client) appsPath() string {
	return "apps"
}

func (c *Client) appPath(appName string) string {
	return fmt.Sprintf("%s/%s", c.appsPath(), appName)
}

func (c *Client) appInstancePath(appName, instanceId string) string {
	return fmt.Sprintf("%s/%s", c.appPath(appName), instanceId)
}

func (c *Client) instancePath(instanceId string) string {
	return fmt.Sprintf("instances/%s", instanceId)
}

func (c *Client) appInstanceStatusPath(appName, instanceId string, status Status) string {
	return fmt.Sprintf("%s/status?value=%s", c.appInstancePath(appName, instanceId), status)
}
