package kcp

import (
	"context"
	"crypto/tls"
	"errors"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/yamakiller/magicLibs/mmath"
	"github.com/yamakiller/magicLibs/net/connection/interfaces"
	"github.com/yamakiller/magicLibs/net/middle"
	"github.com/yamakiller/magicLibs/util"
	"github.com/yamakiller/mgokcp/mkcp"
)

type KCPSeria interface {
	UnSeria([]byte) (interface{}, error)
	Seria(interface{}, *mkcp.KCP) (int, error)
}

type recvData struct {
	_data []byte
	_len  int
}

//Client KCP(UDP)协议客户端
type Client struct {
	WriteWaitQueue int
	ReadWaitQueue  int
	RecvWndSize    int32
	SendWndSize    int32
	NoDelay        int32
	Interval       int32
	Resend         int32
	Nc             int32
	RxMinRto       int32
	FastResend     int32
	Serializer     KCPSeria
	Exception      interfaces.Exception
	Mtu            int
	Middleware     middle.KCMiddleware
	Allocator      func(int) []byte
	Releaser       func([]byte)

	_conn *net.UDPConn
	_id   uint32
	_kcp  *mkcp.KCP
	_sync sync.Mutex
	_addr *net.UDPAddr

	_cancel  context.CancelFunc
	_ctx     context.Context
	_sdQueue chan interface{}
	_rdQueue chan *recvData

	_buffer     []byte
	_wbs        uint64
	_rbs        uint64
	_lastActive int64

	_wg sync.WaitGroup
}

//Connect 连接服务器
func (c *Client) Connect(addr string, timeout time.Duration) error {

	udpAddr, err := net.ResolveUDPAddr("", addr)
	if err != nil {
		return err
	}

	if c.Mtu == 0 {
		c.Mtu = 1400
	}

	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		return err
	}
	c._conn = conn

	if c.Middleware != nil {
		conv, err := c.Middleware.Subscribe(conn, udpAddr, timeout)
		if err != nil {
			return err
		}
		c._id = conv.(uint32)
		conn.SetReadDeadline(time.Time{})
	}

	c._addr = udpAddr
	c._sdQueue = make(chan interface{}, c.WriteWaitQueue)
	c._rdQueue = make(chan *recvData, c.ReadWaitQueue)

	c._buffer = make([]byte, mmath.Align(uint32(c.Mtu), 4))
	c._kcp = mkcp.New(c._id, c)
	c._kcp.WithOutput(output)
	c._kcp.WndSize(c.SendWndSize, c.RecvWndSize)
	c._kcp.NoDelay(c.NoDelay, c.Interval, c.Resend, c.Nc)
	c._kcp.SetMTU(int32(c.Mtu))
	if c.RxMinRto > 0 {
		c._kcp.SetRxMinRto(c.RxMinRto)
	}
	if c.FastResend > 0 {
		c._kcp.SetFastResend(c.FastResend)
	}

	c._wg.Add(2)
	ctx, cancel := context.WithCancel(context.Background())
	c._cancel = cancel
	c._ctx = ctx
	go c.writeServe()
	go c.readServe()

	return nil
}

func (c *Client) ConnectTls(addr string, timeout time.Duration, config *tls.Config) error {
	return nil
}

func (c *Client) writeServe() {
	defer func() {
		c._wg.Done()
	}()

	var current int64
	for {
	active:
		current = util.Timestamp()
		select {
		case <-c._ctx.Done():
			goto exit
		case <-time.After(time.Duration(c.Interval) * time.Millisecond):
		case msg, ok := <-c._sdQueue:
			if !ok {
				goto exit
			}

			c._sync.Lock()
			n, err := c.Serializer.Seria(msg, c._kcp)
			c._sync.Unlock()
			if err != nil {
				if c.Exception != nil {
					c.Exception.Error(err)
				}

				goto active
			}

			atomic.AddUint64(&c._wbs, uint64(n))
		}

		c._sync.Lock()
		c._kcp.Update(uint32(current & 0xFFFFFFFF))
		c._sync.Unlock()
	}
exit:
}

func (c *Client) readServe() {
	defer func() {
		c._wg.Done()
	}()

	for {
	active:
		select {
		case <-c._ctx.Done():
			goto exit
		default:
			n, _, err := c._conn.ReadFromUDP(c._buffer)
			if err != nil {
				if c.Exception != nil {
					c.Exception.Error(err)
				}
				goto active
			}

			c._sync.Lock()
			c._kcp.Input(c._buffer, int32(n))
			c._sync.Unlock()

			atomic.AddUint64(&c._rbs, uint64(n))
			c._lastActive = util.Timestamp()

			for {
				c._sync.Lock()
				n = int(c._kcp.Recv(c._buffer, int32(len(c._buffer))))
				c._sync.Unlock()
				if n < 0 {
					break
				}
				//需要修改池化
				var tmpBuf []byte
				if c.Allocator != nil {
					tmpBuf = c.Allocator(n)
				} else {
					tmpBuf = make([]byte, n)
				}

				copy(tmpBuf, c._buffer[:n])
				c._rdQueue <- &recvData{_data: tmpBuf, _len: n}
			}
		}
	}
exit:
}

//Parse 解析数据, 需要修改
func (c *Client) Parse() (interface{}, error) {
	c._wg.Add(1)
	defer c._wg.Done()

	if c._ctx == nil {
		return nil, errors.New("connect is closed")
	}

	select {
	case <-c._ctx.Done():
		return nil, errors.New("connect is closed")
	case d, ok := <-c._rdQueue:
		if !ok {
			return nil, errors.New("connect is closed")
		}

		defer func() {
			if c.Releaser != nil {
				c.Releaser(d._data)
			}
		}()

		msg, err := c.Serializer.UnSeria(d._data[:d._len])
		if err != nil {
			return nil, err
		}

		return msg, nil
	}
}

//SendTo 发送数据
func (c *Client) SendTo(msg interface{}) error {
	c._wg.Add(1)
	defer c._wg.Done()

	select {
	case <-c._ctx.Done():
		return errors.New("closed")
	case c._sdQueue <- msg:
	}

	return nil
}

func (c *Client) Close() error {
	var err error
	if c._conn != nil {
		c._cancel()
		err = c._conn.Close()
		c._wg.Wait()
		c._conn = nil
		c._ctx = nil
		c._cancel = nil

		close(c._rdQueue)
		for d := range c._rdQueue {
			if d == nil {
				break
			}
			if c.Releaser != nil {
				c.Releaser(d._data)
			}

		}
		c._rdQueue = nil

		close(c._sdQueue)
		c._sdQueue = nil

		mkcp.Free(c._kcp)
		c._kcp = nil
	}

	return err
}

func (c *Client) Shutdown() {
}

func output(buff []byte, user interface{}) int32 {
	client := user.(*Client)
	client._conn.Write(buff)
	return 0
}
