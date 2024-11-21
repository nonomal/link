package forward

import (
    "net"
    "net/url"
    "strings"
    "sync"
)

func Shadow(parsedURL *url.URL, whiteList *sync.Map) error {
    linkAddr, err := net.ResolveUDPAddr("udp", parsedURL.Host)
    if err != nil {
		return err
	}
    targetAddr, err := net.ResolveUDPAddr("udp", strings.TrimPrefix(parsedURL.Path, "/"))
	if err != nil {
		return err
	}
    linkConn, err := net.ListenUDP("udp", linkAddr)
    if err != nil {
        return err
    }
    defer linkConn.Close()
    readBuffer := make([]byte, 4096)
    for {
        n, remoteAddr, err := linkConn.ReadFromUDP(readBuffer)
        if err != nil {
            continue
        }
        if parsedURL.Fragment != "" {
            clientIP := remoteAddr.IP.String()
            if _, exists := whiteList.Load(clientIP); !exists {
                continue
            }
        }
        targetConn, err := net.DialUDP("udp", nil, targetAddr)
        if err != nil {
            targetConn.Close()
            continue
        }
        go func(data []byte, addr *net.UDPAddr) {
            defer targetConn.Close()
            _, err := targetConn.Write(data)
            if err != nil {
                return
            }
            writeBuffer := make([]byte, 4096)
            n, _, err := targetConn.ReadFromUDP(writeBuffer)
            if err == nil {
                linkConn.WriteToUDP(writeBuffer[:n], addr)
            }
        }(readBuffer[:n], remoteAddr)
    }
}
