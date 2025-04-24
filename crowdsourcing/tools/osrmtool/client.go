package osrmtool

import (
	"fmt"
	"io"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

type Client struct {
	HttpClient *http.Client
	Address    string
}

type Config struct {
	Address   string
	TimeoutMS int32
}

func (c *Client) Init(conf Config) {
	c.Address = conf.Address

	timeout := time.Duration(conf.TimeoutMS) * time.Millisecond

	c.HttpClient = &http.Client{
		Transport: &http.Transport{
			ResponseHeaderTimeout: timeout,
			TLSHandshakeTimeout:   timeout,
			IdleConnTimeout:       timeout,
			ExpectContinueTimeout: timeout,
		},
		Timeout: timeout,
	}
}

func (c *Client) HttpGet(address string) ([]byte, error) {
	defer func() {
		if err := recover(); err != nil {
			log.WithError(fmt.Errorf("%s", err)).WithField("address", address).Errorf("recovered from panic when calling httpGet: %s", err)
		}
	}()
	resp, err := c.HttpClient.Get(address)
	if err != nil {
		return []byte{}, err
	}

	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return []byte{}, err
	}

	return body, nil
}
