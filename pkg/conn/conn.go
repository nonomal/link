package conn

import (
	"io"
	"net"
)

func DataExchange(conn1, conn2 net.Conn) {
	done := make(chan struct{}, 2)
	go func() {
		io.Copy(conn1, conn2)
		sendEOF(conn1)
		done <- struct{}{}
	}()
	go func() {
		io.Copy(conn2, conn1)
		sendEOF(conn2)
		done <- struct{}{}
	}()
	<-done
	<-done
	closeConn(conn1)
	closeConn(conn2)
}

func sendEOF(conn net.Conn) {
	switch c := conn.(type) {
	case *net.TCPConn:
		c.CloseWrite()
	default:
		c.Close()
	}
}

func closeConn(conn net.Conn) {
	if conn != nil {
		conn.Close()
	}
}
