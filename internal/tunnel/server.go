package tunnel

import (
	"net"
	"net/url"
	"strings"
	"sync"

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
	linkConn, err := linkListen.AcceptTCP()
	if err != nil {
		log.Error("Unable to accept connections form link address: [%v]", linkAddr)
		return err
	}
	log.Info("Tunnel connection established from: [%v]", linkConn.RemoteAddr().String())
	for {
		targetConn, err := targetListen.AcceptTCP()
		if err != nil {
			log.Error("Unable to accept connections form target address: [%v] %v", targetAddr, err)
			continue
		}
		targetConn.SetNoDelay(true)
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
		go func() {
			if _, err = linkConn.Write([]byte("PASSPORT\n")); err != nil {
				log.Error("Unable to send signal: %v", err)
				targetConn.Close()
				return
			}
			log.Info("Starting data exchange: [%v] <-> [%v]", clientAddr, targetAddr)
			util.HandleConn(linkConn, targetConn)
			log.Info("Connection closed successfully")
		}()
	}
}
