package tunnel

import (
	"net"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/yosebyte/passport/internal/util"
	"github.com/yosebyte/passport/pkg/log"
)

func Server(parsedURL *url.URL, whiteList *sync.Map) error {
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
	targetListen, err := net.ListenTCP("tcp", targetAddr)
	if err != nil {
		log.Error("Unable to listen target address: [%v]", targetAddr)
		return err
	}
	defer targetListen.Close()
	var linkConn *net.TCPConn
	go func() {
		for {
			tempConn, err := linkListen.AcceptTCP()
			if err != nil {
				time.Sleep(1 * time.Second)
				continue
			}
			if linkConn != nil {
				log.Warn("Connection closed by target service: [%v]", targetAddr)
				linkConn.Close()
			}
			linkConn = tempConn
			log.Info("Reconnection complete")
			linkConn.SetNoDelay(true)
		}
	}()
	targetConn, err := targetListen.AcceptTCP()
	if err != nil {
		log.Error("Unable to accept connections form target address: [%v]", targetAddr)
		linkConn.Close()
		return err
	}
	targetConn.SetNoDelay(true)
	clientAddr := linkConn.RemoteAddr().String()
	if parsedURL.Fragment != "" {
		clientIP, _, err := net.SplitHostPort(clientAddr)
		if err != nil {
			log.Error("Unable to extract client IP address: [%v]", clientAddr)
			targetConn.Close()
			linkConn.Close()
			return err
		}
		if _, exists := whiteList.Load(clientIP); !exists && linkConn != nil {
			log.Warn("Unauthorized IP address blocked: [%v]", clientIP)
			targetConn.Close()
			linkConn.Close()
			return nil
		}
	}
	if linkConn == nil {
		targetConn.Close()
		return nil
	}
	log.Info("Starting data exchange: [%v] <-> [%v]", clientAddr, targetAddr)
	util.HandleConn(linkConn, targetConn)
	return nil
}
