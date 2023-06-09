package ttss

const (
	Bus  = "http://91.223.13.70"
	Tram = "http://185.70.182.51"
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
	host string
}

func NewClient(url string) *Client {
	return &Client{host: url}
}
