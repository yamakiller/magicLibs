package listener

import (
	"net"
	"sync"
)

//TCPListener TCP Socket Listener
type TCPListener struct {
	net.Listener
	_wg sync.WaitGroup
}

//Accept 接受链接
func (slf *TCPListener) Accept() (net.Conn, error) {
	c, err := slf.Listener.Accept()
	if err != nil {
		return nil, err
	}

	slf._wg.Add(1)
	return &TCPConn{Conn: c, _wg: &slf._wg}, nil
}

//Wait 等待所有客户端结束
func (slf *TCPListener) Wait() {
	slf._wg.Wait()
}

//TCPConn TCP Accept client
type TCPConn struct {
	net.Conn
	_wg *sync.WaitGroup
}

//Close ...
func (slf *TCPConn) Close() error {
	if slf._wg != nil {
		slf._wg.Done()
	}
	slf._wg = nil
	return slf.Conn.Close()
}
