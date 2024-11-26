package main

import (
	"net/url"
	"sync"
	"time"

	"github.com/yosebyte/passport/internal/util"
	"github.com/yosebyte/passport/pkg/log"
)

func authSetups(parsedURL *url.URL, whiteList *sync.Map) {
	if parsedURL.Fragment == "" {
		log.Info("Authorization disabled")
		return
	}
	parsedAuthURL, err := url.Parse(parsedURL.Fragment)
	if err != nil {
		log.Fatal("Error parsing auth URL: %v", err)
	}
	log.Info("Authorization enabled: %v", parsedAuthURL)
	go func() {
		for {
			if err := util.HandleHTTP(parsedAuthURL, whiteList); err != nil {
				log.Error("Authorization error: %v Restarting in 1s...", err)
				time.Sleep(1 * time.Second)
				continue
			}
		}
	}()
}
