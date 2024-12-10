package conn

import (
	"io"
	"net"
)

func DataExchange(conn1, conn2 net.Conn) {
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
	conn1.Close()
	conn2.Close()
}
