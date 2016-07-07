package network

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/hashicorp/go-cleanhttp"
)

const userAgent = "go-networkapiclient"

// Client is the basic type of this pacGkage. It provides methods for
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

type Vlan struct{}
type ListVlansOptions struct{}

func (c *Client) getURL(path string) string {
	urlStr := strings.TrimRight(c.endpoint.String(), "/")
	return fmt.Sprintf("%s%s", urlStr, path)
}

func (c *Client) do(method, path string) (*http.Response, error) {
	httpClient := c.HTTPClient
	var u string
	u = c.getURL(path)
	// If the user has provided a timeout, apply it.
	if c.timeout != 0 {
		httpClient.Timeout = c.timeout
	}
	req, err := http.NewRequest(method, u, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent)
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		return nil, errors.New("error")
	}
	return resp, nil
}

func (c *Client) ListVlans(opts ListVlansOptions) ([]Vlan, error) {
	path := "/vlan/find"
	resp, err := c.do("GET", path)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var vlans []Vlan
	if err := json.NewDecoder(resp.Body).Decode(&vlans); err != nil {
		return nil, err
	}
	return vlans, nil
}
