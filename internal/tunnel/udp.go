package tunnel

import (
	"net"
	"net/url"
	"sync"

	"github.com/yosebyte/passport/pkg/log"
)

func ServeUDP(parsedURL *url.URL, whiteList *sync.Map, linkAddr *net.TCPAddr, targetAddr *net.UDPAddr, linkListen *net.TCPListener, linkConn *net.TCPConn) error {
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
			break
		}
		if parsedURL.Fragment != "" {
			clientIP := clientAddr.IP.String()
			if _, exists := whiteList.Load(clientIP); !exists {
				log.Warn("Unauthorized IP address blocked: [%v]", clientIP)
				continue
			}
		}
		mu.Lock()
		_, err = linkConn.Write([]byte("[PASSPORT]<UDP>\n"))
		mu.Unlock()
		if err != nil {
			log.Error("Unable to send signal: %v", err)
			break
		}
		semaphore <- struct{}{}
		go func(buffer []byte, n int, clientAddr *net.UDPAddr) {
			defer func() { <-semaphore }()
			remoteConn, err := linkListen.AcceptTCP()
			if err != nil {
				log.Error("Unable to accept connections from link address: [%v] %v", linkAddr, err)
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
	return nil
}
