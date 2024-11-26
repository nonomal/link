package main

import (
	"net/url"
	"os"
	"sync"

	"github.com/yosebyte/passport/pkg/log"
)

func main() {
	if len(os.Args) < 2 {
		log.Info("Usage: server|client|broker://linkAddr/targetAddr#http|https://authAddr/secretPath")
		os.Exit(1)
	}
	rawURL := os.Args[1]
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		log.Fatal("Error parsing core URL: %v", err)
	}
	var whiteList sync.Map
	authSetup(parsedURL, &whiteList)
	coreSelect(parsedURL, rawURL, &whiteList)
}
