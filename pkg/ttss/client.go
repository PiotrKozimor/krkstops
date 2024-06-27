package ttss

import (
	"net/http"
	"time"
)

const (
	Bus  = "http://ttss.mpk.krakow.pl"
	Tram = "http://www.ttss.krakow.pl"
)

type Departure struct {
	PatternText  string
	Direction    string
	PlannedTime  string
	RelativeTime int32
	Predicted    bool
}

type Stop struct {
	Name string
	Id   uint
}

type Client struct {
	host       string
	httpClient *http.Client
}

type opt func(*Client)

func WithTimeout(timeout time.Duration) opt {
	return func(cli *Client) {
		cli.httpClient.Timeout = timeout
	}
}

func NewClient(url string, options ...opt) *Client {
	c := &Client{
		host: url,
		httpClient: &http.Client{
			Timeout: time.Second,
		},
	}
	for _, o := range options {
		o(c)
	}
	return c
}
