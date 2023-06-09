package airly

import (
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

const URL = "https://airapi.airly.eu"

type Measurement struct {
	Caqi        int32
	Humidity    int32
	Temperature float32
	Color       uint32
}

type Installation struct {
	Id        int32
	Latitude  float32
	Longitude float32
}

type Client struct {
	host   string
	apiKey string
	cl     *http.Client
}

func NewClient(host string) *Client {
	return &Client{
		host:   host,
		apiKey: os.Getenv("AIRLY_KEY"),
		cl:     http.DefaultClient,
	}
}

func (c *Client) do(path string, modify func(*http.Request)) (io.Reader, func(), error) {
	req, err := http.NewRequest("GET", strings.Join([]string{c.host, path}, "/"), nil)
	if err != nil {
		return nil, nil, fmt.Errorf("new request: %w", err)
	}
	req.Header.Add("apikey", c.apiKey)
	req.Header.Add("Accept-Encoding", "gzip")

	if modify != nil {
		modify(req)
	}

	resp, err := c.cl.Do(req)
	if err != nil {
		return nil, nil, err
	}
	if resp.StatusCode != 200 {
		return nil, nil, fmt.Errorf("status code: %d", resp.StatusCode)
	}

	gzipR, err := gzip.NewReader(resp.Body)
	if resp.StatusCode != 200 {
		return nil, nil, fmt.Errorf("gzip reader: %d", err)
	}

	return gzipR, func() {
		resp.Body.Close()
		gzipR.Close()
	}, nil
}
