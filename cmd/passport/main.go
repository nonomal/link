package main

import (
	"net/url"
	"os"
	"sync"

	"github.com/yosebyte/passport/pkg/log"
)

func main() {
	if len(os.Args) < 2 {
		log.Info("Usage: server|client|broker://link/target#http|https://auth/path")
		os.Exit(1)
	}
	rawURL := os.Args[1]
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		log.Fatal("Error parsing core URL: %v", err)
	}
	var whiteList sync.Map
	authSetups(parsedURL, &whiteList)
	coreSelect(parsedURL, rawURL, &whiteList)
}
