package tunnel

import (
	"net/url"
	"sync"
)

func Server(parsedURL *url.URL, whiteList *sync.Map) error {
	errChan := make(chan error, 2)
	go func() {
		errChan <- ServeTCP(parsedURL, whiteList)
	}()
	go func() {
		errChan <- ServeUDP(parsedURL, whiteList)
	}()
	return <-errChan
}
