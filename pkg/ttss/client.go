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

func NewClient(url string) *Client {
	return &Client{
		host: url,
		httpClient: &http.Client{
			Timeout: time.Second,
		},
	}
}
