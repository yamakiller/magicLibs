package wss

import (
	"context"
	"crypto/tls"
	"errors"
	"io"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
	"github.com/yamakiller/magicLibs/net/connection/connerrors"
	"github.com/yamakiller/magicLibs/net/connection/interfaces"
)

func New(writeBufferQueueSize int) *Client {
	return &Client{_queue: make(chan interface{}, writeBufferQueueSize)}
}

type Client struct {
	Serializer interfaces.Serialization
	Exception  interfaces.Exception
	T          int

	_conn       *websocket.Conn
	_cancel     context.CancelFunc
	_ctx        context.Context
	_queue      chan interface{}
	_wbs        uint64
	_rbs        uint64
	_lastActive int64
	_connectd   bool
	_wg         sync.WaitGroup
}

//Connect 连接服务器
func (c *Client) Connect(url string, timeout time.Duration) error {
	var d *websocket.Dialer
	if timeout > 0 {
		d.HandshakeTimeout = timeout
	}
	conn, _, err := d.Dial(url, nil)

	if err != nil {
		return nil
	}

	//c._queue = make(chan interface{}, c.WriteWaitQueue)
	c._conn = conn
	c._ctx, c._cancel = context.WithCancel(context.Background())

	c._wg.Add(1)
	go c.writeServe()

	return nil
}

func (c *Client) ConnectTls(addr string, timeout time.Duration, config *tls.Config) error {
	return nil
}

func (c *Client) writeServe() {
	defer func() {
		c._wg.Done()
	}()

	for {
	active:
		select {
		case <-c._ctx.Done():
			goto exit
		case msg := <-c._queue:
			w, err := c._conn.NextWriter(c.T)
			if err != nil {
				if c.Exception != nil {
					c.Exception.Error(err)
				}
			}
			defer w.Close()

			n, err := c.Serializer.Seria(msg, w)
			if err != nil {
				if c.Exception != nil {
					c.Exception.Error(err)
				}

				goto active
			}

			atomic.AddUint64(&c._wbs, uint64(n))
		}
	}
exit:
}

//Parse 解析数据
func (c *Client) Parse() (interface{}, error) {
	c._wg.Add(1)
	defer c._wg.Done()

	var (
		err error
		r   io.Reader
		t   int
		m   interface{}
		n   int
	)

	select {
	case <-c._ctx.Done():
		return nil, connerrors.ErrConnectClosed
	default:
		t, r, err = c._conn.NextReader()
		if err != nil {
			return nil, err
		}

		if t != c.T {
			return nil, errors.New("data type mismatch")
		}

		m, n, err = c.Serializer.UnSeria(r)
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
