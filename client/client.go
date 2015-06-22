package client

import (
	"encoding/json"
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

		config      *Config
		httpClient  *http.Client
		cacheFile   string
		metadataURL string
		contentURL  string
	}

	// Config represents the clients configuration.
	Config struct {
		TokenFile string        `json:"tokenFile"`
		CacheFile string        `json:"cacheFile"`
		Timeout   time.Duration `json:"timeout"`
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
func New(configFile string) (*Client, error) {
	config, err := loadConfig(configFile)
	if err != nil {
		return nil, err
	}

	ts, err := token.New(config.TokenFile)
	if err != nil {
		return nil, err
	}
	c := &Client{
		config:    config,
		cacheFile: config.CacheFile,
		httpClient: &http.Client{
			Timeout: config.Timeout,
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

// Close finalized the client.
func (c *Client) Close() error {
	return c.NodeTree.Close()
}

// Do invokes net/http.Client.Do(). Refer to net/http.Client.Do() for documentation.
func (c *Client) Do(r *http.Request) (*http.Response, error) {
	return c.httpClient.Do(r)
}

func validateFile(file string, checkPerms bool) error {
	stat, err := os.Stat(file)
	if err != nil {
		return err
	}
	if checkPerms && stat.Mode() != os.FileMode(0600) {
		log.Errorf("%s: want 0600 got %s", constants.ErrWrongPermissions, stat.Mode())
		return constants.ErrWrongPermissions
	}

	return nil
}

func loadConfig(configFile string) (*Config, error) {
	// validate the config file
	if err := validateFile(configFile, false); err != nil {
		return nil, err
	}

	cf, err := os.Open(configFile)
	if err != nil {
		log.Errorf("%s: %s", constants.ErrOpenFile, err)
		return nil, err
	}
	defer cf.Close()
	var config Config
	if err := json.NewDecoder(cf).Decode(&config); err != nil {
		log.Errorf("%s: %s", constants.ErrJSONDecoding, err)
		return nil, err
	}

	// validate the token file
	if err := validateFile(config.TokenFile, true); err != nil {
		return nil, err
	}

	return &config, nil
}
