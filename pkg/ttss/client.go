package ttss

import (
	"crypto/tls"
	"crypto/x509"
	_ "embed"
	"log"
	"net/http"
	"time"
)

const (
	Bus  = "https://ttss.mpk.krakow.pl"
	Tram = "https://www.ttss.krakow.pl"
)

var (
	// https://www.certum.pl/CTNCA.pem
	//go:embed CTNCA.pem
	ca []byte
	// https://repository.certum.pl/ovcasha2.pem
	//go:embed ovcasha2.pem
	intermediate []byte
	client       *http.Client
)

func init() {
	caCertPool := x509.NewCertPool()
	for _, cert := range [][]byte{
		ca,
		intermediate,
	} {
		if ok := caCertPool.AppendCertsFromPEM(cert); !ok {
			log.Fatal("Failed to add certificate to pool")
			return
		}
	}
	client = &http.Client{
		Timeout: time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: caCertPool,
			},
		},
	}
}

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
		host:       url,
		httpClient: client,
	}
	for _, o := range options {
		o(c)
	}
	return c
}
