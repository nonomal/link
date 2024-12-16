package tunnel

import (
	"crypto/tls"
	"net"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/yosebyte/passport/internal"
	"github.com/yosebyte/passport/pkg/log"
)

func Server(parsedURL *url.URL, whiteList *sync.Map, tlsConfig *tls.Config) error {
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
	linkListen, err := tls.Listen("tcp", linkAddr.String(), tlsConfig)
	if err != nil {
		log.Error("Unable to listen link address: [%v]", linkAddr)
		return err
	}
	defer linkListen.Close()
	linkConn, err := linkListen.Accept()
	if err != nil {
		log.Error("Unable to accept connections form link address: [%v]", linkAddr)
		return err
	}
	defer linkConn.Close()
	linkTLS, ok := linkConn.(*tls.Conn)
	if !ok {
		log.Error("Non-TLS connection received")
		linkConn.Close()
		return nil
	}
	if err := linkTLS.Handshake(); err != nil {
		linkConn.Close()
		return err
	}
	log.Info("Tunnel connection established from: [%v]", linkConn.RemoteAddr().String())
	var sharedMU sync.Mutex
	errChan := make(chan error, 2)
	done := make(chan struct{})
	go func() {
		errChan <- healthCheck(linkListen, linkTLS, &sharedMU, done)
	}()
	go func() {
		errChan <- ServeTCP(parsedURL, whiteList, linkAddr, targetTCPAddr, linkListen, linkTLS, &sharedMU, done)
	}()
	go func() {
		errChan <- ServeUDP(parsedURL, whiteList, linkAddr, targetUDPAddr, linkListen, linkTLS, &sharedMU, done)
	}()
	return <-errChan
}

func healthCheck(linkListen net.Listener, linkTLS *tls.Conn, sharedMU *sync.Mutex, done chan struct{}) error {
	for {
		time.Sleep(internal.MaxReportInterval * time.Second)
		sharedMU.Lock()
		_, err := linkTLS.Write([]byte("[REPORT]\n"))
		sharedMU.Unlock()
		if err != nil {
			log.Error("Tunnel connection health check failed: %v", err)
			linkTLS.Close()
			linkListen.Close()
			close(done)
			return err
		}
	}
}
