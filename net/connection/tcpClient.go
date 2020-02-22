package connection

import (
	"bufio"
	"errors"
	"io"
	"net"
	"sync"
	"time"
)

//TCPClient tcp 客户端
type TCPClient struct {
	ReadBufferSize  int
	WriteBufferSize int
	WriteWaitQueue  int
	S               Serialization
	E               Exception

	_c          io.ReadWriteCloser
	_closed     chan bool
	_queue      chan interface{}
	_reader     *bufio.Reader
	_writer     *bufio.Writer
	_wTotal     int
	_rTotal     int
	_lastActive int64
	_wwg        sync.WaitGroup
}

//Connect 连接远程地址
func (slf *TCPClient) Connect(addr string, timeout time.Duration) error {
	var err error
	var c net.Conn

	if timeout == 0 {
		tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
		if err != nil {
			return err
		}

		c, err = net.DialTCP("tcp", nil, tcpAddr)
	} else {
		c, err = net.DialTimeout("tcp", addr, timeout)
	}
	if err != nil {
		return err
	}

	slf._c = c.(io.ReadWriteCloser)
	slf._reader = bufio.NewReaderSize(slf._c, slf.ReadBufferSize)
	slf._writer = bufio.NewWriterSize(slf._c, slf.WriteBufferSize)
	slf._queue = make(chan interface{}, slf.WriteWaitQueue)

	if slf._closed == nil {
		slf._closed = make(chan bool, 1)
	}

	slf._wwg.Add(1)
	go slf.writeServe()

	return nil
}

func (slf *TCPClient) writeServe() {
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
			n, err := slf.S.Seria(msg, slf._writer)
			if err != nil {
				if slf.E != nil {
					slf.E.Error(err)
				}

				goto active
			}

			if slf._writer.Buffered() > 0 {
				if err := slf._writer.Flush(); err != nil &&
					slf.E != nil {
					slf.E.Error(err)
				}
			}
			slf._wTotal += n
		}
	}
exit:
}

//Parse 解析数据
func (slf *TCPClient) Parse() (interface{}, error) {
	m, n, err := slf.S.UnSeria(slf._reader)
	if err != nil {
		return nil, err
	}

	slf._rTotal += n
	return m, nil
}

//SendTo 发送数据
func (slf *TCPClient) SendTo(msg interface{}) error {
	select {
	case <-slf._closed:
		return errors.New("closed")
	default:
	}

	slf._queue <- msg
	return nil
}

//Close 关闭连接
func (slf *TCPClient) Close() error {
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
