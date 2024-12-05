package tunnel

import (
	"net"
	"net/url"
	"strings"
	"sync"

	"github.com/yosebyte/passport/pkg/log"
)

func ServeUDP(parsedURL *url.URL, whiteList *sync.Map) error {
	linkAddr, err := net.ResolveTCPAddr("tcp", parsedURL.Host)
	if err != nil {
		log.Error("Unable to resolve link address: %v", parsedURL.Host)
		return err
	}
	targetAddr, err := net.ResolveUDPAddr("udp", strings.TrimPrefix(parsedURL.Path, "/"))
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
	targetConn, err := net.ListenUDP("udp", targetAddr)
	if err != nil {
		log.Error("Unable to listen target address: [%v]", targetAddr)
		return err
	}
	defer targetConn.Close()
	var mu sync.Mutex
	semaphore := make(chan struct{}, 1024)
	for {
		buffer := make([]byte, 8192)
		n, clientAddr, err := targetConn.ReadFromUDP(buffer)
		if err != nil {
			log.Error("Unable to read from client address: [%v] %v", clientAddr, err)
			continue
		}
		if parsedURL.Fragment != "" {
			clientIP := clientAddr.IP.String()
			if _, exists := whiteList.Load(clientIP); !exists {
				log.Warn("Unauthorized IP address blocked: [%v]", clientIP)
				continue
			}
		}
		semaphore <- struct{}{}
		go func(buffer []byte, n int, clientAddr *net.UDPAddr) {
			defer func() { <-semaphore }()
			mu.Lock()
			_, err = linkConn.Write([]byte("[PASSPORT]<UDP>\n"))
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
			log.Info("Starting data transfer: [%v] <-> [%v]", clientAddr, targetAddr)
			_, err = remoteConn.Write(buffer[:n])
			if err != nil {
				log.Error("Unable to write to link address: [%v] %v", linkAddr, err)
				return
			}
			n, err = remoteConn.Read(buffer)
			if err != nil {
				log.Error("Unable to read from link address: [%v] %v", linkAddr, err)
				return
			}
			_, err = targetConn.WriteToUDP(buffer[:n], clientAddr)
			if err != nil {
				log.Error("Unable to write to client address: [%v] %v", clientAddr, err)
				return
			}
			log.Info("Transfer completed successfully")
		}(buffer, n, clientAddr)
	}
}
