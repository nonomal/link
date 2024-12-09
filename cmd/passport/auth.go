package main

import (
	"net/url"
	"sync"
	"time"

	"github.com/yosebyte/passport/internal/auth"
	"github.com/yosebyte/passport/pkg/log"
)

func authSetups(parsedURL *url.URL, whiteList *sync.Map) {
	if parsedURL.Fragment == "" {
		return
	}
	parsedAuthURL, err := url.Parse(parsedURL.Fragment)
	if err != nil {
		log.Fatal("Error parsing auth URL: %v", err)
	}
	log.Info("Auth mode enabled: %v", parsedAuthURL)
	go func() {
		for {
			if err := auth.HandleHTTP(parsedAuthURL, whiteList); err != nil {
				log.Error("Auth mode error: %v", err)
				log.Info("Restarting in 1s...")
				time.Sleep(1 * time.Second)
				continue
			}
		}
	}()
}
