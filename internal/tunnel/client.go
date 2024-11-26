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
		log.Error("Unable to dial link addressL %v", linkAddr)
		return err
	}
	linkConn.SetNoDelay(true)
	targetConn, err := net.DialTCP("tcp", nil, targetAddr)
	if err != nil {
		log.Error("Unable to dial target address: %v", targetAddr)
		linkConn.Close()
		return err
	}
	targetConn.SetNoDelay(true)
	log.Info("Starting data exchange: [%v] <-> [%v]", linkAddr, targetAddr)
	util.HandleConn(linkConn, targetConn)
	return nil
}
