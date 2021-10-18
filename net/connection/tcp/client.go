package tcp

import (
	"bufio"
	"context"
	"crypto/tls"
	"io"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/yamakiller/magicLibs/net/connection/connerrors"
	"github.com/yamakiller/magicLibs/net/connection/interfaces"
)

func New(readWriteBufferSize, writeBufferQueueSize int) *Client {
	return &Client{_readWriteBufferSize: readWriteBufferSize,
		_queue: make(chan interface{}, writeBufferQueueSize)}
}

type Client struct {
	Serializer interfaces.Serialization
	Exception  interfaces.Exception

	_readWriteBufferSize int
	_conn                io.ReadWriteCloser
	_cancel              context.CancelFunc
	_ctx                 context.Context
	_queue               chan interface{}
	_reader              *bufio.Reader
	_writer              *bufio.Writer
	_wbs                 uint64
	_rbs                 uint64
	_lastActive          int64
	_connected           bool
	_wg                  sync.WaitGroup
}

//Connect 连接远程地址
func (c *Client) Connect(address string, timeout time.Duration) error {

	var (
		err  error
		conn net.Conn
	)

	if timeout == 0 {
		var raddr *net.TCPAddr
		raddr, err = net.ResolveTCPAddr("tcp", address)
		if err != nil {
			return err
		}
		conn, err = net.DialTCP("tcp", nil, raddr)
	} else {
		conn, err = net.DialTimeout("tcp", address, timeout)
	}

	if err != nil {
		return err
	}

	c._conn = conn.(io.ReadWriteCloser)
	c._connected = true
	c.init()

	return nil
}

func (c *Client) ConnectTls(addr string, timeout time.Duration, config *tls.Config) error {
	var (
		err  error
		conn net.Conn
	)

	if timeout == 0 {
		conn, err = tls.Dial("tcp", addr, config)
	} else {
		conn, err = tls.DialWithDialer(&net.Dialer{Timeout: timeout}, "tcp", addr, config)
	}
	if err != nil {
		return err
	}

	c._conn = conn.(io.ReadWriteCloser)
	c.init()

	return nil
}

func (c *Client) init() {
	c._reader = bufio.NewReaderSize(c._conn, c._readWriteBufferSize)
	c._writer = bufio.NewWriterSize(c._conn, c._readWriteBufferSize)
	c._ctx, c._cancel = context.WithCancel(context.Background())

	c._wg.Add(1)
	go c.writeServe()
}

func (c *Client) writeServe() {
	defer func() {
		c._connected = false
		c._wg.Done()
	}()

	for {
	active:
		select {
		case <-c._ctx.Done():
			goto exit
		case msg := <-c._queue:
			n, err := c.Serializer.Seria(msg, c._writer)
			if err != nil {
				if c.Exception != nil {
					c.Exception.Error(err)
				}

				goto active
			}

			if c._writer.Buffered() > 0 {
				if err := c._writer.Flush(); err != nil &&
					c.Exception != nil {
					c.Exception.Error(err)
				}
			}

			atomic.AddUint64(&c._wbs, uint64(n))
		}
	}
exit:
}

func (c *Client) GetWriteBytes() uint64 {
	return atomic.LoadUint64(&c._wbs)
}

func (c *Client) GetReadBytes() uint64 {
	return atomic.LoadUint64(&c._rbs)
}

//Parse 解析数据
func (c *Client) Parse() (interface{}, error) {
	c._wg.Add(1)
	defer c._wg.Done()

	select {
	case <-c._ctx.Done():
		return nil, connerrors.ErrConnectClosed
	default:
		m, n, err := c.Serializer.UnSeria(c._reader)
		if err != nil {
			return nil, err
		}

		atomic.AddUint64(&c._rbs, uint64(n))
		return m, nil
	}
}

//SendTo 发送数据
func (c *Client) SendTo(msg interface{}) error {
	c._wg.Add(1)
	defer c._wg.Done()

	if c._queue == nil {
		return connerrors.ErrConnectClosed
	}

	select {
	case <-c._ctx.Done():
		return connerrors.ErrConnectClosed
	case c._queue <- msg:
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
	}

	return err
}

func (c *Client) Shutdown() {
	c.Close()
	if c._queue != nil {
		close(c._queue)
		c._queue = nil
	}
}
