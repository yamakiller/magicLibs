package connection

import (
	"errors"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

//WSSClient Websocket Client
type WSSClient struct {
	WriteWaitQueue int
	S              Serialization
	E              Exception
	T              int

	_c          *websocket.Conn
	_closed     chan bool
	_queue      chan interface{}
	_wTotal     int
	_rTotal     int
	_lastActive int64
	_wwg        sync.WaitGroup
}

//Connect 连接服务器
func (slf *WSSClient) Connect(url string, timeout time.Duration) error {
	var d *websocket.Dialer
	if timeout > 0 {
		d.HandshakeTimeout = timeout
	}
	c, _, err := d.Dial(url, nil)

	if err != nil {
		return nil
	}

	slf._queue = make(chan interface{}, slf.WriteWaitQueue)
	if slf._closed == nil {
		slf._closed = make(chan bool, 1)
	}

	slf._c = c

	slf._wwg.Add(1)
	go slf.writeServe()

	return nil
}

func (slf *WSSClient) writeServe() {
	defer func() {
		close(slf._queue)
		slf._wwg.Done()
	}()

	for {
	active:
		select {
		case <-slf._closed:
			goto exit
		case msg := <-slf._queue:
			w, err := slf._c.NextWriter(slf.T)
			if err != nil {
				if slf.E != nil {
					slf.E.Error(err)
				}
			}
			defer w.Close()

			n, err := slf.S.Seria(msg, w)
			if err != nil {
				if slf.E != nil {
					slf.E.Error(err)
				}

				goto active
			}

			slf._wTotal += n
		}
	}
exit:
}

//Parse 解析数据
func (slf *WSSClient) Parse() (interface{}, error) {
	t, r, err := slf._c.NextReader()
	if err != nil {
		return nil, err
	}

	if t != slf.T {
		return nil, errors.New("data type mismatch")
	}

	m, n, err := slf.S.UnSeria(r)
	if err != nil {
		return nil, err
	}

	slf._rTotal += n
	return m, nil
}

//SendTo 发送数据
func (slf *WSSClient) SendTo(msg interface{}) error {
	select {
	case <-slf._closed:
		return errors.New("closed")
	default:
	}
	slf._queue <- msg
	return nil
}

//Close 关闭连接
func (slf *WSSClient) Close() error {
	if slf._closed != nil {
		select {
		case <-slf._closed:
		default:
			close(slf._closed)
		}
	}
	err := slf._c.Close()
	slf._wwg.Wait()

	return err
}
