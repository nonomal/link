package main

import (
	"net/url"
	"os"
	"sync"

	"github.com/yosebyte/passport/pkg/log"
)

var (
	version = "dev"
	whiteList sync.Map
)

func main() {
	if len(os.Args) < 2 {
		helpInfo()
		os.Exit(1)
	}
	rawURL := os.Args[1]
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		log.Fatal("Error parsing raw URL: %v", err)
	}
	authSetups(parsedURL, &whiteList)
	coreSelect(parsedURL, rawURL, &whiteList)
}
