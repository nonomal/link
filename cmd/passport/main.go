package main

import (
	"log"
	"net/url"
	"os"
	"sync"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("[ERRO] Usage: server|client|broker://linkAddr/targetAddr#http|https://authAddr/secretPath")
	}
	rawURL := os.Args[1]
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		log.Fatalf("[ERRO] URL: %v", err)
	}
	var whiteList sync.Map
	authSetup(parsedURL, &whiteList)
	coreSelect(parsedURL, rawURL, &whiteList)
}
