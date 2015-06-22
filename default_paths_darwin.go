package acd

import (
	"os"
	"path"
)

const tokenFilename = "acd-token.json"
const configFilename = "acd.json"
const cacheFilename = "com.appspot.go-acd.cache"

// DefaultConfigFile returns the default path for the configuration file. This is os-dependent setting.
func DefaultConfigFile() string {
	homePath := os.Getenv("HOME")
	return path.Join(homePath, ".config", configFilename)
}

// DefaultTokenFile returns the default path for the token file. This is os-dependent setting.
func DefaultTokenFile() string {
	homePath := os.Getenv("HOME")
	return path.Join(homePath, ".config", tokenFilename)
}

// DefaultCacheFile returns the default path for the cache file. This is os-dependent setting.
func DefaultCacheFile() string {
	homePath := os.Getenv("HOME")
	return path.Join(homePath, "Library", "Caches", cacheFilename)
}
