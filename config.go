package acd

import (
	"errors"
	"os"

	"gopkg.in/acd.v0/client"
	"gopkg.in/acd.v0/client/token"
)

var (
	// ErrConfigWrongPermissions is returned if the configuration file has permission different than 0600.
	ErrConfigWrongPermissions = errors.New("the config file should have 0600 permissions")
)

// NewClient initialize the token, creates a client and returns it to you.
func NewClient(configFile, cacheFile string) (*client.Client, error) {
	if err := validateConfigFile(configFile); err != nil {
		return nil, err
	}

	ts, err := token.New(configFile)
	if err != nil {
		return nil, err
	}

	c, err := client.New(ts, cacheFile)
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
		return ErrConfigWrongPermissions
	}

	return nil
}
