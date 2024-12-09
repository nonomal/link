package tunnel

import (
	"net"
	"net/url"
	"strings"

	"github.com/yosebyte/passport/internal"
	"github.com/yosebyte/passport/pkg/log"
)

func Client(parsedURL *url.URL) error {
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
	linkConn, err := net.DialTCP("tcp", nil, linkAddr)
	if err != nil {
		log.Error("Unable to dial link address: [%v]", linkAddr)
		return err
	}
	defer linkConn.Close()
	log.Info("Tunnel connection established to: [%v]", linkAddr)
	buffer := make([]byte, internal.MinBufferSize)
	for {
		n, err := linkConn.Read(buffer)
		if err != nil {
			log.Error("Unable to read form link address: [%v] %v", linkAddr, err)
			break
		}
		if string(buffer[:n]) == "[PASSPORT]<TCP>\n" {
			go ClientTCP(linkAddr, targetTCPAddr)
		}
		if string(buffer[:n]) == "[PASSPORT]<UDP>\n" {
			go ClientUDP(linkAddr, targetUDPAddr)
		}
	}
	return nil
}
