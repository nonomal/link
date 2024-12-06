package tunnel

import (
	"net"
	"net/url"
	"sync"

	"github.com/yosebyte/passport/internal/util"
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
		defer targetConn.Close()
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
			defer func() { <-semaphore }()
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
			util.HandleConn(remoteConn, targetConn)
			log.Info("Connection closed successfully")
		}(targetConn)
	}
	return nil
}
