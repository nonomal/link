package main

import (
	"net/url"
	"os"
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
		helpInfo()
		os.Exit(1)
	}
}

func runServer(parsedURL *url.URL, rawURL string, whiteList *sync.Map) {
	log.Info("Server core selected: %v", strings.Split(rawURL, "#")[0])
	for {
		if err := tunnel.Server(parsedURL, whiteList); err != nil {
			log.Error("Server core error: %v", err)
			log.Info("Restarting in 1s...")
			time.Sleep(1 * time.Second)
		}
	}
}

func runClient(parsedURL *url.URL, rawURL string) {
	log.Info("Client core selected: %v", strings.Split(rawURL, "#")[0])
	for {
		if err := tunnel.Client(parsedURL); err != nil {
			log.Error("Client core error: %v", err)
			log.Info("Restarting in 1s...")
			time.Sleep(1 * time.Second)
		}
	}
}

func runBroker(parsedURL *url.URL, rawURL string, whiteList *sync.Map) {
	log.Info("Broker core selected: %v", strings.Split(rawURL, "#")[0])
	for {
		if err := forward.Broker(parsedURL, whiteList); err != nil {
			log.Error("Broker core error: %v", err)
			log.Info("Restarting in 1s...")
			time.Sleep(1 * time.Second)
		}
	}
}
