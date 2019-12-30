package listen

import (
	"net"
	"sync"

	"github.com/yamakiller/magicLibs/mutex"
)

//TCPClient tcp service client
type TCPClient struct {
	_id   int32
	_conn net.Conn
	_addr string
	_ref  int32
}

func (slf *TCPClient) GetID() int32 {
	return slf._id
}

func (slf *TCPClient) GetAddr() string {
	return slf._addr
}

func (slf *TCPClient) Write(b []byte) (int, error) {
	return slf._conn.Write(b)
}

func (slf *TCPClient) Read(b []byte) (int, error) {
	return slf._conn.Read(b)
}

func (slf *TCPClient) Close() {
	slf._conn.Close()
}

//TCPListen tcp listener
type TCPListen struct {
	_net net.Listener
	_map map[int32]*TCPClient
	_sz  int32
	_syn mutex.SpinLock
	_pol sync.Pool
	_cap int
}

//Listen listen tcp network
func (slf *TCPListen) Listen(addr string, cap int) error {
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return err
	}

	l, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return err
	}
	slf._net = l
	slf._cap = cap
	slf._sz = 0
	slf._map = make(map[int32]*TCPClient)

	slf._pol = sync.Pool{
		New: func() interface{} {
			return &TCPClient{}
		},
	}

	return nil
}

func (slf *TCPListen) Accept() (ClientSocket, error) {
	c, err := slf._net.Accept()
	if err != nil {
		return nil, err
	}

	cls := slf._pol.Get().(*TCPClient)
	cls._id = 0
	cls._conn = c

	return nil, nil
}
