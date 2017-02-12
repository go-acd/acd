package acd

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"golang.org/x/oauth2"
	"gopkg.in/acd.v0/internal/constants"
	"gopkg.in/acd.v0/internal/log"
	"gopkg.in/acd.v0/node"
	"gopkg.in/acd.v0/token"
)

type (
	// Config represents the clients configuration.
	Config struct {
		// TokenFile represents the file containing the oauth settings which must
		// be present on disk and has permissions 0600. The file is used by the
		// token package to produce a valid access token by calling the oauthServer
		// with the refresh token.  The default oauth server is hosted at
		// https://go-acd.appspot.com with the source code available at
		// https://github.com/go-acd/oauth-server.  It's currently not possible to
		// change the oauth-server. Please feel free to add this feature if you
		// have a use-case for it.
		TokenFile string `json:"tokenFile"`

		// CacheFile represents the file used by the client to cache the NodeTree.
		// This file is not assumed to be present and will be created on the first
		// run. It is gob-encoded node.Node.
		CacheFile string `json:"cacheFile"`

		// Timeout configures the HTTP Client with a timeout after which the client
		// will cancel the request and return. A timeout of 0 (the default) means
		// no timeout. See http://godoc.org/net/http#Client for more information.
		Timeout time.Duration `json:"timeout"`

		// Oauth2RefreshURL OAuth2 token server
		Oauth2RefreshURL string `json:"oauth2RefreshURL"`
	}

	// Client provides a client for Amazon Cloud Drive.
	Client struct {
		// NodeTree is the tree of nodes as stored on the drive. This tree should
		// be fetched using (*Client).FetchNodeTree() as soon the client is
		// created.
		NodeTree *node.Tree

		config      *Config
		httpClient  *http.Client
		cacheFile   string
		metadataURL string
		contentURL  string
	}

	endpointResponse struct {
		ContentURL     string `json:"contentUrl"`
		MetadataURL    string `json:"metadataUrl"`
		CustomerExists bool   `json:"customerExists"`
	}
)

const (
	endpointURL = "https://drive.amazonaws.com/drive/v1/account/endpoint"
	defaultOauth2RefreshUrl = "https://go-acd.appspot.com/refresh"
)

// New returns a new Amazon Cloud Drive "acd" Client. configFile must exist and must be a valid JSON decodable into Config.
func New(configFile string) (*Client, error) {
	config, err := loadConfig(configFile)
	if err != nil {
		return nil, err
	}

	ts, err := token.New(config.Oauth2RefreshURL, config.TokenFile)
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

// Close finalizes the acd.
func (c *Client) Close() error {
	return c.NodeTree.Close()
}

// Do invokes net/http.Client.Do(). Refer to net/http.Client.Do() for documentation.
func (c *Client) Do(r *http.Request) (*http.Response, error) {
	return c.httpClient.Do(r)
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

	if config.Oauth2RefreshURL == "" {
		config.Oauth2RefreshURL = defaultOauth2RefreshUrl
	}

	// validate the token file
	if err := validateFile(config.TokenFile, true); err != nil {
		return nil, err
	}

	return &config, nil
}

func validateFile(file string, checkPerms bool) error {
	stat, err := os.Stat(file)
	if err != nil {
		if os.IsNotExist(err) {
			log.Errorf("%s: %s -- %s", constants.ErrFileNotFound, err, file)
			return constants.ErrFileNotFound
		}
		log.Errorf("%s: %s -- %s", constants.ErrStatFile, err, file)
		return constants.ErrStatFile
	}
	if checkPerms && stat.Mode() != os.FileMode(0600) {
		log.Errorf("%s: want 0600 got %s", constants.ErrWrongPermissions, stat.Mode())
		return constants.ErrWrongPermissions
	}

	return nil
}
