package conn

import (
	"io"
	"net"
	"sync"
)

func DataExchange(conn1, conn2 net.Conn) error {
	var (
		once1, once2 sync.Once
		wg           sync.WaitGroup
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
	errChan := make(chan error, 2)
	wg.Add(2)
	go func() {
		defer func() {
			wg.Done()
			closeConn1()
			closeConn2()
		}()
		if _, err := io.Copy(conn1, conn2); err != nil {
			closeConn1()
			closeConn2()
			errChan <- err
		}
	}()
	go func() {
		defer func() {
			wg.Done()
			closeConn1()
			closeConn2()
		}()
		if _, err := io.Copy(conn2, conn1); err != nil {
			closeConn1()
			closeConn2()
			errChan <- err
		}
	}()
	wg.Wait()
	select {
	case err := <-errChan:
		return err
	default:
		return nil
	}
}
