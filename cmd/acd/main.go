package main

import (
	"flag"
	"fmt"

	"gopkg.in/acd.v0"
	"gopkg.in/acd.v0/internal/log"
)

var (
	configFile = flag.String("config-file", acd.DefaultConfigFile(), "The path of the configuration file.")
	cacheFile  = flag.String("cache-file", acd.DefaultCacheFile(), "The path of the cache file.")
	logLevel   = flag.Int("log-level", int(log.ErrorLevel), fmt.Sprintf("The log level: possible values: %s.", log.Levels()))
	timeout    = flag.Duration("timeout", 0, "HTTP Client timeout. A timeout of 0 means no timeout.")
)

func main() {
	flag.Parse()
	log.SetLevel(log.Level(*logLevel))
	c, err := acd.NewClient(*configFile, *timeout, *cacheFile)
	if err != nil {
		log.Fatal(err)
	}
	if err := c.FetchNodeTree(); err != nil {
		log.Fatal(err)
	}
	defer c.Close()
}
