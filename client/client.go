package client

import (
	"net/http"
	"time"

	"golang.org/x/oauth2"
	"gopkg.in/acd.v0/client/nodetree"
)

type (
	// Client provides a client for Amazon Cloud Drive.
	Client struct {
		NodeTree *nodetree.NodeTree

		httpClient  *http.Client
		timeout     time.Duration
		cacheFile   string
		metadataURL string
		contentURL  string
		tokenSource oauth2.TokenSource
	}

	endpointResponse struct {
		ContentURL     string `json:"contentUrl"`
		MetadataURL    string `json:"metadataUrl"`
		CustomerExists bool   `json:"customerExists"`
	}
)

var (
	// DefaultTimeout defines the default timeout of the client.
	DefaultTimeout = 10 * time.Hour
	endpointURL    = "https://drive.amazonaws.com/drive/v1/account/endpoint"
)

// New returns a new Amazon Cloud Drive "acd" Client.
func New(src oauth2.TokenSource, cacheFile string) (*Client, error) {
	c := &Client{
		tokenSource: src,
		timeout:     DefaultTimeout,
		cacheFile:   cacheFile,
		httpClient: &http.Client{
			Timeout: DefaultTimeout,
			Transport: &oauth2.Transport{
				Source: oauth2.ReuseTokenSource(nil, src),
			},
		},
	}

	if err := setEndpoints(c); err != nil {
		return nil, err
	}

	return c, nil
}

// GetTimeout returns the client's configured timeout.
func (c *Client) GetTimeout() time.Duration {
	return c.timeout
}

// SetTimeout configures the client's timeout.
func (c *Client) SetTimeout(t time.Duration) {
	c.timeout = t
}

// Close finalized the client.
func (c *Client) Close() error {
	return c.NodeTree.Close()
}

// Do invokes net/http.Client.Do(). Refer to net/http.Client.Do() for documentation.
func (c *Client) Do(r *http.Request) (*http.Response, error) {
	return c.httpClient.Do(r)
}
