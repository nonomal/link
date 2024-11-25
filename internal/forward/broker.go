package forward

import (
	"net/url"
	"sync"
)

func Broker(parsedURL *url.URL, whiteList *sync.Map) error {
	errChan := make(chan error, 2)
	go func() {
		errChan <- HandleTCP(parsedURL, whiteList)
	}()
	go func() {
		errChan <- HandleUDP(parsedURL, whiteList)
	}()
	return <-errChan
}
