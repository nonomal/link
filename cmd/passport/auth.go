package main

import (
	"crypto/tls"
	"net/url"
	"sync"
	"time"

	"github.com/yosebyte/passport/internal"
	"github.com/yosebyte/passport/pkg/log"
)

func authSetups(parsedURL *url.URL, whiteList *sync.Map, tlsConfig *tls.Config) {
	if parsedURL.Fragment == "" {
		return
	}
	parsedAuthURL, err := url.Parse(parsedURL.Fragment)
	if err != nil {
		log.Fatal("Unable to parse auth URL: %v", err)
	}
	log.Info("Auth mode enabled: %v", parsedAuthURL)
	go func() {
		for {
			if err := internal.HandleHTTP(parsedAuthURL, whiteList, tlsConfig); err != nil {
				log.Error("Auth mode error: %v", err)
				log.Info("Restarting in 1s...")
				time.Sleep(1 * time.Second)
				continue
			}
		}
	}()
}
