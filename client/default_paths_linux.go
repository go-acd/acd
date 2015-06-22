package client

import (
	"os"
	"path"
)

const tokenFilename = "acd-token.json"
const configFilename = "acd.json"
const cacheFilename = "com.appspot.go-acd.cache"

// DefaultConfigFile returns the default path for the configuration file.
func DefaultConfigFile() string {
	homePath := os.Getenv("HOME")
	return path.Join(homePath, ".config", configFilename)
}

// DefaultTokenFile returns the default path for the token file.
func DefaultTokenFile() string {
	homePath := os.Getenv("HOME")
	return path.Join(homePath, ".config", tokenFilename)
}

// DefaultCacheFile returns the default path for the cache file.
func DefaultCacheFile() string {
	homePath := os.Getenv("HOME")
	return path.Join(homePath, ".cache", cacheFilename)
}
