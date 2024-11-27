package main

import (
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/yosebyte/passport/internal/forward"
	"github.com/yosebyte/passport/internal/tunnel"
	"github.com/yosebyte/passport/pkg/log"
)

func coreSelect(parsedURL *url.URL, rawURL string, whiteList *sync.Map) {
	switch parsedURL.Scheme {
	case "server":
		runServer(parsedURL, rawURL, whiteList)
	case "client":
		runClient(parsedURL, rawURL)
	case "broker":
		runBroker(parsedURL, rawURL, whiteList)
	default:
		log.Fatal("Invalid running core: use server|client|broker://")
	}
}

func runServer(parsedURL *url.URL, rawURL string, whiteList *sync.Map) {
	log.Info("Server mode enabled: %v", strings.Split(rawURL, "#")[0])
	for {
		if err := tunnel.Server(parsedURL, whiteList); err != nil {
			log.Error("Server core error: %v", err)
			log.Info("Restarting in 1s...")
			time.Sleep(1 * time.Second)
		}
	}
}

func runClient(parsedURL *url.URL, rawURL string) {
	log.Info("Client mode enabled: %v", strings.Split(rawURL, "#")[0])
	for {
		if err := tunnel.Client(parsedURL); err != nil {
			log.Error("Client core error: %v", err)
			log.Info("Restarting in 1s...")
			time.Sleep(1 * time.Second)
		}
	}
}

func runBroker(parsedURL *url.URL, rawURL string, whiteList *sync.Map) {
	log.Info("Broker mode enabled: %v", strings.Split(rawURL, "#")[0])
	for {
		if err := forward.Broker(parsedURL, whiteList); err != nil {
			log.Error("Broker core error: %v", err)
			log.Info("Restarting in 1s...")
			time.Sleep(1 * time.Second)
		}
	}
}
