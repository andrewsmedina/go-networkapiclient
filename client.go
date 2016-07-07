package network

import (
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/hashicorp/go-cleanhttp"
)

// Client is the basic type of this package. It provides methods for
// interaction with the API.
type Client struct {
	HTTPClient *http.Client
	Dialer     *net.Dialer

	endpoint *url.URL

	// A timeout to use when using both the unixHTTPClient and HTTPClient
	timeout time.Duration
}

// NewClient returns a Client instance ready for communication with the given
// server endpoint. It will use the latest remote API version available in the
// server.
func NewClient(endpoint string) (*Client, error) {
	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}
	client := &Client{
		HTTPClient: cleanhttp.DefaultClient(),
		Dialer:     &net.Dialer{},
		endpoint:   u,
	}
	return client, nil
}
