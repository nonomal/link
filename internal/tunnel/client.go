package tunnel

import (
	"net"
	"net/url"
	"strings"
	"time"

	"github.com/yosebyte/passport/internal/util"
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
	buffer := make([]byte, 1024)
	for {
		n, err := linkConn.Read(buffer)
		if err != nil {
			log.Error("Unable to read form link address: [%v] %v", linkAddr, err)
			break
		}
		if string(buffer[:n]) == "[PASSPORT]<TCP>\n" {
			go func() {
				targetConn, err := net.DialTCP("tcp", nil, targetTCPAddr)
				if err != nil {
					log.Error("Unable to dial target address: [%v], %v", targetTCPAddr, err)
					return
				}
				log.Info("Target connection established: [%v]", targetTCPAddr)
				remoteConn, err := net.DialTCP("tcp", nil, linkAddr)
				if err != nil {
					log.Error("Unable to dial target address: [%v], %v", linkAddr, err)
					return
				}
				log.Info("Starting data exchange: [%v] <-> [%v]", linkAddr, targetTCPAddr)
				util.HandleConn(remoteConn, targetConn)
				log.Info("Connection closed successfully")
			}()
		}
		if string(buffer[:n]) == "[PASSPORT]<UDP>\n" {
			go func() {
				remoteConn, err := net.DialTCP("tcp", nil, linkAddr)
				if err != nil {
					log.Error("Unable to dial target address: [%v], %v", linkAddr, err)
					return
				}
				log.Info("Remote connection established: [%v]", linkAddr)
				buffer := make([]byte, 8192)
				n, err := remoteConn.Read(buffer)
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
				err = targetConn.SetDeadline(time.Now().Add(5 * time.Second))
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
				_, err = remoteConn.Write(buffer[:n])
				if err != nil {
					log.Error("Unable to write to remote address: [%v] %v", linkAddr, err)
					return
				}
				log.Info("Transfer completed successfully")
			}()
		}
	}
	return nil
}
