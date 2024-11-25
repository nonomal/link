package main

import (
	"log"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/yosebyte/passport/internal/forward"
	"github.com/yosebyte/passport/internal/tunnel"
)

func coreSelect(parsedURL *url.URL, rawURL string, whiteList *sync.Map) {
	for {
		switch parsedURL.Scheme {
		case "server":
			runServer(parsedURL, rawURL, whiteList)
		case "client":
			runClient(parsedURL, rawURL)
		case "broker":
			runBroker(parsedURL, rawURL, whiteList)
		case "shadow":
			runShadow(parsedURL, rawURL, whiteList)
		default:
			log.Fatalf("[ERRO] Usage: server|client|broker://linkAddr/targetAddr#http|https://authAddr/secretPath")
		}
	}
}

func runServer(parsedURL *url.URL, rawURL string, whiteList *sync.Map) {
	log.Printf("[INFO] Server: %v", strings.Split(rawURL, "#")[0])
	if err := tunnel.Server(parsedURL, whiteList); err != nil {
		log.Printf("[ERRO] Server: %v", err)
		time.Sleep(1 * time.Second)
	}
}

func runClient(parsedURL *url.URL, rawURL string) {
	log.Printf("[INFO] Client: %v", strings.Split(rawURL, "#")[0])
	if err := tunnel.Client(parsedURL); err != nil {
		log.Printf("[ERRO] Client: %v", err)
		time.Sleep(1 * time.Second)
	}
}

func runBroker(parsedURL *url.URL, rawURL string, whiteList *sync.Map) {
	log.Printf("[INFO] Broker: %v", strings.Split(rawURL, "#")[0])
	if err := forward.Broker(parsedURL, whiteList); err != nil {
		log.Printf("[ERRO] Broker: %v", err)
		time.Sleep(1 * time.Second)
	}
}

func runShadow(parsedURL *url.URL, rawURL string, whiteList *sync.Map) {
	log.Printf("[INFO] Shadow: %v", strings.Split(rawURL, "#")[0])
	if err := forward.Shadow(parsedURL, whiteList); err != nil {
		log.Printf("[ERRO] Shadow: %v", err)
		time.Sleep(1 * time.Second)
	}
}
