package forward

import (
	"net"
	"net/url"
	"strings"
	"sync"

	"github.com/yosebyte/passport/pkg/log"
)

func HandleUDP(parsedURL *url.URL, whiteList *sync.Map) error {
	linkAddr, err := net.ResolveUDPAddr("udp", parsedURL.Host)
	if err != nil {
		log.Error("Unable to resolve link address: %v", parsedURL.Host)
		return err
	}
	targetAddr, err := net.ResolveUDPAddr("udp", strings.TrimPrefix(parsedURL.Path, "/"))
	if err != nil {
		log.Error("Unable to resolve target address: %v", strings.TrimPrefix(parsedURL.Path, "/"))
		return err
	}
	linkConn, err := net.ListenUDP("udp", linkAddr)
	if err != nil {
		log.Error("Unable to listen link address: [%v]", linkAddr)
		return err
	}
	defer linkConn.Close()
	readBuffer := make([]byte, 4096)
	for {
		n, remoteAddr, err := linkConn.ReadFromUDP(readBuffer)
		if err != nil {
			log.Error("Unable to read UDP from remote address: [%v]", remoteAddr)
			continue
		}
		if parsedURL.Fragment != "" {
			clientIP := remoteAddr.IP.String()
			if _, exists := whiteList.Load(clientIP); !exists {
				log.Warn("Unauthorized IP address blocked: [%v]", clientIP)
				continue
			}
		}
		targetConn, err := net.DialUDP("udp", nil, targetAddr)
		if err != nil {
			log.Error("Unable to dial target address: [%v]", targetAddr)
			targetConn.Close()
			continue
		}
		go func(data []byte, remoteAddr *net.UDPAddr) {
			defer targetConn.Close()
			_, err := targetConn.Write(data)
			if err != nil {
				log.Error("Unable to write remote data: [%v]", remoteAddr)
				return
			}
			writeBuffer := make([]byte, 4096)
			n, _, err := targetConn.ReadFromUDP(writeBuffer)
			if err == nil {
				log.Info("Starting data transfer: [%v] <-> [%v]", remoteAddr, targetAddr)
				linkConn.WriteToUDP(writeBuffer[:n], remoteAddr)
			}
		}(readBuffer[:n], remoteAddr)
	}
}
