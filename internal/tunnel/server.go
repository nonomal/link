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
		log.Error("Unable to resolve target TCP address: %v", strings.TrimPrefix(parsedURL.Path, "/"))
		return err
	}
	targetUDPAddr, err := net.ResolveUDPAddr("udp", strings.TrimPrefix(parsedURL.Path, "/"))
	if err != nil {
		log.Error("Unable to resolve target UDP address: %v", strings.TrimPrefix(parsedURL.Path, "/"))
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
	targetTCPListen, err := net.ListenTCP("tcp", targetTCPAddr)
	if err != nil {
		log.Error("Unable to listen target TCP address: [%v]", targetTCPAddr)
		return err
	}
	defer targetTCPListen.Close()
	targetUDPConn, err := net.ListenUDP("udp", targetUDPAddr)
	if err != nil {
		log.Error("Unable to listen target UDP address: [%v]", targetUDPAddr)
		return err
	}
	defer targetUDPConn.Close()
	var sharedMU sync.Mutex
	errChan := make(chan error, 2)
	done := make(chan struct{})
	go func() {
		errChan <- healthCheck(linkListen, targetTCPListen, targetUDPConn, linkTLS, &sharedMU, done)
	}()
	go func() {
		errChan <- ServeTCP(parsedURL, whiteList, targetTCPListen, linkListen, linkTLS, &sharedMU, done)
	}()
	go func() {
		errChan <- ServeUDP(parsedURL, whiteList, targetUDPConn, linkListen, linkTLS, &sharedMU, done)
	}()
	return <-errChan
}

func healthCheck(linkListen net.Listener, targetTCPListen *net.TCPListener, targetUDPConn *net.UDPConn, linkTLS *tls.Conn, sharedMU *sync.Mutex, done chan struct{}) error {
	for {
		time.Sleep(internal.MaxReportInterval * time.Second)
		sharedMU.Lock()
		_, err := linkTLS.Write([]byte("[]\n"))
		sharedMU.Unlock()
		if err != nil {
			log.Error("Tunnel connection health check failed")
			if linkListen != nil {
				linkListen.Close()
			}
			if targetTCPListen != nil {
				targetTCPListen.Close()
			}
			if targetUDPConn != nil {
				targetUDPConn.Close()
			}
			if linkTLS != nil {
				linkTLS.Close()
			}
			close(done)
			return err
		}
	}
}
