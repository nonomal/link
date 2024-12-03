package util

import (
	"io"
	"net"
)

func HandleConn(conn1, conn2 *net.TCPConn) {
	done := make(chan struct{}, 2)
	go func() {
		io.Copy(conn1, conn2)
		done <- struct{}{}
	}()
	go func() {
		io.Copy(conn2, conn1)
		done <- struct{}{}
	}()
	<-done
	<-done
}
