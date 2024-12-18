package forward

import (
	"net"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/yosebyte/passport/internal/util"
	"github.com/yosebyte/passport/pkg/log"
)

func HandleTCP(parsedURL *url.URL, whiteList *sync.Map) error {
	linkAddr, err := net.ResolveTCPAddr("tcp", parsedURL.Host)
	if err != nil {
		log.Error("Unable to resolve link address: %v", parsedURL.Host)
		return err
	}
	targetAddr, err := net.ResolveTCPAddr("tcp", strings.TrimPrefix(parsedURL.Path, "/"))
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
	tempSlot := make(chan struct{}, 1024)
	for {
		linkConn, err := linkListen.AcceptTCP()
		if err != nil {
			log.Error("Unable to connect link address: [%v] %v", linkAddr, err)
			time.Sleep(1 * time.Second)
			continue
		}
		linkConn.SetNoDelay(true)
		tempSlot <- struct{}{}
		go func(linkConn net.Conn) {
			defer func() { <-tempSlot }()
			clientAddr := linkConn.RemoteAddr().String()
			log.Info("Client connection established: [%v]", clientAddr)
			if parsedURL.Fragment != "" {
				clientIP, _, err := net.SplitHostPort(clientAddr)
				if err != nil {
					log.Error("Unable to extract client IP address: [%v]", clientAddr)
					linkConn.Close()
					return
				}
				if _, exists := whiteList.Load(clientIP); !exists {
					log.Warn("Unauthorized IP address blocked: [%v]", clientIP)
					linkConn.Close()
					return
				}
			}
			targetConn, err := net.DialTCP("tcp", nil, targetAddr)
			if err != nil {
				log.Error("Unable to dial target address: [%v]", targetAddr)
				linkConn.Close()
				return
			}
			targetConn.SetNoDelay(true)
			log.Info("Target connection established: [%v]", targetAddr)
			log.Info("Starting data exchange: [%v] <-> [%v]", clientAddr, targetAddr)
			util.HandleConn(linkConn, targetConn)
			log.Info("Connection closed successfully")
		}(linkConn)
	}
}
