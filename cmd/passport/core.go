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
	for {
		switch parsedURL.Scheme {
		case "server":
			runServer(parsedURL, rawURL, whiteList)
		case "client":
			runClient(parsedURL, rawURL)
		case "broker":
			runBroker(parsedURL, rawURL, whiteList)
		default:
			log.Fatal("Invalid running core: use server, client or broker before ://")
		}
	}
}

func runServer(parsedURL *url.URL, rawURL string, whiteList *sync.Map) {
	log.Info("Server core enabled: %v", strings.Split(rawURL, "#")[0])
	if err := tunnel.Server(parsedURL, whiteList); err != nil {
		log.Error("Server core error: %v Restarting in 1s...", err)
		time.Sleep(1 * time.Second)
	}
}

func runClient(parsedURL *url.URL, rawURL string) {
	log.Info("Client core enabled: %v", strings.Split(rawURL, "#")[0])
	if err := tunnel.Client(parsedURL); err != nil {
		log.Error("Client core error: %v Restarting in 1s...", err)
		time.Sleep(1 * time.Second)
	}
}

func runBroker(parsedURL *url.URL, rawURL string, whiteList *sync.Map) {
	log.Info("Broker core enabled: %v", strings.Split(rawURL, "#")[0])
	if err := forward.Broker(parsedURL, whiteList); err != nil {
		log.Error("Broker core error: %v Restarting in 1s...", err)
		time.Sleep(1 * time.Second)
	}
}
