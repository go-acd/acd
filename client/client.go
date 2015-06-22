package client

import (
	"net/http"
	"os"
	"time"

	"golang.org/x/oauth2"
	"gopkg.in/acd.v0/client/node"
	"gopkg.in/acd.v0/client/token"
	"gopkg.in/acd.v0/internal/constants"
	"gopkg.in/acd.v0/internal/log"
)

type (
	// Client provides a client for Amazon Cloud Drive.
	Client struct {
		NodeTree *node.Tree

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

const (
	endpointURL = "https://drive.amazonaws.com/drive/v1/account/endpoint"
)

// New returns a new Amazon Cloud Drive "acd" Client. A timeout of 0 means no timeout.
func New(configFile, cacheFile string, timeout time.Duration) (*Client, error) {
	if err := validateConfigFile(configFile); err != nil {
		return nil, err
	}
	ts, err := token.New(configFile)
	if err != nil {
		return nil, err
	}
	c := &Client{
		tokenSource: ts,
		timeout:     timeout,
		cacheFile:   cacheFile,
		httpClient: &http.Client{
			Timeout: timeout,
			Transport: &oauth2.Transport{
				Source: oauth2.ReuseTokenSource(nil, ts),
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

func validateConfigFile(configFile string) error {
	stat, err := os.Stat(configFile)
	if err != nil {
		return err
	}
	if stat.Mode() != os.FileMode(0600) {
		log.Errorf("%s: want 0600 got %s", constants.ErrWrongPermissions, stat.Mode())
		return constants.ErrWrongPermissions
	}

	return nil
}
