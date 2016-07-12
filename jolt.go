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

type client struct {
	endpoints []string
}

func NewClient(endpoints []string) *client {
	rand.Seed(time.Now().Unix())
	return &client{endpoints}
}

func (c *client) endpoint() string {
	return strings.TrimRight(c.endpoints[rand.Intn(len(c.endpoints))], " /")
}

func (c *client) appsURI(instance Instance) string {
	return fmt.Sprintf("%s/apps/%s", c.endpoint(), instance.App)
}

func (c *client) instanceURI(instance Instance) string {
	return fmt.Sprintf("%s/%s", c.appsURI(instance), instance.Id)
}

func (c *client) Register(instance Instance) error {
	data, err := xml.Marshal(instance)
	if err != nil {
		return err
	}

	uri := c.appsURI(instance)
	resp, err := http.Post(uri, "application/xml", bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("Unexpected response code %d", resp.StatusCode)
	}

	return nil
}

func (c *client) Deregister(instance Instance) error {
	uri := c.instanceURI(instance)
	req, err := http.NewRequest("DELETE", uri, nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Unexpected response code %d", resp.StatusCode)
	}

	return nil
}
