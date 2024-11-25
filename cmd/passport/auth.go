package main

import (
	"log"
	"net/url"
	"sync"
	"time"

	"github.com/yosebyte/passport/internal/util"
)

func authSetup(parsedURL *url.URL, whiteList *sync.Map) {
	if parsedURL.Fragment == "" {
		return
	}
	parsedAuthURL, err := url.Parse(parsedURL.Fragment)
	if err != nil {
		log.Fatalf("[ERRO] URL: %v", err)
	}
	log.Printf("[INFO] Auth: %v", parsedAuthURL)
	go func() {
		for {
			if err := util.Auth(parsedAuthURL, whiteList); err != nil {
				log.Printf("[ERRO] Auth: %v", err)
				time.Sleep(1 * time.Second)
				continue
			}
		}
	}()
}
