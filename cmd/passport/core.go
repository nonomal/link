package main

import (
	"crypto/tls"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/yosebyte/passport/internal/forward"
	"github.com/yosebyte/passport/internal/tunnel"
	"github.com/yosebyte/passport/pkg/log"
)

func coreSelect(parsedURL *url.URL, rawURL string, whiteList *sync.Map, tlsConfig *tls.Config) {
	switch parsedURL.Scheme {
	case "server":
		runServer(parsedURL, rawURL, whiteList, tlsConfig)
	case "client":
		runClient(parsedURL, rawURL)
	case "broker":
		runBroker(parsedURL, rawURL, whiteList)
	default:
		helpInfo()
		os.Exit(1)
	}
}

func runServer(parsedURL *url.URL, rawURL string, whiteList *sync.Map, tlsConfig *tls.Config) {
	log.Info("Server core selected: %v", strings.Split(rawURL, "#")[0])
	for {
		if err := tunnel.Server(parsedURL, whiteList, tlsConfig); err != nil {
			log.Error("Server core error: %v", err)
			time.Sleep(1 * time.Second)
			log.Info("Server core restarted")
		}
	}
}

func runClient(parsedURL *url.URL, rawURL string) {
	log.Info("Client core selected: %v", strings.Split(rawURL, "#")[0])
	for {
		if err := tunnel.Client(parsedURL); err != nil {
			log.Error("Client core error: %v", err)
			time.Sleep(1 * time.Second)
			log.Info("Client core restarted")
		}
	}
}

func runBroker(parsedURL *url.URL, rawURL string, whiteList *sync.Map) {
	log.Info("Broker core selected: %v", strings.Split(rawURL, "#")[0])
	for {
		if err := forward.Broker(parsedURL, whiteList); err != nil {
			log.Error("Broker core error: %v", err)
			time.Sleep(1 * time.Second)
			log.Info("Broker core restarted")
		}
	}
}
