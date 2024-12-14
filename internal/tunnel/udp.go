package tunnel

import (
	"crypto/tls"
	"net"
	"net/url"
	"sync"
	"time"

	"github.com/yosebyte/passport/internal"
	"github.com/yosebyte/passport/pkg/log"
)

func ServeUDP(parsedURL *url.URL, whiteList *sync.Map, linkAddr *net.TCPAddr, targetAddr *net.UDPAddr, linkListen net.Listener, linkTLS net.Conn) error {
	targetConn, err := net.ListenUDP("udp", targetAddr)
	if err != nil {
		log.Error("Unable to listen target address: [%v]", targetAddr)
		return err
	}
	defer targetConn.Close()
	var mu sync.Mutex
	sem := make(chan struct{}, internal.MaxSemaphoreLimit)
	for {
		buffer := make([]byte, internal.MaxDataBuffer)
		n, clientAddr, err := targetConn.ReadFromUDP(buffer)
		if err != nil {
			log.Error("Unable to read from client address: [%v] %v", clientAddr, err)
			time.Sleep(1 * time.Second)
			continue
		}
		if parsedURL.Fragment != "" {
			clientIP := clientAddr.IP.String()
			if _, exists := whiteList.Load(clientIP); !exists {
				log.Warn("Unauthorized IP address blocked: [%v]", clientIP)
				continue
			}
		}
		mu.Lock()
		_, err = linkTLS.Write([]byte("[PASSPORT]<UDP>\n"))
		mu.Unlock()
		if err != nil {
			log.Error("Unable to send signal: %v", err)
			time.Sleep(1 * time.Second)
			continue
		}
		remoteConn, err := linkListen.Accept()
		if err != nil {
			log.Error("Unable to accept connections from link address: [%v] %v", linkAddr, err)
			time.Sleep(1 * time.Second)
			continue
		}
		remoteTLS, ok := remoteConn.(*tls.Conn)
		if !ok {
			log.Error("Non-TLS connection received")
			time.Sleep(1 * time.Second)
			continue
		}
		if err := remoteTLS.Handshake(); err != nil {
			log.Error("TLS handshake failed: %v", err)
			time.Sleep(1 * time.Second)
			continue
		}
		sem <- struct{}{}
		go func(buffer []byte, n int, remoteTLS *tls.Conn, clientAddr *net.UDPAddr) {
			defer func() {
				<-sem
				remoteTLS.Close()
			}()
			log.Info("Starting data transfer: [%v] <-> [%v]", clientAddr, targetAddr)
			_, err = remoteTLS.Write(buffer[:n])
			if err != nil {
				log.Error("Unable to write to link address: [%v] %v", linkAddr, err)
				return
			}
			n, err = remoteTLS.Read(buffer)
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
		}(buffer, n, remoteTLS, clientAddr)
	}
}

func ClientUDP(linkAddr *net.TCPAddr, targetUDPAddr *net.UDPAddr) {
	remoteTLS, err := tls.Dial("tcp", linkAddr.String(), &tls.Config{InsecureSkipVerify: true})
	if err != nil {
		log.Error("Unable to dial target address: [%v] %v", linkAddr, err)
		return
	}
	defer remoteTLS.Close()
	if err := remoteTLS.Handshake(); err != nil {
		log.Error("TLS handshake failed: %v", err)
		return
	}
	log.Info("Remote connection established: [%v]", linkAddr)
	buffer := make([]byte, internal.MaxDataBuffer)
	n, err := remoteTLS.Read(buffer)
	if err != nil {
		log.Error("Unable to read from remote address: [%v] %v", linkAddr, err)
		return
	}
	targetConn, err := net.DialUDP("udp", nil, targetUDPAddr)
	if err != nil {
		log.Error("Unable to dial target address: [%v] %v", targetUDPAddr, err)
		return
	}
	defer targetConn.Close()
	log.Info("Target connection established: [%v]", targetUDPAddr)
	err = targetConn.SetDeadline(time.Now().Add(internal.MaxUDPTimeout * time.Second))
	if err != nil {
		log.Error("Unable to set deadline: %v", err)
		return
	}
	log.Info("Starting data transfer: [%v] <-> [%v]", linkAddr, targetUDPAddr)
	_, err = targetConn.Write(buffer[:n])
	if err != nil {
		log.Error("Unable to write to target address: [%v] %v", targetUDPAddr, err)
		return
	}
	n, _, err = targetConn.ReadFromUDP(buffer)
	if err != nil {
		log.Error("Unable to read from target address: [%v] %v", targetUDPAddr, err)
		return
	}
	_, err = remoteTLS.Write(buffer[:n])
	if err != nil {
		log.Error("Unable to write to remote address: [%v] %v", linkAddr, err)
		return
	}
	log.Info("Transfer completed successfully")
}
