package tunnel

import (
	"net"
	"net/url"
	"strings"

	"github.com/yosebyte/passport/internal/util"
	"github.com/yosebyte/passport/pkg/log"
)

func Client(parsedURL *url.URL) error {
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
	linkConn, err := net.DialTCP("tcp", nil, linkAddr)
	if err != nil {
		log.Error("Unable to dial link address: [%v]", linkAddr)
		return err
	}
	linkConn.SetNoDelay(true)
	log.Info("Tunnel connection established to: [%v]", linkAddr)
	buffer := make([]byte, 16)
	for {
		n, err := linkConn.Read(buffer)
		if err != nil {
			log.Error("Error reading form link address: [%v] %v", linkAddr, err)
			continue
		}
		if string(buffer[:n]) == "PASSPORT\n" {
			go func() {
				targetConn, err := net.DialTCP("tcp", nil, targetAddr)
				if err != nil {
					log.Error("Unable to dial target address: [%v], %v", targetAddr, err)
					return
				}
				targetConn.SetNoDelay(true)
				log.Info("Target connection established, starting data exchange: [%v] <-> [%v]", linkAddr, targetAddr)
				util.HandleConn(linkConn, targetConn)
				log.Info("Connection closed successfully")
			}()
		}
	}
}
