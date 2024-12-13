package conn

import (
	"io"
	"net"
	"sync"
	"time"
)

func DataExchange(conn1, conn2 net.Conn) error {
	var (
		once1, once2 sync.Once
		wg           sync.WaitGroup
		timeout      = 30 * time.Second
	)
	closeConn1 := func() {
		once1.Do(func() {
			if conn1 != nil {
				conn1.Close()
			}
		})
	}
	closeConn2 := func() {
		once2.Do(func() {
			if conn2 != nil {
				conn2.Close()
			}
		})
	}
	updateDeadline := func(conn net.Conn) {
		conn.SetDeadline(time.Now().Add(timeout))
	}
	errChan := make(chan error, 2)
	wg.Add(2)
	go func() {
		defer func() {
			wg.Done()
			closeConn1()
			closeConn2()
		}()
		for {
			updateDeadline(conn1)
			updateDeadline(conn2)
			if _, err := io.Copy(conn1, conn2); err != nil {
				errChan <- err
				return
			}
		}
	}()
	go func() {
		defer func() {
			wg.Done()
			closeConn2()
			closeConn1()
		}()
		for {
			updateDeadline(conn2)
			updateDeadline(conn1)
			if _, err := io.Copy(conn2, conn1); err != nil {
				errChan <- err
				return
			}
		}
	}()
	wg.Wait()
	select {
	case err := <-errChan:
		return err
	default:
		return io.EOF
	}
}
