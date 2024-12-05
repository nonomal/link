package tunnel

import (
	"net"
	"net/url"
	"strings"
	"sync"

	"github.com/yosebyte/passport/pkg/log"
)

func Server(parsedURL *url.URL, whiteList *sync.Map) error {
	linkAddr, err := net.ResolveTCPAddr("tcp", parsedURL.Host)
	if err != nil {
		log.Error("Unable to resolve link address: %v", parsedURL.Host)
		return err
	}
	targetTCPAddr, err := net.ResolveTCPAddr("tcp", strings.TrimPrefix(parsedURL.Path, "/"))
	if err != nil {
		log.Error("Unable to resolve target address: %v", strings.TrimPrefix(parsedURL.Path, "/"))
		return err
	}
	targetUDPAddr, err := net.ResolveUDPAddr("udp", strings.TrimPrefix(parsedURL.Path, "/"))
	if err != nil {
		log.Error("Unable to resolve target address: %v", strings.TrimPrefix(parsedURL.Path, "/"))
		return err
	}
	linkListen, err := net.ListenTCP("tcp", linkAddr)
	if err != nil {
		log.Error("Unable to listen link address: [%v]", linkAddr)
		return err
	}
	defer linkListen.Close()
	linkConn, err := linkListen.AcceptTCP()
	if err != nil {
		log.Error("Unable to accept connections form link address: [%v]", linkAddr)
		return err
	}
	defer linkConn.Close()
	log.Info("Tunnel connection established from: [%v]", linkConn.RemoteAddr().String())
	mu := &sync.Mutex{}
	errChan := make(chan error, 2)
	go func() {
		errChan <- ServeTCP(parsedURL, whiteList, linkAddr, targetTCPAddr, linkListen, linkConn, mu)
	}()
	go func() {
		errChan <- ServeUDP(parsedURL, whiteList, linkAddr, targetUDPAddr, linkListen, linkConn, mu)
	}()
	return <-errChan
}
