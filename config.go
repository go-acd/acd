package acd

import (
	"os"
	"time"

	"gopkg.in/acd.v0/client"
	"gopkg.in/acd.v0/client/token"
	"gopkg.in/acd.v0/internal/constants"
	"gopkg.in/acd.v0/internal/log"
)

// NewClient initialize the token, creates a client and returns it to you. A
// timeout of 0 means no timeout.
func NewClient(configFile string, timeout time.Duration, cacheFile string) (*client.Client, error) {
	if err := validateConfigFile(configFile); err != nil {
		return nil, err
	}

	ts, err := token.New(configFile)
	if err != nil {
		return nil, err
	}

	c, err := client.New(ts, timeout, cacheFile)
	if err != nil {
		return nil, err
	}

	return c, nil
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
