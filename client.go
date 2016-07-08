package network

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/hashicorp/go-cleanhttp"
)

// Client is the basic type of this pacGkage. It provides methods for
// interaction with the API.
type Client struct {
	HTTPClient *http.Client
	Dialer     *net.Dialer

	endpoint *url.URL

	// A timeout to use when using both the unixHTTPClient and HTTPClient
	timeout  time.Duration
	user     string
	password string
}

// NewClient returns a Client instance ready for communication with the given
// server endpoint. It will use the latest remote API version available in the
// server.
func NewClient(endpoint, user, password string) (*Client, error) {
	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}
	client := &Client{
		HTTPClient: cleanhttp.DefaultClient(),
		Dialer:     &net.Dialer{},
		endpoint:   u,
		user:       user,
		password:   password,
	}
	return client, nil
}

// ListVlanResult represents a Vlan list
type ListVlanResult struct {
	Vlans []Vlan `xml:"vlan"`
}

// NetworkIPV4 represents a network
type NetworkIPV4 struct {
	Network string `xml:"network"`
}

// Vlan represents a vlan
type Vlan struct {
	Environment int         `xml:"ambiente"`
	Number      int         `xml:"num_vlan"`
	NetworkIPV4 NetworkIPV4 `xml:"redeipv4"`
}

// ListVlansOptions is used to define the ListVlans parameters
type ListVlansOptions struct {
	Name string
}

func (c *Client) getURL(path string) string {
	urlStr := strings.TrimRight(c.endpoint.String(), "/")
	return fmt.Sprintf("%s%s", urlStr, path)
}

func (c *Client) do(method, path string, headers map[string]string, data io.Reader) (*http.Response, error) {
	httpClient := c.HTTPClient
	var u string
	u = c.getURL(path)
	// If the user has provided a timeout, apply it.
	if c.timeout != 0 {
		httpClient.Timeout = c.timeout
	}
	req, err := http.NewRequest(method, u, data)
	if err != nil {
		return nil, err
	}
	req.Header.Set("NETWORKAPI_USERNAME", c.user)
	req.Header.Set("NETWORKAPI_PASSWORD", c.password)
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		return nil, errors.New("error")
	}
	return resp, nil
}

// ListVlans returns a []Vlan and is filtered by ListVlansOptions
func (c *Client) ListVlans(opts ListVlansOptions) ([]Vlan, error) {
	path := "/vlan/find/"
	headers := map[string]string{
		"Content-Type": "text/plain",
	}
	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<networkapi versao="1.0">
<vlan>
<exato>False</exato>
<subrede>0</subrede>
<start_record>0</start_record>
<tipo_rede/>
<nome>%s</nome>
<custom_search/>
<numero/>
<ambiente/>
<versao>0</versao>
<end_record>100</end_record>
<rede/>
<asorting_cols/>
<acl/>
<searchable_columns/>
</vlan>
</networkapi>`
	data := strings.NewReader(fmt.Sprintf(xmlData, opts.Name))
	resp, err := c.do("POST", path, headers, data)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var result ListVlanResult
	if err := xml.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result.Vlans, nil
}
