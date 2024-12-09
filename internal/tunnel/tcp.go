package tunnel

import (
	"net"
	"net/url"
	"sync"

	"github.com/yosebyte/passport/pkg/conn"
	"github.com/yosebyte/passport/pkg/log"
)

func ServeTCP(parsedURL *url.URL, whiteList *sync.Map, linkAddr, targetAddr *net.TCPAddr, linkListen *net.TCPListener, linkConn *net.TCPConn) error {
	targetListen, err := net.ListenTCP("tcp", targetAddr)
	if err != nil {
		log.Error("Unable to listen target address: [%v]", targetAddr)
		return err
	}
	defer targetListen.Close()
	var mu sync.Mutex
	semaphore := make(chan struct{}, 1024)
	for {
		targetConn, err := targetListen.AcceptTCP()
		if err != nil {
			log.Error("Unable to accept connections form target address: [%v] %v", targetAddr, err)
			break
		}
		clientAddr := targetConn.RemoteAddr().String()
		log.Info("Target connection established from: [%v]", clientAddr)
		if parsedURL.Fragment != "" {
			clientIP, _, err := net.SplitHostPort(clientAddr)
			if err != nil {
				log.Error("Unable to extract client IP address: [%v] %v", clientAddr, err)
				targetConn.Close()
				continue
			}
			if _, exists := whiteList.Load(clientIP); !exists {
				log.Warn("Unauthorized IP address blocked: [%v]", clientIP)
				targetConn.Close()
				continue
			}
		}
		semaphore <- struct{}{}
		go func(targetConn *net.TCPConn) {
			defer func() {
				<-semaphore
				targetConn.Close()
			}()
			mu.Lock()
			_, err = linkConn.Write([]byte("[PASSPORT]<TCP>\n"))
			mu.Unlock()
			if err != nil {
				log.Error("Unable to send signal: %v", err)
				targetConn.Close()
				return
			}
			remoteConn, err := linkListen.AcceptTCP()
			if err != nil {
				log.Error("Unable to accept connections form link address: [%v] %v", linkAddr, err)
				return
			}
			defer remoteConn.Close()
			log.Info("Starting data exchange: [%v] <-> [%v]", clientAddr, targetAddr)
			conn.DataExchange(remoteConn, targetConn)
			log.Info("Connection closed successfully")
		}(targetConn)
	}
	return nil
}

func ClientTCP(linkAddr, targetTCPAddr *net.TCPAddr) {
	targetConn, err := net.DialTCP("tcp", nil, targetTCPAddr)
	if err != nil {
		log.Error("Unable to dial target address: [%v], %v", targetTCPAddr, err)
		return
	}
	defer targetConn.Close()
	log.Info("Target connection established: [%v]", targetTCPAddr)
	remoteConn, err := net.DialTCP("tcp", nil, linkAddr)
	if err != nil {
		log.Error("Unable to dial target address: [%v], %v", linkAddr, err)
		return
	}
	defer remoteConn.Close()
	log.Info("Starting data exchange: [%v] <-> [%v]", linkAddr, targetTCPAddr)
	conn.DataExchange(remoteConn, targetConn)
	log.Info("Connection closed successfully")
}
